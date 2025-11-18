package domain

import (
	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// RechargeCalculator calcula el monto a cobrar al usuario basado en el modelo económico
// Fórmula: C = (D + f) / (1 - r)
// Donde:
// - C = Monto a cobrar (Charge)
// - D = Crédito deseado (Desired credit)
// - f = Tarifa fija del procesador (Fixed fee)
// - r = Tasa porcentual del procesador (Rate, como decimal: 3% = 0.03)
type RechargeCalculator struct {
	fixedFee        decimal.Decimal // Tarifa fija del procesador (ej: ₡100)
	processorRate   decimal.Decimal // Tasa porcentual (ej: 0.03 para 3%)
	platformFeeRate decimal.Decimal // Comisión de la plataforma (ej: 0.02 para 2%)
}

// NewRechargeCalculator crea una nueva instancia del calculador
func NewRechargeCalculator(fixedFee, processorRate, platformFeeRate decimal.Decimal) *RechargeCalculator {
	return &RechargeCalculator{
		fixedFee:        fixedFee,
		processorRate:   processorRate,
		platformFeeRate: platformFeeRate,
	}
}

// RechargeBreakdown contiene el desglose completo de un cálculo de recarga
type RechargeBreakdown struct {
	DesiredCredit   decimal.Decimal `json:"desired_credit"`   // Crédito que el usuario recibirá
	FixedFee        decimal.Decimal `json:"fixed_fee"`        // Tarifa fija del procesador
	ProcessorRate   decimal.Decimal `json:"processor_rate"`   // Tasa porcentual del procesador (ej: 0.03)
	ProcessorFee    decimal.Decimal `json:"processor_fee"`    // Comisión calculada del procesador
	PlatformFeeRate decimal.Decimal `json:"platform_fee_rate"` // Tasa de la plataforma (ej: 0.02)
	PlatformFee     decimal.Decimal `json:"platform_fee"`     // Comisión de la plataforma
	TotalFees       decimal.Decimal `json:"total_fees"`       // Total de comisiones
	ChargeAmount    decimal.Decimal `json:"charge_amount"`    // Monto total a cobrar al usuario
}

// CalculateCharge calcula el monto a cobrar para obtener el crédito deseado
// Fórmula: C = (D + f) / (1 - r)
func (rc *RechargeCalculator) CalculateCharge(desiredCredit decimal.Decimal) *RechargeBreakdown {
	// C = (D + f) / (1 - r)
	// Donde r incluye tanto la tasa del procesador como la de la plataforma
	totalRate := rc.processorRate.Add(rc.platformFeeRate)

	// Numerador: D + f
	numerator := desiredCredit.Add(rc.fixedFee)

	// Denominador: 1 - r
	denominator := decimal.NewFromFloat(1.0).Sub(totalRate)

	// C = numerador / denominador
	chargeAmount := numerator.Div(denominator)

	// Redondear a 2 decimales (centavos)
	chargeAmount = chargeAmount.Round(2)

	// Calcular comisiones individuales para el desglose
	// La tarifa fija ya está incluida en el cálculo de C, no se suma aquí
	processorFeePercentage := chargeAmount.Mul(rc.processorRate).Round(2)
	platformFee := chargeAmount.Mul(rc.platformFeeRate).Round(2)

	// Total de comisiones = tarifa fija + comisión porcentual procesador + comisión plataforma
	totalFees := rc.fixedFee.Add(processorFeePercentage).Add(platformFee)

	return &RechargeBreakdown{
		DesiredCredit:   desiredCredit,
		FixedFee:        rc.fixedFee,
		ProcessorRate:   rc.processorRate,
		ProcessorFee:    processorFeePercentage,
		PlatformFeeRate: rc.platformFeeRate,
		PlatformFee:     platformFee,
		TotalFees:       totalFees,
		ChargeAmount:    chargeAmount,
	}
}

// CalculateCredit calcula el crédito que recibirá el usuario dado un monto a cobrar
// Fórmula inversa: D = C * (1 - r) - f
func (rc *RechargeCalculator) CalculateCredit(chargeAmount decimal.Decimal) decimal.Decimal {
	// D = C * (1 - r) - f
	totalRate := rc.processorRate.Add(rc.platformFeeRate)
	oneMinusRate := decimal.NewFromFloat(1.0).Sub(totalRate)

	desiredCredit := chargeAmount.Mul(oneMinusRate).Sub(rc.fixedFee)
	return desiredCredit.Round(2)
}

// GetPredefinedRechargeOptions retorna los rangos predefinidos de recarga
// con sus respectivos desgloses
func (rc *RechargeCalculator) GetPredefinedRechargeOptions() []*RechargeBreakdown {
	// Rangos predefinidos de recarga
	predefinedAmounts := []int64{
		1000,  // ₡1,000
		5000,  // ₡5,000
		10000, // ₡10,000
		15000, // ₡15,000
		20000, // ₡20,000
		30000, // ₡30,000
	}

	breakdowns := make([]*RechargeBreakdown, 0, len(predefinedAmounts))
	for _, amount := range predefinedAmounts {
		desiredCredit := decimal.NewFromInt(amount)
		breakdown := rc.CalculateCharge(desiredCredit)
		breakdowns = append(breakdowns, breakdown)
	}

	return breakdowns
}

// Validate valida que los parámetros del calculador sean válidos
func (rc *RechargeCalculator) Validate() error {
	if rc.fixedFee.LessThan(decimal.Zero) {
		return errors.ErrInvalidConfiguration
	}

	if rc.processorRate.LessThan(decimal.Zero) || rc.processorRate.GreaterThanOrEqual(decimal.NewFromFloat(1.0)) {
		return errors.ErrInvalidConfiguration
	}

	if rc.platformFeeRate.LessThan(decimal.Zero) || rc.platformFeeRate.GreaterThanOrEqual(decimal.NewFromFloat(1.0)) {
		return errors.ErrInvalidConfiguration
	}

	// Validar que la suma de tasas no sea >= 1.0 (evitaría división por cero o negativo)
	totalRate := rc.processorRate.Add(rc.platformFeeRate)
	if totalRate.GreaterThanOrEqual(decimal.NewFromFloat(1.0)) {
		return errors.ErrInvalidConfiguration
	}

	return nil
}

// GetFixedFee retorna la tarifa fija configurada
func (rc *RechargeCalculator) GetFixedFee() decimal.Decimal {
	return rc.fixedFee
}

// GetProcessorRate retorna la tasa del procesador configurada
func (rc *RechargeCalculator) GetProcessorRate() decimal.Decimal {
	return rc.processorRate
}

// GetPlatformFeeRate retorna la tasa de la plataforma configurada
func (rc *RechargeCalculator) GetPlatformFeeRate() decimal.Decimal {
	return rc.platformFeeRate
}
