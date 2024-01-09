package rest

import (
	"context"
	"fmt"
	"github.com/MorZLE/auth/internal/controller"
	"github.com/MorZLE/auth/internal/domain/cerror"
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log/slog"
	"strconv"
	"time"
)

func NewHandler(log *slog.Logger, auth controller.Auth, authAdmin controller.AuthAdmin, port int, ttl time.Duration) *Handler {
	return &Handler{
		log:       log,
		auth:      auth,
		authAdmin: authAdmin,
		port:      port,
		ttl:       ttl,
	}
}

func (h *Handler) Run() {
	app := fiber.New(fiber.Config{ErrorHandler: cerror.ErrorHandler})
	app.Use(recover.New())
	app.Use(logger.New())
	h.Route(app)
	err := app.Listen(fmt.Sprintf(":%d", h.port))

	if err != nil {
		panic(err)
	}
}

type Handler struct {
	auth      controller.Auth
	authAdmin controller.AuthAdmin
	port      int
	ttl       time.Duration
	log       *slog.Logger
}

func (h *Handler) Route(app *fiber.App) {
	app.Post("/api/auth/login", h.Login)
	app.Post("/api/auth/register", h.Register)
	app.Get("/api/auth/checkadmin", h.IsAdmin)
	app.Post("/api/auth/createadmin", h.CreateAdmin)
	app.Delete("/api/auth/deleteadmin", h.DeleteAdmin)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var appID int64

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	login := c.Query("login")
	pass := c.Query("password")
	app := c.Query("app_id")

	appID, err := strconv.ParseInt(app, 10, 32)
	if err != nil || login == "" || pass == "" || appID == 0 {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}

	token, err := h.auth.LoginUser(ctx, login, pass, int32(appID))
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(
		models.Response{
			Status: 200,
			Body:   models.LoginBodyResponse{Token: token}},
	)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var appID int64

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	login := c.Query("login")
	pass := c.Query("password")
	app := c.Query("app_id")

	appID, err := strconv.ParseInt(app, 10, 32)
	if err != nil || login == "" || pass == "" || appID == 0 {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}

	userID, err := h.auth.RegisterNewUser(ctx, login, pass, int32(appID))
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(
		models.Response{
			Status: 200,
			Body:   models.RegisterBodyResponse{UserID: userID}},
	)
}

func (h *Handler) IsAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	userID := c.Query("user_id")
	appID := c.Query("app_id")

	userIDint, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)

	}
	intAPP, err := strconv.ParseInt(appID, 10, 32)

	if err != nil || intAPP == 0 || userIDint == 0 {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}
	admin, err := h.auth.CheckIsAdmin(ctx, int32(userIDint), int32(intAPP))
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(
		models.Response{
			Status: 200,
			Body: models.IsAdminBodyResponse{
				Result: true,
				LVL:    admin.Lvl},
		},
	)
}

func (h *Handler) CreateAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	login := c.Query("login")
	lvlAdmin := c.Query("lvl")
	key := c.Query("key")
	appID := c.Query("app_id")

	intLVL, err := strconv.ParseInt(lvlAdmin, 10, 32)
	if err != nil {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)

	}
	intAPP, err := strconv.ParseInt(appID, 10, 32)

	if err != nil || login == "" || intLVL == 0 || key == "" || intAPP == 0 {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}
	adminID, err := h.authAdmin.CreateAdmin(ctx, login, int32(intLVL), key, int32(intAPP))
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(
		models.Response{
			Status: 200,
			Body:   models.CreateAdminBodyResponse{AdminID: adminID}},
	)
}

func (h *Handler) DeleteAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	login := c.Query("login")
	key := c.Query("key")

	if login == "" || key == "" {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}

	res, err := h.authAdmin.DeleteAdmin(ctx, login, key)
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(
		models.Response{
			Status: 200,
			Body:   models.DeleteAdminBodyResponse{Result: res}},
	)
}

func (h *Handler) AddApp(c *fiber.Ctx) error {

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	name := c.Query("name")
	secret := c.Query("secret")
	key := c.Query("key")

	if name == "" || secret == "" || key == "" {
		return cerror.ErrorHandler(c, cerror.ErrInvalidCredentials)
	}

	appid, err := h.authAdmin.AddApp(ctx, name, secret, key)
	if err != nil {
		return cerror.ErrorHandler(c, err)
	}

	return c.JSON(
		models.Response{
			Status: 200,
			Body:   models.AddAppBodyResponse{AppID: appid},
		})
}
