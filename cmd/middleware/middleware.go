package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	sessions "github.com/MikMuellerDev/QuickClip/sessions"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger(logger *logrus.Logger) {
	log = logger
}

type ResponseStruct struct {
	Success   bool
	ErrorCode int
	Title     string
	Message   string
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		value, ok := session.Values["valid"]
		valid, okParse := value.(bool)

		query := r.URL.Query()
		username := query.Get("username")
		password := query.Get("password")

		// TODO impl checker for api requests that use url params auth instead of session
		if ok && okParse && valid {
			log.Trace(fmt.Sprintf("Valid Session, serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		} else if TestCredentials(username, password, false) {
			log.Trace(fmt.Sprintf("Invalid Session, but authenticated with query. serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		}
		log.Trace(fmt.Sprintf("Invalid Session, redirecting %s to /login", r.URL.Path))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func ApiAuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		value, ok := session.Values["valid"]
		valid, okParse := value.(bool)

		query := r.URL.Query()
		username := query.Get("username")
		password := query.Get("password")

		if ok && okParse && valid {
			log.Trace(fmt.Sprintf("Valid Session, serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		} else if TestCredentials(username, password, false) {

			// UNCOMMENT IF SAVING SESSION WHEN USING API IS WANTED
			// session, _ := sessions.Store.Get(r, "session")
			// session.Values["valid"] = true
			// session.Values["username"] = username
			// session.Save(r, w)

			log.Trace(fmt.Sprintf("Invalid Session, but authenticated with query: Session Saved. Serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		}
		log.Trace(fmt.Sprintf("Invalid Session, redirecting %s to /login", r.URL.Path))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Access denied", "You must be authenticated."})
	}
}

func AdminAuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		value, ok := session.Values["valid"]
		sessionUserTemp := session.Values["username"]
		sessionUser, _ := sessionUserTemp.(string)
		valid, okParse := value.(bool)

		query := r.URL.Query()
		username := query.Get("username")
		password := query.Get("password")

		if ok && okParse && valid {
			if sessionUser != "admin" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Access denied", "You must authenticate as a admin user."})
				return
			}
			log.Trace(fmt.Sprintf("Valid Session, serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		} else if TestCredentials(username, password, false) {
			if username != "admin" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Access denied", "You must authenticate as a admin user."})
				return
			}
			log.Trace(fmt.Sprintf("Invalid Session, but authenticated with query: Session Saved. Serving %s", r.URL.Path))
			handler.ServeHTTP(w, r)
			return
		}
		log.Trace(fmt.Sprintf("Invalid Session, redirecting %s to /login", r.URL.Path))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Access denied", "You must be authenticated as a admin user."})
	}
}

func LogRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// UA stands for user-agent
		log.Trace(fmt.Sprintf("[\x1b[32m%s\x1b[0m] FROM: (\x1b[34m%s\x1b[0m) [%s] Serving path:\x1b[35m%s\x1b[0m, UA:%s", r.Method, r.RemoteAddr, r.Proto, r.URL.Path, r.UserAgent()))
		handler.ServeHTTP(w, r)
	}
}
