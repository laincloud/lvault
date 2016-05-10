package main

import (
	"encoding/json"
	"strings"
	"time"
)

type Data struct {
	Path      string `json:"path,omitempty"`
	Content   string `json:"content"`
	Owner     string `json:"owner,omitempty"`
	Group     string `json:"group,omitempty"`
	Mode      string `json:"mode,omitempty"`
	TimeStamp int64  `json:"timestamp"`
}

func ParseInput(path string, body string) (ret Data) {
	/*
		若 body 的格式为
		{
		    "content": "",
		    "owner": "",
		    "group": "",
		    "mode": ""
		}
		后面三项可以忽略，则将 content owner 等作为 key；
		否则，将 body 整体作为 content.
	*/
	// TODO 不合法的 mode 等边界情况检查
	var data Data
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		data = Data{
			Path:    path,
			Content: body,
		}
	} else {
		data.Path = path
	}
	if strings.EqualFold(data.Owner, "") {
		data.Owner = "root"
	}
	if strings.EqualFold(data.Group, "") {
		data.Group = "root"
	}
	if strings.EqualFold(data.Mode, "") {
		data.Mode = "400"
	}
	//data.TimeStamp = fmt.Sprint(time.Now().Unix())
	data.TimeStamp = time.Now().Unix()
	return data
}
