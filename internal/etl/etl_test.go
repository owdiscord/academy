package etl

import (
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vinovest/sqlx"
)

var e *Etl

func TestMain(m *testing.M) {
	out, err := sqlx.Connect("mysql", "academy:academypass@tcp(localhost:3306)/academy?parseTime=true")
	if err != nil {
		log.Fatalf("could not connect to out database: %v", err)
	}
	mm, err := sqlx.Connect("mysql", "modmail:modmailbot@tcp(localhost:3306)/modmail?parseTime=true")
	if err != nil {
		log.Fatalf("could not connect to modmail database: %v", err)
	}

	start := time.Date(2020, time.January, 1, 1, 1, 1, 1, time.UTC)
	e = &Etl{
		startDate: start,
		mmDB:      mm,
		outDB:     out,
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestImportModmail(t *testing.T) {
	thrds, err := e.FindAllTraineeThreads(t.Context(), []string{"204084691425427466"})
	if err != nil {
		t.Fatalf("could not get trainee threads from modmail database: %v", err)
	}

	if len(thrds) < 1 {
		t.Fatal("no threads were found")
	}

	tx, err := e.OutTx()
	if err != nil {
		t.Fatalf("could not start transaction: %v", err)
	}

	for _, thread := range thrds {
		msgs, err := e.FindThreadMessages(t.Context(), thread.ID.String())
		if err != nil {
			t.Fatalf("could not get thread messages for thread %s: %v", thread.ID.String(), err)
		}
		for _, msg := range msgs {
			if err := e.InsertThreadMessage(t.Context(), tx, msg); err != nil {
				t.Fatalf("could not insert thread message: %v", err)
			}
		}

		if err := e.InsertImportedThread(t.Context(), tx, thread); err != nil {
			t.Fatalf("could not insert thread: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("could not commit transaction: %v", err)
	}
}
