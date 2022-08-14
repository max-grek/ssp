package service

import (
	"context"
	"net/http"
	"test-assignment-cookie-sync/connector"

	satoriuuid "github.com/satori/go.uuid"
)

const COOKIE = "cookie"

type CookieS struct {
	conn connector.DB
}

func NewCookie(conn connector.DB) *CookieS {
	return &CookieS{conn}
}

func (c *CookieS) ProcessNotify(ctx context.Context, cookieId string) error {
	return c.conn.Cookie().PersistNotify(ctx, satoriuuid.NewV4().String(), cookieId)
}

func (c *CookieS) ProcessCookie(ctx context.Context) (*http.Cookie, error) {
	cid := satoriuuid.NewV4().String()
	return &http.Cookie{
		Name:   "ssp",
		Value:  cid,
		MaxAge: 0,
	}, c.conn.Cookie().PersistCookie(ctx, satoriuuid.NewV4().String(), cid)
}
