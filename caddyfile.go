package connectiontracker

import (
    "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
    "github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
    "github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
    httpcaddyfile.RegisterHandlerDirective("connection_tracker", parseCaddyfile)
    httpcaddyfile.RegisterHandlerDirective("connection_tracker_api", parseConnectionTrackerAPI)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
    var m Handler

    for h.Next() {
    }

    return &m, nil
}

func (h *Handler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
    for d.Next() {
    }
    return nil
}

func parseConnectionTrackerAPI(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	return new(API), nil
}
