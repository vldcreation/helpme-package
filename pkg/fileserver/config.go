package fileserver

type FileServerOpt func(*FileServer)

func WithHost(host string) FileServerOpt {
	return func(c *FileServer) {
		c.host = host
	}
}

func WithPort(port string) FileServerOpt {
	return func(c *FileServer) {
		if port != "" {
			if port[0] != ':' {
				port = ":" + port
			}
			c.port = port
		}
	}
}

// WithAuth sets the authentication string for the file server
// Auth type currently only works for basic auth
// Format: username:password
func WithAuth(auth string) FileServerOpt {
	return func(c *FileServer) {
		c.auth = newAuthenticator(authTypeBasic, auth)
	}
}
