package api

import (
	"golift.io/starr/readarr"
)

type BookStatus struct {
	progress int
	state    string
}

func StartSearch(bookID int64) (*readarr.CommandResponse, error) {
	handler := Connect()
	bookIDs := []int64{bookID}

	command := readarr.CommandRequest{
		Name:    "BookSearch",
		BookIDs: bookIDs,
	}

	resp, err := handler.SendCommand(&command)
	return resp, err
}
