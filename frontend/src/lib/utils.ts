import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Utility function to merge Tailwind CSS classes
 * Combines clsx for conditional classes and tailwind-merge for proper Tailwind class merging
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Format currency to CRC (Costa Rican Colón)
 * Usa el símbolo ₡ (colón costarricense)
 */
export function formatCurrency(amount: number): string {
  const formatted = new Intl.NumberFormat("es-CR", {
    style: "decimal",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
  return `₡${formatted}`;
}

/**
 * Format date to Costa Rican format
 */
export function formatDate(date: Date | string): string {
  const d = typeof date === "string" ? new Date(date) : date;
  return new Intl.DateTimeFormat("es-CR", {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(d);
}

/**
 * Format datetime to Costa Rican format
 */
export function formatDateTime(date: Date | string): string {
  const d = typeof date === "string" ? new Date(date) : date;
  return new Intl.DateTimeFormat("es-CR", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(d);
}

/**
 * Validate Colombian cedula format (7-10 digits)
 */
export function isValidCedula(cedula: string): boolean {
  const cleaned = cedula.replace(/\D/g, "");
  return cleaned.length >= 7 && cleaned.length <= 10;
}

/**
 * Format cedula with dots (e.g., 1.234.567)
 */
export function formatCedula(cedula: string): string {
  const cleaned = cedula.replace(/\D/g, "");
  return cleaned.replace(/\B(?=(\d{3})+(?!\d))/g, ".");
}

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email);
}

/**
 * Validate phone format (E.164 format)
 */
export function isValidPhone(phone: string): boolean {
  const phoneRegex = /^\+?[1-9]\d{1,14}$/;
  return phoneRegex.test(phone);
}

/**
 * Get status badge color variant
 */
export function getStatusColor(
  status: string
): "default" | "success" | "warning" | "error" | "info" {
  switch (status) {
    case "active":
      return "success";
    case "draft":
      return "default";
    case "suspended":
      return "warning";
    case "completed":
      return "info";
    case "cancelled":
      return "error";
    default:
      return "default";
  }
}

/**
 * Get status label in Spanish
 */
export function getStatusLabel(status: string): string {
  switch (status) {
    case "draft":
      return "Borrador";
    case "active":
      return "Activo";
    case "suspended":
      return "Suspendido";
    case "completed":
      return "Completado";
    case "cancelled":
      return "Cancelado";
    default:
      return status;
  }
}

/**
 * Get draw method label in Spanish
 */
export function getDrawMethodLabel(method: string): string {
  switch (method) {
    case "loteria_nacional_cr":
      return "Lotería Nacional CR";
    case "manual":
      return "Sorteo Manual";
    case "random":
      return "Sorteo Aleatorio";
    default:
      return method;
  }
}
