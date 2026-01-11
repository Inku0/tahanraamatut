package api

import (
	"context"
	"log"

	"golift.io/starr"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/dotenv"
)

const (
	contentPath           string = "/data/media/books/komga"
	editionsMonitored     bool   = true
	manualAdd             bool   = true
	monitorBook           bool   = true
	monitorAuthor         bool   = true
	searchForNewBook      bool   = true
	searchForMissingBooks bool   = false
	qualityProfileID      int64  = 1
	metadataProfileID     int64  = 1
	authorMonitorType     string = "existing"
)

type ReadarrService struct {
	Client *readarr.Readarr
}

// NewReadarrService creates a connection to the Readarr API
func NewReadarrService() *ReadarrService {
	dotEnvVars, err := dotenv.GetEnv()
	if err != nil {
		log.Fatal("Error reading dot env variables")
		return nil
	}

	starrConfig := starr.New(dotEnvVars.ApiKey, dotEnvVars.ApiURL.String(), 0)
	readarrClient := readarr.New(starrConfig)

	return &ReadarrService{Client: readarrClient}
}

func (service *ReadarrService) GetStatus(ctx context.Context) (*readarr.SystemStatus, error) {
	return service.Client.GetSystemStatusContext(ctx)
}

type BookToAdd struct {
	ForeignBookID,
	AuthorName,
	ForeignAuthorID,
	Title,
	TitleSlug,
	ForeignEditionID string
}

func FormatBookToAdd(add BookToAdd) *readarr.AddBookInput {
	booksToMonitor := []string{add.ForeignBookID}

	editions := []*readarr.AddBookEdition{
		{
			Title:            add.Title,
			TitleSlug:        add.TitleSlug,
			ForeignEditionID: add.ForeignEditionID,
			Monitored:        editionsMonitored,
			ManualAdd:        manualAdd,
		},
	}

	bookToAdd := readarr.AddBookInput{
		Monitored: monitorBook,
		Tags:      []int{},
		AddOptions: &readarr.AddBookOptions{
			SearchForNewBook: searchForNewBook,
		},
		Author: &readarr.AddBookAuthor{
			Monitored:         monitorAuthor,
			QualityProfileID:  qualityProfileID,
			MetadataProfileID: metadataProfileID,
			ForeignAuthorID:   add.ForeignAuthorID,
			RootFolderPath:    contentPath,
			Tags:              []int{},
			AddOptions: &readarr.AddAuthorOptions{
				SearchForMissingBooks: searchForMissingBooks,
				Monitored:             monitorBook,
				Monitor:               authorMonitorType,
				BooksToMonitor:        booksToMonitor,
			},
		},
		Editions:      editions,
		ForeignBookID: add.ForeignBookID,
	}

	return &bookToAdd
}

func (service *ReadarrService) StartSearch(ctx context.Context, bookID int64) (*readarr.CommandResponse, error) {
	bookIDs := []int64{bookID}

	command := readarr.CommandRequest{
		Name:    "BookSearch",
		BookIDs: bookIDs,
	}

	resp, err := service.Client.SendCommandContext(ctx, &command)
	return resp, err
}
