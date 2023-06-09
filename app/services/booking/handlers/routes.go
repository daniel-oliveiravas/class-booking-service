package handlers

import (
	"errors"
	"net/http"

	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	"github.com/daniel-oliveiravas/class-booking-service/foundation/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	MembersUsecase *members.Usecase
	ClassesUsecase *classes.Usecase
	BookingUsecase *bookings.Usecase
	GinMode        string
	Logger         *zap.SugaredLogger
	PgProbe        *postgres.Probe
}

type Handler struct {
	cfg Config
}

func NewHandler(cfg Config) (*Handler, error) {
	if cfg.GinMode == "" {
		cfg.GinMode = "release"
	}

	if cfg.Logger == nil {
		return nil, errors.New("failed to build new handler: missing logger")
	}

	if cfg.MembersUsecase == nil {
		return nil, errors.New("failed to build new handler: missing members usecase")
	}

	if cfg.ClassesUsecase == nil {
		return nil, errors.New("failed to build new handler: missing classes usecase")
	}

	if cfg.BookingUsecase == nil {
		return nil, errors.New("failed to build new handler: missing classes usecase")
	}

	return &Handler{
		cfg: cfg,
	}, nil
}

func (h *Handler) API() http.Handler {
	gin.SetMode(h.cfg.GinMode)
	r := gin.Default()

	//Members routes
	r.POST("/members", h.AddMember)
	r.GET("/members/:id", h.GetMemberByID)
	r.PATCH("/members/:id", h.UpdateMember)
	r.DELETE("/members/:id", h.DeleteMember)
	r.GET("/members/", h.ListMembers)

	//Classes routes
	r.POST("/classes", h.AddClass)
	r.GET("/classes/:id", h.GetClassByID)
	r.PATCH("/classes/:id", h.UpdateClass)
	r.DELETE("/classes/:id", h.DeleteClass)
	r.GET("/classes", h.ListClasses)

	//Booking routes
	r.POST("/bookings", h.BookClass)
	r.GET("/bookings/:id", h.GetBookingByID)
	r.DELETE("/bookings/:id", h.DeleteBooking)
	r.GET("/bookings", h.ListBookings)

	//Health endpoints
	r.GET("/v1/readiness", h.Readiness)
	r.GET("/v1/liveness", h.Liveness)

	return r.Handler()
}
