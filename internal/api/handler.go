package api

import (
	"golift.io/starr"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/dotenv"
)

const (
	contentPath           string = "/data/media/books/komga"
	editionsMonitored     bool   = true
	manualAdd             bool   = true
	monitorBook           bool   = true
	monitorAuthor         bool   = false
	searchForNewBook      bool   = true
	searchForMissingBooks bool   = false
	qualityProfileID      int64  = 1
	metadataProfileID     int64  = 1
	authorMonitorType     string = "existing"
)

// Connect returns a new instance of the starr.readarr API thingamajig
func Connect() *readarr.Readarr {
	dotEnvVars, err := dotenv.GetEnv()
	if err != nil {
		return nil
	}

	starrConfig := starr.New(dotEnvVars.ApiKey, dotEnvVars.ApiURL.String(), 0)
	ReadarrAPI := readarr.New(starrConfig)

	return ReadarrAPI
}

func GetStatus() (*readarr.SystemStatus, error) {
	handler := Connect()
	status, err := handler.GetSystemStatus()
	return status, err
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
