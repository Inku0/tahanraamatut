package api

import (
	"log"
	"net/http"
	"net/url"

	"tahanraamatut/internal/dotenv"
)

type HTTPMethod string

type ReadarrAPI struct {
	dotenv dotenv.Dotenv
	client *http.Client
}

func Connect() ReadarrAPI {
	env := dotenv.GetEnv()
	return ReadarrAPI{
		dotenv: env,
		client: &http.Client{},
	}
}

func (api ReadarrAPI) makeRequest(pathURL *url.URL, method string) (*http.Response, error) {
	if method != "GET" && method != "POST" {
		log.Fatalf("%s is not a supported HTTP method", method)
	}

	req, err := http.NewRequest("GET", pathURL.String(), nil)
	if err != nil {
		log.Fatalf("encountered a fatal error: %s", err)
	}

	req.Header.Add("X-Api-Key", api.dotenv.ApiKey)
	req.Header.Add("Accept", "*/*")

	resp, err := api.client.Do(req)

	return resp, err
}

func (api ReadarrAPI) HealthCheck() (*http.Response, error) {
	apiURL := api.dotenv.ApiURL

	healthURL := apiURL.JoinPath("api", "v1", "health")

	resp, err := api.makeRequest(healthURL, "GET")

	return resp, err
}

func (api ReadarrAPI) Search(term string) (*http.Response, error) {
	apiURL := api.dotenv.ApiURL

	searchURL := apiURL.JoinPath("api", "v1", "search")

	queryValues := searchURL.Query()
	queryValues.Add("term", term)

	searchURL.RawQuery = queryValues.Encode()

	resp, err := api.makeRequest(searchURL, "get")

	return resp, err
}

