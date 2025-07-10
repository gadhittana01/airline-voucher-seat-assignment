package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/airline-voucher-seat-assignment/models"
)

type VoucherRepository struct {
	db *sql.DB
}

func NewVoucherRepository(db *sql.DB) *VoucherRepository {
	return &VoucherRepository{
		db: db,
	}
}

func (r *VoucherRepository) CheckVoucherExists(flightNumber, flightDate string) (bool, error) {
	if r.db == nil {
		return false, fmt.Errorf("database connection is nil")
	}

	if flightNumber == "" || flightDate == "" {
		return false, fmt.Errorf("flight number and date cannot be empty")
	}

	query := `SELECT COUNT(*) FROM vouchers WHERE flight_number = ? AND flight_date = ?`

	var count int
	err := r.db.QueryRow(query, flightNumber, flightDate).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("database query failed: %v", err)
	}

	return count > 0, nil
}

func (r *VoucherRepository) CreateVoucher(voucher *models.VoucherDB) error {
	if r.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if voucher == nil {
		return fmt.Errorf("voucher cannot be nil")
	}

	if voucher.CrewName == "" || voucher.CrewID == "" || voucher.FlightNumber == "" ||
		voucher.FlightDate == "" || voucher.AircraftType == "" {
		return fmt.Errorf("all voucher fields are required")
	}

	if voucher.Seat1 == "" || voucher.Seat2 == "" || voucher.Seat3 == "" {
		return fmt.Errorf("all seat assignments are required")
	}

	query := `INSERT INTO vouchers (crew_name, crew_id, flight_number, flight_date, aircraft_type, seat1, seat2, seat3, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(query, voucher.CrewName, voucher.CrewID, voucher.FlightNumber,
		voucher.FlightDate, voucher.AircraftType, voucher.Seat1, voucher.Seat2, voucher.Seat3, now)
	if err != nil {
		return fmt.Errorf("failed to insert voucher into database: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were inserted")
	}

	return nil
}
