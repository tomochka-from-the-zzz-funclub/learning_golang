package services

import (
	"fmt"
	"math/rand"
	my_errors "shorten_links/internal/errors"
	"sync"
	"time"
)

type HashLink struct {
	ShortLink string
}
type WorkLink struct {
	LongLink     string
	StatRedirect int
	TimeLife     time.Duration
	Create       time.Time
}
type SetHashLink struct {
	SetLink map[HashLink]WorkLink
	mutex   sync.Mutex
}

type DataLink struct {
	ShortLink    string    `json:"short_link"`
	LongLink     string    `json:"long_link"`
	StatRedirect int       `json:"stat_redirect"`
	Death        time.Time `json:"death"`
}

func NewSetHashLink() SetHashLink {
	m := SetHashLink{SetLink: make(map[HashLink]WorkLink)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.mutex.Lock()
			for k, v := range m.SetLink {
				if now.Sub(v.Create) > v.TimeLife {
					fmt.Println("Удалил")
					delete(m.SetLink, k)
				}
			}
			m.mutex.Unlock()
		}
	}()

	return m
}

func HashingLink(LongLink string) HashLink {
	var hlink HashLink
	hash := make([]byte, 6)
	for i := range hash {
		hash[i] = LongLink[rand.Intn(len(LongLink))]
	}
	hlink.ShortLink = string(hash)
	return hlink
}

func (s *SetHashLink) CreateShortLink(llink string, timelife time.Duration) {
	l := WorkLink{
		LongLink:     llink,
		StatRedirect: 0,
		Create:       time.Now(),
		TimeLife:     timelife,
	}
	s.mutex.Lock()
	s.SetLink[HashingLink(l.LongLink)] = l
	s.mutex.Unlock()
}

func (s *SetHashLink) GetStatLink(llink string) (int, error) {
	s.mutex.Lock()
	for slink := range s.SetLink {
		if s.SetLink[slink].LongLink == llink {
			s.mutex.Unlock()
			return s.SetLink[slink].StatRedirect, nil
		}
	}
	s.mutex.Unlock()
	return 0, my_errors.ErrNoLlink
}

func (s *SetHashLink) GetLongLink(slink HashLink) (WorkLink, error) {
	s.mutex.Lock()
	if llink, ok := s.SetLink[slink]; !ok {
		s.mutex.Unlock()
		return llink, my_errors.ErrNoSlink
	} else {
		s.mutex.Unlock()
		return llink, nil
	}
}

func (s *SetHashLink) GetShortLink(llink string) (HashLink, error) {
	var slink HashLink
	l := WorkLink{
		LongLink:     llink,
		StatRedirect: 0,
	}
	s.mutex.Lock()
	for slink := range s.SetLink {
		if s.SetLink[slink].LongLink == l.LongLink {
			s.mutex.Unlock()
			return slink, nil
		}
	}
	s.mutex.Unlock()
	return slink, my_errors.ErrNoLlink
}

func (s *SetHashLink) SetRedirect(llink string) error {
	s.mutex.Lock()
	for slink := range s.SetLink {
		if s.SetLink[slink].LongLink == llink {
			last := s.SetLink[slink].StatRedirect
			s.SetLink[slink] = WorkLink{
				LongLink:     llink,
				StatRedirect: last + 1,
				TimeLife:     s.SetLink[slink].TimeLife,
				Create:       s.SetLink[slink].Create,
			}
			s.mutex.Unlock()
			return nil
		}
	}
	s.mutex.Unlock()
	return my_errors.ErrNoLlink
}

func (s *SetHashLink) GetAllStat() []DataLink {
	var data []DataLink
	var d DataLink
	s.mutex.Lock()
	//var t time.Time
	//t = s.SetLink[slink].Create + int64(s.SetLink[slink].TimeLife.Seconds())
	for slink := range s.SetLink {
		d.LongLink = s.SetLink[slink].LongLink
		d.ShortLink = slink.ShortLink
		d.StatRedirect = s.SetLink[slink].StatRedirect
		d.Death = s.SetLink[slink].Create.Add(s.SetLink[slink].TimeLife)
		data = append(data, d)
	}
	s.mutex.Unlock()
	return data
}
