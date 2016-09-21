package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/hashicorp/vault/api"
	"github.com/mijia/sweb/form"
	"github.com/mijia/sweb/log"
)

const (
	SECRET_SHARES    = 1
	SECRET_THRESHOLD = 1
)

func (l *Lvault) init(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {

	initReq := api.InitRequest{}
	if err := form.ParamBodyJson(req, &initReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return ctx
	}

	if initReq.SecretShares == 0 || initReq.SecretThreshold == 0 || initReq.SecretThreshold > initReq.SecretShares {
		initReq.SecretShares = SECRET_SHARES
		initReq.SecretThreshold = SECRET_THRESHOLD
	}

	log.Debug(initReq)
	initResp, err := l.Vault.InitVault(&initReq)
	log.Debug(initResp)

	if err != nil {
		log.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return ctx
	} else {
		l.RootToken = initResp.RootToken
		l.UnsealKey = initResp.Keys
		l.SetMissToken(false)
		tk_keys, _ := json.Marshal(initResp)
		w.Write(tk_keys)
		return ctx
	}
}
