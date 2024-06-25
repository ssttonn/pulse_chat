package http

import (
	"encoding/json"
	"net/http"
	"pulse/src/api-core/internal/domain"
	"pulse/src/api-core/internal/usecase"
	"pulse/src/pkg/auth"
)

type UserHandler struct {
	usecase usecase.UserUseCase
}

// NewUserHandler creates a new handler injecting the usecase
func NewUserHandler(u usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: u,
	}
}

type registerRequest struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
}

type registerResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input request", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.RegisterUser(r.Context(), req.Username, req.AvatarURL)
	if err != nil {
		http.Error(w, "Cannot create new user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	secretKey := "super-secret-key"
	token, err := auth.GenerateToken(user.ID, secretKey)
	if err != nil {
		http.Error(w, "Cannot create new token", http.StatusInternalServerError)
		return
	}

	resp := registerResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
