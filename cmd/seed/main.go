package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/auth"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/infra/postgres"
	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/converter"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func main() {
	ctx := context.Background()

	// Подключение к БД
	db, err := postgres.Connect(ctx, &postgres.Config{
		Host:               "localhost",
		Port:               5432, //nolint:mnd // seeding
		User:               "postgres",
		Password:           "postgres",
		DBName:             "myapp",
		SSLMode:            "disable",
		ConnectTimeout:     10 * time.Second, //nolint:mnd // seeding
		MaxConnections:     5,                //nolint:mnd // seeding
		MaxIdleConnections: 5,                //nolint:mnd // seeding
		MaxConnLifetime:    time.Minute,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Очищаем существующие данные
	log.Println("Cleaning existing data...")
	tables := []string{"bookings", "slots", "schedules", "rooms", "users"}
	for _, table := range tables {
		_, err = db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("Failed to clean table %s: %v", table, err)
		}
	}

	// Создаем пользователей
	log.Println("Creating users...")
	users := createUsers(db)

	// Создаем переговорки
	log.Println("Creating rooms...")
	rooms := createRooms(db)

	// Создаем расписания
	log.Println("Creating schedules...")
	schedules := createSchedules(db, rooms)

	// Создаем слоты
	log.Println("Creating slots...")
	slots := createSlots(db, schedules)

	// Создаем бронирования
	log.Println("Creating bookings...")
	createBookings(db, slots, users)

	log.Println("Seeding completed!")
}

func createUsers(db *sql.DB) []*domain.User {
	users := []*domain.User{}

	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           uuid.MustParse(auth.DefaultAdminUID),
		Email:        "admin@example.com",
		PasswordHash: string(adminPassword),
		Role:         domain.AdminRole,
		CreatedAt:    time.Now().UTC(),
	}))
	users = append(users, admin)

	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           uuid.MustParse(auth.DefaultUserUID),
		Email:        "user@example.com",
		PasswordHash: string(userPassword),
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))
	users = append(users, user)

	for i := 1; i <= 5; i++ {
		pass, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
		u := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
			ID:           uuid.New(),
			Email:        fmt.Sprintf("user%d@example.com", i),
			PasswordHash: string(pass),
			Role:         domain.UserRole,
			CreatedAt:    time.Now().UTC(),
		}))
		users = append(users, u)
	}

	for _, u := range users {
		query := `INSERT INTO users (id, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.ExecContext(context.Background(), query, u.ID(), u.Email(),
			u.PasswordHash(), u.Role().String(), u.CreatedAt())
		if err != nil {
			log.Printf("Failed to create user %s: %v", u.Email(), err)
		}
	}

	return users
}

func createRooms(db *sql.DB) []*domain.Room {
	rooms := []*domain.Room{
		domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
			ID:          uuid.New(),
			Name:        "Конференц-зал А",
			Description: ptr.To("Большой зал с проектором, до 20 человек"),
			Capacity:    ptr.To(20), //nolint:mnd // seeding
			CreatedAt:   time.Now().UTC(),
		})),
		domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
			ID:          uuid.New(),
			Name:        "Переговорная Б",
			Description: ptr.To("Средний зал для совещаний, до 10 человек"),
			Capacity:    ptr.To(10), //nolint:mnd // seeding
			CreatedAt:   time.Now().UTC(),
		})),
		domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
			ID:          uuid.New(),
			Name:        "Кабинет В",
			Description: ptr.To("Маленькая комната для переговоров, до 5 человек"),
			Capacity:    ptr.To(5), //nolint:mnd // seeding
			CreatedAt:   time.Now().UTC(),
		})),
		domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
			ID:          uuid.New(),
			Name:        "Зал для презентаций",
			Description: ptr.To("Специальное оборудование для презентаций"),
			Capacity:    ptr.To(30), //nolint:mnd // seeding
			CreatedAt:   time.Now().UTC(),
		})),
		domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
			ID:          uuid.New(),
			Name:        "Виртуальная комната",
			Description: nil,
			Capacity:    nil,
			CreatedAt:   time.Now().UTC(),
		})),
	}

	for _, r := range rooms {
		query := `INSERT INTO rooms (id, name, description, capacity, created_at) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.ExecContext(context.Background(), query, r.ID(), r.Name(), r.Description(), r.Capacity(), r.CreatedAt())
		if err != nil {
			log.Printf("Failed to create room %s: %v", r.Name(), err)
		}
	}

	return rooms
}

