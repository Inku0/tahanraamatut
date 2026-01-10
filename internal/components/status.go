package components

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"

	"tahanraamatut/internal/api"
)

func Status(ctx *h.RequestContext) *h.Element {
	locator := ctx.ServiceLocator()
	readarr := service.Get[api.ReadarrService](locator)
	status, err := readarr.GetStatus(ctx.Request.Context())

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
