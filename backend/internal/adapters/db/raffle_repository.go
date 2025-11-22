package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// RaffleRepository define los métodos de acceso a datos para Raffle
type RaffleRepository interface {
	Create(raffle *domain.Raffle) error
	Update(raffle *domain.Raffle) error
	Delete(id int64) error
	SoftDelete(id int64) error
	FindByID(id int64) (*domain.Raffle, error)
	FindByUUID(uuid string) (*domain.Raffle, error)
	FindByUserID(userID int64, offset, limit int) ([]*domain.Raffle, int64, error)
	List(offset, limit int, filters map[string]interface{}) ([]*domain.Raffle, int64, error)
	ListActive(offset, limit int) ([]*domain.Raffle, int64, error)
	CountByStatus(status domain.RaffleStatus) (int64, error)
	UpdateStatus(id int64, status domain.RaffleStatus) error
	SetWinner(id int64, winnerNumber string, winnerUserID *int64) error
	IncrementSoldCount(id int64) error
	DecrementSoldCount(id int64) error

	// Earnings methods
	GetUserEarningsSummary(userID int64) (*domain.UserEarnings, error)
	GetUserCompletedRaffles(userID int64, limit, offset int) ([]domain.RaffleEarning, error)
}

// RaffleRepositoryImpl implementa RaffleRepository
type RaffleRepositoryImpl struct {
	db *gorm.DB
}

// NewRaffleRepository crea una nueva instancia del repositorio
func NewRaffleRepository(db *gorm.DB) RaffleRepository {
	return &RaffleRepositoryImpl{db: db}
}

