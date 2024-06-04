package main

import (
	"booking/internal/models"
	"github.com/gin-gonic/gin"
	"time"
)

func (app *application) bookComputer(c *gin.Context) {
	var bookingRequest struct {
		ComputerID uint      `json:"computer_id"`
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
	}

	if err := c.BindJSON(&bookingRequest); err != nil {
		app.logger.Error("Invalid request format: " + err.Error())
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	userID := c.GetString("userID")

	var computer models.Computer
	if err := app.db.First(&computer, bookingRequest.ComputerID).Error; err != nil {
		app.logger.Error("Computer not found: " + err.Error())
		c.JSON(404, gin.H{"error": "Computer not found"})
		return
	}

	if computer.Status == "booked" {
		c.JSON(400, gin.H{"error": "Computer is already booked"})
		return
	}

	booking := models.Booking{
		UserID:     userID,
		ComputerID: bookingRequest.ComputerID,
		StartTime:  bookingRequest.StartTime,
		EndTime:    bookingRequest.EndTime,
	}

	if err := app.db.Create(&booking).Error; err != nil {
		app.logger.Error("Failed to create booking: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to create booking"})
		return
	}

	computer.Status = "booked"
	if err := app.db.Save(&computer).Error; err != nil {
		app.logger.Error("Failed to update computer status: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to update computer status"})
		return
	}

	c.JSON(200, gin.H{"message": "Computer booked successfully"})
}

func (app *application) getAvailableComputers(c *gin.Context) {
	var computers []models.Computer
	if err := app.db.Where("status = ?", "available").Find(&computers).Error; err != nil {
		app.logger.Error("Failed to retrieve computers: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to retrieve computers"})
		return
	}

	c.JSON(200, computers)
}

func (app *application) cancelBooking(c *gin.Context) {
	bookingID := c.Param("id")
	userID := c.GetString("userID")

	var booking models.Booking
	if err := app.db.First(&booking, bookingID).Error; err != nil {
		app.logger.Error("Booking not found: " + err.Error())
		c.JSON(404, gin.H{"error": "Booking not found"})
		return
	}

	if booking.UserID != userID {
		c.JSON(403, gin.H{"error": "You can only cancel your own bookings"})
		return
	}

	var computer models.Computer
	if err := app.db.First(&computer, booking.ComputerID).Error; err != nil {
		app.logger.Error("Computer not found: " + err.Error())
		c.JSON(404, gin.H{"error": "Computer not found"})
		return
	}

	computer.Status = "available"
	if err := app.db.Save(&computer).Error; err != nil {
		app.logger.Error("Failed to update computer status: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to update computer status"})
		return
	}

	if err := app.db.Delete(&booking).Error; err != nil {
		app.logger.Error("Failed to delete booking: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(200, gin.H{"message": "Booking cancelled successfully"})
}
