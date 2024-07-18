package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	my_errors "shorten_links/internal/errors"
	"time"

	"github.com/valyala/fasthttp"
)

func ParseJsonL(ctx *fasthttp.RequestCtx) (string, time.Duration, error) {
	var data struct {
		LongLink string `json:"long_link"`
		TimeLife string `json:"time_life"`
	}
	err := json.NewDecoder(bytes.NewReader(ctx.Request.Body())).Decode(&data)
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

func WriteJson(ctx *fasthttp.RequestCtx, s string) error {
	ctx.SetContentType("application/json")
	ctx.Response.BodyWriter()
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(s)
	if err != nil {
		return my_errors.ErrWriteJSON
	}
	return nil
}

func WriteJsonErr(ctx *fasthttp.RequestCtx, mes string) error {
	var err_mes struct {
		ErrorMessage string
	}
	err_mes.ErrorMessage = mes
	ctx.SetContentType("application/json")
	switch mes {
	case my_errors.ErrNoSlink.Error():
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	case my_errors.ErrNoLlink.Error():
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	case my_errors.ErrMethodNotAllowed.Error():
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	case my_errors.ErrEqualJSON.Error():
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(err_mes)
	if err != nil {
		return my_errors.ErrWriteJSONerr
	}
	return nil
}