// Create crea un nuevo sorteo
func (r *RaffleRepositoryImpl) Create(raffle *domain.Raffle) error {
	if err := r.db.Create(raffle).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Update actualiza un sorteo existente
func (r *RaffleRepositoryImpl) Update(raffle *domain.Raffle) error {
	if err := r.db.Save(raffle).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Delete elimina físicamente un sorteo
func (r *RaffleRepositoryImpl) Delete(id int64) error {
	if err := r.db.Delete(&domain.Raffle{}, id).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// SoftDelete elimina lógicamente un sorteo
func (r *RaffleRepositoryImpl) SoftDelete(id int64) error {
	now := time.Now()
	if err := r.db.Model(&domain.Raffle{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByID busca un sorteo por ID
func (r *RaffleRepositoryImpl) FindByID(id int64) (*domain.Raffle, error) {
	var raffle domain.Raffle
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &raffle, nil
}

// FindByUUID busca un sorteo por UUID
func (r *RaffleRepositoryImpl) FindByUUID(uuid string) (*domain.Raffle, error) {
	var raffle domain.Raffle
	if err := r.db.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &raffle, nil
}

// FindByUserID busca sorteos de un usuario específico
func (r *RaffleRepositoryImpl) FindByUserID(userID int64, offset, limit int) ([]*domain.Raffle, int64, error) {
	var raffles []*domain.Raffle
	var total int64

	query := r.db.Model(&domain.Raffle{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&raffles).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return raffles, total, nil
}

// List retorna una lista paginada de sorteos con filtros
func (r *RaffleRepositoryImpl) List(offset, limit int, filters map[string]interface{}) ([]*domain.Raffle, int64, error) {
	var raffles []*domain.Raffle
	var total int64

	query := r.db.Model(&domain.Raffle{}).Where("deleted_at IS NULL")

	// Aplicar filtros
	if status, ok := filters["status"].(domain.RaffleStatus); ok {
		query = query.Where("status = ?", status)
	}

	if userID, ok := filters["user_id"].(int64); ok {
		query = query.Where("user_id = ?", userID)
	}

	if categoryID, ok := filters["category_id"].(int64); ok {
		query = query.Where("category_id = ?", categoryID)
	}

	if drawMethod, ok := filters["draw_method"].(domain.DrawMethod); ok {
		query = query.Where("draw_method = ?", drawMethod)
	}

	// Filtro por fecha de sorteo
	if drawDateFrom, ok := filters["draw_date_from"].(time.Time); ok {
		query = query.Where("draw_date >= ?", drawDateFrom)
	}

	if drawDateTo, ok := filters["draw_date_to"].(time.Time); ok {
		query = query.Where("draw_date <= ?", drawDateTo)
	}

	// Filtro por texto (búsqueda en título)
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	// Filtro para excluir sorteos de un usuario específico (para /explore)
	if excludeUserID, ok := filters["exclude_user_id"].(int64); ok {
		query = query.Where("user_id != ?", excludeUserID)
	}

	// Filtro por disponibilidad
	if onlyAvailable, ok := filters["only_available"].(bool); ok && onlyAvailable {
		query = query.Where("sold_count < total_numbers")
	}

	// Ordenamiento
	orderBy := "created_at DESC"
	if order, ok := filters["order_by"].(string); ok {
		orderBy = order
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener página
	if err := query.Order(orderBy).Offset(offset).Limit(limit).Find(&raffles).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return raffles, total, nil
}

// ListActive retorna sorteos activos paginados
func (r *RaffleRepositoryImpl) ListActive(offset, limit int) ([]*domain.Raffle, int64, error) {
	var raffles []*domain.Raffle
	var total int64

	query := r.db.Model(&domain.Raffle{}).
		Where("status = ? AND deleted_at IS NULL AND draw_date > ?",
			domain.RaffleStatusActive, time.Now())

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("draw_date ASC").Offset(offset).Limit(limit).Find(&raffles).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return raffles, total, nil
}

// CountByStatus cuenta sorteos por estado
func (r *RaffleRepositoryImpl) CountByStatus(status domain.RaffleStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Raffle{}).
		Where("status = ? AND deleted_at IS NULL", status).
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return count, nil
}

// UpdateStatus actualiza solo el estado del sorteo
func (r *RaffleRepositoryImpl) UpdateStatus(id int64, status domain.RaffleStatus) error {
	if err := r.db.Model(&domain.Raffle{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// SetWinner establece el ganador del sorteo
func (r *RaffleRepositoryImpl) SetWinner(id int64, winnerNumber string, winnerUserID *int64) error {
	now := time.Now()
	updates := map[string]interface{}{
		"winner_number":   winnerNumber,
		"winner_user_id":  winnerUserID,
		"status":          domain.RaffleStatusCompleted,
		"completed_at":    now,
		"updated_at":      now,
	}

	if err := r.db.Model(&domain.Raffle{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// IncrementSoldCount incrementa el contador de vendidos
func (r *RaffleRepositoryImpl) IncrementSoldCount(id int64) error {
	if err := r.db.Model(&domain.Raffle{}).
		Where("id = ?", id).
		UpdateColumn("sold_count", gorm.Expr("sold_count + ?", 1)).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// DecrementSoldCount decrementa el contador de vendidos
func (r *RaffleRepositoryImpl) DecrementSoldCount(id int64) error {
	if err := r.db.Model(&domain.Raffle{}).
		Where("id = ? AND sold_count > 0", id).
		UpdateColumn("sold_count", gorm.Expr("sold_count - ?", 1)).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// GetUserEarningsSummary obtiene el resumen total de ganancias de un usuario
func (r *RaffleRepositoryImpl) GetUserEarningsSummary(userID int64) (*domain.UserEarnings, error) {
	type Summary struct {
		TotalCollected     float64
		PlatformCommission float64
		NetEarnings        float64
		CompletedRaffles   int64
	}

	var summary Summary
	err := r.db.Model(&domain.Raffle{}).
		Select(`
			COALESCE(SUM(total_revenue), 0) as total_collected,
			COALESCE(SUM(platform_fee_amount), 0) as platform_commission,
			COALESCE(SUM(net_amount), 0) as net_earnings,
			COUNT(*) as completed_raffles
		`).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL",
			userID, domain.RaffleStatusActive).
		Scan(&summary).Error

	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Convertir a UserEarnings domain
	earnings := &domain.UserEarnings{
		TotalCollected:     domain.NewDecimalFromFloat(summary.TotalCollected),
		PlatformCommission: domain.NewDecimalFromFloat(summary.PlatformCommission),
		NetEarnings:        domain.NewDecimalFromFloat(summary.NetEarnings),
		CompletedRaffles:   int(summary.CompletedRaffles),
		Raffles:            []domain.RaffleEarning{},
	}

	return earnings, nil
}

// GetUserCompletedRaffles obtiene los sorteos completados con desglose
func (r *RaffleRepositoryImpl) GetUserCompletedRaffles(userID int64, limit, offset int) ([]domain.RaffleEarning, error) {
	var raffles []domain.Raffle

	query := r.db.Where("user_id = ? AND status = ? AND deleted_at IS NULL",
		userID, domain.RaffleStatusActive).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&raffles).Error; err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Convertir a RaffleEarning
	earnings := make([]domain.RaffleEarning, 0, len(raffles))
	for _, r := range raffles {
		earnings = append(earnings, domain.RaffleEarning{
			RaffleID:           r.ID,
			RaffleUUID:         r.UUID.String(),
			Title:              r.Title,
			DrawDate:           r.DrawDate,
			CompletedAt:        r.CompletedAt,
			TotalRevenue:       r.TotalRevenue,
			PlatformFeePercent: r.PlatformFeePercentage,
			PlatformFeeAmount:  r.PlatformFeeAmount,
			NetAmount:          r.NetAmount,
			SettlementStatus:   string(r.SettlementStatus),
			SettledAt:          r.SettledAt,
		})
	}

	return earnings, nil
}
