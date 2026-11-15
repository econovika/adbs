package realty

import (
	"encoding/json"
	"energosbyt/service/internal/handlers"
	"log/slog"
	"net/http"
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
	realtyID, err := handlers.ParsePathID(r, "realty_id")
	if err != nil {
		handlers.JsonError(w, http.StatusBadRequest, "invalid realty id")
		return
	}

	realty, err := h.domain.GetRealty(r.Context(), realtyID)
	if err != nil {
		h.log.Error("failed to get realty", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if realty == nil {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(realty)
}
