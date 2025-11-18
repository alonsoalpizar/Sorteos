package user

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ListUsersInput datos de entrada para listar usuarios
type ListUsersInput struct {
	Page      int
	PageSize  int
	Role      *domain.UserRole
	Status    *string // 'active', 'suspended', 'banned'
	KYCLevel  *domain.KYCLevel
	Search    string // Buscar en name, email, cedula
	OrderBy   string
	DateFrom  *string
	DateTo    *string
}

// ListUsersOutput resultado del listado
type ListUsersOutput struct {
	Users      []*domain.User
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListUsersUseCase caso de uso para listar usuarios (admin)
type ListUsersUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListUsersUseCase crea una nueva instancia
func NewListUsersUseCase(db *gorm.DB, log *logger.Logger) *ListUsersUseCase {
	return &ListUsersUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListUsersUseCase) Execute(ctx context.Context, input *ListUsersInput, adminID int64) (*ListUsersOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Construir query
	query := uc.db.Model(&domain.User{})

	// Aplicar filtros
	if input.Role != nil {
		query = query.Where("role = ?", *input.Role)
	}

	if input.Status != nil {
		switch *input.Status {
		case "active":
			query = query.Where("is_active = ? AND suspended_at IS NULL", true)
		case "suspended":
			query = query.Where("suspended_at IS NOT NULL")
		case "banned":
			query = query.Where("is_active = ?", false)
		}
	}

	if input.KYCLevel != nil {
		query = query.Where("kyc_level = ?", *input.KYCLevel)
	}

	if input.Search != "" {
		searchPattern := "%" + input.Search + "%"
		query = query.Where(
			"name ILIKE ? OR email ILIKE ? OR cedula ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}

	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting users", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Aplicar ordenamiento
	orderBy := "created_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Obtener usuarios con paginación
	var users []*domain.User
	if err := query.Offset(offset).Limit(input.PageSize).Find(&users).Error; err != nil {
		uc.log.Error("Error listing users", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Log auditoría
	uc.log.Info("Admin listed users",
		logger.Int64("admin_id", adminID),
		logger.Int("total_results", len(users)),
		logger.String("action", "admin_list_users"))

	return &ListUsersOutput{
		Users:      users,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
