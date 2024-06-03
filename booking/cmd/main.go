package main

import "booking/internal/logger"

func main() {
	log := logger.NewLogger("../logs/booking.log")

	log.Info("Starting booking service")
}
