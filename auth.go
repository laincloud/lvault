package main

import (
	"io/ioutil"
	"net/http"

	"github.com/mijia/sweb/log"
	"golang.org/x/net/context"
)

const SCOPE = ""

func (l *Lvault) auth(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	r.ParseForm()
	redirect_uri := r.Form.Get("redirect_uri")
	response_type := r.Form.Get("response_type")
	state := r.Form.Get("state")
	client_id := l.ClientId
	scope := SCOPE
	url := l.SSOSite + r.URL.Path + "?response_type=" + response_type + "&redirect_uri=" + redirect_uri + "&client_id=" + client_id + "&scope=" + scope + "&state=" + state
	log.Debug(url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusSeeOther)
	return ctx
}

func (l *Lvault) token(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	r.ParseForm()
	redirect_uri := r.Form.Get("redirect_uri")
	code := r.Form.Get("code")
	grant_type := r.Form.Get("grant_type")
	client_id := l.ClientId
	client_secret := l.ClientSecret
	url := l.SSOSite + r.URL.Path + "?grant_type=" + grant_type + "&redirect_uri=" + redirect_uri + "&client_id=" + client_id + "&code=" + code + "&client_secret=" + client_secret
	log.Debug(url)
	req, _ := http.NewRequest(r.Method, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
	} else {
		log.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return ctx
	}
	log.Debug(resp)
	b, _ := ioutil.ReadAll(resp.Body)
	log.Debug(string(b))
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
	return ctx
}
