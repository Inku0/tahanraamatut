package api

import (
	"golift.io/starr/readarr"
)

func CleanFailedAdd(grab *readarr.Book) error {
	handler := Connect()

	err := handler.DeleteBook(grab.ID, true, false)
	if err != nil {
		return err
	}

	return nil
}
