package reports

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ExportDataInput datos de entrada
type ExportDataInput struct {
	EntityType string  // users, raffles, payments, settlements
	Format     string  // csv, xlsx, pdf (por ahora solo CSV)
	DateFrom   *string
	DateTo     *string
	Filters    map[string]interface{} // Filtros adicionales según entidad
}

// ExportDataOutput resultado
type ExportDataOutput struct {
	FilePath     string `json:"file_path"`
	FileName     string `json:"file_name"`
	DownloadURL  string `json:"download_url"`
	RecordCount  int    `json:"record_count"`
	FileSize     int64  `json:"file_size"`
	ExpiresAt    string `json:"expires_at"` // Timestamp de expiración
}

// ExportDataUseCase caso de uso para exportar datos
type ExportDataUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewExportDataUseCase crea una nueva instancia
func NewExportDataUseCase(db *gorm.DB, log *logger.Logger) *ExportDataUseCase {
	return &ExportDataUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ExportDataUseCase) Execute(ctx context.Context, input *ExportDataInput, adminID int64) (*ExportDataOutput, error) {
	// Validar entity_type
	validEntities := map[string]bool{
		"users":       true,
		"raffles":     true,
		"payments":    true,
		"settlements": true,
		"audit_logs":  true,
	}

	if !validEntities[input.EntityType] {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("invalid entity_type: %s", input.EntityType), 400, nil)
	}

	// Validar format
	if input.Format == "" {
		input.Format = "csv"
	}
	if input.Format != "csv" {
		// TODO: Implementar xlsx y pdf
		return nil, errors.New("VALIDATION_FAILED",
			"only CSV format is currently supported", 400, nil)
	}

	// Crear directorio temporal si no existe
	exportDir := "/tmp/exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		uc.log.Error("Error creating export directory", logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
	}

	// Generar nombre de archivo único
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_export_%s.csv", input.EntityType, timestamp)
	filePath := filepath.Join(exportDir, fileName)

	// Crear archivo
	file, err := os.Create(filePath)
	if err != nil {
		uc.log.Error("Error creating export file", logger.String("path", filePath), logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var recordCount int

	// Exportar según entity_type
	switch input.EntityType {
	case "users":
		recordCount, err = uc.exportUsers(writer, input)
	case "raffles":
		recordCount, err = uc.exportRaffles(writer, input)
	case "payments":
		recordCount, err = uc.exportPayments(writer, input)
	case "settlements":
		recordCount, err = uc.exportSettlements(writer, input)
	case "audit_logs":
		recordCount, err = uc.exportAuditLogs(writer, input)
	default:
		return nil, errors.New("VALIDATION_FAILED", "unsupported entity_type", 400, nil)
	}

	if err != nil {
		uc.log.Error("Error exporting data", logger.String("entity_type", input.EntityType), logger.Error(err))
		return nil, err
	}

	// Obtener tamaño del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		uc.log.Error("Error getting file info", logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
	}

	// Calcular expiración (24 horas)
	expiresAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	// TODO: Generar URL de descarga pública (presigned URL o endpoint de download)
	downloadURL := fmt.Sprintf("/api/v1/admin/exports/download/%s", fileName)

	// Log auditoría crítica (exportación de datos sensibles)
	uc.log.Error("Admin exported data",
		logger.Int64("admin_id", adminID),
		logger.String("entity_type", input.EntityType),
		logger.String("format", input.Format),
		logger.Int("record_count", recordCount),
		logger.String("file_name", fileName),
		logger.String("action", "admin_export_data"),
		logger.String("severity", "critical"))

	return &ExportDataOutput{
		FilePath:    filePath,
		FileName:    fileName,
		DownloadURL: downloadURL,
		RecordCount: recordCount,
		FileSize:    fileInfo.Size(),
		ExpiresAt:   expiresAt,
	}, nil
}

