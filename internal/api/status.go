package api

import (
	"github.com/maddalax/htmgo/framework/h"
	"golift.io/starr/readarr"
)

func Status() *h.Element {
	status, err := getStatus()
	if err != nil {
		return h.Div(
			h.H4(
				h.Class("flex dark:text-white flex-col gap-2 mx-auto"),
				h.Text("Failed to get Readarr version, maybe the API key is wrong? Proceed with caution."),
				h.TextF("Error: %s", err),
			),
		)
	}
	return h.Div(
		h.H4(
			h.Class("flex dark:text-white flex-col gap-2 mx-auto"),
			h.TextF("Readarr version: %s", status.Version),
		),
	)
}

func getStatus() (*readarr.SystemStatus, error) {
	handler := Connect()
	status, err := handler.GetSystemStatus()
	return status, err
}