// {
//  "title": "Monkey: The Journey to the West",
//  "authorTitle": "cheng'en, wu Monkey: The Journey to the West",
//  "seriesTitle": "",
//  "disambiguation": "",
//  "overview": "Probably the most popular book in the history of the Far East, this classic sixteenth century novel is a combination of picaresque novel and folk epic that mixes satire, allegory, and history into a rollicking adventure. It is the story of the roguish Monkey and his encounters with major and minor spirits, gods, demigods, demons, ogres, monsters, and fairies. This translation, by the distinguished scholar Arthur Waley, is the first accurate English version; it makes available to the Western reader a faithful reproduction of the spirit and meaning of the original.",
//  "authorId": 0,
//  "foreignBookId": "96649",
//  "foreignEditionId": "100237",
//  "titleSlug": "96649",
//  "monitored": true,
//  "anyEditionOk": true,
//  "ratings": {
//    "votes": 5684,
//    "value": 4.06,
//    "popularity": 23077.039999999997
//  },
//  "releaseDate": "1592-01-01T07:52:58Z",
//  "pageCount": 306,
//  "genres": [
//    "Classics",
//    "Fiction",
//    "Fantasy",
//    "China",
//    "Mythology",
//    "Chinese Literature",
//    "Asia",
//    "Literature",
//    "Novels",
//    "Adventure"
//  ],
//  "author": {
//    "authorMetadataId": 0,
//    "status": "continuing",
//    "ended": false,
//    "authorName": "Wu Cheng'en",
//    "authorNameLastFirst": "Cheng'en, Wu",
//    "foreignAuthorId": "92018",
//    "titleSlug": "92018",
//    "overview": "Librarian Note: There is more than one author by this name in the Goodreads database.Wu Cheng'en (simplified Chinese: 吴承恩; traditional Chinese: 吳承恩; pinyin: Wú Chéng'ēn, ca. 1505–1580 or 1500–1582, courtesy name Ruzhong (汝忠), pen name \"Sheyang Hermit,\" was a Chinese novelist and poet of the Ming Dynasty, best known for being the probable author of one of the Four Great Classical Novels of Chinese literature, Journey to the West, also called Monkey.",
//    "links": [
//      {
//        "url": "https://www.goodreads.com/author/show/92018.Wu_Cheng_en",
//        "name": "Goodreads"
//      }
//    ],
//    "images": [
//      {
//        "url": "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/authors/1282356359i/92018._UY200_CR33,0,200,200_.jpg",
//        "coverType": "poster",
//        "extension": ".jpg"
//      }
//    ],
//    "qualityProfileId": 1,
//    "metadataProfileId": 1,
//    "monitored": true,
//    "monitorNewItems": "none",
//    "folder": "Wu Cheng'en",
//    "genres": [],
//    "cleanName": "wuchengen",
//    "sortName": "wu cheng'en",
//    "sortNameLastFirst": "cheng'en, wu",
//    "tags": [],
//    "added": "0001-01-01T00:00:00Z",
//    "ratings": {
//      "votes": 0,
//      "value": 0,
//      "popularity": 0
//    },
//    "statistics": {
//      "bookFileCount": 0,
//      "bookCount": 0,
//      "availableBookCount": 0,
//      "totalBookCount": 0,
//      "sizeOnDisk": 0,
//      "percentOfBooks": 0
//    },
//    "addOptions": {
//      "searchForMissingBooks": false,
//      "booksToMonitor": [
//        "96649"
//      ]
//    },
//    "rootFolderPath": "/data/media/books/komga"
//  },
//  "images": [
//    {
//      "url": "/MediaCoverProxy/e8ba25a7073f54b50efe5012dc51f2dac6362d601198e43b4bc69f095c3a89bf/100237.jpg",
//      "coverType": "cover",
//      "extension": ".jpg",
//      "remoteUrl": "https://m.media-amazon.com/images/S/compressed.photo.goodreads.com/books/1347431752i/100237.jpg"
//    }
//  ],
//  "links": [
//    {
//      "url": "https://www.goodreads.com/work/96649-x-y-u-j",
//      "name": "Goodreads Editions"
//    },
//    {
//      "url": "https://www.goodreads.com/book/show/100237.Monkey",
//      "name": "Goodreads Book"
//    }
//  ],
//  "added": "0001-01-01T00:00:00Z",
//  "remoteCover": "https://m.media-amazon.com/images/S/compressed.photo.goodreads.com/books/1347431752i/100237.jpg",
//  "editions": [
//    {
//      "bookId": 0,
//      "foreignEditionId": "100237",
//      "titleSlug": "100237",
//      "isbn13": "9780802130860",
//      "asin": "0802130860",
//      "title": "Monkey: The Journey to the West",
//      "language": "eng",
//      "overview": "Probably the most popular book in the history of the Far East, this classic sixteenth century novel is a combination of picaresque novel and folk epic that mixes satire, allegory, and history into a rollicking adventure. It is the story of the roguish Monkey and his encounters with major and minor spirits, gods, demigods, demons, ogres, monsters, and fairies. This translation, by the distinguished scholar Arthur Waley, is the first accurate English version; it makes available to the Western reader a faithful reproduction of the spirit and meaning of the original.",
//      "format": "Paperback",
//      "isEbook": false,
//      "disambiguation": "",
//      "publisher": "Grove Press",
//      "pageCount": 306,
//      "releaseDate": "1994-01-12T08:00:00Z",
//      "images": [
//        {
//          "url": "/MediaCoverProxy/e8ba25a7073f54b50efe5012dc51f2dac6362d601198e43b4bc69f095c3a89bf/100237.jpg",
//          "coverType": "cover",
//          "extension": ".jpg",
//          "remoteUrl": "https://m.media-amazon.com/images/S/compressed.photo.goodreads.com/books/1347431752i/100237.jpg"
//        }
//      ],
//      "links": [
//        {
//          "url": "https://www.goodreads.com/book/show/100237.Monkey",
//          "name": "Goodreads Book"
//        }
//      ],
//      "ratings": {
//        "votes": 5684,
//        "value": 4.06,
//        "popularity": 23077.039999999997
//      },
//      "monitored": true,
//      "manualAdd": false,
//      "grabbed": false
//    }
//  ],
//  "grabbed": false,
//  "addOptions": {
//    "searchForNewBook": false
//  }
//}
//{
//  "title": "Monkey",
//  "authorTitle": "cheng'en, wu Monkey",
//  "seriesTitle": "",
//  "disambiguation": "",
//  "authorId": 299,
//  "foreignBookId": "96649",
//  "foreignEditionId": "100237",
//  "titleSlug": "96649",
//  "monitored": true,
//  "anyEditionOk": true,
//  "ratings": {
//    "votes": 5684,
//    "value": 4.06,
//    "popularity": 23077.039999999997
//  },
//  "releaseDate": "1592-01-01T06:13:58Z",
//  "pageCount": 306,
//  "genres": [
//    "Classics",
//    "Fiction",
//    "Fantasy",
//    "China",
//    "Mythology",
//    "Chinese Literature",
//    "Asia",
//    "Literature",
//    "Novels",
//    "Adventure"
//  ],
//  "author": {
//    "authorMetadataId": 299,
//    "status": "continuing",
//    "ended": false,
//    "authorName": "Wu Cheng'en",
//    "authorNameLastFirst": "Cheng'en, Wu",
//    "foreignAuthorId": "92018",
//    "titleSlug": "92018",
//    "overview": "Librarian Note: There is more than one author by this name in the Goodreads database.Wu Cheng'en (simplified Chinese: 吴承恩; traditional Chinese: 吳承恩; pinyin: Wú Chéng'ēn, ca. 1505–1580 or 1500–1582, courtesy name Ruzhong (汝忠), pen name \"Sheyang Hermit,\" was a Chinese novelist and poet of the Ming Dynasty, best known for being the probable author of one of the Four Great Classical Novels of Chinese literature, Journey to the West, also called Monkey.",
//    "links": [
//      {
//        "url": "https://www.goodreads.com/author/show/92018.Wu_Cheng_en",
//        "name": "Goodreads"
//      }
//    ],
//    "images": [
//      {
//        "url": "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/authors/1282356359i/92018._UY200_CR33,0,200,200_.jpg",
//        "coverType": "poster",
//        "extension": ".jpg"
//      }
//    ],
//    "path": "/data/media/books/komga/Wu Cheng'en",
//    "qualityProfileId": 1,
//    "metadataProfileId": 1,
//    "monitored": true,
//    "monitorNewItems": "none",
//    "genres": [],
//    "cleanName": "wuchengen",
//    "sortName": "wu cheng'en",
//    "sortNameLastFirst": "cheng'en, wu",
//    "tags": [],
//    "added": "2025-11-23T13:10:04Z",
//    "addOptions": {
//      "searchForMissingBooks": false,
//      "monitor": "all",
//      "booksToMonitor": [
//        "96649"
//      ],
//      "monitored": false
//    },
//    "ratings": {
//      "votes": 17228,
//      "value": 4.145983,
//      "popularity": 71426.99512400001
//    },
//    "statistics": {
//      "bookFileCount": 0,
//      "bookCount": 0,
//      "availableBookCount": 0,
//      "totalBookCount": 0,
//      "sizeOnDisk": 0,
//      "percentOfBooks": 0
//    },
//    "id": 299
//  },
//  "images": [
//    {
//      "url": "/MediaCover/Books/15964/cover.jpg",
//      "coverType": "cover",
//      "extension": ".jpg",
//      "remoteUrl": "https://m.media-amazon.com/images/S/compressed.photo.goodreads.com/books/1347431752i/100237.jpg"
//    }
//  ],
//  "links": [
//    {
//      "url": "https://www.goodreads.com/work/96649-x-y-u-j",
//      "name": "Goodreads Editions"
//    },
//    {
//      "url": "https://www.goodreads.com/book/show/100237.Monkey",
//      "name": "Goodreads Book"
//    }
//  ],
//  "statistics": {
//    "bookFileCount": 0,
//    "bookCount": 1,
//    "totalBookCount": 1,
//    "sizeOnDisk": 0,
//    "percentOfBooks": 0
//  },
//  "added": "2025-11-23T13:10:04Z",
//  "grabbed": false,
//  "id": 15964
//}
