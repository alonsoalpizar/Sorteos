package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// getAdminIDFromContext obtiene el ID del admin desde el contexto
func getAdminIDFromContext(c *gin.Context) (int64, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("UNAUTHORIZED", "admin not authenticated", 401, nil)
	}

	adminID, ok := userID.(int64)
	if !ok {
		return 0, errors.New("UNAUTHORIZED", "invalid admin ID", 401, nil)
	}

	// TODO: Validar que el usuario tiene rol super_admin o admin
	// role, exists := c.Get("user_role")
	// if !exists || (role != "super_admin" && role != "admin") {
	//     return 0, errors.New("FORBIDDEN", "insufficient permissions", 403, nil)
	// }

	return adminID, nil
}

// stringPtr convierte un string a puntero, retorna nil si está vacío
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// handleError maneja errores de manera centralizada
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Status, gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	// Error genérico
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "An internal error occurred",
		},
	})
}
