package periodic

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/owdiscord/academy/internal/config"
	"github.com/owdiscord/academy/internal/database"
	"github.com/owdiscord/academy/internal/etl"
	"github.com/vinovest/sqlx"
)

var e *etl.Etl

func TestMain(m *testing.M) {
	// Load from dotenv so I can stay lazy for test running
	_ = godotenv.Load()
	config, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	dsn, err := database.URLtoDSN(config.DatabaseURI)
	if err != nil {
		log.Fatalf("could not parse database dsn: %v", err)
	}

	out, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("could not connect to out database: %s, %v", config.DatabaseURI, err)
	}

	// Only used for testing so we hard-code this value.
	mockDB, err := sqlx.Connect("mysql", "academy:academypass@tcp(localhost:3306)/modmail2?parseTime=true")
	if err != nil {
		log.Fatalf("could not connect to modmail database: %v", err)
	}

	start := time.Date(2020, time.January, 1, 1, 1, 1, 1, time.UTC)
	staff := database.Staff{
		ID:          1,
		WaveID:      67,
		Snowflake:   "999000000000000002",
		Username:    "isaac",
		DisplayName: "testisaac",
		Role:        "trainee",
	}
	e = etl.New(67, start, mockDB, mockDB, out, []database.Staff{staff}, config.PrivateChannels)

	out.MustExec("DELETE FROM case_notes;")
	out.MustExec("DELETE FROM cases;")
	out.MustExec("DELETE FROM thread_messages;")
	out.MustExec("DELETE FROM threads;")
	out.MustExec("DELETE FROM threads;")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRunFullImportCycle(t *testing.T) {
	if err := ImportData(t.Context(), e); err != nil {
		t.Fatal(err)
	}
}
