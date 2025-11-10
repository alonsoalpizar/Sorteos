package raffle

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// CreateRaffleInput representa los datos de entrada para crear un sorteo
type CreateRaffleInput struct {
	UserID int64

	// Basic info
	Title       string
	Description string

	// Pricing
	PricePerNumber decimal.Decimal
	TotalNumbers   int
	MinNumber      int
	MaxNumber      int

	// Draw info
	DrawDate   time.Time
	DrawMethod domain.DrawMethod

	// Platform fee (opcional, usa default 10%)
	PlatformFeePercentage *decimal.Decimal
}

// CreateRaffleOutput representa el resultado de crear un sorteo
type CreateRaffleOutput struct {
	Raffle  *domain.Raffle
	Numbers []*domain.RaffleNumber
}

// CreateRaffleUseCase maneja la lógica de crear un nuevo sorteo
type CreateRaffleUseCase struct {
	raffleRepo       db.RaffleRepository
	raffleNumberRepo db.RaffleNumberRepository
	userRepo         domain.UserRepository
	auditRepo        domain.AuditLogRepository
	logger           *logger.Logger
}

// NewCreateRaffleUseCase crea una nueva instancia del caso de uso
func NewCreateRaffleUseCase(
	raffleRepo db.RaffleRepository,
	raffleNumberRepo db.RaffleNumberRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	logger *logger.Logger,
) *CreateRaffleUseCase {
	return &CreateRaffleUseCase{
		raffleRepo:       raffleRepo,
		raffleNumberRepo: raffleNumberRepo,
		userRepo:         userRepo,
		auditRepo:        auditRepo,
		logger:           logger,
	}
}

// Execute ejecuta el caso de uso de crear un sorteo
func (uc *CreateRaffleUseCase) Execute(ctx context.Context, input *CreateRaffleInput) (*CreateRaffleOutput, error) {
	// Validar que el usuario existe y tiene permisos
	user, err := uc.userRepo.FindByID(input.UserID)
	if err != nil {
		uc.logger.Error("Usuario no encontrado", logger.Int64("user_id", input.UserID), logger.Error(err))
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Verificar que el usuario tenga al menos email verificado
	if user.KYCLevel == domain.KYCLevelNone {
		return nil, fmt.Errorf("debes verificar tu email antes de crear un sorteo")
	}

	// Verificar que el usuario esté activo
	if !user.IsActive() {
		return nil, fmt.Errorf("tu cuenta no está activa")
	}

	// Validar inputs básicos
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Crear entidad Raffle
	raffle := domain.NewRaffle(
		input.UserID,
		input.Title,
		input.PricePerNumber,
		input.TotalNumbers,
		input.DrawDate,
	)

	// Aplicar valores opcionales
	if input.Description != "" {
		raffle.Description = input.Description
	}

	if input.MinNumber != 0 || input.MaxNumber != 0 {
		raffle.MinNumber = input.MinNumber
		raffle.MaxNumber = input.MaxNumber
	}

	if input.DrawMethod != "" {
		raffle.DrawMethod = input.DrawMethod
	}

	if input.PlatformFeePercentage != nil {
		raffle.PlatformFeePercentage = *input.PlatformFeePercentage
	}

	// Validar la entidad
	if err := raffle.Validate(); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Guardar el sorteo
	if err := uc.raffleRepo.Create(raffle); err != nil {
		uc.logger.Error("Error creando sorteo", logger.Error(err))
		return nil, fmt.Errorf("error al crear el sorteo")
	}

	// Generar números del sorteo
	numbers, err := uc.generateNumbers(raffle)
	if err != nil {
		uc.logger.Error("Error generando números", logger.Int64("raffle_id", raffle.ID), logger.Error(err))
		return nil, fmt.Errorf("error al generar los números")
	}

	// Guardar números en batch
	if err := uc.raffleNumberRepo.CreateBatch(numbers); err != nil {
		uc.logger.Error("Error guardando números", logger.Int64("raffle_id", raffle.ID), logger.Error(err))
		return nil, fmt.Errorf("error al guardar los números")
	}

	// Registrar en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionRaffleCreated).
		WithUser(user.ID).
		WithEntity("raffle", raffle.ID).
		WithDescription(fmt.Sprintf("Sorteo creado: %s", raffle.Title)).
		WithMetadata(map[string]interface{}{
			"title":         raffle.Title,
			"total_numbers": raffle.TotalNumbers,
			"price":         raffle.PricePerNumber.String(),
		}).
		Build()

	uc.auditRepo.Create(auditLog)

	uc.logger.Info("Sorteo creado exitosamente",
		logger.Int64("raffle_id", raffle.ID),
		logger.Int64("user_id", user.ID),
		logger.String("title", raffle.Title),
	)

	return &CreateRaffleOutput{
		Raffle:  raffle,
		Numbers: numbers,
	}, nil
}

// validateInput valida los inputs del caso de uso
func (uc *CreateRaffleUseCase) validateInput(input *CreateRaffleInput) error {
	if input.UserID <= 0 {
		return fmt.Errorf("user_id es requerido")
	}

	if input.Title == "" {
		return fmt.Errorf("el título es requerido")
	}

	if len(input.Title) < 5 {
		return fmt.Errorf("el título debe tener al menos 5 caracteres")
	}

	if input.PricePerNumber.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("el precio debe ser mayor a 0")
	}

	if input.TotalNumbers <= 0 {
		return fmt.Errorf("el total de números debe ser mayor a 0")
	}

	if input.TotalNumbers > 10000 {
		return fmt.Errorf("el total de números no puede exceder 10,000")
	}

	if input.DrawDate.Before(time.Now()) {
		return fmt.Errorf("la fecha de sorteo debe ser futura")
	}

	// Si se especifican min/max, validar que sean coherentes
	if input.MinNumber != 0 || input.MaxNumber != 0 {
		if input.MaxNumber <= input.MinNumber {
			return fmt.Errorf("el número máximo debe ser mayor al mínimo")
		}

		if (input.MaxNumber - input.MinNumber + 1) != input.TotalNumbers {
			return fmt.Errorf("el rango de números no coincide con el total")
		}
	}

	return nil
}

// generateNumbers genera los números del sorteo
func (uc *CreateRaffleUseCase) generateNumbers(raffle *domain.Raffle) ([]*domain.RaffleNumber, error) {
	numbers := make([]*domain.RaffleNumber, 0, raffle.TotalNumbers)

	// Generar números del min al max
	for i := raffle.MinNumber; i <= raffle.MaxNumber; i++ {
		// Formatear número con ceros a la izquierda si es necesario
		var numberStr string
		if raffle.MaxNumber < 100 {
			numberStr = fmt.Sprintf("%02d", i) // 00-99
		} else if raffle.MaxNumber < 1000 {
			numberStr = fmt.Sprintf("%03d", i) // 000-999
		} else {
			numberStr = fmt.Sprintf("%04d", i) // 0000-9999
		}

		number := domain.NewRaffleNumber(raffle.ID, numberStr)
		numbers = append(numbers, number)
	}

	if len(numbers) != raffle.TotalNumbers {
		return nil, fmt.Errorf("error generando números: esperados %d, generados %d", raffle.TotalNumbers, len(numbers))
	}

	return numbers, nil
}
