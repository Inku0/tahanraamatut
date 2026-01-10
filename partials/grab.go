package partials

import (
	"context"
	"time"

	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/js"
	"github.com/maddalax/htmgo/framework/service"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/api"
)

func GrabBookForm(book *readarr.Book, edition *readarr.Edition) *h.Element {
	buttonClasses := "rounded border-2 border-black dark:border-white-200 items-center px-3 py-2 bg-white-800 dark:text-white w-full text-center"

	return h.Form(
		h.PostPartial(Grab), // HTMX will call Grab() and swap the returned partial
		h.Class("flex justify-center w-xs"),
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
			h.Text("Grabbing..."),
		),
		h.Button(
			h.Type("submit"),
			h.Class("submit", buttonClasses),
			h.Text("Grab!"),
		),
	)
}

// Grab is the HTMX handler for the grab action. It reads the necessary fields from the request context,
func Grab(rctx *h.RequestContext) *h.Partial {
	ForeignBookID := rctx.FormValue("ForeignBookID")
	AuthorName := rctx.FormValue("AuthorName")
	ForeignAuthorID := rctx.FormValue("ForeignAuthorID")
	Title := rctx.FormValue("Title")
	TitleSlug := rctx.FormValue("TitleSlug")
	ForeignEditionID := rctx.FormValue("ForeignEditionID")

	handler := service.Get[api.ReadarrService](rctx.ServiceLocator())

	// upper bound for time
	ctx, cancel := context.WithTimeout(rctx.Request.Context(), 45*time.Second)
	defer cancel()

	bookToAdd := api.FormatBookToAdd(api.BookToAdd{
		ForeignBookID:    ForeignBookID,
		AuthorName:       AuthorName,
		ForeignAuthorID:  ForeignAuthorID,
		Title:            Title,
		TitleSlug:        TitleSlug,
		ForeignEditionID: ForeignEditionID,
	})

	grab, err := handler.Client.AddBookContext(ctx, bookToAdd)
	if err != nil || grab == nil {
		return h.NewPartial(
			h.Div(
				h.P(h.TextF("Failed to start grab for %s by %s because: %s", Title, AuthorName, err)),
			),
		)
	}

	_, err = handler.StartSearch(ctx, grab.ID)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.P(h.TextF("Failed to start search for id %d by %s: %s", grab.ID, AuthorName, err)),
			),
		)
	}

	grabbed, err := waitForGrabbed(ctx, handler, grab, 30*time.Second, 2*time.Second)
	if err != nil {
		return h.NewPartial(
			h.Div(
				h.P(h.TextF("Failed to check whether book with id %d by %s got grabbed because: %s", grab.ID, AuthorName, err)),
			),
		)
	}

	if !grabbed {
		err = handler.CleanFailedAdd(ctx, grab)
		if err != nil {
			return h.NewPartial(
				h.Div(
					h.P(h.TextF("Failed to delete book with id %d by %s because: %s", grab.ID, AuthorName, err)),
				),
			)
		}
	}

	return h.NewPartial(
		h.Div(
			h.IfElseE(
				grabbed,
				h.P(h.TextF("Found and downloading %s!", grab.Title)),
				h.P(h.TextF("Couldn't find %s...", grab.Title)),
			),
		),
	)
}

func waitForGrabbed(ctx context.Context, handler *api.ReadarrService, grab *readarr.Book, maxWait, interval time.Duration) (bool, error) {
	waitCtx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		grabbed, err := handler.GotGrabbed(waitCtx, grab)
		if err != nil {
			return false, err
		}
		if grabbed {
			return true, nil
		}

		select {
		case <-waitCtx.Done():
			return false, waitCtx.Err()
		case <-t.C:
		}
	}
}
