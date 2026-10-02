package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var Store *sessions.CookieStore

const cookieName = "session"
const sessionLifetime = 30 * 24 * time.Hour

type activeSession struct {
	username string
	expires  time.Time
}

// The cookie only holds a random session id, the session itself lives on the server.
// This allows logging out, and revoking sessions on password changes or user deletion.
var activeSessions = make(map[string]activeSession)
var sessionsMutex sync.Mutex

func InitLogger(logger *logrus.Logger) {
	log = logger
}

func Init() {
	hashKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)
	if hashKey == nil || encryptionKey == nil {
		log.Fatal("Could not generate session keys.")
	}
	Store = sessions.NewCookieStore(hashKey, encryptionKey)
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(sessionLifetime.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func newSessionId() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Starts a new session for the user, replacing the session the request might already have
func Login(w http.ResponseWriter, r *http.Request, username string) error {
	id, err := newSessionId()
	if err != nil {
		return err
	}
	session, _ := Store.Get(r, cookieName)
	oldId, _ := session.Values["id"].(string)

	sessionsMutex.Lock()
	now := time.Now()
	for key, value := range activeSessions {
		if now.After(value.expires) {
			delete(activeSessions, key)
		}
	}
	delete(activeSessions, oldId)
	activeSessions[id] = activeSession{username: username, expires: now.Add(sessionLifetime)}
	sessionsMutex.Unlock()

	session.Values = map[interface{}]interface{}{"id": id}
	return session.Save(r, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, cookieName)
	if id, ok := session.Values["id"].(string); ok {
		sessionsMutex.Lock()
		delete(activeSessions, id)
		sessionsMutex.Unlock()
	}
	session.Values = map[interface{}]interface{}{}
	session.Options.MaxAge = -1
	session.Save(r, w)
}

// Returns the user of the request's session, if the session is valid
func User(r *http.Request) (string, bool) {
	session, _ := Store.Get(r, cookieName)
	id, ok := session.Values["id"].(string)
	if !ok {
		return "", false
	}
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	active, ok := activeSessions[id]
	if !ok {
		return "", false
	}
	if time.Now().After(active.expires) {
		delete(activeSessions, id)
		return "", false
	}
	return active.username, true
}

// Invalidates every session of the user
func RevokeUser(username string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	for key, value := range activeSessions {
		if value.username == username {
			delete(activeSessions, key)
		}
	}
}
