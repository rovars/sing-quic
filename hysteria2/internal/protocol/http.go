package protocol

import (
	"strconv"

	"github.com/metacubex/http"
)

const (
	URLPath = "/auth"

	StatusAuthOK = 233
)

type prefixSet struct {
	Auth       string
	UDPEnabled string
	CCRX       string
	Padding    string
}

var (
	prefixHysteria = prefixSet{
		Auth:       "Hysteria-Auth",
		UDPEnabled: "Hysteria-UDP",
		CCRX:       "Hysteria-CC-RX",
		Padding:    "Hysteria-Padding",
	}
	prefixZivpn = prefixSet{
		Auth:       "Zivpnudp-Auth",
		UDPEnabled: "Zivpnudp-UDP",
		CCRX:       "Zivpnudp-CC-RX",
		Padding:    "Zivpnudp-Padding",
	}
)

func GetURLHost(mode string) string {
	if mode == "zudp" {
		return "zivpnudp"
	}
	return "hysteria"
}

func IsValidHost(host string) bool {
	return host == "hysteria" || host == "zivpnudp"
}

func getPrefix(mode string) prefixSet {
	if mode == "zudp" {
		return prefixZivpn
	}
	return prefixHysteria
}

func DetectRequestMode(h http.Header) string {
	if h.Get("Zivpnudp-Auth") != "" || h.Get("Zivpnudp-CC-RX") != "" {
		return "zudp"
	}
	return "hudp"
}

// AuthRequest is what client sends to server for authentication.
type AuthRequest struct {
	Auth string
	Rx   uint64 // 0 = unknown, client asks server to use bandwidth detection
}

// AuthResponse is what server sends to client when authentication is passed.
type AuthResponse struct {
	UDPEnabled bool
	Rx         uint64 // 0 = unlimited
	RxAuto     bool   // true = server asks client to use bandwidth detection
}

func AuthRequestFromHeader(h http.Header) AuthRequest {
	p := getPrefix(DetectRequestMode(h))
	rx, _ := strconv.ParseUint(h.Get(p.CCRX), 10, 64)
	return AuthRequest{
		Auth: h.Get(p.Auth),
		Rx:   rx,
	}
}

func AuthRequestToHeader(h http.Header, req AuthRequest, mode string) {
	p := getPrefix(mode)
	h.Set(p.Auth, req.Auth)
	h.Set(p.CCRX, strconv.FormatUint(req.Rx, 10))
	h.Set(p.Padding, authRequestPadding.String())
}

func AuthResponseFromHeader(h http.Header) AuthResponse {
	p := getPrefix(DetectRequestMode(h))
	resp := AuthResponse{}
	resp.UDPEnabled, _ = strconv.ParseBool(h.Get(p.UDPEnabled))
	rxStr := h.Get(p.CCRX)
	if rxStr == "auto" {
		// Special case for server requesting client to use bandwidth detection
		resp.RxAuto = true
	} else {
		resp.Rx, _ = strconv.ParseUint(rxStr, 10, 64)
	}
	return resp
}

func AuthResponseToHeader(h http.Header, resp AuthResponse, mode string) {
	p := getPrefix(mode)
	h.Set(p.UDPEnabled, strconv.FormatBool(resp.UDPEnabled))
	if resp.RxAuto {
		h.Set(p.CCRX, "auto")
	} else {
		h.Set(p.CCRX, strconv.FormatUint(resp.Rx, 10))
	}
	h.Set(p.Padding, authResponsePadding.String())
}
