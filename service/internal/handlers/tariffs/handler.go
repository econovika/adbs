package tariffs

import (
	"encoding/json"
	"energosbyt/service/internal/handlers"
	"log/slog"
	"net/http"
	"strings"
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"date\" parameter")
		return
	}

	tariff, err := h.domain.GetTariffByDate(r.Context(), dateStr)
	if err != nil {
		if strings.Contains(err.Error(), "invalid date") {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"date\" parameter")
			return
		}
		h.log.Error("failed to get tariff", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if tariff == nil {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tariff)
}
