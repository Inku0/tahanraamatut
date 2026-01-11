package partials

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"golift.io/starr"

	"tahanraamatut/internal/api"
	"tahanraamatut/internal/components"
)

func GrabStatusWidget(ctx *h.RequestContext, id int64) *h.Element {
	return h.Div(
		h.GetPartialWithQs(GrabStatusPartial, h.NewQs("id", strconv.FormatInt(id, 10)), "load, every 3s"),
	)
}

func GrabStatusPartial(ctx *h.RequestContext) *h.Partial {
	bookID := ctx.QueryParam("id")
	if bookID == "" {
		return h.NewPartial(
			components.GrabError("No bookID provided, unable to query status..."),
		)
	}

	locator := ctx.ServiceLocator()
	handler := service.Get[api.ReadarrService](locator)

	intBookID, _ := strconv.ParseInt(bookID, 10, 64)
	record, err := handler.GetBookQueueStatus(ctx.Request.Context(), intBookID)
	if err != nil {
		if errors.Is(err, api.ErrNoResults) && api.WasGrabbed(intBookID) {
			api.ClearGrabbed(intBookID)
			return h.NewPartial(h.P(h.Text("Grab completed!")))
		}
		errStr := fmt.Sprintf("Unable to get queue status: %v", err)
		return h.NewPartial(
			components.GrabError(errStr),
		)
	}

	return h.NewPartial(
		h.Ul(
			h.List(record.StatusMessages, func(msg *starr.StatusMessage, _ int) *h.Element {
				return h.Li(
					h.Class("status-message"),
					h.P(h.TextF(msg.Title)),
					h.List(msg.Messages, func(item string, _ int) *h.Element {
						return h.P(h.TextF(item))
					}),
				)
			}),
			h.Li(h.TextF("Status: %s", record.Status)),
			h.Li(h.TextF("Messages: %+v", record.StatusMessages)),
			h.Li(h.TextF("Download ID: %s", record.DownloadID)),
			h.Li(h.TextF("Time left: %s", record.Timeleft)),
			h.Li(h.TextF("ETA: %s", record.EstimatedCompletionTime)),
			h.ElementIf(
				len(strings.TrimSpace(record.ErrorMessage)) != 0,
				h.Li(h.TextF("Errors: %s", record.ErrorMessage))),
		),
	)
}
