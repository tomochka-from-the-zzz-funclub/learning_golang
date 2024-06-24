package transport

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"net/http"
	my_errors "shorten_links/internal/errors"
	"shorten_links/internal/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type HandlersBuilder struct {
	s  services.SetHashLink
	lg zerolog.Logger
}

func HandleCreate() {
	hb := HandlersBuilder{
		s:  services.NewSetHashLink(),
		lg: zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.UnixDate}),
	}
	http.HandleFunc("/shortlink/get", hb.GetShortLink())
	http.HandleFunc("/shortlink/redirect", hb.Redirect())
	http.HandleFunc("/shortlink/stat", hb.GetStat())
	http.HandleFunc("/shortlink/allstat", hb.GetAllStat())
}

func (hb *HandlersBuilder) GetShortLink() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			longlink, timelife, err := ParseJsonL(r)
			if err != nil {
				//err_ :=
				WriteJsonErr(&w, err.Error())

				w.WriteHeader(http.StatusBadRequest)
				// if err_ != nil {
				// 	hb.lg.Warn().
				// 		Msgf("message from func GetShortLink %v", err_.Error())
				// 	return
				// }
				hb.lg.Warn().
					Msgf("message from func GetShortLink %v", err.Error())
			} else {
				hb.s.CreateShortLink(longlink, timelife)
				slink, err := hb.s.GetShortLink(longlink)
				if err != nil {
					WriteJsonErr(&w, err.Error())
					w.WriteHeader(http.StatusNotFound)
					hb.lg.Warn().
						Msgf("message from func GetShortLink %v", err.Error())
					return
				}
				err = WriteJson(&w, slink)
				if err != nil {
					WriteJsonErr(&w, err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					hb.lg.Warn().
						Msgf("message from func GetShortLink %v", err.Error())
					return
				}
			}
		} else {
			WriteJsonErr(&w, my_errors.ErrMethodNotAllowed.Error())
			w.WriteHeader(http.StatusMethodNotAllowed)
			hb.lg.Warn().
				Msgf("message from func GetShortLink %v", my_errors.ErrMethodNotAllowed.Error())

		}
	}
}

func (hb *HandlersBuilder) Redirect() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			findhashlink := services.HashLink{
				ShortLink: r.URL.Query().Get("url"),
			}
			llink, err := hb.s.GetLongLink(findhashlink)
			if err != nil {
				WriteJsonErr(&w, err.Error())
				w.WriteHeader(http.StatusNotFound)
				hb.lg.Warn().
					Msgf("message from func Redirect %v", err.Error())
				return
			}
			http.Redirect(w, r, llink.LongLink, http.StatusSeeOther)
			hb.s.SetRedirect(llink.LongLink)
		} else {
			WriteJsonErr(&w, my_errors.ErrMethodNotAllowed.Error())
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Warn().
				Msgf("message from func Redirect %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}
}

func (hb *HandlersBuilder) GetStat() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			findhashlink := services.HashLink{
				ShortLink: r.URL.Query().Get("url"),
			}
			llink, err := hb.s.GetLongLink(findhashlink)
			if err != nil {
				WriteJsonErr(&w, err.Error())
				w.WriteHeader(http.StatusNotFound)
				log.Warn().
					Msgf("message from func GetStat %v", err.Error())
				return
			}
			red, err := hb.s.GetStatLink(llink.LongLink)
			if err != nil {
				WriteJsonErr(&w, err.Error())
				w.WriteHeader(http.StatusNotFound)
				log.Warn().
					Msgf("message from func GetStat %v", err.Error())
				return
			}
			fmt.Fprint(w, "Переходов по ссылке: ", red)
		} else {
			WriteJsonErr(&w, my_errors.ErrMethodNotAllowed.Error())
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Warn().
				Msgf("message from func GetStat %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}
}

func (hb *HandlersBuilder) GetAllStat() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			array := hb.s.GetAllStat()
			if len(array) == 0 {
				fmt.Fprint(w, "There are no links in the collection")
			} else {
				w.Header().Set("Content-Type", "application/json") //проставляем заголовок для json
				err := json.NewEncoder(w).Encode(array)
				if err != nil {
					WriteJsonErr(&w, err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					log.Warn().
						Msgf("message from func GetAllStat %v", err.Error())
				}
			}
		} else {
			WriteJsonErr(&w, my_errors.ErrMethodNotAllowed.Error())
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Warn().
				Msgf("message from func GetAllStat %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}
}
