package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/airline-voucher-seat-assignment/models"
	"github.com/airline-voucher-seat-assignment/services"
	"github.com/airline-voucher-seat-assignment/utils"
)

type VoucherHandler struct {
	voucherService *services.VoucherService
}

func NewVoucherHandler(voucherService *services.VoucherService) *VoucherHandler {
	return &VoucherHandler{
		voucherService: voucherService,
	}
}

func (h *VoucherHandler) validateJSONRequest(r *http.Request, req interface{}) *utils.AppError {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewBadRequestError("Invalid JSON format")
	}
	return nil
}

func (h *VoucherHandler) validateMethod(r *http.Request, method string) *utils.AppError {
	if r.Method != method {
		return utils.NewMethodNotAllowedError("Method not allowed")
	}
	return nil
}

func (h *VoucherHandler) CheckVoucher(w http.ResponseWriter, r *http.Request) {
	if err := h.validateMethod(r, http.MethodPost); err != nil {
		utils.HandleError(w, err)
		return
	}

	var req models.VoucherRequest
	if err := h.validateJSONRequest(r, &req); err != nil {
		utils.HandleError(w, err)
		return
	}

	if req.FlightNumber == "" || req.Date == "" {
		utils.HandleError(w, utils.NewBadRequestError("Flight number and date are required"))
		return
	}

	response, err := h.voucherService.CheckVoucherExists(&req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SendSuccessResponse(w, response)
}

func (h *VoucherHandler) GenerateVoucher(w http.ResponseWriter, r *http.Request) {
	if err := h.validateMethod(r, http.MethodPost); err != nil {
		utils.HandleError(w, err)
		return
	}

	var req models.VoucherRequest
	if err := h.validateJSONRequest(r, &req); err != nil {
		utils.HandleError(w, err)
		return
	}

	if req.Name == "" || req.ID == "" || req.FlightNumber == "" ||
		req.Date == "" || req.Aircraft == "" {
		utils.HandleError(w, utils.NewBadRequestError("All fields (name, id, flightNumber, date, aircraft) are required"))
		return
	}

	if req.Aircraft != "ATR" && req.Aircraft != "Airbus 320" && req.Aircraft != "Boeing 737 Max" {
		utils.HandleError(w, utils.NewBadRequestError("Invalid aircraft type. Supported types: ATR, Airbus 320, Boeing 737 Max"))
		return
	}

	response, err := h.voucherService.GenerateVoucher(&req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SendSuccessResponse(w, response)
}
