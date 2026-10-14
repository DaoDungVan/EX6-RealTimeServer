package user

import (
	"encoding/json"
	"net/http"
	"realtime-server/internal/common"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).
			Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).
			Encode(map[string]string{"error": "invalid JSON body"})
		return
	}

	if hashed, err := common.
		HashPassword(input.Password); err == nil {
		input.Password = hashed
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).
			Encode(map[string]string{"error": "Cannot hash password"})
		return
	}

	if err := h.Repo.Create(r.Context(), &input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).
			Encode(map[string]string{"error": "cannot create user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).
		Encode(map[string]string{"message": "Register successfully!"})
}


