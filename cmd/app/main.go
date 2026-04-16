package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"

	"github.com/mmms914/test-avito-internship/cmd/config"
	"github.com/mmms914/test-avito-internship/internal/api"
	authH "github.com/mmms914/test-avito-internship/internal/api/handlers/auth"
	bookingH "github.com/mmms914/test-avito-internship/internal/api/handlers/booking"
	roomH "github.com/mmms914/test-avito-internship/internal/api/handlers/room"
	scheduleH "github.com/mmms914/test-avito-internship/internal/api/handlers/schedule"
	slotH "github.com/mmms914/test-avito-internship/internal/api/handlers/slot"
	"github.com/mmms914/test-avito-internship/internal/infra/hasher"
	"github.com/mmms914/test-avito-internship/internal/infra/postgres"
	"github.com/mmms914/test-avito-internship/internal/infra/zoom"
	bookingR "github.com/mmms914/test-avito-internship/internal/repository/booking"
	roomR "github.com/mmms914/test-avito-internship/internal/repository/room"
	scheduleR "github.com/mmms914/test-avito-internship/internal/repository/schedule"
	slotR "github.com/mmms914/test-avito-internship/internal/repository/slot"
	userR "github.com/mmms914/test-avito-internship/internal/repository/user"
	"github.com/mmms914/test-avito-internship/internal/service/booking"
	"github.com/mmms914/test-avito-internship/internal/service/room"
	"github.com/mmms914/test-avito-internship/internal/service/schedule"
	"github.com/mmms914/test-avito-internship/internal/service/slot"
	"github.com/mmms914/test-avito-internship/internal/service/user"
)

// @title           Room Booking Service API
// @version         1.0.0
// @description     Сервис бронирования переговорок

// @host      localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	appConfig, err := config.LoadConfig("./cmd/config/config.yaml")
	if err != nil {
		log.Fatalf("Can't load config file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	validate := validator.New()
	hash := hasher.NewHasher()

	db := initDatabase(appConfig, logger)
	defer db.Close()

	repos := initRepositories(db)
	services := initServices(repos, logger, hash)
	handlers := initHandlers(services, appConfig, logger, validate)

	router := api.NewRouter(&api.Config{
		Logger:          logger,
		JWTSecret:       appConfig.Auth.JWTSecret,
		AuthHandler:     handlers.auth,
		RoomHandler:     handlers.room,
		ScheduleHandler: handlers.schedule,
		SlotHandler:     handlers.slot,
		BookingHandler:  handlers.booking,
	})

	runServer(appConfig, router, logger)
}

func initDatabase(cfg *config.Config, logger *slog.Logger) *sql.DB {
	ctx := context.Background()
	db, err := postgres.Connect(ctx, &postgres.Config{
		Host:               cfg.Database.Host,
		Port:               cfg.Database.Port,
		User:               cfg.Database.User,
		Password:           cfg.Database.Password,
		DBName:             cfg.Database.Name,
		SSLMode:            cfg.Database.SSLMode,
		ConnectTimeout:     cfg.Database.ConnectTimeout,
		MaxIdleConnections: cfg.Database.MaxIdleConns,
		MaxConnections:     cfg.Database.MaxOpenConns,
		MaxConnLifetime:    cfg.Database.MaxConnLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger.Info("Connected to database")

	if err = postgres.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	logger.Info("Migrations completed")
	return db
}

type Repositories struct {
	room     *roomR.Repository
	user     *userR.Repository
	schedule *scheduleR.Repository
	slot     *slotR.Repository
	booking  *bookingR.Repository
}

func initRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		room:     roomR.NewRepository(db),
		user:     userR.NewRepository(db),
		schedule: scheduleR.NewRepository(db),
		slot:     slotR.NewRepository(db),
		booking:  bookingR.NewRepository(db),
	}
}

type Services struct {
	user     *user.Service
	room     *room.Service
	schedule *schedule.Service
	slot     *slot.Service
	booking  *booking.Service
}

func initServices(repos *Repositories, logger *slog.Logger, hash *hasher.Hasher) *Services {
	userService := user.NewService(&user.Config{
		UserRepo: repos.user,
		Hasher:   hash,
	})

	roomService := room.NewService(repos.room)

	scheduleService := schedule.NewService(&schedule.Config{
		RoomRepo:     repos.room,
		ScheduleRepo: repos.schedule,
	})

	slotService := slot.NewService(&slot.Config{
		SlotRepo:     repos.slot,
		ScheduleRepo: repos.schedule,
		RoomRepo:     repos.room,
	})

	bookingService := booking.NewService(&booking.Config{
		Logger:      logger,
		BookingRepo: repos.booking,
		SlotRepo:    repos.slot,
		LinkManager: zoom.NewZoom(),
	})

	return &Services{
		user:     userService,
		room:     roomService,
		schedule: scheduleService,
		slot:     slotService,
		booking:  bookingService,
	}
}

type Handlers struct {
	auth     *authH.Handler
	room     *roomH.Handler
	schedule *scheduleH.Handler
	slot     *slotH.Handler
	booking  *bookingH.Handler
}

func initHandlers(services *Services, cfg *config.Config, logger *slog.Logger, validate *validator.Validate) *Handlers {
	authHandler := authH.NewHandler(&authH.HandlerConfig{
		Logger:          logger,
		Service:         services.user,
		Validate:        validate,
		JwtSecret:       cfg.Auth.JWTSecret,
		ExpirationHours: cfg.Auth.TokenDurationHours,
	})

	roomHandler := roomH.NewHandler(&roomH.HandlerConfig{
		Logger:   logger,
		Validate: validate,
		Service:  services.room,
	})

	scheduleHandler := scheduleH.NewHandler(&scheduleH.HandlerConfig{
		Service:  services.schedule,
		Validate: validate,
		Logger:   logger,
	})

	slotHandler := slotH.NewHandler(&slotH.HandlerConfig{
		Service: services.slot,
		Logger:  logger,
	})

	bookingHandler := bookingH.NewHandler(&bookingH.HandlerConfig{
		Service:  services.booking,
		Logger:   logger,
		Validate: validate,
	})

	return &Handlers{
		auth:     authHandler,
		room:     roomHandler,
		schedule: scheduleHandler,
		slot:     slotHandler,
		booking:  bookingHandler,
	}
}

func runServer(cfg *config.Config, router *api.Router, logger *slog.Logger) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("Server started", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error", "error", err)
		return
	}

	logger.Info("Server stopped successfully")
}
