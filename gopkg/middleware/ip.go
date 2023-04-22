package middleware

import (
	"net"
	"net/http"
)

//获取客户端IP
func GetRemoteIp(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(req.RemoteAddr)
	}

	return ip
}
