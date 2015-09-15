package session

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Unknwon/macaron"
)

// Options represents a struct for specifying configuration options for the session middleware.
type Options struct {
	// Name of provider. Default is "memory".
	Provider string
	// Provider configuration, it's corresponding to provider.
	ProviderConfig string
	// Cookie name to save session ID. Default is "MacaronSession".
	CookieName string
	// Cookie path to store. Default is "/".
	CookiePath string
	// GC interval time in seconds. Default is 3600.
	Gclifetime int64
	// Max life time in seconds. Default is whatever GC interval time is.
	Maxlifetime int64
	// Use HTTPS only. Default is false.
	Secure bool
	// Cookie life time. Default is 0.
	CookieLifeTime int
	// Cookie domain name. Default is empty.
	Domain string
	// Session ID length. Default is 16.
	IDLength int
	// Configuration section name. Default is "session".
	Section string
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Section) == 0 {
		opt.Section = "session"
	}
	sec := macaron.Config().Section(opt.Section)

	if len(opt.Provider) == 0 {
		opt.Provider = sec.Key("PROVIDER").MustString("memory")
	}
	if len(opt.ProviderConfig) == 0 {
		opt.ProviderConfig = sec.Key("PROVIDER_CONFIG").MustString("data/sessions")
	}
	if len(opt.CookieName) == 0 {
		opt.CookieName = sec.Key("COOKIE_NAME").MustString("MacaronSession")
	}
	if len(opt.CookiePath) == 0 {
		opt.CookiePath = sec.Key("COOKIE_PATH").MustString("/")
	}
	if opt.Gclifetime == 0 {
		opt.Gclifetime = sec.Key("GC_INTERVAL_TIME").MustInt64(3600)
	}
	if opt.Maxlifetime == 0 {
		opt.Maxlifetime = sec.Key("MAX_LIFE_TIME").MustInt64(opt.Gclifetime)
	}
	if !opt.Secure {
		opt.Secure = sec.Key("SECURE").MustBool()
	}
	if opt.CookieLifeTime == 0 {
		opt.CookieLifeTime = sec.Key("COOKIE_LIFE_TIME").MustInt()
	}
	if len(opt.Domain) == 0 {
		opt.Domain = sec.Key("DOMAIN").String()
	}
	if opt.IDLength == 0 {
		opt.IDLength = sec.Key("ID_LENGTH").MustInt(16)
	}

	return opt
}
