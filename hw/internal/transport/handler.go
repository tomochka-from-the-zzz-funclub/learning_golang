package transport

import (
	"fmt"
	my_errors "hw/internal/errors"
	"hw/internal/services"
	"net/http"
)

func HandleCreate() {
	hb := HandlersBuilder{
		s: services.NewSet(),
	}
	http.HandleFunc("/add", hb.AddElemBuild())
	http.HandleFunc("/delete/elem", hb.DeleteElemBuild())
	http.HandleFunc("/delete/all", hb.DeleteAllBuild())
	http.HandleFunc("/check", hb.CheckElemBuild())
}

func (hb *HandlersBuilder) AddElemBuild() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := hb.s.Add(r.URL.Query().Get("element"))
			if err != nil {
				if err == my_errors.ErrRepeat {
					w.WriteHeader(http.StatusFound)
				} else {
					w.WriteHeader(http.StatusUnsupportedMediaType)
				}
			} else {
				fmt.Fprint(w, "Элемент ", r.URL.Query().Get("element"), " успешно добавлен")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (hb *HandlersBuilder) DeleteElemBuild() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			err := hb.s.DeleteElem(r.URL.Query().Get("element"))
			if err != nil {
				if err == my_errors.ErrNoElem {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusUnsupportedMediaType)
				}
			} else {
				fmt.Fprint(w, "Элемент ", r.URL.Query().Get("element"), " успешно удален")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (hb *HandlersBuilder) DeleteAllBuild() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			hb.s.DeleteAll()
			fmt.Fprint(w, "Коллекция успешно удалена")
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (hb *HandlersBuilder) CheckElemBuild() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := hb.s.Check(r.URL.Query().Get("element"))
			if err == nil {
				fmt.Fprint(w, "Элемент ", r.URL.Query().Get("element"), " существует в коллекции")
			} else if err == my_errors.ErrNoElem {
				fmt.Fprint(w, "Элемента ", r.URL.Query().Get("element"), " нет в коллекции")
			} else {
				w.WriteHeader(http.StatusUnsupportedMediaType)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

type HandlersBuilder struct {
	s services.Set
}
