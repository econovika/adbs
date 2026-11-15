package bind_gauge

import (
	"encoding/json"
	"energosbyt/service/internal/handlers"
	"energosbyt/service/internal/models"
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

	var req models.BindGaugeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GaugeID == 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"gauge_id\" field")
		return
	}

	gExists, err := h.domain.GaugeExists(ctx, req.GaugeID)
	if err != nil || !gExists {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}

	if err := h.domain.BindGauge(ctx, realtyID, req.GaugeID); err != nil {
		h.log.Error("failed to bind gauge", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
