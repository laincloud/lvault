package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/vault/api"
	"github.com/mijia/sweb/log"
)

const (
	sendTokenInterval = 1
)

type etcdResponse struct {
	Node Instances `json:"node"`
}

type Instances struct {
	Nodes []Instance `json:"nodes"`
}

type Instance struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (l *Lvault) lvaultStatus(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {
	url := etcdURL + etcdLvaultTokenStatus + "?recursive=true"
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return ctx
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	var test etcdResponse
	err = json.Unmarshal(b, &test)
	if err != nil {
		log.Error(err)
	}
	ret := make([]status, 0, 5)
	for _, vContainerStatus := range test.Node.Nodes {
		var sta status
		json.Unmarshal([]byte(vContainerStatus.Value), &sta)
		ret = append(ret, sta)
	}
	lsta, _ := json.Marshal(ret)
	log.Debug(string(lsta))
	w.Write(lsta)
	return ctx
}

func (l *Lvault) SendToken() {
	for {
		time.Sleep(time.Duration(sendTokenInterval) * time.Second)
		if !l.GetMissToken() {
			shouldClose := false
			url := etcdURL + etcdLvaultTokenStatus + "?recursive=true"
			resp, err := http.Get(url)
			if resp != nil {
				shouldClose = true
			}
			if err != nil {
				log.Error(err)
				if shouldClose {
					resp.Body.Close()
				}
				continue
			}
			b, _ := ioutil.ReadAll(resp.Body)
			//	b = []byte(`{"action":"get","node":{"key":"/vault/lvault/tokenstatus","dir":true,"nodes":[{"key":"/vault/lvault/tokenstatus/1","value":"{\"Host\":\"http://web-1.lvault.lain:8001\",\"IsMiss\":true}","modifiedIndex":115892,"createdIndex":115892}],"modifiedIndex":111172,"createdIndex":111172}}`)
			var test etcdResponse
			err = json.Unmarshal(b, &test)
			if err != nil {
				log.Error(err)
			}
			for _, vContainerStatus := range test.Node.Nodes {
				var sta status
				json.Unmarshal([]byte(vContainerStatus.Value), &sta)
				if !l.IsHostExists(sta.Host) {
					go l.AddHost(sta.Host)
				}
				if sta.IsMiss {
					l.sendToken(sta.Host)
				}
			}
			if shouldClose {
				resp.Body.Close()
			}
		}
	}
}

func (l *Lvault) sendToken(host string) {
	c := http.DefaultClient
	req, _ := http.NewRequest("PUT", host+"/reset", nil)
	tk_keys := api.InitResponse{
		Keys:      l.UnsealKey,
		RootToken: l.RootToken,
	}
	b, _ := json.Marshal(tk_keys)
	req.Body = ioutil.NopCloser(strings.NewReader(string(b)))
	res, err := c.Do(req)
	if err == nil {
		defer res.Body.Close()
	} else {
		log.Error(err)
	}
}

func (l *Lvault) vaultStatus(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {
	ret := l.Vault.GetAllStatus()
	mapdata := make(map[string]interface{})
	json.Unmarshal(ret, &mapdata)
	sl := make([]interface{}, 0, 5)
	for _, v := range mapdata {
		sl = append(sl, v)
	}
	ret, _ = json.Marshal(sl)
	w.Write(ret)

	return ctx
}
