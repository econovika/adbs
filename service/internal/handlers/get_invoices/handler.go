package get_invoices

import (
	"encoding/json"
	"energosbyt/service/internal/handlers"
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

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0
	var errP error

	if limitStr != "" {
		limit, errP = strconv.Atoi(limitStr)
		if errP != nil || limit < 0 {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"limit\" parameter")
			return
		}
	}
	if offsetStr != "" {
		offset, errP = strconv.Atoi(offsetStr)
		if errP != nil || offset < 0 {
			handlers.JsonError(w, http.StatusBadRequest, "invalid \"offset\" parameter")
			return
		}
	}

	invoices, total, err := h.domain.GetInvoices(ctx, realtyID, limit, offset)
	if err != nil {
		h.log.Error("failed to get invoices", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"invoices": invoices,
		"count":    total,
	})
}
