package components

import (
	"github.com/maddalax/htmgo/framework/h"
)

func GrabError(error string) *h.Element {
	return h.Div(
		h.P(h.TextF(error)),
	)
}

func SearchError(error string) *h.Element {
	return h.Div(
		h.Class("flex gap-1"),
		h.P(
			h.TextF(error),
		),
	)
}
