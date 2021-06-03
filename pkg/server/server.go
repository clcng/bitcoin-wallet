package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/clcng/bitcoin-wallet/pkg/log"
	"github.com/go-chi/chi"
)

type Server struct {
}

func WalletServ(addr string) error {
	s := &Server{}

	defer func() {
		if r := recover(); r != nil {
			log.Fatal().Msg(string(debug.Stack()))
		}
	}()

	r := chi.NewRouter()

	r.Use(loggingHandler)

	//health check path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/wallet", func(r chi.Router) {
		r.Get("/mnemonic", s.handleGenerateBIP39Mnemonic)
		r.Post("/hdAddress", s.handleGenerateHDAddress)
	})

	log.Info().Msgf("wallet service running at %v", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal().Err(err)
	}
	return nil
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ctx := context.WithValue(r.Context(), "id", now.Nanosecond())

		param := ""
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				writeRst(w, err)
				return
			}
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			s := fmt.Sprintf("%q", rdr1)
			s = strings.Replace(s, "\\n", "", -1)
			param = strings.Replace(s, "\\", "", -1)
			r.Body = rdr1
		} else {
			param = r.URL.Query().Encode()
		}
		log.Info().Interface("ctxId", ctx.Value("id")).Str("method", r.Method).Interface("header", r.Header).Str("param", param).Str("requestUrl", r.URL.Path).Msg("http request")

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeRst(w http.ResponseWriter, v interface{}) {
	if err, ok := v.(error); ok {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if jenc, ok := v.(json.Marshaler); ok {
			ba, _ := jenc.MarshalJSON()
			w.Write(ba)
		} else {
			msg := fmt.Sprintf(`{"message":"%s"`, err.Error())
			w.Write([]byte(msg))
		}
	} else {
		ba, err := json.Marshal(v)
		if err == nil {
			w.Header().Add("Content-Type", "application/json")
			w.Write(ba)
		} else {
			msg := fmt.Sprintf(`{"message":"%s"`, err.Error())
			w.Write([]byte(msg))
		}
	}
}
