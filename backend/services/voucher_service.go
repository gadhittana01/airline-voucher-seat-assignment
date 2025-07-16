package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/airline-voucher-seat-assignment/models"
	"github.com/airline-voucher-seat-assignment/repository"
	"github.com/airline-voucher-seat-assignment/utils"
)

type VoucherService struct {
	voucherRepo *repository.VoucherRepository
}

func NewVoucherService(voucherRepo *repository.VoucherRepository) *VoucherService {
	return &VoucherService{
		voucherRepo: voucherRepo,
	}
}

func (s *VoucherService) validateDate(dateStr string) *utils.AppError {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return utils.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	return nil
}

func (s *VoucherService) CheckVoucherExists(req *models.VoucherRequest) (*models.CheckResponse, *utils.AppError) {
	if err := s.validateDate(req.Date); err != nil {
		return nil, err
	}

	exists, err := s.voucherRepo.CheckVoucherExists(req.FlightNumber, req.Date)
	if err != nil {
		return nil, utils.NewInternalServerError(fmt.Sprintf("failed to check voucher existence: %v", err))
	}

	return &models.CheckResponse{
		Exists: exists,
	}, nil
}

func (s *VoucherService) checkSeatExists(seats []string, seat string) int {
	for idx, s := range seats {
		if s == seat {
			return idx
		}
	}

	return -1
}

func (s *VoucherService) GenerateVoucher(req *models.VoucherRequest) (*models.GenerateResponse, *utils.AppError) {
	if err := s.validateDate(req.Date); err != nil {
		return nil, err
	}

	var seats []string

	if req.IsRegenerate {
		currentSeats, err := s.voucherRepo.GetSeatByFlightNumberAndDate(req.FlightNumber, req.Date)
		if err != nil {
			return nil, utils.NewInternalServerError(err.Error())
		}

		if len(currentSeats) == 0 {
			return nil, utils.NewNotFoundError("seats not found")
		}

		seats = currentSeats

		for _, seat := range req.UpdatedSeat {
			idx := s.checkSeatExists(seats, seat)
			if idx != -1 {
				generatedSeat := s.generateRandomSeats(1, req.Aircraft)
				if generatedSeat == nil {
					return nil, utils.NewBadRequestError("invalid aircraft type")
				}

				seats[idx] = generatedSeat[0]
			}
		}

		err = s.voucherRepo.UpdateVoucher(&models.UpdateVoucherDB{
			FlightNumber: req.FlightNumber,
			FlightDate:   req.Date,
			AircraftType: req.Aircraft,
			Seat1:        seats[0],
			Seat2:        seats[1],
			Seat3:        seats[2],
		})
		if err != nil {
			return nil, utils.NewInternalServerError(err.Error())
		}
	} else {
		seats = s.generateRandomSeats(3, req.Aircraft)
		if seats == nil {
			return nil, utils.NewBadRequestError("invalid aircraft type")
		}

		if len(seats) != 3 {
			return nil, utils.NewInternalServerError("failed to generate required number of seats")
		}

		exists, err := s.voucherRepo.CheckVoucherExists(req.FlightNumber, req.Date)
		if err != nil {
			return nil, utils.NewInternalServerError(err.Error())
		}

		if exists {
			return nil, utils.NewConflictError("vouchers already exist for this flight date")
		}

		voucherDB := &models.VoucherDB{
			CrewName:     req.Name,
			CrewID:       req.ID,
			FlightNumber: req.FlightNumber,
			FlightDate:   req.Date,
			AircraftType: req.Aircraft,
			Seat1:        seats[0],
			Seat2:        seats[1],
			Seat3:        seats[2],
		}

		if err := s.voucherRepo.CreateVoucher(voucherDB); err != nil {
			return nil, utils.NewInternalServerError(fmt.Sprintf("failed to create voucher: %v", err))
		}
	}

	return &models.GenerateResponse{
		Success: true,
		Seats:   seats,
	}, nil
}

func (s *VoucherService) generateRandomSeats(count int, aircraftType string) []string {
	return GetRandomSeats(count, aircraftType)
}

func GetRandomSeats(count int, aircraft string) []string {
	layout := map[string]struct {
		rows int
		cols []string
	}{
		"ATR":            {18, []string{"A", "C", "D", "F"}},
		"Airbus 320":     {32, []string{"A", "B", "C", "D", "E", "F"}},
		"Boeing 737 Max": {32, []string{"A", "B", "C", "D", "E", "F"}},
	}

	cfg, ok := layout[aircraft]
	if !ok {
		return nil
	}

	var seats []string
	for r := 1; r <= cfg.rows; r++ {
		for _, c := range cfg.cols {
			seats = append(seats, fmt.Sprintf("%d%s", r, c))
		}
	}

	if len(seats) < count {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(seats), func(i, j int) {
		seats[i], seats[j] = seats[j], seats[i]
	})

	return seats[:count]
}
