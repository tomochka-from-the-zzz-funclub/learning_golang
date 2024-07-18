package database

import "time"

type InfoLLink interface {
	SetLongLink(llink string)
	SetDeath(death time.Time)
	SetStatRedirect(redirect int)
	GetLongLink() string
	GetStatRedirect() int
	GetDeath() time.Time
}

type DataLLink struct {
	LongLink     string    `json:"long_link"`
	StatRedirect int       `json:"stat_redirect"`
	Death        time.Time `json:"death"`
}

func (d DataLLink) SetLongLink(llink string) {
	d.LongLink = llink
}
func (d DataLLink) SetDeath(death time.Time) {
	d.Death = death
}
func (d DataLLink) SetStatRedirect(redirect int) {
	d.StatRedirect = redirect
}
func (d DataLLink) GetLongLink() string {
	return d.LongLink
}
func (d DataLLink) GetDeath() time.Time {
	return d.Death
}
func (d DataLLink) GetStatRedirect() int {
	return d.StatRedirect
}
