/**
 * Formatea un número como moneda de Costa Rica (Colones)
 * @param amount - Monto a formatear
 * @returns String formateado con símbolo ₡ y separadores de miles
 */
export function formatCurrency(amount: number): string {
  return `₡${amount.toLocaleString('es-CR', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  })}`;
}

/**
 * Formatea un número como moneda de Costa Rica sin símbolo
 * @param amount - Monto a formatear
 * @returns String formateado con separadores de miles
 */
export function formatAmount(amount: number): string {
  return amount.toLocaleString('es-CR', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  });
}
