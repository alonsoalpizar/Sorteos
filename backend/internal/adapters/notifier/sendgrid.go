package notifier

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// SendGridNotifier implementa el envío de emails con SendGrid
type SendGridNotifier struct {
	client   *sendgrid.Client
	config   *config.SendGridConfig
	logger   *logger.Logger
	fromMail *mail.Email
}

// NewSendGridNotifier crea una nueva instancia del notifier
func NewSendGridNotifier(cfg *config.SendGridConfig, logger *logger.Logger) *SendGridNotifier {
	client := sendgrid.NewSendClient(cfg.APIKey)
	fromMail := mail.NewEmail(cfg.FromName, cfg.FromEmail)

	return &SendGridNotifier{
		client:   client,
		config:   cfg,
		logger:   logger,
		fromMail: fromMail,
	}
}

// SendVerificationEmail envía un email de verificación
func (n *SendGridNotifier) SendVerificationEmail(email, code string) error {
	to := mail.NewEmail("", email)
	subject := "Verifica tu cuenta - Sorteos Platform"

	plainTextContent := fmt.Sprintf(`
Hola,

Gracias por registrarte en Sorteos Platform.

Tu código de verificación es: %s

Este código expirará en 15 minutos.

Si no solicitaste este código, puedes ignorar este email.

Saludos,
Equipo de Sorteos Platform
	`, code)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verifica tu cuenta</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #3B82F6;">Verifica tu cuenta</h2>
        <p>Hola,</p>
        <p>Gracias por registrarte en <strong>Sorteos Platform</strong>.</p>
        <div style="background-color: #EFF6FF; border-left: 4px solid #3B82F6; padding: 15px; margin: 20px 0;">
            <p style="margin: 0; font-size: 14px; color: #64748B;">Tu código de verificación es:</p>
            <p style="margin: 10px 0 0 0; font-size: 32px; font-weight: bold; color: #3B82F6; letter-spacing: 5px;">%s</p>
        </div>
        <p style="color: #64748B; font-size: 14px;">Este código expirará en <strong>15 minutos</strong>.</p>
        <p style="color: #64748B; font-size: 14px;">Si no solicitaste este código, puedes ignorar este email.</p>
        <hr style="border: none; border-top: 1px solid #E2E8F0; margin: 30px 0;">
        <p style="color: #94A3B8; font-size: 12px;">
            Saludos,<br>
            <strong>Equipo de Sorteos Platform</strong>
        </p>
    </div>
</body>
</html>
	`, code)

	message := mail.NewSingleEmail(n.fromMail, subject, to, plainTextContent, htmlContent)

	response, err := n.client.Send(message)
	if err != nil {
		n.logger.Error("Error sending verification email",
			logger.String("email", email),
			logger.Error(err),
		)
		return err
	}

	if response.StatusCode >= 400 {
		n.logger.Error("SendGrid returned error",
			logger.String("email", email),
			logger.Int("status_code", response.StatusCode),
			logger.String("body", response.Body),
		)
		return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
	}

	n.logger.Info("Verification email sent",
		logger.String("email", email),
		logger.Int("status_code", response.StatusCode),
	)

	return nil
}

// SendPasswordResetEmail envía un email de reset de contraseña
func (n *SendGridNotifier) SendPasswordResetEmail(email, token string) error {
	to := mail.NewEmail("", email)
	subject := "Restablecer contraseña - Sorteos Platform"

	resetURL := fmt.Sprintf("https://sorteos.com/reset-password?token=%s", token)

	plainTextContent := fmt.Sprintf(`
Hola,

Recibimos una solicitud para restablecer tu contraseña.

Haz clic en el siguiente enlace para crear una nueva contraseña:
%s

Este enlace expirará en 1 hora.

Si no solicitaste restablecer tu contraseña, puedes ignorar este email de forma segura.

Saludos,
Equipo de Sorteos Platform
	`, resetURL)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Restablecer contraseña</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #3B82F6;">Restablecer contraseña</h2>
        <p>Hola,</p>
        <p>Recibimos una solicitud para restablecer tu contraseña en <strong>Sorteos Platform</strong>.</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #3B82F6; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">Restablecer Contraseña</a>
        </div>
        <p style="color: #64748B; font-size: 14px;">Este enlace expirará en <strong>1 hora</strong>.</p>
        <p style="color: #64748B; font-size: 14px;">Si no solicitaste restablecer tu contraseña, puedes ignorar este email de forma segura.</p>
        <hr style="border: none; border-top: 1px solid #E2E8F0; margin: 30px 0;">
        <p style="color: #94A3B8; font-size: 12px;">
            Saludos,<br>
            <strong>Equipo de Sorteos Platform</strong>
        </p>
    </div>
</body>
</html>
	`, resetURL)

	message := mail.NewSingleEmail(n.fromMail, subject, to, plainTextContent, htmlContent)

	response, err := n.client.Send(message)
	if err != nil {
		n.logger.Error("Error sending password reset email",
			logger.String("email", email),
			logger.Error(err),
		)
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
	}

	n.logger.Info("Password reset email sent",
		logger.String("email", email),
	)

	return nil
}

// SendWelcomeEmail envía un email de bienvenida
func (n *SendGridNotifier) SendWelcomeEmail(email, firstName string) error {
	to := mail.NewEmail(firstName, email)
	subject := "¡Bienvenido a Sorteos Platform!"

	plainTextContent := fmt.Sprintf(`
Hola %s,

¡Bienvenido a Sorteos Platform!

Tu cuenta ha sido verificada exitosamente y ya puedes empezar a participar en sorteos o crear los tuyos propios.

Explora nuestra plataforma y encuentra sorteos increíbles.

Saludos,
Equipo de Sorteos Platform
	`, firstName)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Bienvenido</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #3B82F6;">¡Bienvenido a Sorteos Platform!</h2>
        <p>Hola <strong>%s</strong>,</p>
        <p>Tu cuenta ha sido <strong>verificada exitosamente</strong> y ya puedes empezar a participar en sorteos o crear los tuyos propios.</p>
        <div style="background-color: #F0FDF4; border-left: 4px solid #10B981; padding: 15px; margin: 20px 0;">
            <p style="margin: 0; color: #065F46;">✓ Cuenta verificada</p>
            <p style="margin: 10px 0 0 0; color: #065F46;">✓ Listo para participar</p>
        </div>
        <div style="text-align: center; margin: 30px 0;">
            <a href="https://sorteos.com" style="background-color: #3B82F6; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">Explorar Sorteos</a>
        </div>
        <hr style="border: none; border-top: 1px solid #E2E8F0; margin: 30px 0;">
        <p style="color: #94A3B8; font-size: 12px;">
            Saludos,<br>
            <strong>Equipo de Sorteos Platform</strong>
        </p>
    </div>
</body>
</html>
	`, firstName)

	message := mail.NewSingleEmail(n.fromMail, subject, to, plainTextContent, htmlContent)

	response, err := n.client.Send(message)
	if err != nil {
		n.logger.Error("Error sending welcome email",
			logger.String("email", email),
			logger.Error(err),
		)
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
	}

	n.logger.Info("Welcome email sent",
		logger.String("email", email),
	)

	return nil
}
