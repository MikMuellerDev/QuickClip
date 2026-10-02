package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/MikMuellerDev/QuickClip/utils"
)

const maxFailedLogins = 10
const failedLoginWindow = 15 * time.Minute

type failedLogins struct {
	count int
	first time.Time
}

var failedLoginsByIp = make(map[string]*failedLogins)
var failedLoginsMutex sync.Mutex

func clientIp(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func isRateLimited(ip string) bool {
	failedLoginsMutex.Lock()
	defer failedLoginsMutex.Unlock()
	entry, ok := failedLoginsByIp[ip]
	if !ok {
		return false
	}
	if time.Since(entry.first) > failedLoginWindow {
		delete(failedLoginsByIp, ip)
		return false
	}
	return entry.count >= maxFailedLogins
}

func recordFailedLogin(ip string) {
	failedLoginsMutex.Lock()
	defer failedLoginsMutex.Unlock()
	now := time.Now()
	for key, value := range failedLoginsByIp {
		if now.Sub(value.first) > failedLoginWindow {
			delete(failedLoginsByIp, key)
		}
	}
	entry, ok := failedLoginsByIp[ip]
	if !ok {
		entry = &failedLogins{first: now}
		failedLoginsByIp[ip] = entry
	}
	entry.count++
}

// Checks the credentials, limiting failed attempts per client IP.
// limited is true if the client has too many failed attempts, the credentials are not checked in that case.
func TestCredentials(r *http.Request, user string, password string) (success bool, limited bool) {
	ip := clientIp(r)
	if isRateLimited(ip) {
		log.Warn(fmt.Sprintf("Login for User %q from %s blocked: too many failed attempts", user, ip))
		return false, true
	}
	if utils.CheckCredentials(user, password) {
		log.Debug(fmt.Sprintf("Login successful for User %q", user))
		return true, false
	}
	recordFailedLogin(ip)
	log.Warn(fmt.Sprintf("Login failed for User %q from %s", user, ip))
	return false, false
}
