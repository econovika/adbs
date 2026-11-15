package get_readings

import (
	"encoding/json"
	"energosbyt/service/internal/handlers"
	"log/slog"
	"net/http"
	"strconv"
	"time"
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

	realtyIDStr := r.PathValue("realty_id")
	realtyID, err := strconv.ParseUint(realtyIDStr, 10, 64)
	if err != nil || realtyID == 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid realty id")
		return
	}

	realty, err := h.domain.GetRealty(ctx, realtyID)
	if err != nil {
		h.log.Error("failed to check realty", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if realty == nil {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}

	from := r.URL.Query().Get("date_from")
	to := r.URL.Query().Get("date_to")

	if from == "" {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"date_from\" parameter")
		return
	}
	if to == "" {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"date_to\" parameter")
		return
	}
	if _, err := time.Parse("2006-01-02", from); err != nil {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"date_from\" parameter")
		return
	}
	if _, err := time.Parse("2006-01-02", to); err != nil {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"date_to\" parameter")
		return
	}

	readings, err := h.domain.GetMeterReadings(ctx, realtyID, from, to)
	if err != nil {
		h.log.Error("failed to get readings", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"readings": readings,
		"count":    len(readings),
	})
}
