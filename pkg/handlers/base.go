package handlers

import (
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

type baseHandler struct {
	lg *logrus.Entry
}

func newBaseHandler(lg *logrus.Logger, name string) baseHandler {
	return baseHandler{lg.WithField("handler", name)}
}

func (h *baseHandler) recover(lg *logrus.Entry) {
	if r := recover(); r != nil {
		lg.WithFields(logrus.Fields{
			"stack_trace": debug.Stack(),
			"error":       r,
		}).Error("panic was recovered")
	}
}

func (h *baseHandler) newLoggerFromRequest(r *http.Request) *logrus.Entry {
	lg := h.lg.WithField("request_id", xid.New())

	lg.WithFields(logrus.Fields{
		"ua":        r.UserAgent(),
		"method":    r.Method,
		"path":      r.URL.Path,
		"query":     r.URL.Query(),
		"host":      r.Host,
		"proto":     r.Proto,
		"origin":    r.Header.Get("origin"),
		"forwarded": r.Header.Get("X-FORWARDED-FOR"),
		"addr":      r.RemoteAddr,
	}).Trace("incoming request")

	return lg
}
