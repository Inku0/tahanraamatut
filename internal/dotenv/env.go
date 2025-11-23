package dotenv

import (
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Dotenv struct {
	ApiKey string
	ApiURL url.URL
}

func GetEnv() Dotenv {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	URL := os.Getenv("API_URL")

	parsedURL, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("error parsing API_URL: %s", err)
	}

	return Dotenv{
		ApiKey: apiKey,
		ApiURL: *parsedURL,
	}
}
