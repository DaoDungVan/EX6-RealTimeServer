package friend

import (
	"encoding/json"
	"net/http"
	"realtime-server/internal/auth"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(r *Repository) *Handler {
	return &Handler{Repo: r}
}

func (h *Handler) SendRequest(
	w http.ResponseWriter, r *http.Request) {
	var input struct {
		ToUserID string `json:"to_user_id"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	fromID := r.Context().Value(auth.UserIDKey).(string)
	fromObjID, _ := bson.ObjectIDFromHex(fromID)
	toObjID, _ := bson.ObjectIDFromHex(input.ToUserID)

	err := h.Repo.SendRequest(r.Context(), fromObjID, toObjID)
	if err != nil {
		http.Error(w, "Can not send invitation",
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).
		Encode(
			map[string]string{"message": "Friend request sent."})
}

func (h *Handler) AcceptRequest(
	w http.ResponseWriter, r *http.Request) {
	var input struct {
		RequestID string `json:"request_id"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	requestObjID, _ := bson.ObjectIDFromHex(input.RequestID)
	err := h.Repo.AcceptRequest(r.Context(), requestObjID)
	if err != nil {
		http.Error(w, "Can not accept invitation",
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{"message": "Friend request accepted."})
}
