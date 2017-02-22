package proxy

import "errors"
import "strings"

var (
	ErrUnknownProxy      = errors.New("unknown proxy")
	ErrApacheProxy       = errors.New("apache? O RLY?")
	ErrInvalidProxyParam = errors.New("invalid proxy parameter")
)

type ProxyType int

const (
	UndefinedProxy ProxyType = iota
	Caddy
	Nginx
)

type ProxyInit struct {
	Proxy        ProxyType
	DefaultProxy bool
	ConfigFile   string
}

func ParseProxyArgument(proxyArg string) (*ProxyInit, error) {
	var proxyStr string
	var proxyConfig string
	// no proxy
	if proxyArg == "" || proxyArg == ":" {
		return nil, nil
	}
	splittedArg := strings.SplitN(proxyArg, ":", 2)
	switch len(splittedArg) {
	case 1:
		proxyStr = splittedArg[0]
	case 2:
		proxyStr = splittedArg[0]
		proxyConfig = splittedArg[1]
	default:
		return nil, ErrInvalidProxyParam
	}
	var proxy ProxyType
	switch proxyStr {
	case "nginx":
		proxy = Nginx
	case "caddy":
		proxy = Caddy
	case "apache":
		// >:(
		return nil, ErrApacheProxy
	default:
		return nil, ErrUnknownProxy
	}
	p := &ProxyInit{
		Proxy:      proxy,
		ConfigFile: proxyConfig,
	}
	if p.ConfigFile == "" {
		p.DefaultProxy = true
	}
	return p, nil
}
