//go:build integration
// +build integration

package integration_test

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/booking"
	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/room"
	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/schedule"
	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/slot"
	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/user"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx          context.Context
	db           *sql.DB
	userRepo     *user.Repository
	roomRepo     *room.Repository
	scheduleRepo *schedule.Repository
	slotRepo     *slot.Repository
	bookingRepo  *booking.Repository
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()

	connStr := "host=localhost port=5433 user=test_user password=test_password dbname=test_db sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	require.NoError(s.T(), err)

	err = db.PingContext(s.ctx)
	require.NoError(s.T(), err)

	s.db = db
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *IntegrationTestSuite) SetupTest() {
	s.userRepo = user.NewRepository(s.db)
	s.roomRepo = room.NewRepository(s.db)
	s.scheduleRepo = schedule.NewRepository(s.db)
	s.slotRepo = slot.NewRepository(s.db)
	s.bookingRepo = booking.NewRepository(s.db)

	s.cleanupTables()
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.cleanupTables()
}

func (s *IntegrationTestSuite) cleanupTables() {
	tables := []string{"bookings", "slots", "schedules", "rooms", "users"}
	for _, table := range tables {
		s.db.ExecContext(s.ctx, "DELETE FROM "+table)
	}
}
