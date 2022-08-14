package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-assignment-cookie-sync/service"
)

type cookieHandler struct {
	*service.CookieS
}

func NewCookieHandler(sync *service.CookieS) cookieHandler {
	return cookieHandler{sync}
}

func (sync cookieHandler) Notify(w http.ResponseWriter, r *http.Request) {
	m, err := readData(r, "cookie_id")
	if err != nil {
		responseBuilder(w, http.StatusBadRequest, fmt.Errorf("parse body: %v", err))
		return
	}
	err = sync.CookieS.ProcessNotify(r.Context(), m["cookie_id"].(string))
	if err != nil {
		responseBuilder(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("success"))
}

func (sync cookieHandler) Cookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := sync.CookieS.ProcessCookie(r.Context())
	if err != nil {
		responseBuilder(w, http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
