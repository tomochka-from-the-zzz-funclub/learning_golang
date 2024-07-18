package transport

import (
	"shorten_links/internal/services"
	database "shorten_links/internal/storages/redis"
	"time"
)

type Set interface {
	CreateShortLink(llink string, timelife time.Duration) (string, error)
	GetAllStat() ([]services.InfoLink, error)
	Set(shortlink string, data database.InfoLLink) error
	GetLongL(shortlink string) (string, error)
	GetAllData(shortlink string) (database.InfoLLink, error)
	GetRedirect(shortlink string) (int, error)
	Increment(key string) error
}
