// Package periodic contains our cron backed functions for pulling and updating data
package periodic

import (
	"context"

	"github.com/owdiscord/academy/internal/etl"
)

func ImportData(ctx context.Context, e *etl.Etl) error {
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
	threads, err := e.FindAllTraineeThreads(ctx)
	if err != nil {
		return err
	}

	for _, thread := range threads {
		if thread.ClosedByID != nil {
			e.IncreaseCloseStat(*thread.ClosedByID, 1)
		}
	}

	return nil
}

// --- Utils ----
