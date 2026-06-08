package internal

import (
	"crypto/tls"
	"fmt"
	"time"
)

func CheckSSLExpiry(host string, port int) (time.Time, error) {
	conn, err := tls.DialWithDialer(&tls.Dialer{Timeout: 5 * time.Second}, "tcp", fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		return time.Time{}, err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return time.Time{}, fmt.Errorf("no certificates found")
	}

	return certs[0].NotAfter, nil
}
