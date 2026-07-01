// Package periodic contains our cron backed functions for pulling and updating data
package periodic

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/owdiscord/academy/internal/config"
	"github.com/owdiscord/academy/internal/database"
	"github.com/owdiscord/academy/internal/etl"
	"github.com/vinovest/sqlx"
)

// I need to, in order
//
// 01. Get new threads & messages
// 02. Insert new threads & messages and updated values, cache count per-user???
// 03. Get new cases & notes
// 04. Insert new cases & case notes
// 05. Get message count in time since last check (?)
// 06. Set stats for $NOW
//     - Messages sent privately
//     - Messages sent publicly
//     - Cases created
//     - ^ This means created_at > $start
//     - Thread messages sent
//     	 - Replies
//     	 - Chat
// 		 - ^ This means created_at > $start
//     - Threads closed
//       - WHERE closed_by_id = $id AND (created_at > $START OR updated_at > $START)

type Manager struct {
	cfg       config.Config
	scheduler gocron.Scheduler
}

func NewManager(cfg config.Config) (*Manager, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Manager{
		cfg:       cfg,
		scheduler: s,
	}, nil
}

func (m *Manager) AddImportJob(athDB *sqlx.DB, mmDB *sqlx.DB, outDB *database.DB) {
	job, err := m.scheduler.NewJob(gocron.CronJob("0 * */3 * *", false), gocron.NewTask(func() {
		ctx := context.Background()
		wave, err := outDB.GetLatestWave(ctx)
		if err != nil {
			slog.Default().Error("could not get latest wave in import job", "err", err)
		}
		staff, err := outDB.GetWaveTrainees(ctx, wave.ID)
		if err != nil {
			slog.Default().Error("could not get trainees in import job", "err", err)
		}

		start := time.Now().Add(3 * (time.Minute * -1))
		e := etl.New(67, start, athDB, mmDB, outDB.Conn(), staff, m.cfg.PrivateChannels)
		ImportData(ctx, e)
	}))
	if err != nil {
		fmt.Printf("[task] could not start order cleanup job: %v\n", err)
		return
	}

	slog.Default().Info("adding import job to task queue", "cron", "0 * */3 * *", "id", job.ID())
}

func (m *Manager) Start() {
	m.scheduler.Start()
}

func ImportData(ctx context.Context, e *etl.Etl) error {
	tx, err := e.OutTx()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	//
	// # Modmail threads & thread replies / messages
	//

	threads, err := e.FindAllTraineeThreads(ctx)
	if err != nil {
		return fmt.Errorf("ImportData: fetching threads: %w", err)
	}

	for _, thread := range threads {
		if err := e.InsertImportedThread(ctx, tx, thread); err != nil {
			return fmt.Errorf("ImportData: inserting thread %s: %w", thread.ID, err)
		}

		// Track thread closures per closer
		if thread.ClosedByID != nil {
			e.IncreaseCloseStat(*thread.ClosedByID, 1)
		}

		// Fetch and insert all messages for this thread
		threadIDStr := thread.ID.String()
		messages, err := e.FindThreadMessages(ctx, threadIDStr)
		if err != nil {
			return fmt.Errorf("could not fetch messages for thread %s: %w", thread.ID, err)
		}

		for _, msg := range messages {
			if err := e.InsertThreadMessage(ctx, tx, msg); err != nil {
				return fmt.Errorf("could not insert message %d: %w", msg.ID, err)
			}

			// kind 1 = inbound (user reply), 2 = outbound (staff reply), 3 = internal chat
			// Accumulate per-user stats for messages created after the wave start.
			// FindThreadMessages returns all messages for the thread; the created_at
			// filter is intentionally loose here — SaveDateStats uses CURRENT_DATE so
			// double-counting on re-runs is prevented by the ON DUPLICATE KEY logic.
			switch msg.Kind {
			case 2:
				e.IncreaseThreadReplyStat(msg.UserID, 1)
			case 3:
				e.IncreaseThreadChatStat(msg.UserID, 1)
			}
		}

		// Recalculate denormalised counts now that messages are inserted
		if err := e.RecalculateThreadMessageCounts(ctx, tx, thread.ID); err != nil {
			return fmt.Errorf("could not recalculate counts for thread %s: %w", thread.ID, err)
		}
	}

	//
	// # Athena cases and case notes
	//

	cases, err := e.FindAllTraineeCases(ctx)
	if err != nil {
		return fmt.Errorf("could not fetch cases: %w", err)
	}

	for _, c := range cases {
		if err := e.InsertImportedCase(ctx, tx, c); err != nil {
			return fmt.Errorf("could not insert case %d: %w", c.ID, err)
		}

		// ModID is a *uint64 but stat keys are string snowflakes — convert carefully
		if c.ModID != nil {
			modSnowflake := fmt.Sprintf("%d", *c.ModID)
			e.IncreaseCasesStat(modSnowflake, 1)
		}

		notes, err := e.FindCaseNotes(ctx, c.ID)
		if err != nil {
			return fmt.Errorf("could not fetch notes for case %d: %w", c.ID, err)
		}

		for _, note := range notes {
			if err := e.InsertCaseNote(ctx, tx, note); err != nil {
				return fmt.Errorf("could not insert case note %d: %w", note.ID, err)
			}
		}
	}

	//
	// # Message counts (from Athena, for public / private messages)
	//

	msgStats, err := e.GetMessageStats(ctx, tx)
	if err != nil {
		return fmt.Errorf("could not fetch message stats: %w", err)
	}

	for _, stat := range msgStats {
		e.IncreasePublicMsgStat(stat.UserID, stat.Public)
		e.IncreasePrivateMsgStat(stat.UserID, stat.Private)
	}

	//
	// # Save all stats
	//

	if err := e.SaveAllDateStats(ctx, tx); err != nil {
		return fmt.Errorf("could not save stats: %w", err)
	}

	//
	// # Commit the transaction
	//

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
