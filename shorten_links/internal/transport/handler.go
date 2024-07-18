package transport

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"net/http"
	my_errors "shorten_links/internal/errors"
	"shorten_links/internal/services"

	"github.com/fasthttp/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// type Handler interface {
// 	GetShortLink()
// 	Redirect()
// 	GetStat()
// 	GetAllStat()
// }

type HandlersBuilder struct {
	s    Set
	lg   zerolog.Logger
	rout *router.Router
}

func HandleCreate() {
	hb := HandlersBuilder{
		s:    services.NewSetHashLink(),
		lg:   zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.UnixDate}),
		rout: router.New(),
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.POST("/shortlink/get", hb.GetShortLink())
	hb.rout.GET("/shortlink/redirect", hb.Redirect())
	hb.rout.GET("/shortlink/stat", hb.GetStat())
	hb.rout.GET("/shortlink/allstat", hb.GetAllStat())
	fmt.Println(fasthttp.ListenAndServe(":80", hb.rout.Handler))
}

func (hb *HandlersBuilder) GetShortLink() func(ctx *fasthttp.RequestCtx) {
	hb.lg.Info().
		Msgf("Start func GetShortLink")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsPost() {
			longlink, timelife, err := ParseJsonL(ctx)
			if err != nil {
				err_ := WriteJsonErr(ctx, err.Error())

				if err_ != nil {
					hb.lg.Warn().
						Msgf("message from func GetShortLink %v", err_.Error())
				}

				hb.lg.Warn().
					Msgf("message from func GetShortLink %v", err.Error())
			} else {
				slink, err := hb.s.CreateShortLink(longlink, timelife)
				if err != nil {
					WriteJsonErr(ctx, err.Error())
					hb.lg.Warn().
						Msgf("message from func GetShortLink with Set %v", err.Error())
				}

				err = WriteJson(ctx, slink)
				if err != nil {
					WriteJsonErr(ctx, err.Error())
					hb.lg.Warn().
						Msgf("message from func GetShortLink %v", err.Error())
				}
			}
		} else {
			WriteJsonErr(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			hb.lg.Warn().
				Msgf("message from func GetShortLink %v", my_errors.ErrMethodNotAllowed.Error())

		}
	}, "GetShortLink")
}

func (hb *HandlersBuilder) Redirect() func(ctx *fasthttp.RequestCtx) {
	hb.lg.Info().
		Msgf("Start func Redirect")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			findhashlink := string(ctx.QueryArgs().Peek("url"))
			llink, err := hb.s.GetLongL(findhashlink)
			if err != nil {
				WriteJsonErr(ctx, err.Error())
				hb.lg.Warn().
					Msgf("message from func Redirect s %v", err.Error())
				return
			}
			ctx.Redirect(llink, http.StatusSeeOther)
			err = hb.s.Increment(findhashlink)
			if err != nil {
				WriteJsonErr(ctx, err.Error())
				hb.lg.Warn().
					Msgf("message from func Redirect s %v", err.Error())
			}
		} else {
			WriteJsonErr(ctx, my_errors.ErrMethodNotAllowed.Error())
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			log.Warn().
				Msgf("message from func Redirect %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}, "Redirect")
}

func (hb *HandlersBuilder) GetStat() func(ctx *fasthttp.RequestCtx) {
	hb.lg.Info().
		Msgf("Start func GetStat")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			findhashlink := string(ctx.QueryArgs().Peek("url"))
			redirect, err := hb.s.GetRedirect(findhashlink)
			if err != nil {
				WriteJsonErr(ctx, err.Error())
				log.Warn().
					Msgf("message from func GetStat %v", err.Error())
				return
			}
			fmt.Fprint(ctx, "Переходов по ссылке: ", redirect)
		} else {
			WriteJsonErr(ctx, my_errors.ErrMethodNotAllowed.Error())
			log.Warn().
				Msgf("message from func GetStat %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}, "statistic_one")
}

func (hb *HandlersBuilder) GetAllStat() func(ctx *fasthttp.RequestCtx) {
	hb.lg.Info().
		Msgf("Start func GetAllStat")
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if ctx.IsGet() {
			array, err := hb.s.GetAllStat()
			if err != nil {
				WriteJsonErr(ctx, err.Error())
				log.Warn().
					Msgf("message from func GetAllStat %v", err.Error())
			}
			if len(array) == 0 {
				fmt.Fprint(ctx, "There are no links in the collection")
			} else {
				ctx.SetContentType("application/json")
				ctx.Response.BodyWriter()
				err := json.NewEncoder((*ctx).Response.BodyWriter()).Encode(array)
				if err != nil {
					WriteJsonErr(ctx, err.Error())
					log.Warn().
						Msgf("message from func GetAllStat %v", err.Error())
				}
			}
		} else {
			WriteJsonErr(ctx, my_errors.ErrMethodNotAllowed.Error())
			log.Warn().
				Msgf("message from func GetAllStat %v", my_errors.ErrMethodNotAllowed.Error())
		}
	}, "all_statistic")
}
