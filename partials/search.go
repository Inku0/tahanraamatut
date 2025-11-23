package partials

import (
	"io"
	"strconv"

	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/hx"
	"github.com/maddalax/htmgo/framework/js"

	"tahanraamatut/internal/api"
)

func Search() *h.Partial {
	return h.NewPartial(
		h.Form(
			h.TriggerChildren(),
			h.PostPartial(SubmitForm),
			h.Class("flex dark:text-white flex-col gap-2 mx-auto"),
			h.LabelFor("name", "raamatu nimi:"),
			h.Input(
				"text",
				h.Required(),
				h.Type("search"),
				h.Role("search"),
				h.Class("p-4 dark:text-black rounded-md border border-slate-200"),
				h.Name("name"),
				h.Attribute("autocomplete", "off"),
				h.Placeholder("nimi"),
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

	readarr := api.Connect()

	health, err := readarr.HealthCheck()
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.Text("Encountered a fatal error while doing a health check. Wrong API key? error:"),
					h.P(h.TextF("%s", err)),
				),
			),
		)
	} else if health.StatusCode != 200 {
		body, err := io.ReadAll(health.Body)
		if err != nil {
			return h.NewPartial(
				h.Div(
					h.Class("text-base/8"),
					h.P(h.TextF("Encountered a fatal error while reading the body of the health check:"),
						h.P(h.TextF("error: %s", err)),
					),
				))
		}
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.Text("API health check returned a non-successful status code:"),
					h.P(h.TextF("%s", body)),
				),
			))
	}

	searchResults, err := readarr.Search(name)
	if err != nil || searchResults.StatusCode != 200 {
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.TextF("Encountered a fatal error while searching for a book named: %s", name)),
				h.P(h.TextF("error: %s", err)),
			),
		)
	}

	body, err := io.ReadAll(searchResults.Body)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.TextF("Encountered a fatal error while reading the search results for a book named: %s", name)),
				h.P(h.TextF("error: %s", err)),
			),
		)
	}

	parsedResults, err := api.ParseSearch(body)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.Class("text-base/8"),
				h.P(h.TextF("Encountered a fatal error while parsing the search results for a book named: %s", name)),
				h.P(h.TextF("error: %s", err)),
			),
		)
	}

	return h.NewPartial(
		h.Div(
			h.Class("flex flex-col mx-12"),
			h.P(h.TextF("Searched for: %s", name)),
			h.Hr(),
			h.Div(
				h.Class("mt-[1rem] flex flex-col gap-y-4"),
				h.List(parsedResults, func(searchResult api.SearchResult, index int) *h.Element {
					if searchResult.Book == nil || searchResult.Book.Editions == nil || searchResult.ID == 0 {
						return h.Empty()
					}

					return h.Div(
						h.P(
							h.TextF("\"%s\" by %s", searchResult.Book.Title, searchResult.Book.Author.AuthorName),
						),
						h.P(
							h.TextF("page count: %d", searchResult.Book.PageCount),
						),
						h.P(
							h.TextF("release date: %s", searchResult.Book.ReleaseDate),
						),
						h.P(
							h.TextF("editions: %s", searchResult.Book.Editions[0].Title),
						),
						h.P(
							h.TextF("already monitored? %s", strconv.FormatBool(searchResult.Book.Monitored)),
						),
						h.P(
							h.TextF("grabbed (downloading)?: %s", strconv.FormatBool(searchResult.Book.Grabbed)),
						),
						h.Div(
							h.Class("flex"),
							h.P(
								h.Class("p-4 text-justify flex-[10]"),
								h.TextF("%s", searchResult.Book.Overview),
							),
							h.If(searchResult.Book.RemoteCover != "",
								h.Div(
									h.Class("p-8 flex-[2]"),
									h.Img(
										h.Src(searchResult.Book.RemoteCover),
										h.Width(200),
									),
								),
							),
						),
						h.Hr(),
					)
				}),
			),
		),
	)
}
