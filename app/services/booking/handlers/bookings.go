package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	"github.com/gin-gonic/gin"
)

func (h *Handler) BookClass(c *gin.Context) {
	var bookClass bookings.BookClass
	err := c.BindJSON(&bookClass)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing body").Error()})
			return
		}
		h.cfg.Logger.Debugw("failed to bind booking", "error", err.Error())
		return
	}

	ctx := c.Request.Context()

	booking, err := h.cfg.BookingUsecase.BookClass(ctx, bookClass)
	if err != nil {
		if isInvalidBookingDataErr(err) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		h.cfg.Logger.Errorw("failed to add new booking", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add new booking"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func isInvalidBookingDataErr(err error) bool {
	return errors.Is(err, bookings.ErrMemberNotFound) || errors.Is(err, bookings.ErrClassNotFound) || errors.Is(err, bookings.ErrInvalidClassDate)
}

func (h *Handler) GetBookingByID(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()
	booking, err := h.cfg.BookingUsecase.GetByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, bookings.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("booking with ID %s not found", bookingID)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get booking. Retry later"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *Handler) DeleteBooking(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()

	err := h.cfg.BookingUsecase.DeleteBooking(ctx, bookingID)
	if err != nil {
		h.cfg.Logger.Debugw("failed to delete booking: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete booking"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListBookings(c *gin.Context) {
	pageInfo, err := h.extractBookingsPageInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query params are invalid"})
		return
	}

	ctx := c.Request.Context()

	allBookings, err := h.cfg.BookingUsecase.ListBookings(ctx, pageInfo)
	if err != nil {
		h.cfg.Logger.Debugw("failed to list bookings: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}

	c.JSON(http.StatusOK, allBookings)
}

func (h *Handler) extractBookingsPageInfo(c *gin.Context) (bookings.PageInfo, error) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var page int
	var limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			h.cfg.Logger.Debugw("failed to parse page param: %w", err)
			return bookings.PageInfo{}, err
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.cfg.Logger.Debugw("failed to parse limit param: %w", err)
			return bookings.PageInfo{}, err
		}
	}

	return bookings.PageInfo{
		Limit: limit,
		Page:  page,
	}, nil
}
