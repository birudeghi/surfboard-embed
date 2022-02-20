package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifetime int64
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error //given an sid, deleted the corresponding session
	SessionGC(maxLifetime int64)     // deletes expired session variables according to maxLifetime
}

type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}

// creating sessionId
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// Check the existence of any sessions related to the current user based on cookies, and creating a new session if none is found.
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {

	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)

	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifetime)}
		http.SetCookie(w, &cookie)

	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}

	return
}

// destroy session
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)

	if cookie.Value == "" || err != nil {
		return

	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		sid, _ := url.QueryUnescape(cookie.Value)
		manager.provider.SessionDestroy(sid)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

// session garbage collection
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifetime)
	time.AfterFunc(time.Duration(manager.maxLifetime), func() { manager.GC() })
}

var provides = make(map[string]Provider)

// Register makes a session provider available by the provided name.
// Panic if Register is called twice with the same name or driver is nil.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}

	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

func NewManager(provideName, cookieName string, maxLifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, errors.New("session: unknown provide " + provideName + "(forgotten import?)")
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifetime: maxLifetime}, nil
}
