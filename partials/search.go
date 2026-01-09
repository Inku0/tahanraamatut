package partials

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/hx"
	"github.com/maddalax/htmgo/framework/js"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/api"
)

func Search() *h.Partial {
	return h.NewPartial(
		h.Form(
			h.TriggerChildren(),
			h.PostPartial(SubmitForm),
			h.Class("flex dark:text-white flex-col gap-2 mx-auto"),
			h.LabelFor("name", "Raamatu nimi:"),
			h.Input(
				"text",
				h.Required(),
				h.Type("search"),
				h.Role("search"),
				h.Class("p-4 dark:text-black rounded-md border border-slate-200"),
				h.Name("name"),
				h.Attribute("autocomplete", "off"),
				h.Placeholder("East of Eden...?"),
				h.OnEvent(
					hx.KeyDownEvent,
					js.SubmitFormOnEnter(),
				),
			),
			SearchSubmitButton(),
		),
	)
}

func SearchSubmitButton() *h.Element {
	buttonClasses := "rounded dark:border-2 dark:border-white-200 items-center px-3 py-2 bg-white-800 dark:text-white w-full text-center"
	return h.Div(
		h.HxBeforeRequest(
			js.RemoveClassOnChildren(".loading", "hidden"),
			js.SetClassOnChildren(".submit", "hidden"),
		),
		h.HxAfterRequest(
			js.SetClassOnChildren(".loading", "hidden"),
			js.RemoveClassOnChildren(".submit", "hidden"),
		),
		h.Class("flex gap-2 justify-center"),
		h.Button(
			h.Class("loading hidden relative text-center", buttonClasses),
			Spinner(),
			h.Disabled(),
			h.Text("searching..."),
		),
		h.Button(
			h.Type("submit"),
			h.Class("submit", buttonClasses),
			h.Text("search"),
		),
	)
}

func Spinner(children ...h.Ren) *h.Element {
	return h.Div(
		h.Children(children...),
		h.Class("absolute left-1 spinner spinner-border animate-spin "+
			"inline-block w-6 h-6 border-4 rounded-full dark:border-slate-200 border-black-200 border-t-transparent"),
		h.Attribute("role", "status"),
	)
}

func SubmitForm(ctx *h.RequestContext) *h.Partial {
	name := ctx.FormValue("name")

	handler := api.Connect()

	searchResults, err := handler.Search(name)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.TextF("Encountered a fatal error while searching for a book named: %s", name)),
				h.P(h.TextF("error: %s", err)),
			),
		)
	} else if len(searchResults) == 0 {
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.TextF("No results for a book named: \"%s\"", name)),
			),
		)
	}

	return h.NewPartial(
		h.Div(
			h.Class("flex flex-col mx-12"),
			h.P(h.TextF("Searched for: %s", name)),
			h.Hr(),
			h.List(searchResults, func(book *readarr.SearchResult, index int) *h.Element {
				if book.Book == nil {
					return h.Empty()
				}

				formattedTime := book.Book.ReleaseDate.Format("2006-01-02")

				if len(book.Book.Editions) == 0 {
					return h.Div(
						h.Class("flex gap-1"),
						h.P(
							h.TextF("Found no editions for a book named %s", book.Book.Title),
						),
					)
				}

				return h.Div(
					h.Div(
						h.Class("mt-2 flex flex-row gap-y-4"),
						h.Div(
							h.Class("max-w-[70rem]"),
							h.P(
								h.TextF("\"%s\" by %s", book.Book.Title, book.Book.Author.AuthorName),
							),
							h.Ul(
								h.Class("list-disc list-inside"),
								h.Li(
									h.TextF("Page count: %d", book.Book.PageCount),
								),
								h.Li(
									h.TextF("Release date: %s", formattedTime),
								),
								h.Li(
									h.TextF("Already exists?: %t", book.Book.Monitored),
								),
								h.Li(
									h.TextF("Grabbed (aka already downloading)?: %t", book.Book.Grabbed),
								),
							),
							h.P(
								h.Class("p-4 text-justify flex-[10]"),
								h.TextF("%s", book.Book.Overview),
							),
						),

						h.Button(),

						h.List(book.Book.Editions, func(edition *readarr.Edition, index int) *h.Element {
							var cover *h.Element
							if len(edition.Images) > 0 {
								image := edition.Images[0]
								cover = h.Img(
									h.Src(image.RemoteURL),
									h.Width(320),
								)
							} else {
								cover = h.P(
									h.Class("italic text-sm"),
									h.Text("No cover available"),
								)
							}

							return h.Div(
								h.Class("flex flex-col gap-1"),
								cover,
								h.IfElseE(
									book.Book.Monitored,
									h.P(h.Text("This book is already available!"), h.Class("font-bold text-center")),
									GrabBookForm(book.Book, edition),
								),
							)
						}),
					),
					h.Hr(h.Class("mt-2 mb-2")),
				)
			}),
		),
	)
}
