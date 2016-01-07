package pecelserver

import (
	"github.com/eaciit/toolkit"
	"time"
)

type Session struct {
	SessionID   string
	ReferenceID string
	Secret      string
	Created     time.Time
	ExpireOn    time.Time
}

var _defaultSessionLifetime time.Duration

func SetSesionLifetime(t time.Duration) {
	_defaultSessionLifetime = t
}

func SessionLifetime() time.Duration {
	if _defaultSessionLifetime == 0 {
		_defaultSessionLifetime = 90 * time.Minute
	}
	return _defaultSessionLifetime
}

func NewSession(referenceid string) *Session {
	s := new(Session)
	s.SessionID = toolkit.RandomString(32)
	s.ReferenceID = referenceid
	s.Created = time.Now()
	s.ExpireOn = s.Created.Add(SessionLifetime())
	s.Secret = toolkit.RandomString(32)
	return s
}

func (s *Session) IsValid() bool {
	return time.Now().Before(s.ExpireOn)
}
