package main

import (
	"encoding/json"
	"github.com/mijia/sweb/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ConsoleResponse struct {
	Role MaintainerRole `json:"role"`
}

type MaintainerRole struct {
	Role string `json:"role"`
}

func checkMaintainerAuth(req *http.Request) bool {
	//	return true
	var ap string
	if len(req.Form[APPNAME]) > 0 {
		ap = req.Form[APPNAME][0]
	} else {
		return false
	}

	actk := req.Header.Get("access-token")
	log.Debug(actk)

	if strings.EqualFold(actk, "") {
		actk = req.Header.Get("Access-Token")
		log.Debug(actk)
		if actk == "" {
			return false
		}
	}

	domain := os.Getenv("LAIN_DOMAIN")
	url := "http://console." + domain + "/api/v1/repos/" + ap + "/roles/"
	conReq, _ := http.NewRequest("GET", url, nil)
	conReq.Header.Set("access-token", actk)
	resp, err := http.DefaultClient.Do(conReq)
	if err != nil {
		log.Error(err)
		return false
	} else {
		defer resp.Body.Close()
		var tmp ConsoleResponse
		resBody, _ := ioutil.ReadAll(resp.Body)
		log.Debug(string(resBody))
		err = json.Unmarshal(resBody, &tmp)
		if err != nil {
			log.Error(err)
			return false
		}
		log.Debug(tmp)
		if strings.EqualFold(tmp.Role.Role, "") {
			return false
		} else {
			return true
		}
	}
}
