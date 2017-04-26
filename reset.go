package main

import (
	"net/http"

	"github.com/hashicorp/vault/api"
	"github.com/mijia/sweb/form"
	"github.com/mijia/sweb/log"
	"golang.org/x/net/context"
)

func (l *Lvault) settokenkeys(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {
	initResponse := api.InitResponse{}
	if err := form.ParamBodyJson(req, &initResponse); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		if l.GetMissToken() {
			if l.Vault.CheckRootToken(initResponse.RootToken) {
				l.SetMissToken(false)
				l.RootToken = initResponse.RootToken
			} else {
				log.Debug("invalid root token")
				l.RootToken = initResponse.RootToken
			}
		}
		l.UnsealKey = initResponse.Keys
	}
	return ctx
}
