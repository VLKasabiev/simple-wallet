package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/internal/service"
	"github.com/VLKasabiev/simple-wallet/internal/utils/validator"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}


func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": validator.FormatError(err),
		})
	}

	slog.Info("Login attempt", "email", req.Email)

	token, err := h.userService.Login(
		c.Request().Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		slog.Warn("Failed login attempt", "email", req.Email, "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	slog.Info("Successful login", "email", req.Email)

	return c.JSON(http.StatusCreated, echo.Map{"token": token})
}


func (h *UserHandler) Create(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": validator.FormatError(err),
		})
	}

	user, err := h.userService.CreateUser(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, echo.Map{
				"errors": echo.Map{
					"email": "User with such email already exist",
				},
			})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) List(c echo.Context) error {
	users, err := h.userService.ListUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0{
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id"})
	}

	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	user, err := h.userService.GetUserByID(c.Request().Context(), id, userID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		if errors.Is(err, model.ErrNotUserProfileOwner) {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "You can't watch another user profile"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}