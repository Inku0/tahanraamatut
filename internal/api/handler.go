package api

import (
	"golift.io/starr"
	"golift.io/starr/readarr"

	"tahanraamatut/internal/dotenv"
)

// Connect returns a new instance of the starr.readarr API thingamajig
func Connect() *readarr.Readarr {
	dotEnvVars := dotenv.GetEnv()

	starrConfig := starr.New(dotEnvVars.ApiKey, dotEnvVars.ApiURL.String(), 0)
	ReadarrAPI := readarr.New(starrConfig)

	return ReadarrAPI
}
