package connectiontracker

import (
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type API struct{}

func init() {
	caddy.RegisterModule(API{})
}

func (API) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.connection_tracker_api",
		New: func() caddy.Module {
			return new(API)
		},
	}
}

func (a API) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	next caddyhttp.Handler,
) error {

	w.Header().Set("Content-Type", "application/json")

	users := storage.ListUsers()

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

var (
	_ caddy.Module                = (*API)(nil)
	_ caddyhttp.MiddlewareHandler = (*API)(nil)
)