// exportUsers exporta usuarios a CSV
func (uc *ExportDataUseCase) exportUsers(writer *csv.Writer, input *ExportDataInput) (int, error) {
	// Escribir header
	header := []string{"ID", "Email", "First Name", "Last Name", "Role", "Status", "KYC Level", "Created At", "Last Login"}
	if err := writer.Write(header); err != nil {
		return 0, errors.Wrap(errors.ErrInternalServer, err)
	}

	// Query usuarios
	query := uc.db.Table("users")

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}

	rows, err := query.Rows()
	if err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int64
		var email, role, status, kycLevel string
		var firstName, lastName *string
		var createdAt, lastLoginAt *time.Time

		if err := rows.Scan(&id, &email, &firstName, &lastName, &role, &status, &kycLevel, &createdAt, &lastLoginAt); err != nil {
			continue
		}

		firstNameStr := ""
		if firstName != nil {
			firstNameStr = *firstName
		}
		lastNameStr := ""
		if lastName != nil {
			lastNameStr = *lastName
		}
		createdAtStr := ""
		if createdAt != nil {
			createdAtStr = createdAt.Format(time.RFC3339)
		}
		lastLoginStr := ""
		if lastLoginAt != nil {
			lastLoginStr = lastLoginAt.Format(time.RFC3339)
		}

		record := []string{
			fmt.Sprintf("%d", id),
			email,
			firstNameStr,
			lastNameStr,
			role,
			status,
			kycLevel,
			createdAtStr,
			lastLoginStr,
		}

		if err := writer.Write(record); err != nil {
			return count, errors.Wrap(errors.ErrInternalServer, err)
		}
		count++
	}

	return count, nil
}

// exportRaffles exporta rifas a CSV
func (uc *ExportDataUseCase) exportRaffles(writer *csv.Writer, input *ExportDataInput) (int, error) {
	// Escribir header
	header := []string{"ID", "Title", "Organizer ID", "Status", "Total Numbers", "Sold Count", "Price", "Created At", "Completed At"}
	if err := writer.Write(header); err != nil {
		return 0, errors.Wrap(errors.ErrInternalServer, err)
	}

	// Query raffles
	query := uc.db.Table("raffles").Where("deleted_at IS NULL")

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}

	rows, err := query.Select("id, title, user_id, status, total_numbers, sold_count, price_per_number, created_at, completed_at").Rows()
	if err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, userID, totalNumbers, soldCount int64
		var title, status string
		var price float64
		var createdAt time.Time
		var completedAt *time.Time

		if err := rows.Scan(&id, &title, &userID, &status, &totalNumbers, &soldCount, &price, &createdAt, &completedAt); err != nil {
			continue
		}

		completedAtStr := ""
		if completedAt != nil {
			completedAtStr = completedAt.Format(time.RFC3339)
		}

		record := []string{
			fmt.Sprintf("%d", id),
			title,
			fmt.Sprintf("%d", userID),
			status,
			fmt.Sprintf("%d", totalNumbers),
			fmt.Sprintf("%d", soldCount),
			fmt.Sprintf("%.2f", price),
			createdAt.Format(time.RFC3339),
			completedAtStr,
		}

		if err := writer.Write(record); err != nil {
			return count, errors.Wrap(errors.ErrInternalServer, err)
		}
		count++
	}

	return count, nil
}

