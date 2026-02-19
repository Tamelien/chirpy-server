package api

import (
	"sync/atomic"

	"github.com/tamelien/chirpy-server/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
	PLATFORM       string
}