func createSchedules(db *sql.DB, rooms []*domain.Room) []*domain.Schedule {
	schedules := []*domain.Schedule{}

	weekdays := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
	}

	for _, room := range rooms[:3] {
		schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
			ID:         uuid.New(),
			RoomID:     room.ID(),
			DaysOfWeek: weekdays,
			StartTime:  9 * time.Hour,  //nolint:mnd // seeding
			EndTime:    18 * time.Hour, //nolint:mnd // seeding
		}))
		schedules = append(schedules, schedule)
	}

	scheduleWeekend := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         uuid.New(),
		RoomID:     rooms[3].ID(),
		DaysOfWeek: []time.Weekday{time.Saturday, time.Sunday},
		StartTime:  10 * time.Hour, //nolint:mnd // seeding
		EndTime:    16 * time.Hour, //nolint:mnd // seeding
	}))
	schedules = append(schedules, scheduleWeekend)

	for _, s := range schedules {
		query := `INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`
		daysOfWeek := converter.WeekdaysToIntArray(s.DaysOfWeek())
		_, err := db.ExecContext(context.Background(), query, s.ID(), s.RoomID(), daysOfWeek,
			converter.DurationToStringTime(s.StartTime()), converter.DurationToStringTime(s.EndTime()))
		if err != nil {
			log.Printf("Failed to create schedule for room %s: %v", s.RoomID(), err)
		}
	}

	return schedules
}

func createSlots(db *sql.DB, schedules []*domain.Schedule) []*domain.Slot {
	allSlots := []*domain.Slot{}

	for i := 1; i <= 14; i++ {
		date := time.Now().UTC().Add(time.Duration(i) * 24 * time.Hour).Truncate(24 * time.Hour) //nolint:mnd // 1 day

		for _, schedule := range schedules {
			if !schedule.IsAppliedForDate(date) {
				continue
			}

			startTime := date.Add(schedule.StartTime())
			endTime := date.Add(schedule.EndTime())

			current := startTime
			for current.Before(endTime) {
				slot := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    schedule.RoomID(),
					StartTime: current,
					EndTime:   current.Add(domain.SlotDuration),
				}))
				allSlots = append(allSlots, slot)
				current = current.Add(domain.SlotDuration)
			}
		}
	}

	for _, slot := range allSlots {
		query := `INSERT INTO slots (id, room_id, start_time, end_time) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
		_, err := db.ExecContext(context.Background(), query, slot.ID(), slot.RoomID(), slot.StartTime(), slot.EndTime())
		if err != nil {
			log.Printf("Failed to create slot: %v", err)
		}
	}

	log.Printf("Created %d slots", len(allSlots))
	return allSlots
}

func createBookings(db *sql.DB, slots []*domain.Slot, users []*domain.User) {
	user := users[1] // обычный пользователь
	userID := user.ID()

	bookingsCreated := 0

	for i := 0; i < 10 && i < len(slots); i++ {
		slot := slots[i*3]

		booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         slot.ID(),
			UserID:         userID,
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: nil,
			CreatedAt:      time.Now().UTC(),
		}))

		query := `INSERT INTO bookings (id, slot_id, user_id, status, created_at) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.ExecContext(context.Background(), query, booking.ID(), booking.SlotID(),
			booking.UserID(), booking.Status().String(), booking.CreatedAt())
		if err != nil {
			log.Printf("Failed to create booking: %v", err)
		} else {
			bookingsCreated++
		}
	}

	log.Printf("Created %d bookings", bookingsCreated)
}
