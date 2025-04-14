package fileserver

import (
	"net/http"
	"strings"
)

const (
	authTypeBasic = "basic"
)

type authenticator struct {
	_type     string
	basicAuth *basicAuth
}

type basicAuth struct {
	username string
	password string
}

func (b *basicAuth) authenticate(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok || username != b.username || password != b.password {
		return false
	}
	return true
}

func (a *authenticator) authenticate(r *http.Request) bool {
	if a._type == authTypeBasic {
		return a.basicAuth.authenticate(r)
	}

	return false
}

func newAuthenticator(authType string, auth string) *authenticator {
	auhtenticator := &authenticator{
		_type: authType,
	}
	if authType == authTypeBasic {
		authSplites := strings.Split(auth, ":")
		if len(authSplites) != 2 {
			panic("invalid basic auth, format => username:password")
		}

		username, password := authSplites[0], authSplites[1]
		auhtenticator.basicAuth = &basicAuth{
			username: username,
			password: password,
		}
	}

	return auhtenticator
}
