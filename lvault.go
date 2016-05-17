package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/mijia/sweb/log"

	"github.com/laincloud/lvault/vaultcli"
)

type Lvault struct {

	//store the ip or something else to contact with other instances
	Others []string
	others map[string]string

	RootToken string

	//for unseal keys, can be an array, but often single
	UnsealKey []string

	//the struct for call vault's api for lvault
	Vault *vaultcli.VaultClient

	L sync.RWMutex

	ContainerAddr string

	InstanceNum string

	missToken bool

	SSOSite      string
	ClientId     string
	ClientSecret string

	HTTPS bool
}

const (
	etcdURL               = "http://etcd.lain:4001"
	etcdLvaultTokenStatus = "/v2/keys/vault/lvault/tokenstatus/"
)

type status struct {
	Host   string
	IsMiss bool
}

func (l *Lvault) Init() error {
	l.Others = make([]string, 0, 5)
	l.others = make(map[string]string)
	l.Vault = &vaultcli.VaultClient{}

	l.Vault.InitClient(l.HTTPS)

	procname := os.Getenv("LAIN_PROCNAME")
	insnum := os.Getenv("DEPLOYD_POD_INSTANCE_NO")
	appname := os.Getenv("LAIN_APPNAME")
	l.InstanceNum = insnum
	l.ContainerAddr = "http://" + procname + "-" + insnum + "." + appname + ".lain:8001"

	l.InitMissToken()
	return nil
}
func (l *Lvault) InitMissToken() {
	//加锁的必要性不大
	l.L.Lock()
	l.missToken = true
	l.setTokenStatus(true)
	l.L.Unlock()
}

func (l *Lvault) SetMissToken(miss bool) {
	l.L.Lock()
	if l.missToken != miss {
		l.missToken = miss
		l.setTokenStatus(miss)
	}
	l.L.Unlock()
}

func (l *Lvault) GetMissToken() (miss bool) {
	l.L.RLock()
	miss = l.missToken
	l.L.RUnlock()
	return miss
}

func (l *Lvault) setTokenStatus(miss bool) {

	tmpStatus := status{l.ContainerAddr, miss}
	b, _ := json.Marshal(tmpStatus)
	url := etcdURL + etcdLvaultTokenStatus + l.InstanceNum + "?value=" + string(b)
	req, _ := http.NewRequest("PUT", url, nil)
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err == nil {
		resp.Body.Close()
	} else {
		log.Error(err)
	}
}
