package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	sessions "github.com/MikMuellerDev/QuickClip/sessions"
	"github.com/MikMuellerDev/QuickClip/utils"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Maximum size of a request body, larger requests are rejected
const MaxBodyBytes = 2 << 20

func InitLogger(logger *logrus.Logger) {
	log = logger
}

type ResponseStruct struct {
	Success   bool
	ErrorCode int
	Title     string
	Message   string
}

func writeJsonError(w http.ResponseWriter, status int, title string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ResponseStruct{false, status, title, message})
}

// Returns the authenticated user of the request.
// Browsers authenticate using the session cookie, API clients can use HTTP Basic Auth.
func CurrentUser(r *http.Request) (string, bool) {
	if username, ok := sessions.User(r); ok && utils.DoesUserExist(username) {
		return username, true
	}
	if username, password, ok := r.BasicAuth(); ok {
		if success, _ := TestCredentials(r, username, password); success {
			return username, true
		}
	}
	return "", false
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := CurrentUser(r); ok {
			log.Trace(fmt.Sprintf("Authenticated, serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		}
		log.Trace(fmt.Sprintf("Not authenticated, redirecting %s to /login", r.URL.Path))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func ApiAuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := CurrentUser(r); ok {
			log.Trace(fmt.Sprintf("Authenticated, serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		}
		log.Trace(fmt.Sprintf("Not authenticated, denying %s", r.URL.Path))
		writeJsonError(w, http.StatusUnauthorized, "Access denied", "You must be authenticated.")
	}
}

func AdminAuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := CurrentUser(r)
		if !ok {
			writeJsonError(w, http.StatusUnauthorized, "Access denied", "You must be authenticated as a admin user.")
			return
		}
		if username != "admin" {
			writeJsonError(w, http.StatusForbidden, "Access denied", "You must authenticate as a admin user.")
			return
		}
		log.Trace(fmt.Sprintf("Authenticated as admin, serving %s", r.URL.Path))
		handler.ServeHTTP(w, r)
	}
}

// Rejects state-changing requests sent by browsers from other origins (CSRF protection)
func CheckOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		origin := r.Header.Get("Origin")
		if origin != "" {
			parsed, err := url.Parse(origin)
			if err != nil || origin == "null" || parsed.Host != r.Host {
				log.Warn(fmt.Sprintf("Blocked cross-origin %s request to %s from origin %q", r.Method, r.URL.Path, origin))
				writeJsonError(w, http.StatusForbidden, "Access denied", "Cross-origin requests are not allowed.")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func LimitBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		next.ServeHTTP(w, r)
	})
}

func LogRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// UA stands for user-agent
		log.Trace(fmt.Sprintf("[\x1b[32m%s\x1b[0m] FROM: (\x1b[34m%s\x1b[0m) [%s] Serving path:\x1b[35m%s\x1b[0m, UA:%q", r.Method, r.RemoteAddr, r.Proto, r.URL.Path, r.UserAgent()))
		handler.ServeHTTP(w, r)
	}
}
