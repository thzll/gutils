package socketStream

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
)

const rootPEM = `-----BEGIN CERTIFICATE-----
MIIDxzCCAq+gAwIBAgIJAJOa+w+gbMi0MA0GCSqGSIb3DQEBBQUAMHkxCzAJBgNV
BAYTAkNOMRAwDgYDVQQIDAdqaWFuZ3hpMREwDwYDVQQHDAhuYW5jaGFuZzENMAsG
A1UECgwEbmNidDETMBEGA1UEAwwKZG5zZHVuLmNvbTEhMB8GCSqGSIb3DQEJARYS
a2VlbmdvOTlAZ21haWwuY29tMCAXDTE0MDUxMzA3MTcwNVoYDzIxMTQwNDE5MDcx
NzA1WjB5MQswCQYDVQQGEwJDTjEQMA4GA1UECAwHamlhbmd4aTERMA8GA1UEBwwI
bmFuY2hhbmcxDTALBgNVBAoMBG5jYnQxEzARBgNVBAMMCmRuc2R1bi5jb20xITAf
BgkqhkiG9w0BCQEWEmtlZW5nbzk5QGdtYWlsLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAMhLDFLe3qySRPm59X7L9lItg9gVhqr3NoabwRYvB8tE
Wqes+vY1TC9ad0Lnju77ZjQb7exAV09sxckETKuQ3uuIRsy6muUrAN8OV7XTUDUT
wT+vY3UvwT3KmEHoCw7riVIpxt/pFwIvyu7AkYGmVgSbjXxLMnQ27tuG0AqS0EVF
xmu+Vj4EimY9vi4i7aO4gDzufwYgEtzTtSJCduvU7ii/erh9T+40QR81fwOjVsFk
SVMcjzof27C0b7Ievr1B4QVkyOlq3DQukHoQCM38i5reU2WdXw/WUONCVL7Mb+ZQ
GbTD3+i2BX160s9ddhLgsT067n93rl+HI7nUrCZ7QAECAwEAAaNQME4wHQYDVR0O
BBYEFAaDQp9sbXlJz+Kwdc4xIeTWO9uXMB8GA1UdIwQYMBaAFAaDQp9sbXlJz+Kw
dc4xIeTWO9uXMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAKbJ2/SJ
zlkWvWzBMQgEvDsUM8QeehnzvK+5NITvXYgJBUcZ0uv0QUwvCrueXaCH0TD2CiWQ
ugbr25/XxsD3bQg4+sNXyKmrVB9cqlP9wxRhvIJ5DARJWqsGEGvAhzALJNEGa9Fd
It5v4wHSOxNDQI0qL64w5Mwl9fN/x+5wLSirkdsisNcE+E6vz52K3CLj6VcJy+fE
+bdaYsP0yN8uZFxPGg+k+bONrVAffDYjOp51vfi6HEHtyIPsvzIxcBrHQI9jqsif
0McsBJRWIh5obvHfB5b+BElEk8qVLHoeRtFYndJVhuueINTkQc7osOajLDxjffoX
0af/ywcJOMem4a8=
-----END CERTIFICATE-----`

func createTlsConn(conn net.Conn) (net.Conn, error) {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		return nil, fmt.Errorf("failed to parse root certificate")
	}
	tlsconn := tls.Client(conn, &tls.Config{RootCAs: roots, InsecureSkipVerify: true})
	if tlsconn == nil {
		return nil, fmt.Errorf("cann't call tls client")
	}
	return tlsconn, nil
}
