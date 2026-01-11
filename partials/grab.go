package partials

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/api"
	"tahanraamatut/internal/components"
)

func GrabBookForm(book *readarr.Book, edition *readarr.Edition) *h.Element {
	return h.Form(
		h.PostPartial(Grab), // HTMX will call Grab() and swap the returned partial
		h.Class("flex justify-center w-full"),
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
		components.SpinnerButton("grab", "grabbing..."),
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
		errStr := fmt.Sprintf("Failed to start grab for %s by %s because: %s", Title, AuthorName, err)
		return h.NewPartial(components.GrabError(errStr))
	}

	_, err = handler.StartSearch(ctx, grab.ID)
	if err != nil {
		errStr := fmt.Sprintf("Failed to start search for id %d by %s: %s", grab.ID, AuthorName, err)
		return h.NewPartial(components.GrabError(errStr))
	}

	grabbed, err := waitForGrabbed(ctx, handler, grab, 30*time.Second, 2*time.Second)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		errStr := fmt.Sprintf("Failed to check whether book with id %d by %s got grabbed because: %s", grab.ID, AuthorName, err)
		return h.NewPartial(components.GrabError(errStr))
	}

	if !grabbed {
		err = handler.CleanFailedAdd(ctx, grab)
		if err != nil {
			return h.NewPartial(
				h.Div(
					h.IfElseE(
						err != nil,
						h.P(h.TextF("Failed to delete unfound book entry with id %d by %s because: %s", grab.ID, AuthorName, err)),
						h.P(h.TextF("Couldn't find %s or couldn't parse a name. Try Readarr's interactive search instead.", grab.Title)),
					),
				),
			)
		}
	}

	api.MarkGrabbed(grab.ID)

	// success! now to track the status
	return h.NewPartial(
		h.Div(
			h.Id("grab-status"),
			h.Class("my-2"),
			GrabStatusWidget(rctx, grab.ID),
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
