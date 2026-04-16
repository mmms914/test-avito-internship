package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mmms914/test-avito-internship/internal/api/handlers/auth"
	"github.com/mmms914/test-avito-internship/internal/api/handlers/booking"
	"github.com/mmms914/test-avito-internship/internal/api/handlers/room"
	"github.com/mmms914/test-avito-internship/internal/api/handlers/schedule"
	"github.com/mmms914/test-avito-internship/internal/api/handlers/slot"
	appmiddleware "github.com/mmms914/test-avito-internship/internal/api/middleware"
)

type Router struct {
	router    *chi.Mux
	logger    *slog.Logger
	jwtSecret string

	authHandler     *auth.Handler
	roomHandler     *room.Handler
	scheduleHandler *schedule.Handler
	slotHandler     *slot.Handler
	bookingHandler  *booking.Handler
}

type Config struct {
	Logger    *slog.Logger
	JWTSecret string

	AuthHandler     *auth.Handler
	RoomHandler     *room.Handler
	ScheduleHandler *schedule.Handler
	SlotHandler     *slot.Handler
	BookingHandler  *booking.Handler
}

func NewRouter(c *Config) *Router {
	r := &Router{
		router:    chi.NewRouter(),
		logger:    c.Logger,
		jwtSecret: c.JWTSecret,

		authHandler:     c.AuthHandler,
		roomHandler:     c.RoomHandler,
		scheduleHandler: c.ScheduleHandler,
		slotHandler:     c.SlotHandler,
		bookingHandler:  c.BookingHandler,
	}

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

	rt.router.Post("/dummyLogin", rt.authHandler.DummyLogin)
	rt.router.Post("/login", rt.authHandler.Login)
	rt.router.Post("/register", rt.authHandler.Register)

	rt.router.Group(func(r chi.Router) {
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
}

func (rt *Router) handleInfo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (rt *Router) Handler() http.Handler {
	return rt.router
}
