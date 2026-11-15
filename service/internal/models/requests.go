package models

type PostReadingRequest struct {
	UserID      uint64 `json:"user_id"`
	GaugeID     uint64 `json:"gauge_id"`
	TariffID    uint64 `json:"tariff_id"`
	ReadingDate string `json:"reading_date"`
	DayKW       uint32 `json:"day_kw"`
	NightKW     uint32 `json:"night_kw"`
}

type PatchUserRequest struct {
	FIO   *string `json:"fio"`
	Phone *string `json:"phone"`
}

type BindGaugeRequest struct {
	GaugeID uint64 `json:"gauge_id"`
}
