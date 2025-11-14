package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/adapters/http/middleware"
	redisAdapter "github.com/sorteos-platform/backend/internal/adapters/redis"
	"github.com/sorteos-platform/backend/internal/infrastructure/payment"
	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/internal/usecases"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
	"github.com/sorteos-platform/backend/internal/infrastructure/websocket"
)

// Helper function to convert int64 userID to uuid.UUID
func getUserUUID(userRepo domain.UserRepository, userID int64) (uuid.UUID, error) {
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return uuid.Nil, err
	}
	if user == nil {
		return uuid.Nil, err
	}
	return uuid.Parse(user.UUID)
}

// setupReservationAndPaymentRoutes configura las rutas de reservas y pagos
func setupReservationAndPaymentRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, wsHub *websocket.Hub, cfg *config.Config, log *logger.Logger) {
	// Inicializar repositorios existentes
	raffleRepo := db.NewRaffleRepository(gormDB)
	userRepo := db.NewUserRepository(gormDB)
	raffleNumberRepo := db.NewRaffleNumberRepository(gormDB)

	// Inicializar nuevos repositorios
	reservationRepo := db.NewReservationRepository(gormDB)
	paymentRepo := db.NewPaymentRepository(gormDB)
	idempotencyKeyRepo := db.NewIdempotencyKeyRepository(gormDB)

	// Inicializar servicios de infraestructura
	lockService := redisinfra.NewLockService(rdb)

	// Inicializar payment provider basado en configuración
	var paymentProvider payment.PaymentProvider
	if cfg.Payment.Provider == "paypal" {
		var err error
		paymentProvider, err = payment.NewPayPalProvider(
			cfg.Payment.ClientID,
			cfg.Payment.Secret,
			cfg.Payment.Sandbox,
		)
		if err != nil {
			log.Fatal("Failed to initialize PayPal provider", logger.Error(err))
		}
		log.Info("Using PayPal as payment provider")
	} else {
		// Fallback a Stripe si está configurado
		paymentProvider = payment.NewStripeProvider(cfg.Stripe.SecretKey)
		log.Info("Using Stripe as payment provider")
	}

	// Inicializar use cases
	reservationUseCases := usecases.NewReservationUseCases(
		reservationRepo,
		raffleRepo,
		raffleNumberRepo,
		userRepo,
		lockService,
		wsHub,
	)

	paymentUseCases := usecases.NewPaymentUseCases(
		paymentRepo,
		reservationRepo,
		raffleRepo,
		idempotencyKeyRepo,
		paymentProvider,
		reservationUseCases,
	)

	// Inicializar middlewares
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)
	rateLimiter := middleware.NewRateLimiter(rdb, log)

	// Grupo de rutas de reservas
	reservationsGroup := router.Group("/api/v1/reservations")
	reservationsGroup.Use(authMiddleware.Authenticate())
	reservationsGroup.Use(authMiddleware.RequireMinKYC("email_verified"))
	{
		// POST /api/v1/reservations - Crear reserva
		reservationsGroup.POST("",
			rateLimiter.LimitByUser(cfg.Business.RateLimitReservePerMinute, time.Minute),
			func(c *gin.Context) {
				var req struct {
					RaffleID  string   `json:"raffle_id" binding:"required"`
					NumberIDs []string `json:"number_ids" binding:"required,min=1"`
					SessionID string   `json:"session_id" binding:"required"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_INPUT", "message": err.Error()})
					return
				}

				userIDInt, _ := middleware.GetUserID(c)
				userUUID, err := getUserUUID(userRepo, userIDInt)
				if err != nil {
					log.Error("Failed to get user UUID", logger.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
					return
				}

				raffleID, err := uuid.Parse(req.RaffleID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_RAFFLE_ID", "message": "invalid raffle_id"})
					return
				}

				reservation, err := reservationUseCases.CreateReservation(c.Request.Context(), usecases.CreateReservationInput{
					RaffleID:  raffleID,
					UserID:    userUUID,
					NumberIDs: req.NumberIDs,
					SessionID: req.SessionID,
				})

				if err != nil {
					log.Error("Failed to create reservation", logger.Error(err))
					c.JSON(http.StatusConflict, gin.H{"code": "RESERVATION_FAILED", "message": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"reservation": reservation,
				})
			},
		)

		// GET /api/v1/reservations/:id - Ver reserva
		reservationsGroup.GET("/:id", func(c *gin.Context) {
			reservationID, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "invalid reservation id"})
				return
			}

			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			reservation, err := reservationUseCases.GetReservation(c.Request.Context(), reservationID)
			if err != nil || reservation.UserID != userUUID {
				c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "reservation not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": reservation})
		})

		// GET /api/v1/reservations/me - Mis reservas
		reservationsGroup.GET("/me", func(c *gin.Context) {
			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			reservations, err := reservationUseCases.GetUserReservations(c.Request.Context(), userUUID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": "FETCH_FAILED", "message": "failed to fetch reservations"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": reservations, "count": len(reservations)})
		})

		// POST /api/v1/reservations/:id/cancel - Cancelar reserva
		reservationsGroup.POST("/:id/cancel", func(c *gin.Context) {
			reservationID, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "invalid reservation id"})
				return
			}

			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			// Verify ownership
			reservation, err := reservationUseCases.GetReservation(c.Request.Context(), reservationID)
			if err != nil || reservation.UserID != userUUID {
				c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "reservation not found"})
				return
			}

			// Cancel reservation
			if err := reservationUseCases.CancelReservation(c.Request.Context(), reservationID); err != nil {
				log.Error("Failed to cancel reservation", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "CANCEL_FAILED", "message": "failed to cancel reservation"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "reservation cancelled"})
		})

		// POST /api/v1/reservations/:id/add-number - Agregar número a reserva existente
		reservationsGroup.POST("/:id/add-number", func(c *gin.Context) {
			reservationID, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "invalid reservation id"})
				return
			}

			var req struct {
				NumberID string `json:"number_id" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_INPUT", "message": err.Error()})
				return
			}

			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			// Verify ownership
			reservation, err := reservationUseCases.GetReservation(c.Request.Context(), reservationID)
			if err != nil || reservation.UserID != userUUID {
				c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "reservation not found"})
				return
			}

			// Verificar que el número no esté ya en la reserva
			for _, num := range reservation.NumberIDs {
				if num == req.NumberID {
					c.JSON(http.StatusConflict, gin.H{"code": "NUMBER_ALREADY_IN_RESERVATION", "message": "number already selected"})
					return
				}
			}

			// Usar el use case para agregar el número (esto actualiza raffle_numbers y envía WebSocket)
			if err := reservationUseCases.AddNumberToReservation(c.Request.Context(), reservationID, req.NumberID); err != nil {
				log.Error("Failed to add number to reservation", logger.Error(err))

				// Manejar errores específicos
				if err.Error() == "number is already reserved" {
					c.JSON(http.StatusConflict, gin.H{"code": "NUMBER_ALREADY_RESERVED", "message": "number is already reserved"})
					return
				}

				c.JSON(http.StatusInternalServerError, gin.H{"code": "ADD_NUMBER_FAILED", "message": err.Error()})
				return
			}

			// Obtener la reserva actualizada
			updatedReservation, err := reservationUseCases.GetReservation(c.Request.Context(), reservationID)
			if err != nil {
				log.Error("Failed to get updated reservation", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "FETCH_FAILED", "message": "failed to fetch updated reservation"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"reservation": updatedReservation})
		})

		// POST /api/v1/reservations/:id/confirm - Confirmar reserva (pago completado)
		reservationsGroup.POST("/:id/confirm", func(c *gin.Context) {
			reservationID, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "invalid reservation id"})
				return
			}

			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			// Verify ownership
			reservation, err := reservationUseCases.GetReservation(c.Request.Context(), reservationID)
			if err != nil || reservation.UserID != userUUID {
				c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "reservation not found"})
				return
			}

			// Confirm reservation
			if err := reservationUseCases.ConfirmReservation(c.Request.Context(), reservationID); err != nil {
				log.Error("Failed to confirm reservation", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "CONFIRM_FAILED", "message": "failed to confirm reservation"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "reservation confirmed"})
		})
	}

	// GET /api/v1/raffles/:id/my-reservation - Obtener reserva activa del usuario para un sorteo
	router.GET("/api/v1/raffles/:id/my-reservation",
		authMiddleware.Authenticate(),
		func(c *gin.Context) {
			raffleID := c.Param("id")
			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusNotFound, gin.H{"message": "no active reservation"})
				return
			}

			reservation, err := reservationUseCases.GetActiveReservation(
				c.Request.Context(),
				userUUID,
				raffleID,
			)

			if err != nil || reservation == nil {
				// 200 con null para indicar que no hay reserva activa (no es un error)
				c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
				return
			}

			// Check if expired
			if reservation.IsExpired() {
				// Auto-cancel expired reservation
				reservationUseCases.CancelReservation(c.Request.Context(), reservation.ID)
				// 200 con null porque la reserva ya expiró
				c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": reservation})
		},
	)

	// Grupo de rutas de pagos
	paymentsGroup := router.Group("/api/v1/payments")
	paymentsGroup.Use(authMiddleware.Authenticate())
	paymentsGroup.Use(authMiddleware.RequireMinKYC("email_verified"))
	{
		// POST /api/v1/payments/intent - Crear payment intent
		paymentsGroup.POST("/intent",
			rateLimiter.LimitByUser(cfg.Business.RateLimitPaymentPerMinute, time.Minute),
			func(c *gin.Context) {
				var req struct {
					ReservationID  string `json:"reservation_id" binding:"required"`
					IdempotencyKey string `json:"idempotency_key"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_INPUT", "message": err.Error()})
					return
				}

				userIDInt, _ := middleware.GetUserID(c)
				userUUID, err := getUserUUID(userRepo, userIDInt)
				if err != nil {
					log.Error("Failed to get user UUID", logger.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
					return
				}

				reservationID, err := uuid.Parse(req.ReservationID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_RESERVATION_ID", "message": "invalid reservation_id"})
					return
				}

				// Get idempotency key from header or request body
				idempotencyKey := c.GetHeader("Idempotency-Key")
				if idempotencyKey == "" {
					idempotencyKey = req.IdempotencyKey
				}

				output, err := paymentUseCases.CreatePaymentIntent(c.Request.Context(), usecases.CreatePaymentIntentInput{
					ReservationID:  reservationID,
					UserID:         userUUID,
					IdempotencyKey: idempotencyKey,
				})

				if err != nil {
					log.Error("Failed to create payment intent", logger.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"code": "PAYMENT_FAILED", "message": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"success": true,
					"data": gin.H{
						"payment_id":    output.PaymentID.String(),
						"client_secret": output.ClientSecret,
						"amount":        output.Amount,
						"currency":      output.Currency,
					},
				})
			},
		)

		// GET /api/v1/payments/:id - Ver pago
		paymentsGroup.GET("/:id", func(c *gin.Context) {
			paymentID, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "invalid payment id"})
				return
			}

			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			payment, err := paymentUseCases.GetPayment(c.Request.Context(), paymentID)
			if err != nil || payment.UserID != userUUID {
				c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "payment not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": payment})
		})

		// GET /api/v1/payments/me - Mis pagos
		paymentsGroup.GET("/me", func(c *gin.Context) {
			userIDInt, _ := middleware.GetUserID(c)
			userUUID, err := getUserUUID(userRepo, userIDInt)
			if err != nil {
				log.Error("Failed to get user UUID", logger.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"code": "USER_NOT_FOUND", "message": "user not found"})
				return
			}

			payments, err := paymentUseCases.GetUserPayments(c.Request.Context(), userUUID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": "FETCH_FAILED", "message": "failed to fetch payments"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": payments, "count": len(payments)})
		})
	}

	// Webhook de Stripe (sin autenticación - Stripe firma los requests)
	router.POST("/api/v1/webhooks/stripe", func(c *gin.Context) {
		payload, err := c.GetRawData()
		if err != nil {
			log.Error("Failed to read webhook payload", logger.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PAYLOAD", "message": "invalid payload"})
			return
		}

		signature := c.GetHeader("Stripe-Signature")
		if signature == "" {
			log.Warn("Missing Stripe signature header")
			c.JSON(http.StatusBadRequest, gin.H{"code": "MISSING_SIGNATURE", "message": "missing signature"})
			return
		}

		event, err := paymentProvider.ConstructWebhookEvent(payload, signature, cfg.Payment.WebhookSecret)
		if err != nil {
			log.Error("Webhook signature verification failed", logger.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_SIGNATURE", "message": "invalid signature"})
			return
		}

		log.Info("Received Stripe webhook event", logger.String("type", event.Type))

		// Parse event data to extract payment intent ID
		var eventData struct {
			Object struct {
				ID string `json:"id"`
			} `json:"object"`
		}

		rawData, ok := event.Data.([]byte)
		if !ok {
			log.Error("Failed to cast event data to bytes")
			c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "internal error"})
			return
		}

		if err := json.Unmarshal(rawData, &eventData); err != nil {
			log.Error("Failed to parse event data", logger.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_EVENT_DATA", "message": "invalid event data"})
			return
		}

		if eventData.Object.ID == "" {
			log.Error("Payment intent ID not found in event data")
			c.JSON(http.StatusBadRequest, gin.H{"code": "MISSING_INTENT_ID", "message": "missing payment intent id"})
			return
		}

		if err := paymentUseCases.ProcessPaymentWebhook(c.Request.Context(), event.Type, eventData.Object.ID); err != nil {
			log.Error("Failed to process webhook", logger.Error(err), logger.String("event_type", event.Type))
			c.JSON(http.StatusInternalServerError, gin.H{"code": "WEBHOOK_FAILED", "message": "webhook processing failed"})
			return
		}

		log.Info("Successfully processed webhook", logger.String("event_type", event.Type))
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "webhook processed"})
	})
}
