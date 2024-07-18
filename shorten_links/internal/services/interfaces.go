package services

import (
	database "shorten_links/internal/storages/redis"
	"time"
)

type LinkStorage interface { // база в сервисе
	GetAll() (map[string]database.InfoLLink, error)
	GetAllData(shortlink string) (database.InfoLLink, error)
	GetLongL(shortlink string) (string, error)
	GetRedirect(shortlink string) (int, error)
	Increment(key string) error
	Set(shortlink string, data database.InfoLLink) error
}

type InfoLink interface {
	SetShortLink(string)
	SetLongLink(string)
	SetStatRedirect(int)
	SetDeath(time.Time)
	GetShortLink() string
	GetLongLink() string
	GetStatRedirect() int
	GetDeath() time.Time
}

type DataLink struct {
	ShortLink    string    `json:"short_link"`
	LongLink     string    `json:"long_link"`
	StatRedirect int       `json:"stat_redirect"`
	Death        time.Time `json:"death"`
}

func (d DataLink) SetShortLink(slink string) {
	d.ShortLink = slink
}
func (d DataLink) SetLongLink(llink string) {
	d.LongLink = llink
}
func (d DataLink) SetDeath(death time.Time) {
	d.Death = death
}
func (d DataLink) SetStatRedirect(redirect int) {
	d.StatRedirect = redirect
}

func (d DataLink) GetShortLink() string {
	return d.ShortLink
}
func (d DataLink) GetLongLink() string {
	return d.LongLink
}
func (d DataLink) GetDeath() time.Time {
	return d.Death
}
func (d DataLink) GetStatRedirect() int {
	return d.StatRedirect
}
