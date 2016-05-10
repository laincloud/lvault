package main

import (
	"net/http"

	"golang.org/x/net/context"
)

const (
	unsealInterval = 1 //second
)

func (l *Lvault) unseal(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {
	l.Unseal()
	return ctx
}

func (l *Lvault) Unseal() {
	err := l.Vault.Unseal(l.UnsealKey)
	if err != nil {
		panic(err)
	}
}
