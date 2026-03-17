package users

import (
	"kwadw0/gocommerce/internal/json"
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
		slog.Error("Error creating user", "error", err)
		json.Error(w, http.StatusInternalServerError, "Failed to create user", "Internal server error")
		return
	}

	json.Success(w, http.StatusCreated, "User created successfully!", user)
}
