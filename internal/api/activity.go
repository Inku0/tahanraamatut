package api

import (
	"log"

	"golift.io/starr/readarr"
)

// IsInQueue checks whether a grabbed book got placed into the queue
func isInQueue(grab *readarr.Book) (bool, error) {
	handler := Connect()

	queue, err := handler.GetQueue(0, 0)
	if err != nil {
		return false, err
	}

	// log.Printf("queue length: %d", len(queue.Records))

	for _, record := range queue.Records {
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

func isInHistory(grab *readarr.Book) (bool, error) {
	handler := Connect()

	history, err := handler.GetHistory(100, 0)
	if err != nil {
		return false, err
	}

	// log.Printf("history length: %d", len(history.Records))

	for _, record := range history.Records {
		// log.Printf("record[%d] Event=%s AuthorID=%d BookID=%d\n", i, record.EventType, record.AuthorID, record.BookID)
		if record.BookID == grab.ID && record.AuthorID == grab.AuthorID {
			log.Printf("match found for bookID=%d authorID=%d\n", grab.ID, grab.AuthorID)
			return true, nil
		}
	}

	return false, nil
}

func GotGrabbed(grab *readarr.Book) (bool, error) {
	isQueued, err := isInQueue(grab)
	if err != nil {
		return false, err
	}

	isHistory, err := isInHistory(grab)
	if err != nil {
		return false, err
	}
	return isHistory || isQueued, nil
}
