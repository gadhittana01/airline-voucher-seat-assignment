package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/airline-voucher-seat-assignment/database"
	"github.com/airline-voucher-seat-assignment/handlers"
	"github.com/airline-voucher-seat-assignment/repository"
	"github.com/airline-voucher-seat-assignment/services"
	"github.com/airline-voucher-seat-assignment/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config := utils.CheckAndSetConfig(".", "config")

	db, err := database.InitDB(config)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	voucherRepo := repository.NewVoucherRepository(db)

	voucherService := services.NewVoucherService(voucherRepo)

	voucherHandler := handlers.NewVoucherHandler(voucherService)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/check", voucherHandler.CheckVoucher).Methods("POST")
	api.HandleFunc("/generate", voucherHandler.GenerateVoucher).Methods("POST")

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(router)

	port := strconv.Itoa(config.ServerPort)

	log.Printf("Environment: %s", config.Environment)
	log.Printf("Database Path: %s", config.DBPath)
	log.Printf("Server starting on port %s", port)
	log.Printf("Log level: %s", config.LogLevel)
	log.Printf("Available endpoints:")
	log.Printf("  POST /api/check - Check if vouchers exist")
	log.Printf("  POST /api/generate - Generate new vouchers")
	log.Printf("  GET /api/health - Health check")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
