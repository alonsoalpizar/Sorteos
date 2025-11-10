package raffle

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// UpdateRaffleInput datos de entrada
type UpdateRaffleInput struct {
	RaffleID    int64
	UserID      int64
	UserRole    domain.UserRole
	Title       *string
	Description *string
	DrawDate    *time.Time
	DrawMethod  *domain.DrawMethod
}

// UpdateRaffleOutput resultado de la actualización
type UpdateRaffleOutput struct {
	Raffle *domain.Raffle
}

// UpdateRaffleUseCase caso de uso para actualizar un sorteo
type UpdateRaffleUseCase struct {
	raffleRepo db.RaffleRepository
	auditRepo  domain.AuditLogRepository
}

// NewUpdateRaffleUseCase crea una nueva instancia
func NewUpdateRaffleUseCase(
	raffleRepo db.RaffleRepository,
	auditRepo domain.AuditLogRepository,
) *UpdateRaffleUseCase {
	return &UpdateRaffleUseCase{
		raffleRepo: raffleRepo,
		auditRepo:  auditRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateRaffleUseCase) Execute(ctx context.Context, input *UpdateRaffleInput) (*UpdateRaffleOutput, error) {
	// 1. Buscar el sorteo
	raffle, err := uc.raffleRepo.FindByID(input.RaffleID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 2. Verificar permisos (owner o admin)
	if raffle.UserID != input.UserID && input.UserRole != domain.UserRoleAdmin {
		return nil, errors.ErrForbidden
	}

	// 3. No permitir edición si está completado o cancelado
	if raffle.Status == domain.RaffleStatusCompleted || raffle.Status == domain.RaffleStatusCancelled {
		return nil, errors.New("RAFFLE_CLOSED", "No se puede editar un sorteo completado o cancelado", 400, nil)
	}

	// 4. Si está activo y tiene ventas, solo permitir cambios limitados
	if raffle.Status == domain.RaffleStatusActive && raffle.SoldCount > 0 {
		// Solo permitir cambiar descripción y fecha (si es futura)
		if input.Title != nil {
			return nil, errors.New("CANNOT_CHANGE_TITLE", "No se puede cambiar el título de un sorteo con ventas", 400, nil)
		}
	}

	// 5. Aplicar cambios
	if input.Title != nil && *input.Title != "" {
		raffle.Title = *input.Title
	}

	if input.Description != nil {
		raffle.Description = *input.Description
	}

	if input.DrawDate != nil {
		// Validar que sea futura
		if input.DrawDate.Before(time.Now()) {
			return nil, errors.New("INVALID_DRAW_DATE", "La fecha del sorteo debe ser en el futuro", 400, nil)
		}
		raffle.DrawDate = *input.DrawDate
	}

	if input.DrawMethod != nil {
		raffle.DrawMethod = *input.DrawMethod
	}

	// 6. Validar
	if err := raffle.Validate(); err != nil {
		return nil, errors.New("VALIDATION_FAILED", err.Error(), 400, nil)
	}

	// 7. Guardar
	if err := uc.raffleRepo.Update(raffle); err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 8. Audit log
	newValues := map[string]interface{}{
		"title":       raffle.Title,
		"description": raffle.Description,
		"draw_date":   raffle.DrawDate,
		"draw_method": raffle.DrawMethod,
	}

	auditLog := domain.NewAuditLog(domain.AuditActionRaffleCreated). // Will use a generic action
		WithUser(input.UserID).
		WithEntity("raffle", raffle.ID).
		WithDescription("Sorteo actualizado").
		WithMetadata(newValues).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		fmt.Printf("Error creating audit log: %v\n", err)
	}

	return &UpdateRaffleOutput{
		Raffle: raffle,
	}, nil
}

// SuspendRaffleInput datos de entrada para suspender sorteo
type SuspendRaffleInput struct {
	RaffleID int64
	UserID   int64
	UserRole domain.UserRole
	Reason   string
}

// SuspendRaffleOutput resultado de la suspensión
type SuspendRaffleOutput struct {
	Raffle *domain.Raffle
}

// SuspendRaffleUseCase caso de uso para suspender un sorteo (admin only)
type SuspendRaffleUseCase struct {
	raffleRepo db.RaffleRepository
	auditRepo  domain.AuditLogRepository
}

// NewSuspendRaffleUseCase crea una nueva instancia
func NewSuspendRaffleUseCase(
	raffleRepo db.RaffleRepository,
	auditRepo domain.AuditLogRepository,
) *SuspendRaffleUseCase {
	return &SuspendRaffleUseCase{
		raffleRepo: raffleRepo,
		auditRepo:  auditRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *SuspendRaffleUseCase) Execute(ctx context.Context, input *SuspendRaffleInput) (*SuspendRaffleOutput, error) {
	// 1. Verificar que sea admin
	if input.UserRole != domain.UserRoleAdmin {
		return nil, errors.ErrForbidden
	}

	// 2. Buscar el sorteo
	raffle, err := uc.raffleRepo.FindByID(input.RaffleID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 3. Suspender
	if err := raffle.Suspend(); err != nil {
		return nil, errors.New("SUSPEND_FAILED", err.Error(), 400, nil)
	}

	// 4. Guardar
	if err := uc.raffleRepo.Update(raffle); err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 5. Audit log
	auditLog := domain.NewAuditLog(domain.AuditActionRaffleSuspended).
		WithUser(input.UserID).
		WithEntity("raffle", raffle.ID).
		WithDescription(fmt.Sprintf("Sorteo suspendido: %s", input.Reason)).
		WithMetadata(map[string]interface{}{
			"status": string(raffle.Status),
			"reason": input.Reason,
		}).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		fmt.Printf("Error creating audit log: %v\n", err)
	}

	return &SuspendRaffleOutput{
		Raffle: raffle,
	}, nil
}

// DeleteRaffleInput datos de entrada para eliminar sorteo
type DeleteRaffleInput struct {
	RaffleID int64
	UserID   int64
	UserRole domain.UserRole
}

// DeleteRaffleUseCase caso de uso para eliminar un sorteo (soft delete)
type DeleteRaffleUseCase struct {
	raffleRepo db.RaffleRepository
	auditRepo  domain.AuditLogRepository
}

// NewDeleteRaffleUseCase crea una nueva instancia
func NewDeleteRaffleUseCase(
	raffleRepo db.RaffleRepository,
	auditRepo domain.AuditLogRepository,
) *DeleteRaffleUseCase {
	return &DeleteRaffleUseCase{
		raffleRepo: raffleRepo,
		auditRepo:  auditRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *DeleteRaffleUseCase) Execute(ctx context.Context, input *DeleteRaffleInput) error {
	// 1. Buscar el sorteo
	raffle, err := uc.raffleRepo.FindByID(input.RaffleID)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrRaffleNotFound
		}
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 2. Verificar permisos (owner o admin)
	if raffle.UserID != input.UserID && input.UserRole != domain.UserRoleAdmin {
		return errors.ErrForbidden
	}

	// 3. No permitir eliminar si tiene ventas
	if raffle.SoldCount > 0 {
		return errors.New("HAS_SALES", "No se puede eliminar un sorteo con números vendidos", 400, nil)
	}

	// 4. Solo permitir eliminar si está en draft o suspended
	if raffle.Status != domain.RaffleStatusDraft && raffle.Status != domain.RaffleStatusSuspended {
		return errors.New("INVALID_STATUS", "Solo se pueden eliminar sorteos en estado borrador o suspendidos", 400, nil)
	}

	// 5. Soft delete
	if err := uc.raffleRepo.SoftDelete(raffle.ID); err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 6. Audit log
	auditLog := domain.NewAuditLog(domain.AuditActionRaffleDeleted).
		WithUser(input.UserID).
		WithEntity("raffle", raffle.ID).
		WithDescription(fmt.Sprintf("Sorteo eliminado: %s", raffle.Title)).
		WithMetadata(map[string]interface{}{
			"title":      raffle.Title,
			"status":     string(raffle.Status),
			"sold_count": raffle.SoldCount,
		}).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		fmt.Printf("Error creating audit log: %v\n", err)
	}

	return nil
}

// Helper para evitar warning de variable no usada
var _ = decimal.Decimal{}
