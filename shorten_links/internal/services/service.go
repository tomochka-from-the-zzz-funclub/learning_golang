package services

import (
	"math/rand"
	//"shorten_links/internal/services"
	database "shorten_links/internal/storages/redis"
	"time"
)

type SetHashLink struct {
	Base LinkStorage
}

// type DataLink struct {
// 	ShortLink    string    `json:"short_link"`
// 	LongLink     string    `json:"long_link"`
// 	StatRedirect int       `json:"stat_redirect"`
// 	Death        time.Time `json:"death"`
// }

func NewSetHashLink() *SetHashLink {
	m := SetHashLink{
		Base: database.NewRedis(),
	}
	return &m
}

func HashingLink(LongLink string) string {
	var hlink string
	hash := make([]byte, 6)
	for i := range hash {
		hash[i] = LongLink[rand.Intn(len(LongLink))]
	}
	hlink = string(hash)
	return hlink
}

func (s *SetHashLink) CreateShortLink(llink string, timelife time.Duration) (string, error) {
	var l database.InfoLLink
	l.SetLongLink(llink)
	l.SetStatRedirect(0)
	l.SetDeath(time.Now().Add(timelife))

	sslink := HashingLink(l.GetLongLink())
	err := s.Base.Set(sslink, l)
	return sslink, err
}

func (s *SetHashLink) GetAllStat() ([]InfoLink, error) {
	var data []InfoLink
	var d InfoLink
	mapp, err := s.Base.GetAll()
	if err != nil {
		return data, err
	}
	for slink := range mapp {
		d.SetShortLink(slink)
		d.SetDeath(mapp[slink].GetDeath())
		d.SetLongLink(mapp[slink].GetLongLink())
		d.SetStatRedirect(mapp[slink].GetStatRedirect())
		data = append(data, d)
	}
	return data, err
}

func (s *SetHashLink) Set(shortlink string, data database.InfoLLink) error {
	return s.Base.Set(shortlink, data)
}

func (s *SetHashLink) GetLongL(shortlink string) (string, error) {
	return s.Base.GetLongL(shortlink)
}

func (s *SetHashLink) GetAllData(shortlink string) (database.InfoLLink, error) {
	return s.Base.GetAllData(shortlink)
}

func (s *SetHashLink) GetRedirect(shortlink string) (int, error) {
	return s.Base.GetRedirect(shortlink)
}

func (s *SetHashLink) Increment(key string) error {
	return s.Base.Increment(key)
}
