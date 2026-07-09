package connectiontracker

import (
	"log"
	"encoding/base64"
	"strings"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/google/uuid"
)

type Handler struct{}

func init() {
	caddy.RegisterModule(Handler{})
}

func (Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.connection_tracker",
		New: func() caddy.Module {
			return new(Handler)
		},
	}
}

func (h Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	next caddyhttp.Handler,
) error {

	if r.Method == http.MethodConnect {
		log.Println("=== REQUEST ===")
		log.Println("RemoteAddr:", r.RemoteAddr)
		log.Println("Host:", r.Host)
		log.Println("Method:", r.Method)
		log.Println("Proto:", r.Proto)
		log.Println("RequestURI:", r.RequestURI)
		log.Println("URL:", r.URL.String())
		log.Println("Forwarded:", r.Header.Get("X-Forwarded-For"))
		log.Printf("Headers: %+v", r.Header)
		log.Printf("Context: %+v", r.Context())

		for k, v := range r.Header {
			log.Printf("Header %s: %v\n", k, v)
		}

		if r.TLS != nil {
			log.Println("TLS ServerName:", r.TLS.ServerName)
			log.Println("TLS Version:", r.TLS.Version)
		}

		ip := r.RemoteAddr
		repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
		log.Println("Caddy remote:", repl.ReplaceAll("{http.request.remote}", ""))
		log.Println("Caddy host:", repl.ReplaceAll("{http.request.remote.host}", ""))
		host := r.Host

		id := uuid.New().String()

		conn := &UserConnection{
			IP:      ip,
			Host:    host,
			Started: time.Now(),
		}

		username := getUsername(r)

		storage.AddConnection(
			username,
			id,
			conn,
		)

		defer func() {
			storage.RemoveConnection(
				username,
				id,
			)
		}()
	}

	return next.ServeHTTP(w, r)
}

func getUsername(r *http.Request) string {

	auth := r.Header.Get("Proxy-Authorization")

	if auth == "" {
		return "unknown"
	}

	encoded, ok := strings.CutPrefix(auth, "Basic ")

	if !ok {
		return "unknown"
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return "unknown"
	}

	username, _, ok := strings.Cut(string(decoded), ":")

	if !ok {
		return "unknown"
	}

	return username
}

var (
	_ caddy.Module                = (*Handler)(nil)
	_ caddyhttp.MiddlewareHandler = (*Handler)(nil)
)