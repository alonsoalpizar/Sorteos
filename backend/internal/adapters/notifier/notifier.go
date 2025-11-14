package notifier

// Notifier es la interface para envío de notificaciones por email
// Permite usar diferentes implementaciones (SendGrid, SMTP, etc.)
// manteniendo el mismo contrato
type Notifier interface {
	// SendVerificationEmail envía un email con código de verificación
	SendVerificationEmail(email, code string) error

	// SendPasswordResetEmail envía un email con token de reset de contraseña
	SendPasswordResetEmail(email, token string) error

	// SendWelcomeEmail envía un email de bienvenida post-verificación
	SendWelcomeEmail(email, firstName string) error
}
