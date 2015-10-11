// Copyright 2013 Beego Authors
// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package session a middleware that provides the session management of Macaron.
package session

import (
	"encoding/hex"
	"fmt"
	"github.com/kolonse/KolonseWeb/HttpLib"
	"github.com/kolonse/KolonseWeb/Type"
	"net/http"
	"time"
)

const _VERSION = "0.1.7"

func Version() string {
	return _VERSION
}

// Sessioner is a middleware that maps a session.SessionStore service into the Macaron handler chain.
// An single variadic session.Options struct can be optionally provided to configure.
func Sessioner(options ...Options) Type.DoStep {
	opt := prepareOptions(options)
	manager, err := NewManager(opt.Provider, opt)
	if err != nil {
		panic(err)
	}
	go manager.startGC()

	return func(req *HttpLib.Request, res *HttpLib.Response, next Type.Next) {
		sess, err := manager.Start(req, res)
		if err != nil {
			panic("session(start): " + err.Error())
		}
		s := store{
			RawStore: sess,
			Manager:  manager,
		}

		req.MapTo(s, (*Store)(nil))
		if err = sess.Release(); err != nil {
			panic("session(release): " + err.Error())
		}
		next()
	}
}

//    _____
//   /     \ _____    ____ _____     ____   ___________
//  /  \ /  \\__  \  /    \\__  \   / ___\_/ __ \_  __ \
// /    Y    \/ __ \|   |  \/ __ \_/ /_/  >  ___/|  | \/
// \____|__  (____  /___|  (____  /\___  / \___  >__|
//         \/     \/     \/     \//_____/      \/

// Manager represents a struct that contains session provider and its configuration.
type Manager struct {
	provider Provider
	opt      Options
}

// NewManager creates and returns a new session manager by given provider name and configuration.
// It panics when given provider isn't registered.
func NewManager(name string, opt Options) (*Manager, error) {
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider '%s'(forgotten import?)", name)
	}
	return &Manager{p, opt}, p.Init(opt.Maxlifetime, opt.ProviderConfig)
}

// sessionId generates a new session ID with rand string, unix nano time, remote addr by hash function.
func (m *Manager) sessionId() string {
	return hex.EncodeToString(generateRandomKey(m.opt.IDLength / 2))
}

// Start starts a session by generating new one
// or retrieve existence one by reading session ID from HTTP request if it's valid.
func (m *Manager) Start(req *HttpLib.Request, res *HttpLib.Response) (RawStore, error) {
	sid := req.GetCookie(m.opt.CookieName)
	if len(sid) > 0 && m.provider.Exist(sid) {
		return m.provider.Read(sid)
	}

	sid = m.sessionId()
	sess, err := m.provider.Read(sid)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     m.opt.CookieName,
		Value:    sid,
		Path:     m.opt.CookiePath,
		HttpOnly: true,
		Secure:   m.opt.Secure,
		Domain:   m.opt.Domain,
	}
	if m.opt.CookieLifeTime >= 0 {
		cookie.MaxAge = m.opt.CookieLifeTime
	}
	http.SetCookie(res, cookie)
	req.AddCookie(cookie)
	return sess, nil
}

// Read returns raw session store by session ID.
func (m *Manager) Read(sid string) (RawStore, error) {
	return m.provider.Read(sid)
}

// Destory deletes a session by given ID.
func (m *Manager) Destory(req *HttpLib.Request, res *HttpLib.Response) error {
	sid := req.GetCookie(m.opt.CookieName)
	if len(sid) == 0 {
		return nil
	}

	if err := m.provider.Destory(sid); err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     m.opt.CookieName,
		Path:     m.opt.CookiePath,
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
	}
	http.SetCookie(res, cookie)
	return nil
}

// RegenerateId regenerates a session store from old session ID to new one.
func (m *Manager) RegenerateId(req *HttpLib.Request, res *HttpLib.Response) (sess RawStore, err error) {
	sid := m.sessionId()
	oldsid := req.GetCookie(m.opt.CookieName)
	sess, err = m.provider.Regenerate(oldsid, sid)
	if err != nil {
		return nil, err
	}
	ck := &http.Cookie{
		Name:     m.opt.CookieName,
		Value:    sid,
		Path:     m.opt.CookiePath,
		HttpOnly: true,
		Secure:   m.opt.Secure,
		Domain:   m.opt.Domain,
	}
	if m.opt.CookieLifeTime >= 0 {
		ck.MaxAge = m.opt.CookieLifeTime
	}
	http.SetCookie(res, ck)
	req.AddCookie(ck)
	return sess, nil
}

// Count counts and returns number of sessions.
func (m *Manager) Count() int {
	return m.provider.Count()
}

// GC starts GC job in a certain period.
func (m *Manager) GC() {
	m.provider.GC()
}

// startGC starts GC job in a certain period.
func (m *Manager) startGC() {
	m.GC()
	time.AfterFunc(time.Duration(m.opt.Gclifetime)*time.Second, func() { m.startGC() })
}

// SetSecure indicates whether to set cookie with HTTPS or not.
func (m *Manager) SetSecure(secure bool) {
	m.opt.Secure = secure
}
