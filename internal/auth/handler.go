package auth

import (
	"encoding/json"
	"net/http"
	"realtime-server/internal/common"
	"realtime-server/internal/user"
)

type Handler struct {
	UserRepo *user.Repository
}

var input struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(repo *user.Repository) *Handler {
	return &Handler{UserRepo: repo}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).
		Decode(&input); err != nil {
		http.
			Error(w, "Invalid input",
				http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.
		FindByEmail(r.Context(), input.Email)
	if err != nil || !common.
		CheckPassword(u.Password, input.Password) {
		http.
			Error(w, "Invalid email or password.",
				http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(u.ID.Hex())
	if err != nil {
		http.
			Error(w, "Create token failed!",
				http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
