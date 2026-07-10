package connectiontracker

import (
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

	if r.Method != http.MethodConnect {
		return next.ServeHTTP(w, r)
	}

	username := getUsername(r)

	id := uuid.New().String()

	conn := &UserConnection{
		Host:    r.Host,
		Started: time.Now(),
	}

	// проверяем лимит и добавляем подключение
	ok := storage.AddConnection(
		username,
		id,
		conn,
	)

	if !ok {
		http.Error(
			w,
			"Too many active connections",
			http.StatusTooManyRequests,
		)

		return nil
	}

	// удаляем подключение после закрытия CONNECT
	defer func() {
		storage.RemoveConnection(
			username,
			id,
		)
	}()

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