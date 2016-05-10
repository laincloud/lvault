package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/mijia/sweb/log"

	"golang.org/x/net/context"
)

var (
	slugPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*(.[a-zA-Z0-9]+)*$`)
	coder       = base64.StdEncoding
)

const (
	APPNAME        string = "app"
	PROCNAME              = "proc"
	PATHNAME              = "path"
	DATAPATHPREFIX        = "secret/lvault/"
)

func ValidateSlug(slug string) error {
	if slug == "" {
		return errors.New("Empty slug")
	}

	if len(slug) > 256 {
		return errors.New("Slug too long")
	}

	if !slugPattern.MatchString(slug) {
		log.Debug("don't match")
		return errors.New("Invalid slug")
	}

	return nil
}

func parseErr(err error) (int, []byte) {
	if strings.Contains(err.Error(), "Vault is sealed") {
		return http.StatusServiceUnavailable, []byte(err.Error())
	} else {
		panic(err)
	}
}

func (l *Lvault) SecretsEndpoint(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {
	log.Debug(req)
	body, _ := ioutil.ReadAll(req.Body)
	req.ParseForm()

	// Auth
	if ok := checkMaintainerAuth(req); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"authentication failed"}`))
		return ctx
	}

	app, proc, path, err := getCheckedParams(req)
	log.Debug(app)
	log.Debug(proc)
	log.Debug(err)
	log.Debug(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return ctx
	}
	switch req.Method {
	case "GET":
		var totalPath string
		if proc == "" {
			totalPath = DATAPATHPREFIX + app
		} else {
			totalPath = DATAPATHPREFIX + app + "/" + proc
		}
		secrets, err := l.Vault.ListSecrets(l.RootToken, totalPath)
		if err != nil {
			log.Error(err)
			code, retD := parseErr(err)
			w.WriteHeader(code)
			w.Write(retD)
			return ctx
		}
		var data []Data
		data = []Data{}
		for i := 0; i < len(secrets); i++ {
			var singleData Data
			plainData, _ := coder.DecodeString(secrets[i])
			err = json.Unmarshal(plainData, &singleData)
			if err != nil {
				log.Error(err)
				panic(err)
			}
			if !strings.EqualFold(singleData.Path, "") {
				data = append(data, singleData)
			}
		}
		retData, err := json.Marshal(data)
		if err != nil {
			log.Error(err)
			panic(err)
		}
		if len(data) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Secrets Not Found"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(retData)
		}
	case "PUT":
		if path == "" || proc == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("require parameters"))
			return ctx
		}
		totalPath := DATAPATHPREFIX + app + "/" + proc + path
		putData := ParseInput((app + "/" + proc + path), string(body))
		byteData, _ := json.Marshal(putData)
		byteDataBase64 := coder.EncodeToString(byteData)
		log.Debug(string(byteData))
		err = l.Vault.PutSecret(l.RootToken, totalPath, byteDataBase64)
		if err != nil {
			log.Error(err)
			code, retD := parseErr(err)
			w.WriteHeader(code)
			w.Write(retD)
			return ctx
		}
	case "DELETE":
		if path == "" || proc == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("require parameters"))
			return ctx
		}
		totalPath := app + "/" + proc + path
		err = l.Vault.DeleteSecret(l.RootToken, DATAPATHPREFIX+totalPath)
		if err != nil {
			log.Error(err)
			code, retD := parseErr(err)
			w.WriteHeader(code)
			w.Write(retD)
			return ctx
		}
	}
	return ctx
}

func getCheckedParams(req *http.Request) (app string, proc string, path string, err error) {
	log.Debug(req.Form)
	if len(req.Form[APPNAME]) > 0 {
		app = req.Form[APPNAME][0]
		if err = ValidateSlug(app); err != nil {
			return
		}
	} else {
		app = ""
		err = errors.New("can't find appname")
		return
	}
	if len(req.Form[PROCNAME]) > 0 {
		proc = req.Form[PROCNAME][0]
		if err = ValidateSlug(proc); err != nil {
			return
		}
	} else {
		proc = ""
	}
	if len(req.Form[PATHNAME]) > 0 {
		path = req.Form[PATHNAME][0]
		if !strings.HasPrefix(path, "/") {
			err = errors.New("path should begin with /")
			return
		}
		if strings.Index(path, "..") != -1 {
			err = errors.New("path should not contain ..")
			return
		}
	} else {
		path = ""
	}
	return
}
