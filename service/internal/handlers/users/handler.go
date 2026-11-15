package users

import (
	"database/sql"
	"encoding/json"
	"energosbyt/service/internal/handlers"
	"energosbyt/service/internal/models"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	domain domain
	log    *slog.Logger
}

func New(domain domain, log *slog.Logger) *Handler {
	return &Handler{
		domain: domain,
		log:    log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("user_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.domain.GetUser(ctx, id)
		if err != nil {
			h.log.Error("failed to get user", "error", err)
			handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		if user == nil {
			handlers.JsonError(w, http.StatusNotFound, "Not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)

	case http.MethodPatch:
		var req models.PatchUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"fio\" field")
			return
		}
		if req.FIO == nil && req.Phone == nil {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"fio\" field")
			return
		}
		if req.FIO != nil && *req.FIO == "" {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"fio\" field")
			return
		}
		if req.Phone != nil && *req.Phone == "" {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"phone\" field")
			return
		}

		err := h.domain.UpdateUser(ctx, id, req.FIO, req.Phone)
		if errors.Is(err, sql.ErrNoRows) {
			handlers.JsonError(w, http.StatusNotFound, "Not found")
			return
		}
		if err != nil {
			h.log.Error("failed to update user", "error", err)
			handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
