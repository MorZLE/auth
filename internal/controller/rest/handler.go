package rest

import (
	"context"
	"github.com/MorZLE/auth/internal/controller"
	"github.com/MorZLE/auth/internal/domain/cerror"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log/slog"
	"strconv"
)

func NewHandler(log *slog.Logger, auth controller.Auth, authAdmin controller.AuthAdmin) *Handler {
	return &Handler{
		log:       log,
		auth:      auth,
		authAdmin: authAdmin,
	}
}

func (h *Handler) Run() {
	app := fiber.New(fiber.Config{ErrorHandler: cerror.ErrorHandler})
	app.Use(recover.New())
	app.Use(logger.New())
	h.Route(app)
	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}

type Handler struct {
	auth      controller.Auth
	authAdmin controller.AuthAdmin
	log       *slog.Logger
}

func (h *Handler) Route(app *fiber.App) {
	app.Post("/login", h.Login)
	app.Post("/register", h.Register)
	app.Get("/admin", h.IsAdmin)
	app.Post("/admin", h.CreateAdmin)
	app.Delete("/admin", h.DeleteAdmin)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var appID int64

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	login := c.Query("login")
	pass := c.Query("password")
	app := c.Query("app_ID")

	appID, err := strconv.ParseInt(app, 10, 32)
	if err != nil || login == "" || pass == "" || appID == 0 {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}

	token, err := h.auth.LoginUser(ctx, login, pass, int32(appID))
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(models.LoginResponse{Status: 200,
		Body: models.LoginBodyResponse{Token: token}})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var appID int64

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	login := c.Query("login")
	pass := c.Query("password")
	app := c.Query("app_ID")

	appID, err := strconv.ParseInt(app, 10, 32)
	if err != nil || login == "" || pass == "" || appID == 0 {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}

	userID, err := h.auth.RegisterNewUser(ctx, login, pass, int32(appID))
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(models.RegisterResponse{Status: 200,
		Body: models.RegisterBodyResponse{UserID: userID}})
}

func (h *Handler) IsAdmin(c *fiber.Ctx) error {
	return nil
}
func (h *Handler) CreateAdmin(c *fiber.Ctx) error {
	return nil
}
func (h *Handler) DeleteAdmin(c *fiber.Ctx) error {
	return nil
}
func (h *Handler) AddApp(c *fiber.Ctx) error {
	return nil
}
