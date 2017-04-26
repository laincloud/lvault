package vaultcli

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mijia/sweb/log"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (c *VaultClient) GetAllStatus() []byte {
	return c.status.GetStatus()
}

const (
	AppName        = "lvault"
	ProcName       = "lvault.web.web"
	UpdateInterval = 2 // 更新 vault 状态的周期，单位为秒
)

type VaultStatus struct {
	Lock sync.RWMutex

	//The key is container_ip+container_port
	Containers map[string]ContainerInfoWithStatus

	Cli *LainletClient

	unsealedAddress string
}

type LainletClient struct {
	c    *http.Client
	resp *http.Response
}

type ContainerStatus struct {
	Sealed          bool `json:"sealed"`
	Threshold       int  `json:"t"`
	TotalSharedKeys int  `json:"n"`
	Progress        int  `json:"progress"`
}

type ContainerInfoWithStatus struct {
	Info   ContainerInfo
	Status ContainerStatus
}

func (v *VaultStatus) UnsealedURL() string {
	v.Lock.RLock()
	defer v.Lock.RUnlock()
	return v.unsealedAddress
}

func (v *VaultStatus) AllUnsealed() bool {
	return v.sealStatus(true)
}

func (v *VaultStatus) AllSealed() bool {
	return v.sealStatus(false)
}

func (v *VaultStatus) sealStatus(all bool) bool {
	v.Lock.RLock()
	defer v.Lock.RUnlock()

	for _, status := range v.Containers {
		if status.Status.Sealed == all {
			return false
		}
	}
	return true
}

// return a copy of the corresponding slice
func (v *VaultStatus) GetContainers() map[string]ContainerInfoWithStatus {
	v.Lock.RLock()
	ret := make(map[string]ContainerInfoWithStatus)
	for k, d := range v.Containers {
		ret[k] = d
	}
	v.Lock.RUnlock()
	return ret
}

func (v *VaultStatus) GetStatus() []byte {
	containers := v.GetContainers()
	b, err := json.Marshal(containers)
	if err != nil {
		log.Error(err)
	}
	return b
}

func (v *VaultStatus) updateContainers() {
	for {
		time.Sleep(time.Duration(UpdateInterval) * time.Second)
		v.Lock.Lock()
		for k, d := range v.Containers {
			status := getContainerStatus(d.Info)
			v.Containers[k] = ContainerInfoWithStatus{d.Info, status}
		}
		v.Lock.Unlock()
	}
}

func (v *VaultStatus) UpdateStatus() {
	v.Containers = make(map[string]ContainerInfoWithStatus)
	v.Cli = &LainletClient{}
	v.Cli.c = &http.Client{}
	//	go v.doQueryLainlet()
	go v.updateContainers()
	for {
		time.Sleep(1 * time.Second)
		data, err := v.doQueryLainlet()
		if err == nil {
			//此时 data 应该是一个包含所有 proc 数据的 map
			if _, ok := data[ProcName]; !ok {
				message, err := json.Marshal(data)
				log.Error(string(message), err)
				continue
			}
			CInfo := data[ProcName].Containers
			//	fmt.Println("num instances:", data[ProcName].NumInstances)
			// always 0 bug of lainlet?
			for _, d := range CInfo {
				ip_port := d.ContainerIp + ":" + strconv.Itoa(d.ContainerPort)
				tmp := getContainerStatus(d)
				v.Lock.Lock()
				v.Containers[ip_port] = ContainerInfoWithStatus{d, tmp}
				if tmp.Sealed == false {
					v.unsealedAddress = "http://" + ip_port
					if isUsingHTTPS() {
						v.unsealedAddress = "https://" + ip_port
					}
				}
				v.Lock.Unlock()
			}
		}
	}
}

func isUsingHTTPS() bool {
	return strings.HasPrefix(VaultURL, "https")
}

func getContainerStatus(info ContainerInfo) ContainerStatus {
	ip_port := info.ContainerIp + ":" + strconv.Itoa(info.ContainerPort)
	url := "http://" + ip_port + "/v1/sys/seal-status"
	if isUsingHTTPS() {
		url = "https://" + ip_port + "/v1/sys/seal-status"
	}
	var tmp ContainerStatus
	resp, geterr := http.Get(url)
	for geterr != nil {
		time.Sleep(1 * time.Second)
		resp, geterr = http.Get(url)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jsonerr := json.Unmarshal(body, &tmp)
	if jsonerr != nil {
	}
	return tmp
}

func (v *VaultStatus) doQueryLainlet() (data map[string]*ProcInfo, err error) {
	v.Cli.resp, err = v.Cli.c.Get(ProcWatcherUrl(AppName))
	if err != nil {
		log.Error(err)
		return nil, err
	} else if v.Cli.resp == nil {
		log.Error("nil response")
		return nil, err
	} else {
		defer v.Cli.resp.Body.Close()
		return GetEventFrom(v.Cli.resp)
	}
}
