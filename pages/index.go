package pages

import (
	"github.com/maddalax/htmgo/framework/h"

	"tahanraamatut/internal/components"
	"tahanraamatut/partials"
)

func IndexPage(ctx *h.RequestContext) *h.Page {
	return RootPage(
		h.Div(
			h.Class("flex-auto flex flex-col gap-4 dark:bg-gray-800 items-center justify-center min-h-screen bg-neutral-100"),
			h.H3(
				h.Id("intro-text"),
				h.Text("okei, millist?"),
				h.Class("text-5xl dark:text-white flex-3"),
			),
			partials.Search(),
			h.Div(
				h.Class("my-2"),
				components.Status(),
			),
		),
	)
}
