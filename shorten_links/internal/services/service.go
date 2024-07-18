package services

import (
	"math/rand"
	my_errors "shorten_links/internal/errors"
	database "shorten_links/internal/storages/redis"
	"sync"
	"time"
)

type SetHashLink struct {
	Base  *database.Redis
	mutex sync.Mutex
}

type DataLink struct {
	ShortLink    string    `json:"short_link"`
	LongLink     string    `json:"long_link"`
	StatRedirect int       `json:"stat_redirect"`
	Death        time.Time `json:"death"`
}

func NewSetHashLink() SetHashLink {
	m := SetHashLink{
		Base: database.NewRedis(),
	}
	return m
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
	l := database.DataBase{
		LongLink:     llink,
		StatRedirect: 0,
		Death:        time.Now().Add(timelife),
	}
	s.mutex.Lock()
	sslink := HashingLink(l.LongLink)
	err := s.Base.Set(sslink, l)
	s.mutex.Unlock()
	return sslink, err
}

func (s *SetHashLink) SetRedirect(slink string, llink database.DataBase) error {
	s.mutex.Lock()

	new := database.DataBase{
		LongLink:     llink.LongLink,
		StatRedirect: llink.StatRedirect + 1,
		Death:        llink.Death,
	}
	s.Base.Set(slink, new)

	s.mutex.Unlock()
	return my_errors.ErrNoLlink
}

func (s *SetHashLink) GetAllStat() ([]DataLink, error) {
	var data []DataLink
	var d DataLink
	s.mutex.Lock()
	mapp, err := s.Base.GetAll()
	if err != nil {
		return data, err
	}
	for slink := range mapp {
		d.ShortLink = slink
		d.Death = mapp[slink].Death
		d.LongLink = mapp[slink].LongLink
		d.StatRedirect = mapp[slink].StatRedirect
		data = append(data, d)
	}
	s.mutex.Unlock()
	return data, err
}
