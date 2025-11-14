package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/adapters/notifier"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// TestEmailHandler maneja el envío de emails de prueba
func setupTestEmailRoute(router *gin.Engine, emailNotifier notifier.Notifier, log *logger.Logger) {
	router.POST("/api/v1/test/send-email", func(c *gin.Context) {
		var req struct {
			To      string `json:"to" binding:"required,email"`
			Type    string `json:"type" binding:"required"` // "verification", "welcome", "reset"
			Code    string `json:"code"`                     // Para verification
			Token   string `json:"token"`                    // Para reset
			Name    string `json:"name"`                     // Para welcome
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request: " + err.Error(),
			})
			return
		}

		var err error

		switch req.Type {
		case "verification":
			code := req.Code
			if code == "" {
				code = "123456" // Código por defecto para pruebas
			}
			err = emailNotifier.SendVerificationEmail(req.To, code)

		case "welcome":
			name := req.Name
			if name == "" {
				name = "Usuario" // Nombre por defecto
			}
			err = emailNotifier.SendWelcomeEmail(req.To, name)

		case "reset":
			token := req.Token
			if token == "" {
				token = "test-token-12345" // Token por defecto para pruebas
			}
			err = emailNotifier.SendPasswordResetEmail(req.To, token)

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid type. Use: verification, welcome, or reset",
			})
			return
		}

		if err != nil {
			log.Error("Failed to send test email",
				logger.String("to", req.To),
				logger.String("type", req.Type),
				logger.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to send email",
				"details": err.Error(),
			})
			return
		}

		log.Info("Test email sent successfully",
			logger.String("to", req.To),
			logger.String("type", req.Type),
		)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Email sent successfully",
			"to":      req.To,
			"type":    req.Type,
		})
	})
}
