package users

import (
	"errors"
	"kwadw0/gocommerce/internal/json"
	"kwadw0/gocommerce/utils"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

var validate = validator.New()

// CreateUser godoc
// @Summary      Create a new user
// @Description  Takes a JSON payload and registers a new user in our database
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserDTO  true  "User to register"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Router       /signup [post]
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserDTO

	if err := json.ReadJson(w, r, &payload); err != nil {
		slog.Error("Error reading json", "error", err)
		json.Error(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		slog.Error("Validation failed", "error", err)
		json.Error(w, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	user, err := h.service.CreateUser(r.Context(), payload)
	if err != nil {
		if errors.Is(err, utils.ErrDuplicateEmail) {
			json.Error(w, http.StatusConflict, "Signup failed", err.Error())
			return
		}
		slog.Error("Error creating user", "error", err)
		json.Error(w, http.StatusInternalServerError, "Failed to create user", "Internal server error")
		return
	}

	json.Success(w, http.StatusCreated, "User created successfully!", user)
}

// Login godoc
// @Summary      Login
// @Description  Takes a JSON payload and logs in a user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      LoginDto  true  "User to login"
// @Success      200   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Router       /login [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginDto

	if err := json.ReadJson(w, r, &payload); err != nil {
		slog.Error("Error reading json", "error", err)
		json.Error(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		slog.Error("Validation failed", "error", err)
		json.Error(w, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	user, token, err := h.service.Login(r.Context(), payload)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) || errors.Is(err, utils.ErrInvalidCredentials) {
			json.Error(w, http.StatusUnauthorized, "Login failed", "Invalid email or password")
			return
		}
		slog.Error("Error logging in", "error", err)
		json.Error(w, http.StatusInternalServerError, "Failed to login", "Internal server error")
		return
	}

	json.Success(w, http.StatusOK, "User logged in successfully!", map[string]any{
		"user":  user,
		"token": token,
	})
}
