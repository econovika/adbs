package models

// Domain Models

type User struct {
	ID        uint64 `json:"id"`
	FIO       string `json:"fio"`
	Phone     string `json:"phone"`
	PhonePrev string `json:"phone_prev"`
	CreatedAt string `json:"-"`
}

type Realty struct {
	ID      uint64 `json:"id"`
	Address string `json:"address"`
	Square  uint32 `json:"square"`
	Type    string `json:"type"`
	Number  string `json:"number"`
	Rooms   uint8  `json:"rooms"`
	Floor   uint8  `json:"floor"`
	Porch   uint8  `json:"porch"`
}

type MeterReading struct {
	ReadingDate string `json:"reading_date"`
	UserID      uint64 `json:"user_id"`
	TariffID    uint64 `json:"tariff_id"`
	GaugeID     uint64 `json:"gauge_id"`
	DayKW       uint32 `json:"day_kw"`
	NightKW     uint32 `json:"night_kw"`
}

type Invoice struct {
	TariffID  uint64 `json:"tariff_id"`
	GaugeID   uint64 `json:"gauge_id"`
	DayKW     uint32 `json:"day_kw"`
	NightKW   uint32 `json:"night_kw"`
	Price     uint64 `json:"price"`
	Deadline  string `json:"deadline"`
	CreatedAt string `json:"created_at"`
}

type Tariff struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	PriceDayKWPH   uint32 `json:"price_day_kwph"`
	PriceNightKWPH uint32 `json:"price_night_kwph"`
	ValidFrom      string `json:"valid_from"`
	ValidTo        string `json:"valid_to"`
}