// exportPayments exporta pagos a CSV
func (uc *ExportDataUseCase) exportPayments(writer *csv.Writer, input *ExportDataInput) (int, error) {
	header := []string{"ID", "User ID", "Raffle ID", "Amount", "Currency", "Status", "Payment Method", "Created At", "Paid At"}
	if err := writer.Write(header); err != nil {
		return 0, errors.Wrap(errors.ErrInternalServer, err)
	}

	query := uc.db.Table("payments")

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}

	rows, err := query.Select("id, user_id, raffle_id, amount, currency, status, payment_method, created_at, paid_at").Rows()
	if err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, userID, raffleID string
		var amount float64
		var currency, status string
		var paymentMethod *string
		var createdAt time.Time
		var paidAt *time.Time

		if err := rows.Scan(&id, &userID, &raffleID, &amount, &currency, &status, &paymentMethod, &createdAt, &paidAt); err != nil {
			continue
		}

		paymentMethodStr := ""
		if paymentMethod != nil {
			paymentMethodStr = *paymentMethod
		}
		paidAtStr := ""
		if paidAt != nil {
			paidAtStr = paidAt.Format(time.RFC3339)
		}

		record := []string{
			id,
			userID,
			raffleID,
			fmt.Sprintf("%.2f", amount),
			currency,
			status,
			paymentMethodStr,
			createdAt.Format(time.RFC3339),
			paidAtStr,
		}

		if err := writer.Write(record); err != nil {
			return count, errors.Wrap(errors.ErrInternalServer, err)
		}
		count++
	}

	return count, nil
}

// exportSettlements exporta settlements a CSV
func (uc *ExportDataUseCase) exportSettlements(writer *csv.Writer, input *ExportDataInput) (int, error) {
	header := []string{"ID", "Organizer ID", "Raffle ID", "Total Revenue", "Platform Fee", "Net Amount", "Status", "Calculated At", "Paid At"}
	if err := writer.Write(header); err != nil {
		return 0, errors.Wrap(errors.ErrInternalServer, err)
	}

	query := uc.db.Table("settlements")

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("calculated_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("calculated_at <= ?", *input.DateTo+" 23:59:59")
	}

	rows, err := query.Select("id, organizer_id, raffle_id, total_revenue, platform_fee, net_amount, status, calculated_at, paid_at").Rows()
	if err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, organizerID, raffleID int64
		var totalRevenue, platformFee, netAmount float64
		var status string
		var calculatedAt time.Time
		var paidAt *time.Time

		if err := rows.Scan(&id, &organizerID, &raffleID, &totalRevenue, &platformFee, &netAmount, &status, &calculatedAt, &paidAt); err != nil {
			continue
		}

		paidAtStr := ""
		if paidAt != nil {
			paidAtStr = paidAt.Format(time.RFC3339)
		}

		record := []string{
			fmt.Sprintf("%d", id),
			fmt.Sprintf("%d", organizerID),
			fmt.Sprintf("%d", raffleID),
			fmt.Sprintf("%.2f", totalRevenue),
			fmt.Sprintf("%.2f", platformFee),
			fmt.Sprintf("%.2f", netAmount),
			status,
			calculatedAt.Format(time.RFC3339),
			paidAtStr,
		}

		if err := writer.Write(record); err != nil {
			return count, errors.Wrap(errors.ErrInternalServer, err)
		}
		count++
	}

	return count, nil
}

// exportAuditLogs exporta audit logs a CSV
func (uc *ExportDataUseCase) exportAuditLogs(writer *csv.Writer, input *ExportDataInput) (int, error) {
	header := []string{"ID", "Admin ID", "Action", "Entity Type", "Entity ID", "Description", "Severity", "Created At"}
	if err := writer.Write(header); err != nil {
		return 0, errors.Wrap(errors.ErrInternalServer, err)
	}

	query := uc.db.Table("audit_logs")

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}

	rows, err := query.Select("id, admin_id, action, entity_type, entity_id, description, severity, created_at").Rows()
	if err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, adminID int64
		var action, entityType, description, severity string
		var entityID *int64
		var createdAt time.Time

		if err := rows.Scan(&id, &adminID, &action, &entityType, &entityID, &description, &severity, &createdAt); err != nil {
			continue
		}

		entityIDStr := ""
		if entityID != nil {
			entityIDStr = fmt.Sprintf("%d", *entityID)
		}

		record := []string{
			fmt.Sprintf("%d", id),
			fmt.Sprintf("%d", adminID),
			action,
			entityType,
			entityIDStr,
			description,
			severity,
			createdAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return count, errors.Wrap(errors.ErrInternalServer, err)
		}
		count++
	}

	return count, nil
}
