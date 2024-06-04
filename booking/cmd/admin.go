package main

import (
	"booking/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) createComputer(c *gin.Context) {
	var computer models.Computer
	if err := c.BindJSON(&computer); err != nil {
		app.logger.Error("Invalid request format: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := app.db.Create(&computer).Error; err != nil {
		app.logger.Error("Failed to create computer: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create computer"})
		return
	}

	c.JSON(http.StatusCreated, computer)
}

func (app *application) getComputers(c *gin.Context) {
	var computers []models.Computer
	if err := app.db.Find(&computers).Error; err != nil {
		app.logger.Error("Failed to retrieve computers: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve computers"})
		return
	}

	c.JSON(http.StatusOK, computers)
}

func (app *application) getComputerByID(c *gin.Context) {
	id := c.Param("id")
	var computer models.Computer
	if err := app.db.First(&computer, id).Error; err != nil {
		app.logger.Error("Computer not found: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Computer not found"})
		return
	}

	c.JSON(http.StatusOK, computer)
}

func (app *application) updateComputer(c *gin.Context) {
	id := c.Param("id")
	var computer models.Computer
	if err := app.db.First(&computer, id).Error; err != nil {
		app.logger.Error("Computer not found: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Computer not found"})
		return
	}

	if err := c.BindJSON(&computer); err != nil {
		app.logger.Error("Invalid request format: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := app.db.Save(&computer).Error; err != nil {
		app.logger.Error("Failed to update computer: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update computer"})
		return
	}

	c.JSON(http.StatusOK, computer)
}

func (app *application) deleteComputer(c *gin.Context) {
	id := c.Param("id")
	if err := app.db.Delete(&models.Computer{}, id).Error; err != nil {
		app.logger.Error("Failed to delete computer: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete computer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Computer deleted successfully"})
}

func (app *application) createBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.BindJSON(&booking); err != nil {
		app.logger.Error("Invalid request format: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := app.db.Create(&booking).Error; err != nil {
		app.logger.Error("Failed to create booking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (app *application) getBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := app.db.Find(&bookings).Error; err != nil {
		app.logger.Error("Failed to retrieve bookings: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (app *application) getBookingByID(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	if err := app.db.First(&booking, id).Error; err != nil {
		app.logger.Error("Booking not found: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (app *application) updateBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	if err := app.db.First(&booking, id).Error; err != nil {
		app.logger.Error("Booking not found: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if err := c.BindJSON(&booking); err != nil {
		app.logger.Error("Invalid request format: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := app.db.Save(&booking).Error; err != nil {
		app.logger.Error("Failed to update booking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (app *application) deleteBooking(c *gin.Context) {
	id := c.Param("id")
	if err := app.db.Delete(&models.Booking{}, id).Error; err != nil {
		app.logger.Error("Failed to delete booking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
