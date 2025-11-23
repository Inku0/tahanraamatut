package api

import (
	"encoding/json"
)

// SearchResult models each top-level element in the API search results array
type SearchResult struct {
	ForeignID string  `json:"foreignId,omitempty"`
	Author    *Author `json:"author,omitempty"`
	Book      *Book   `json:"book,omitempty"`
	ID        int     `json:"id,omitempty"`
}

// Author models the author metadata blocks
type Author struct {
	AuthorMetadataId    int         `json:"authorMetadataId,omitempty"`
	Status              string      `json:"status,omitempty"`
	Ended               bool        `json:"ended,omitempty"`
	AuthorName          string      `json:"authorName,omitempty"`
	AuthorNameLastFirst string      `json:"authorNameLastFirst,omitempty"`
	ForeignAuthorId     string      `json:"foreignAuthorId,omitempty"`
	TitleSlug           string      `json:"titleSlug,omitempty"`
	Overview            string      `json:"overview,omitempty"`
	Links               []Link      `json:"links,omitempty"`
	Images              []Image     `json:"images,omitempty"`
	RemotePoster        string      `json:"remotePoster,omitempty"`
	Path                string      `json:"path,omitempty"`
	QualityProfileId    int         `json:"qualityProfileId,omitempty"`
	MetadataProfileId   int         `json:"metadataProfileId,omitempty"`
	Monitored           bool        `json:"monitored,omitempty"`
	MonitorNewItems     string      `json:"monitorNewItems,omitempty"`
	Folder              string      `json:"folder,omitempty"`
	Genres              []string    `json:"genres,omitempty"`
	CleanName           string      `json:"cleanName,omitempty"`
	SortName            string      `json:"sortName,omitempty"`
	SortNameLastFirst   string      `json:"sortNameLastFirst,omitempty"`
	Tags                []string    `json:"tags,omitempty"`
	Added               string      `json:"added,omitempty"`
	AddOptions          interface{} `json:"addOptions,omitempty"`
	Ratings             *Ratings    `json:"ratings,omitempty"`
	Statistics          *Statistics `json:"statistics,omitempty"`
	ID                  int         `json:"id,omitempty"`
}

type Edition struct {
	BookID           int      `json:"bookId,omitempty"`
	ForeignEditionId string   `json:"foreignEditionId,omitempty"`
	TitleSlug        string   `json:"titleSlug,omitempty"`
	ISBN13           string   `json:"isbn13,omitempty"`
	ASIN             string   `json:"asin,omitempty"`
	Title            string   `json:"title,omitempty"`
	Language         string   `json:"language,omitempty"`
	Overview         string   `json:"overview,omitempty"`
	Format           string   `json:"format,omitempty"`
	IsEbook          bool     `json:"isEbook,omitempty"`
	Disambiguation   string   `json:"disambiguation,omitempty"`
	Publisher        string   `json:"publisher,omitempty"`
	PageCount        int      `json:"pageCount,omitempty"`
	ReleaseDate      string   `json:"releaseDate,omitempty"`
	Images           []Image  `json:"images,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Ratings          *Ratings `json:"ratings,omitempty"`
	Monitored        bool     `json:"monitored,omitempty"`
	ManualAdd        bool     `json:"manualAdd,omitempty"`
	Grabbed          bool     `json:"grabbed,omitempty"`
	ID               int      `json:"id,omitempty"`
}

type Book struct {
	Title            string    `json:"title,omitempty"`
	AuthorTitle      string    `json:"authorTitle,omitempty"`
	SeriesTitle      string    `json:"seriesTitle,omitempty"`
	Disambiguation   string    `json:"disambiguation,omitempty"`
	Overview         string    `json:"overview,omitempty"`
	AuthorId         int       `json:"authorId,omitempty"`
	ForeignBookId    string    `json:"foreignBookId,omitempty"`
	ForeignEditionId string    `json:"foreignEditionId,omitempty"`
	TitleSlug        string    `json:"titleSlug,omitempty"`
	Monitored        bool      `json:"monitored,omitempty"`
	AnyEditionOk     bool      `json:"anyEditionOk,omitempty"`
	Ratings          *Ratings  `json:"ratings,omitempty"`
	ReleaseDate      string    `json:"releaseDate,omitempty"`
	PageCount        int       `json:"pageCount,omitempty"`
	Genres           []string  `json:"genres,omitempty"`
	Author           *Author   `json:"author,omitempty"`
	Images           []Image   `json:"images,omitempty"`
	Links            []Link    `json:"links,omitempty"`
	Added            string    `json:"added,omitempty"`
	RemoteCover      string    `json:"remoteCover,omitempty"`
	Editions         []Edition `json:"editions,omitempty"`
	Grabbed          bool      `json:"grabbed,omitempty"`
	ID               int       `json:"id,omitempty"`
}

// Link and Image helper types
type Link struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type Image struct {
	URL       string `json:"url,omitempty"`
	CoverType string `json:"coverType,omitempty"`
	Extension string `json:"extension,omitempty"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// Ratings and Statistics
type Ratings struct {
	Votes      int     `json:"votes,omitempty"`
	Value      float64 `json:"value,omitempty"`
	Popularity float64 `json:"popularity,omitempty"`
}

type Statistics struct {
	BookFileCount      int `json:"bookFileCount,omitempty"`
	BookCount          int `json:"bookCount,omitempty"`
	AvailableBookCount int `json:"availableBookCount,omitempty"`
	TotalBookCount     int `json:"totalBookCount,omitempty"`
	SizeOnDisk         int `json:"sizeOnDisk,omitempty"`
	PercentOfBooks     int `json:"percentOfBooks,omitempty"`
}

// ParseSearch unmarshals the provided JSON bytes into SearchResult slice
func ParseSearch(data []byte) ([]SearchResult, error) {
	var results []SearchResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
