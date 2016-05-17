package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof"

	"github.com/getsentry/raven-go"
	"github.com/mijia/sweb/log"
	"github.com/mijia/sweb/render"
	"github.com/mijia/sweb/server"

	"golang.org/x/net/context"
)

const (
	kContentCharset = "; charset=UTF-8"
	kContentJson    = "application/json"
)

type Server struct {
	*server.Server
	render *render.Render

	isDebug bool
}

type ApiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	var ssoserver, ssoid, ssosecret string
	flag.StringVar(&ssoserver, "ssoserver", "https://sso.example.com", "Base URL of SSO site")
	flag.StringVar(&ssoid, "ssoid", "1", "Client ID of Lvault on SSO site")
	flag.StringVar(&ssosecret, "ssosecret", "http://localhost:8011", "Client secret of Lvault")

	var isDebug, isHttps bool
	flag.BoolVar(&isDebug, "debug", false, "Debug mode switch")
	flag.BoolVar(&isHttps, "https", false, "Scheme for communicating with the vault backend")

	flag.Parse()

	s := &Server{
		isDebug: isDebug,
	}
	if isDebug {
		log.EnableDebug()
	}

	var l Lvault
	l.HTTPS = isHttps
	l.Init()
	l.SSOSite = ssoserver
	l.ClientId = ssoid
	l.ClientSecret = ssosecret

	go l.SendToken()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "lvault", &l)

	s.Server = server.New(ctx, s.isDebug)

	sentryClient, _ := raven.NewClient("http://5:7@sentry.example.com/5", nil)

	_ = sentryClient
	//	s.Middleware(server.NewSentryRecoveryWare(sentryClient, s.isDebug))

	s.render = s.initRender()
	s.EnableExtraAssetsJson("assets_map.json")
	s.RestfulHandlerAdapter(s.adaptResourceHandler)

	s.Get("/", "Home", s.Home)

	s.Get("/spa/*lochash", "lvaultSpa", s.PageApplication)

	s.Files("/static/*filepath", http.Dir("public"))
	s.Files("/apidoc/*filepath", http.Dir("apidoc"))
	s.Files("/assets/*filepath", http.Dir("public"))

	s.Get("/secrets", "GetSecrets", l.SecretsEndpoint)
	s.Put("/secrets", "PutSecrets", l.SecretsEndpoint)
	s.Delete("/secrets", "DeleteSecrets", l.SecretsEndpoint)

	s.Put("/init", "Init", l.init)
	s.Put("/reset", "Reset", l.settokenkeys)

	s.Put("/unsealall", "UnsealAll", l.unseal)

	s.Get("/status", "lvaultStatus", l.lvaultStatus)
	s.Get("/vaultstatus", "vaultstatus", l.vaultStatus)

	s.Get("/oauth2/auth", "auth", l.auth)
	s.Get("/oauth2/token", "token", l.token)

	s.Run(":8001")
}

func (s *Server) Home(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	log.Debug(r)
	http.Redirect(w, r, "/v2/spa/", http.StatusSeeOther)
	return ctx
}

func (s *Server) PageApplication(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	log.Debug(r)
	s.renderHtmlOr500(w, http.StatusOK, "spa", nil)
	return ctx
}

func (s *Server) initRender() *render.Render {
	tSets := []*render.TemplateSet{
		render.NewTemplateSet("spa", "spa.html", "spa.html"),
	}
	r := render.New(render.Options{
		Directory:     "templates",
		Funcs:         s.renderFuncMaps(),
		Delims:        render.Delims{"{{", "}}"},
		IndentJson:    true,
		UseBufPool:    true,
		IsDevelopment: s.isDebug,
	}, tSets)
	return r
}

func (s *Server) renderHtmlOr500(w http.ResponseWriter, status int, name string, binding interface{}) {
	if err := s.render.Html(w, status, name, binding); err != nil {
		log.Errorf("Server got a rendering error, %q", err)
		if s.isDebug {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "500, Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) adaptResourceHandler(handler server.ResourceHandler) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		code, data := handler(ctx, r)
		if code < 400 {
			s.renderJsonOr500(w, code, data)
		} else {
			errMessage := ""
			if msg, ok := data.(string); ok {
				errMessage = msg
			} else if msg, ok := data.(error); ok {
				errMessage = msg.Error()
			}
			switch code {
			case http.StatusMethodNotAllowed:
				if errMessage == "" {
					errMessage = fmt.Sprintf("Method %q is not allowed", r.Method)
				}
				s.renderError(w, code, errMessage, data)
			case http.StatusNotFound:
				if errMessage == "" {
					errMessage = "Cannot find the resource"
				}
				s.renderError(w, code, errMessage, data)
			case http.StatusBadRequest:
				if errMessage == "" {
					errMessage = "Invalid request get or post params"
				}
				s.renderError(w, code, errMessage, data)
			default:
				if errMessage == "" {
					errMessage = fmt.Sprintf("HTTP Error Code: %d", code)
				}
				s.renderError(w, code, errMessage, data)
			}
		}
		return ctx
	}
}

func (s *Server) renderJsonOr500(w http.ResponseWriter, status int, v interface{}) {
	if err := s.renderJson(w, status, v); err != nil {
		s.renderError(w, http.StatusInternalServerError, err.Error(), "")
	}
}

func (s *Server) renderJson(w http.ResponseWriter, status int, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	data = append(data, '\n')
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", kContentJson+kContentCharset)
	w.WriteHeader(status)
	if status != http.StatusNoContent {
		_, err = w.Write(data)
	}
	return err
}

func (s *Server) renderFuncMaps() []template.FuncMap {
	funcs := make([]template.FuncMap, 0)
	funcs = append(funcs, s.DefaultRouteFuncs())
	return funcs
}

func (s *Server) renderError(w http.ResponseWriter, status int, msg string, data interface{}) {
	apiError := ApiError{msg, data}
	if err := s.renderJson(w, status, apiError); err != nil {
		log.Errorf("Server got a json rendering error, %s", err)
		// we fallback to the http.Error instead return a json formatted error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
