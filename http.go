package main

import (
	"fmt"
	"log"
	"net/http"
	"test-assignment-cookie-sync/config"
	"test-assignment-cookie-sync/handler"
	"test-assignment-cookie-sync/service"
	"time"

	"github.com/rs/cors"
)

type httpService struct {
	*config.HTTPConfig
	router *http.ServeMux
}

func newHTTPService(cfg *config.HTTPConfig) *httpService {
	return &httpService{
		HTTPConfig: cfg,
		router:     http.NewServeMux(),
	}
}

func (h *httpService) registerRoutes(ss *service.CookieS) {
	h.registerSync(ss)
}

func (h *httpService) registerSync(ss *service.CookieS) {
	sh := handler.NewCookieHandler(ss)
	h.register(http.MethodGet, "/api/cookie", sh.Cookie)
	h.register(http.MethodPost, "/api/notify", sh.Notify)
}

func (h *httpService) register(method, path string, handler http.HandlerFunc) {
	timeout, _ := time.ParseDuration(h.HTTPConfig.Timeout)
	h.router.Handle(path, http.TimeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming connection from %s to %s", r.RemoteAddr, path)
		handler.ServeHTTP(w, r)
	}), timeout, "request canceled by timeout"))
}

func (h *httpService) run() error {
	var listenAddress = fmt.Sprintf("%s:%d", h.HTTPConfig.Host, h.HTTPConfig.Port)
	log.Printf("service listening on %s", listenAddress)
	return http.ListenAndServe(listenAddress, cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:9000"},
		AllowedHeaders: []string{"*"},
	}).Handler(h.router))
}
