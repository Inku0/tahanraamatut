package partials

import (
	"log"

	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/api"
)

const (
	contentPath string = "/data/media/books/komga"
)

func GrabBookForm(book *readarr.Book, edition *readarr.Edition) *h.Element {
	buttonClasses := "rounded dark:border-2 dark:border-white-200 items-center px-3 py-2 bg-white-800 dark:text-white w-full text-center"

	return h.Form(
		h.PostPartial(Grab), // HTMX will call Book() and swap the returned partial
		h.Class("flex justify-center w-[200px]"),
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
				h.Value(edition.ForeignEditionID),
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

	handler := api.Connect()

	booksToMonitor := make([]string, 0)
	booksToMonitor = append(booksToMonitor, ForeignBookID)

	log.Println(booksToMonitor, ForeignBookID, ForeignEditionID, ForeignAuthorID)

	editions := []*readarr.AddBookEdition{
		{
			Title:            Title,
			TitleSlug:        TitleSlug,
			ForeignEditionID: ForeignEditionID,
			Monitored:        true,
			ManualAdd:        true,
		},
	}

	bookToAdd := readarr.AddBookInput{
		Monitored: true,
		Tags:      make([]int, 0),
		AddOptions: &readarr.AddBookOptions{
			SearchForNewBook: true, // change this to download
		},
		Author: &readarr.AddBookAuthor{
			Monitored:         true,
			QualityProfileID:  1,
			MetadataProfileID: 0,
			ForeignAuthorID:   ForeignAuthorID,
			RootFolderPath:    contentPath,
			Tags:              make([]int, 0),
			AddOptions: &readarr.AddAuthorOptions{
				SearchForMissingBooks: true, // change this to download
				Monitored:             true,
				Monitor:               "none",
				BooksToMonitor:        booksToMonitor,
			},
		},
		Editions:      editions,
		ForeignBookID: ForeignBookID,
	}

	grab, err := handler.AddBook(&bookToAdd)
	if err != nil || grab == nil {
		return h.NewPartial(
			h.Div(
				h.P(h.TextF("failed to start grab for id %s by %s: %s", ForeignBookID, AuthorName, err)),
			),
		)
	}

	resp, err := api.StartSearch(grab.ID)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.P(h.TextF("failed to start search for id %s by %s: %s", grab.ID, AuthorName, err)),
			),
		)
	}

	return h.NewPartial(
		h.Div(
			h.P(h.TextF("Started grab for book id %d by %s", grab.ID, AuthorName)),
			h.P(h.TextF("response: %s", resp)),
		),
	)
}
