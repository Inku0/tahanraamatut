package api

import (
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"golift.io/starr/readarr"
)

const (
	contentPath string = "/data/media/books/komga"
)

type BookStatus struct {
	progress int
	state    string
}

func StartSearch(bookID int64) (*readarr.CommandResponse, error) {
	handler := Connect()
	bookIDs := make([]int64, 0)
	bookIDs = append(bookIDs, bookID)

	command := readarr.CommandRequest{
		Name:    "BookSearch",
		BookIDs: bookIDs,
	}

	resp, err := handler.SendCommand(&command)
	return resp, err
}

func Spinner(children ...h.Ren) *h.Element {
	return h.Div(
		h.Children(children...),
		h.Class("absolute left-1 spinner spinner-border animate-spin "+
			"inline-block w-6 h-6 border-4 rounded-full dark:border-slate-200 border-black-200 border-t-transparent"),
		h.Attribute("role", "status"),
	)
}

func GrabBookForm(book *readarr.Book, edition *readarr.Edition) *h.Element {
	buttonClasses := "rounded dark:border-2 dark:border-white-200 items-center px-3 py-2 bg-white-800 dark:text-white w-full text-center"

	return h.Form(
		h.PostPartial(Grab), // HTMX will call Book() and swap the returned partial
		h.Class("flex gap-2 justify-center"),
		// hidden id so Book() can read it via ctx.FormValue("id")
		h.Div(
			h.Input(
				"hidden",
				h.Name("ForeignBookID"),
				h.Value(book.ForeignBookID),
			),
			h.Input(
				"hidden",
				h.Name("AuthorName"),
				h.Value(book.Author.AuthorName),
			),
			h.Input(
				"hidden",
				h.Name("ForeignAuthorID"),
				h.Value(book.Author.ForeignAuthorID),
			),
			h.Input(
				"hidden",
				h.Name("Title"),
				h.Value(book.Title),
			),
			h.Input(
				"hidden",
				h.Name("TitleSlug"),
				h.Value(book.TitleSlug),
			),
			h.Input(
				"hidden",
				h.Name("ForeignEditionID"),
				h.Value(edition.ForeignEditionID), // gets the first edition
			),
		),
		// keep the before/after hooks for spinner / submit toggling
		h.HxBeforeRequest(
			js.RemoveClassOnChildren(".loading", "hidden"),
			js.SetClassOnChildren(".submit", "hidden"),
		),
		h.HxAfterRequest(
			js.SetClassOnChildren(".loading", "hidden"),
			js.RemoveClassOnChildren(".submit", "hidden"),
		),
		h.Button(
			h.Class("loading hidden relative text-center", buttonClasses),
			Spinner(),
			h.Disabled(),
			h.Text("grabbing..."),
		),
		h.Button(
			h.Type("submit"),
			h.Class("submit", buttonClasses),
			h.Text("grab"),
		),
	)
}

// Grab is the HTMX handler for the grab action. It reads the `id` form value,
// performs the grab (replace the placeholder logic with your real API call),
// and returns a partial fragment that HTMX will inject where the trigger element is targeted.
func Grab(ctx *h.RequestContext) *h.Partial {
	ForeignBookID := ctx.FormValue("ForeignBookID")
	AuthorName := ctx.FormValue("AuthorName")
	ForeignAuthorID := ctx.FormValue("ForeignAuthorID")
	Title := ctx.FormValue("Title")
	TitleSlug := ctx.FormValue("TitleSlug")
	ForeignEditionID := ctx.FormValue("ForeignEditionID")

	handler := Connect()

	booksToMonitor := make([]string, 0)
	booksToMonitor = append(booksToMonitor, ForeignBookID)

	editions := []*readarr.AddBookEdition{
		{
			Title:            Title,
			TitleSlug:        TitleSlug,
			Images:           nil,
			ForeignEditionID: ForeignEditionID,
			Monitored:        true,
			ManualAdd:        false,
		},
	}

	bookToAdd := readarr.AddBookInput{
		Monitored: true,
		Tags:      make([]int, 0),
		AddOptions: &readarr.AddBookOptions{
			AddType:          "automatic",
			SearchForNewBook: true, // change this to download
		},
		Author: &readarr.AddBookAuthor{
			Monitored:         true,
			QualityProfileID:  1,
			MetadataProfileID: 1,
			ForeignAuthorID:   ForeignAuthorID,
			RootFolderPath:    contentPath,
			Tags:              make([]int, 0),
			AddOptions: &readarr.AddAuthorOptions{
				SearchForMissingBooks: false, // change this to download
				Monitored:             true,
				Monitor:               "none",
				BooksToMonitor:        booksToMonitor,
			},
		},
		Editions:      editions,
		ForeignBookID: ForeignBookID,
	}

	grab, err := handler.AddBook(&bookToAdd)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.P(h.TextF("failed to start grab for id %s by %s: %s", ForeignBookID, AuthorName, err)),
			),
		)
	}

	return h.NewPartial(
		h.Div(
			h.P(h.TextF("Started grab for book id %s by %s", ForeignBookID, AuthorName)),
			h.P(h.TextF("response: %s", grab)),
		),
	)
}
