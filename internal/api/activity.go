package api

import (
	"context"
	"log"

	"golift.io/starr/readarr"
)

// CleanFailedAdd deletes an added book if no suitable sources for it were found
func (service *ReadarrService) CleanFailedAdd(ctx context.Context, grab *readarr.Book) error {
	return service.Client.DeleteBookContext(ctx, grab.ID, true, false)
}

// isInQueue checks whether a grabbed book got placed into the queue
func (service *ReadarrService) isInQueue(ctx context.Context, grab *readarr.Book) (bool, error) {
	queue, err := service.Client.GetQueueContext(ctx, 0, 0)
	if err != nil {
		return false, err
	}

	// log.Printf("queue length: %d", len(queue.Records))

	for _, record := range queue.Records {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		if record == nil {
			// log.Printf("record[%d] is nil", i)
			continue
		}

		// log.Printf("record[%d] Title=%s AuthorID=%d BookID=%d\n", i, record.Title, record.AuthorID, record.BookID)
		if record.BookID == grab.ID && record.AuthorID == grab.AuthorID {
			log.Printf("match found for bookID=%d authorID=%d\n", grab.ID, grab.AuthorID)
			return true, nil
		}
	}

	return false, nil
}

// isInHistory checks whether a book was mentioned in the last 100 elements of history
func (service *ReadarrService) isInHistory(ctx context.Context, grab *readarr.Book) (bool, error) {
	history, err := service.Client.GetHistoryContext(ctx, 100, 0)
	if err != nil {
		return false, err
	}

	// log.Printf("history length: %d", len(history.Records))

	for _, record := range history.Records {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		// log.Printf("record[%d] Event=%s AuthorID=%d BookID=%d\n", i, record.EventType, record.AuthorID, record.BookID)
		if record.BookID == grab.ID && record.AuthorID == grab.AuthorID {
			log.Printf("match found for bookID=%d authorID=%d\n", grab.ID, grab.AuthorID)
			return true, nil
		}
	}

	return false, nil
}

// GotGrabbed returns a heuristic for determining whether a book was grabbed and/or is downloading/ed
func (service *ReadarrService) GotGrabbed(ctx context.Context, grab *readarr.Book) (bool, error) {
	isQueued, err := service.isInQueue(ctx, grab)
	if err != nil {
		return false, err
	}

	isHistory, err := service.isInHistory(ctx, grab)
	if err != nil {
		return false, err
	}

	return isHistory || isQueued, nil
}
