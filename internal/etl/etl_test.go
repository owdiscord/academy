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
	mm, err := sqlx.Connect("mysql", "modmail:modmailbot@tcp(localhost:3306)/modmail2?parseTime=true")
	if err != nil {
		log.Fatalf("could not connect to modmail database: %v", err)
	}

	start := time.Date(2020, time.January, 1, 1, 1, 1, 1, time.UTC)
	e = &Etl{
		startDate: start,
		waveID:    67,
		mmDB:      mm,
		athDB:     mm,
		outDB:     out,
	}

	out.MustExec("DELETE FROM case_notes;")
	out.MustExec("DELETE FROM cases;")
	out.MustExec("DELETE FROM thread_messages;")
	out.MustExec("DELETE FROM threads;")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestImportModmail(t *testing.T) {
	threads, err := e.FindAllTraineeThreads(t.Context(), []string{"204084691425427466"})
	if err != nil {
		t.Fatalf("could not get trainee threads from modmail database: %v", err)
	}

	if len(threads) < 1 {
		t.Fatal("no threads were found")
	}

	tx, err := e.OutTx()
	if err != nil {
		t.Fatalf("could not start transaction: %v", err)
	}

	for _, thread := range threads {
		if err := e.InsertImportedThread(t.Context(), tx, thread); err != nil {
			t.Fatalf("could not insert thread: %v", err)
		}

		msgs, err := e.FindThreadMessages(t.Context(), thread.ID.String())
		if err != nil {
			t.Fatalf("could not get thread messages for thread %s: %v", thread.ID.String(), err)
		}

		for _, msg := range msgs {
			if err := e.InsertThreadMessage(t.Context(), tx, msg); err != nil {
				t.Fatalf("could not insert thread message: %v", err)
			}
		}

		if err := e.RecalculateThreadMessageCounts(t.Context(), tx, thread.ID); err != nil {
			t.Fatalf("could not recalculate thread message counts: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("could not commit transaction: %v", err)
	}
}

func TestImportAthena(t *testing.T) {
	cases, err := e.FindAllTraineeCases(t.Context(), []string{"204084691425427466"})
	if err != nil {
		t.Fatalf("could not get trainee cases from athena database: %v", err)
	}

	if len(cases) < 1 {
		t.Fatal("no cases were found")
	}

	tx, err := e.OutTx()
	if err != nil {
		t.Fatalf("could not start transaction: %v", err)
	}

	for _, modCase := range cases {
		if err := e.InsertImportedCase(t.Context(), tx, modCase); err != nil {
			t.Fatalf("could not insert case: %v", err)
		}

		notes, err := e.FindCaseNotes(t.Context(), modCase.ID)
		if err != nil {
			t.Fatalf("could not get notes on case #%d: %v", modCase.ID, err)
		}
		for _, note := range notes {
			if err := e.InsertCaseNote(t.Context(), tx, note); err != nil {
				t.Fatalf("could not insert case note: %v", err)
			}
		}

	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("could not commit transaction: %v", err)
	}
}
