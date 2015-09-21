package session

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Unknwon/macaron"
)

// RawStore is the interface that operates the session data.
type RawStore interface {
	// Set sets value to given key in session.
	Set(interface{}, interface{}) error
	// Get gets value by given key in session.
	Get(interface{}) interface{}
	// Delete deletes a key from session.
	Delete(interface{}) error
	// ID returns current session ID.
	ID() string
	// Release releases session resource and save data to provider.
	Release() error
	// Flush deletes all session data.
	Flush() error
}

// Store is the interface that contains all data for one session process with specific ID.
type Store interface {
	RawStore
	// Read returns raw session store by session ID.
	Read(string) (RawStore, error)
	// Destory deletes a session.
	Destory(*macaron.Context) error
	// RegenerateId regenerates a session store from old session ID to new one.
	RegenerateId(*macaron.Context) (RawStore, error)
	// Count counts and returns number of sessions.
	Count() int
	// GC calls GC to clean expired sessions.
	GC()
}

type store struct {
	RawStore
	*Manager
}

var _ Store = &store{}
