package transport

import (
	"encoding/json"
	"net/http"
	my_errors "shorten_links/internal/errors"
	"shorten_links/internal/services"
	"time"
)

func ParseJsonL(r *http.Request) (string, time.Duration, error) {
	var data struct {
		LongLink string `json:"long_link"`
		TimeLife string `json:"time_life"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	var empty time.Duration
	if err != nil {
		return data.LongLink, empty, my_errors.ErrParseJSON
	}
	if data.LongLink == "" {
		return data.LongLink, empty, my_errors.ErrEqualJSON
	}
	d, err := time.ParseDuration(data.TimeLife)
	if err != nil {
		return data.LongLink, empty, my_errors.ErrParseDuration
	}

	return data.LongLink, d, err
}

func ParseJsonS(r *http.Request) (string, error) {
	var data struct {
		ShortLink string `json:"short_link"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return data.ShortLink, my_errors.ErrParseJSON
	}
	return data.ShortLink, err
}

func WriteJson(w *http.ResponseWriter, s services.HashLink) error {
	(*w).Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(*w).Encode(s)
	if err != nil {
		return my_errors.ErrWriteJSON
	}
	return nil
}
func WriteJsonErr(w *http.ResponseWriter, mes string) error {
	var err_mes struct {
		ErrorMessage string
	}
	err_mes.ErrorMessage = mes
	(*w).Header().Set("Content-Type", "application/json") //проставляем заголовок для json
	(*w).WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(*w).Encode(err_mes)
	if err != nil {
		return my_errors.ErrWriteJSONerr
	}
	return nil
}
