package models

import "time"

type VoucherRequest struct {
	IsRegenerate bool     `json:"is_regenerate"`
	UpdatedSeat  []string `json:"updated_seat"`
	Name         string   `json:"name" db:"crew_name"`
	ID           string   `json:"id" db:"crew_id"`
	FlightNumber string   `json:"flightNumber" db:"flight_number"`
	Date         string   `json:"date" db:"flight_date"`
	Aircraft     string   `json:"aircraft" db:"aircraft_type"`
}

type CheckResponse struct {
	Exists bool `json:"exists"`
}

type GenerateResponse struct {
	Success bool     `json:"success"`
	Seats   []string `json:"seats"`
}

type VoucherSeat struct {
	Seat       string `json:"seat"`
	AssignedTo string `json:"assigned_to"`
}

type VoucherDB struct {
	ID           int       `db:"id"`
	CrewName     string    `db:"crew_name"`
	CrewID       string    `db:"crew_id"`
	FlightNumber string    `db:"flight_number"`
	FlightDate   string    `db:"flight_date"`
	AircraftType string    `db:"aircraft_type"`
	Seat1        string    `db:"seat1"`
	Seat2        string    `db:"seat2"`
	Seat3        string    `db:"seat3"`
	CreatedAt    time.Time `db:"created_at"`
}

type UpdateVoucherDB struct {
	FlightNumber string `db:"flight_number"`
	FlightDate   string `db:"flight_date"`
	AircraftType string `db:"aircraft_type"`
	Seat1        string `db:"seat1"`
	Seat2        string `db:"seat2"`
	Seat3        string `db:"seat3"`
}
