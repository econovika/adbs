package post_reading

import (
	"encoding/json"
	"energosbyt/service/internal/handlers"
	"energosbyt/service/internal/models"
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

type readingInput struct {
	UserID      *int64  `json:"user_id"`
	GaugeID     *int64  `json:"gauge_id"`
	TariffID    *int64  `json:"tariff_id"`
	ReadingDate *string `json:"reading_date"`
	DayKW       *int64  `json:"day_kw"`
	NightKW     *int64  `json:"night_kw"`
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

	var input readingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handlers.JsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.UserID == nil || *input.UserID <= 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"user_id\" field")
		return
	}
	if input.GaugeID == nil || *input.GaugeID <= 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"gauge_id\" field")
		return
	}
	if input.TariffID == nil || *input.TariffID <= 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"tariff_id\" field")
		return
	}
	if input.ReadingDate == nil || *input.ReadingDate == "" {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"reading_date\" field")
		return
	}
	if _, err := time.Parse("2006-01-02", *input.ReadingDate); err != nil {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"reading_date\" field")
		return
	}
	if input.DayKW == nil || *input.DayKW < 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"day_kw\" field")
		return
	}
	if input.NightKW == nil || *input.NightKW < 0 {
		handlers.JsonError(w, http.StatusBadRequest, "invalid \"night_kw\" field")
		return
	}

	user, err := h.domain.GetUser(ctx, uint64(*input.UserID))
	if err != nil || user == nil {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}
	gExists, err := h.domain.GaugeExists(ctx, uint64(*input.GaugeID))
	if err != nil || !gExists {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}
	tExists, err := h.domain.TariffExists(ctx, uint64(*input.TariffID))
	if err != nil || !tExists {
		handlers.JsonError(w, http.StatusNotFound, "Not found")
		return
	}

	req := models.PostReadingRequest{
		UserID:      uint64(*input.UserID),
		GaugeID:     uint64(*input.GaugeID),
		TariffID:    uint64(*input.TariffID),
		ReadingDate: *input.ReadingDate,
		DayKW:       uint32(*input.DayKW),
		NightKW:     uint32(*input.NightKW),
	}

	if err := h.domain.AddMeterReading(ctx, realtyID, &req); err != nil {
		h.log.Error("failed to add reading", "error", err)
		handlers.JsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
