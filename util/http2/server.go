/*
 * HTTP2 Server
 */

package http2

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func NewServer(bindAddr string, tlskeylog string, handler http.Handler) (server *http.Server, err error) {
	keylogFile, err := os.OpenFile(tlskeylog, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, errors.New("server need handler")
	}
	server = &http.Server{
		Addr: bindAddr,
		TLSConfig: &tls.Config{
			KeyLogWriter: keylogFile,
		},
		Handler: handler,
	}
	return
}
