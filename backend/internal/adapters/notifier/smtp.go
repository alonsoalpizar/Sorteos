package notifier

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// SMTPNotifier implementa el envío de emails con SMTP estándar
type SMTPNotifier struct {
	config   *config.SMTPConfig
	logger   *logger.Logger
	auth     smtp.Auth
	fromMail string
	fromName string
}

// NewSMTPNotifier crea una nueva instancia del notifier SMTP
func NewSMTPNotifier(cfg *config.SMTPConfig, logger *logger.Logger) *SMTPNotifier {
	// Configurar autenticación SMTP
	var auth smtp.Auth
	if cfg.Username != "" && cfg.Password != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	return &SMTPNotifier{
		config:   cfg,
		logger:   logger,
		auth:     auth,
		fromMail: cfg.FromEmail,
		fromName: cfg.FromName,
	}
}

// SendVerificationEmail envía un email de verificación
func (n *SMTPNotifier) SendVerificationEmail(email, code string) error {
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

	return n.sendEmail(email, subject, plainTextContent, htmlContent)
}

// SendPasswordResetEmail envía un email de reset de contraseña
func (n *SMTPNotifier) SendPasswordResetEmail(email, token string) error {
	subject := "Restablecer contraseña - Sorteos Platform"

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", n.config.FrontendURL, token)

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

	return n.sendEmail(email, subject, plainTextContent, htmlContent)
}

// SendWelcomeEmail envía un email de bienvenida
func (n *SMTPNotifier) SendWelcomeEmail(email, firstName string) error {
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
            <a href="%s" style="background-color: #3B82F6; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">Explorar Sorteos</a>
        </div>
        <hr style="border: none; border-top: 1px solid #E2E8F0; margin: 30px 0;">
        <p style="color: #94A3B8; font-size: 12px;">
            Saludos,<br>
            <strong>Equipo de Sorteos Platform</strong>
        </p>
    </div>
</body>
</html>
	`, firstName, n.config.FrontendURL)

	return n.sendEmail(email, subject, plainTextContent, htmlContent)
}

// sendEmail es el método interno que envía el email usando SMTP
func (n *SMTPNotifier) sendEmail(to, subject, plainText, html string) error {
	// Construir el mensaje MIME multipart/alternative
	from := fmt.Sprintf("%s <%s>", n.fromName, n.fromMail)

	// Headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/alternative; boundary=\"boundary123\""

	// Construir mensaje
	var message strings.Builder

	// Agregar headers
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	// Parte texto plano
	message.WriteString("--boundary123\r\n")
	message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(plainText)
	message.WriteString("\r\n")

	// Parte HTML
	message.WriteString("--boundary123\r\n")
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(html)
	message.WriteString("\r\n")

	// Fin del mensaje
	message.WriteString("--boundary123--\r\n")

	// Preparar dirección del servidor
	addr := fmt.Sprintf("%s:%d", n.config.Host, n.config.Port)

	// Enviar email
	var err error
	if n.config.UseTLS {
		// Conexión TLS directa (puerto 465)
		err = n.sendEmailWithTLS(addr, to, message.String())
	} else if n.config.UseSTARTTLS {
		// STARTTLS (puerto 587)
		err = n.sendEmailWithSTARTTLS(addr, to, message.String())
	} else {
		// Sin cifrado (puerto 25) - NO RECOMENDADO
		err = smtp.SendMail(addr, n.auth, n.fromMail, []string{to}, []byte(message.String()))
	}

	if err != nil {
		n.logger.Error("Error sending email via SMTP",
			logger.String("to", to),
			logger.String("subject", subject),
			logger.Error(err),
		)
		return fmt.Errorf("smtp error: %w", err)
	}

	n.logger.Info("Email sent successfully via SMTP",
		logger.String("to", to),
		logger.String("subject", subject),
	)

	return nil
}

// sendEmailWithTLS envía email con TLS directo (puerto 465)
func (n *SMTPNotifier) sendEmailWithTLS(addr, to, message string) error {
	// Configuración TLS
	tlsConfig := &tls.Config{
		ServerName:         n.config.Host,
		InsecureSkipVerify: n.config.SkipVerify, // Solo para desarrollo
	}

	// Conectar con TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial error: %w", err)
	}
	defer conn.Close()

	// Crear cliente SMTP
	client, err := smtp.NewClient(conn, n.config.Host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Quit()

	// Autenticar
	if n.auth != nil {
		if err := client.Auth(n.auth); err != nil {
			return fmt.Errorf("smtp auth error: %w", err)
		}
	}

	// Enviar MAIL FROM
	if err := client.Mail(n.fromMail); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}

	// Enviar RCPT TO
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to error: %w", err)
	}

	// Enviar DATA
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data error: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}

	return nil
}

// sendEmailWithSTARTTLS envía email con STARTTLS (puerto 587)
func (n *SMTPNotifier) sendEmailWithSTARTTLS(addr, to, message string) error {
	// Conectar sin TLS
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial error: %w", err)
	}
	defer client.Quit()

	// Decir HELO
	if err := client.Hello(n.config.Host); err != nil {
		return fmt.Errorf("smtp hello error: %w", err)
	}

	// Iniciar STARTTLS
	tlsConfig := &tls.Config{
		ServerName:         n.config.Host,
		InsecureSkipVerify: n.config.SkipVerify,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("smtp starttls error: %w", err)
	}

	// Autenticar
	if n.auth != nil {
		if err := client.Auth(n.auth); err != nil {
			return fmt.Errorf("smtp auth error: %w", err)
		}
	}

	// Enviar MAIL FROM
	if err := client.Mail(n.fromMail); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}

	// Enviar RCPT TO
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to error: %w", err)
	}

	// Enviar DATA
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data error: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}

	return nil
}
