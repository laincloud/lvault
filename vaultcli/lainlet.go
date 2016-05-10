package vaultcli

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	//	"github.com/mijia/sweb/log"
)

const (
	lainletServiceAddr = "http://lainlet.lain:9001"
	// LAINLET_SERVICE_TIMEOUT = 1 * time.Minute
)

type ProcInfo struct {
	NumInstances int             `json:"num_instances"`
	Containers   []ContainerInfo `json:"containers"`
}

type ContainerInfo struct {
	ContainerIp   string `json:"container_ip"`
	ContainerPort int    `json:"container_port"`
}

func ProcWatcherUrl(appName string) (url string) {
	lainletAddr := os.Getenv("LAINLET_WEBSERVICE")
	if lainletAddr == "" {
		lainletAddr = lainletServiceAddr
		//	log.Warn("environmental variable LAINLEI_WEBSERVICE is null")
	}
	url = lainletAddr + "/v2/proxywatcher/?appname=" + appName + "&watch=0"
	//log.Debug("procWatcherUrl is: ", url)
	return url
}

func GetEventFrom(resp *http.Response) (data map[string]*ProcInfo, err error) {
	bdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error parse data from http.Response: " + err.Error())
	}
	//log.Debug(string(bdata))

	err = json.Unmarshal(bdata, &data)
	if err != nil {
		return nil, errors.New("Error Unmarshal data to ProcInfo: " + err.Error())
	}
	return data, nil
}
