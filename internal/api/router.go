package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers"
	appmiddleware "github.com/avito-internships/test-backend-1-mmms914/internal/api/middleware"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/booking"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/room"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/schedule"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/slot"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/user"
)

type Router struct {
	router    *chi.Mux
	logger    *slog.Logger
	jwtSecret string

	authHandler     *handlers.AuthHandler
	roomHandler     *handlers.RoomHandler
	scheduleHandler *handlers.ScheduleHandler
	slotHandler     *handlers.SlotHandler
	bookingHandler  *handlers.BookingHandler
}

type Config struct {
	Logger          *slog.Logger
	JWTSecret       string
	ExpirationHours int

	UserService     *user.Service
	RoomService     *room.Service
	ScheduleService *schedule.Service
	SlotService     *slot.Service
	BookingService  *booking.Service
}

func NewRouter(c *Config) *Router {
	r := &Router{
		router:    chi.NewRouter(),
		jwtSecret: c.JWTSecret,
		logger:    c.Logger,
	}

	validate := validator.New()

	r.authHandler = handlers.NewAuthHandler(&handlers.AuthHandlerConfig{
		Logger:          c.Logger,
		Service:         c.UserService,
		Validate:        validate,
		JwtSecret:       c.JWTSecret,
		ExpirationHours: c.ExpirationHours,
	})
	r.roomHandler = handlers.NewRoomHandler(&handlers.RoomHandlerConfig{
		Logger:   c.Logger,
		Validate: validate,
		Service:  c.RoomService,
	})
	r.scheduleHandler = handlers.NewScheduleHandler(&handlers.ScheduleHandlerConfig{
		Service:  c.ScheduleService,
		Validate: validate,
		Logger:   c.Logger,
	})
	r.slotHandler = handlers.NewSlotHandler(&handlers.SlotHandlerConfig{
		Service: c.SlotService,
		Logger:  c.Logger,
	})
	r.bookingHandler = handlers.NewBookingHandler(&handlers.BookingHandlerConfig{
		Service:  c.BookingService,
		Logger:   c.Logger,
		Validate: validate,
	})

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

func (rt *Router) setupMiddleware() {
	rt.router.Use(middleware.Recoverer)
	rt.router.Use(appmiddleware.Logging(rt.logger))
}

func (rt *Router) setupRoutes() {
	rt.router.Get("/_info", rt.handleInfo)

	rt.router.Route("", func(r chi.Router) {
		r.Post("/dummyLogin", rt.authHandler.DummyLogin)
		r.Post("/login", rt.authHandler.Login)
		r.Post("/register", rt.authHandler.Register)

		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Auth(rt.jwtSecret))

			r.Route("/rooms", func(r chi.Router) {
				r.Get("/list", rt.roomHandler.List)
				r.Post("/create", rt.roomHandler.Create)

				r.Route("/{roomId}", func(r chi.Router) {
					r.Post("/schedule/create", rt.scheduleHandler.Create)

					r.Get("/slots/list", rt.slotHandler.List)
				})
			})

			r.Route("/bookings", func(r chi.Router) {
				r.Post("/create", rt.bookingHandler.Create)
				r.Get("/my", rt.bookingHandler.ListMy)
				r.Get("/list", rt.bookingHandler.List)
				r.Post("/{bookingId}/cancel", rt.bookingHandler.Cancel)
			})
		})
	})
}

func (rt *Router) handleInfo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
