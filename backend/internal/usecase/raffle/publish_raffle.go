package raffle

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// PublishRaffleInput datos de entrada
type PublishRaffleInput struct {
	RaffleID int64
	UserID   int64
}

// PublishRaffleOutput resultado de la publicación
type PublishRaffleOutput struct {
	Raffle *domain.Raffle
}

// PublishRaffleUseCase caso de uso para publicar un sorteo
type PublishRaffleUseCase struct {
	raffleRepo       db.RaffleRepository
	raffleImageRepo  db.RaffleImageRepository
	raffleNumberRepo db.RaffleNumberRepository
	auditRepo        domain.AuditLogRepository
}

// NewPublishRaffleUseCase crea una nueva instancia
func NewPublishRaffleUseCase(
	raffleRepo db.RaffleRepository,
	raffleImageRepo db.RaffleImageRepository,
	raffleNumberRepo db.RaffleNumberRepository,
	auditRepo domain.AuditLogRepository,
) *PublishRaffleUseCase {
	return &PublishRaffleUseCase{
		raffleRepo:       raffleRepo,
		raffleImageRepo:  raffleImageRepo,
		raffleNumberRepo: raffleNumberRepo,
		auditRepo:        auditRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *PublishRaffleUseCase) Execute(ctx context.Context, input *PublishRaffleInput) (*PublishRaffleOutput, error) {
	// 1. Buscar el sorteo
	raffle, err := uc.raffleRepo.FindByID(input.RaffleID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 2. Verificar que el usuario sea el owner
	if raffle.UserID != input.UserID {
		return nil, errors.ErrForbidden
	}

	// 3. Validar que esté en estado draft
	if raffle.Status != domain.RaffleStatusDraft {
		return nil, errors.New("RAFFLE_NOT_DRAFT", "Solo se pueden publicar sorteos en estado borrador", 400, nil)
	}

	// 4. Verificar que tenga al menos una imagen
	// TODO: TEMPORARILY DISABLED - Image upload not implemented yet (Sprint 4 pending)
	// This validation will be re-enabled once image upload functionality is complete
	// See: roadmap.md Sprint 4 - Image upload implementation
	/*
	imageCount, err := uc.raffleImageRepo.CountByRaffleID(raffle.ID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	if imageCount == 0 {
		return nil, errors.New("NO_IMAGES", "El sorteo debe tener al menos una imagen", 400, nil)
	}

	// 5. Verificar que tenga una imagen principal
	_, err = uc.raffleImageRepo.FindPrimaryByRaffleID(raffle.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.New("NO_PRIMARY_IMAGE", "El sorteo debe tener una imagen principal", 400, nil)
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	*/

	// 6. Verificar que la fecha del sorteo sea futura
	if raffle.DrawDate.Before(time.Now()) {
		return nil, errors.New("INVALID_DRAW_DATE", "La fecha del sorteo debe ser en el futuro", 400, nil)
	}

	// 7. Verificar que tenga números generados
	numberCount, err := uc.raffleNumberRepo.CountByStatus(raffle.ID, domain.RaffleNumberStatusAvailable)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	if numberCount == 0 {
		return nil, errors.New("NO_NUMBERS", "El sorteo no tiene números generados", 500, nil)
	}

	// 8. Publicar el sorteo
	if err := raffle.Publish(); err != nil {
		return nil, errors.New("PUBLISH_FAILED", err.Error(), 400, nil)
	}

	// 9. Guardar cambios
	if err := uc.raffleRepo.Update(raffle); err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 10. Registrar en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionRafflePublished).
		WithUser(input.UserID).
		WithEntity("raffle", raffle.ID).
		WithDescription(fmt.Sprintf("Sorteo publicado: %s", raffle.Title)).
		WithMetadata(map[string]interface{}{
			"status":       string(raffle.Status),
			"published_at": raffle.PublishedAt,
		}).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		// No falla si el audit log falla, solo registra el error
		fmt.Printf("Error creating audit log: %v\n", err)
	}

	return &PublishRaffleOutput{
		Raffle: raffle,
	}, nil
}
