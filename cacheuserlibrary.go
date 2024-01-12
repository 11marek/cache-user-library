package cacheuserlibrary

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// DatabaseConfig stores the configuration for the database.
type DatabaseConfig struct {
	// Add fields for the database configuration, e.g., address, port, password, etc.
}

// SafeMap is a map that is safe for concurrent use.
type SafeMap struct {
	m  map[string]time.Time
	mu sync.RWMutex
}

// NewSafeMap initializes a new SafeMap.
func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[string]time.Time)}
}

// Set sets the value for the given key.
func (sm *SafeMap) Set(key string, value time.Time) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

// Get gets the value for the given key.
func (sm *SafeMap) Get(key string) (time.Time, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.m[key]
	return value, ok
}

// UserCacheHandler handles the caching logic.
type UserCacheHandler struct {
	uniqueIDsMap   *SafeMap
	uniqueIDsLimit int
	DatabaseConfig *DatabaseConfig
}

// NewUserCacheHandler initializes a new caching handler.
func NewUserCacheHandler(uniqueIDsLimit int, dbConfig *DatabaseConfig) *UserCacheHandler {
	return &UserCacheHandler{
		uniqueIDsMap:   NewSafeMap(),
		uniqueIDsLimit: uniqueIDsLimit,
		DatabaseConfig: dbConfig,
	}
}

// IsUserCached checks if the user is cached.
func (h *UserCacheHandler) IsUserCached(userID string) bool {
	_, exists := h.uniqueIDsMap.Get(userID)
	return exists
}

// HandleCaching is the core caching logic.
func (h *UserCacheHandler) HandleCaching(userID string) bool {
	h.uniqueIDsMap.mu.Lock()
	defer h.uniqueIDsMap.mu.Unlock()

	// Check if the userID is already in the group
	if _, exists := h.uniqueIDsMap.Get(userID); !exists {
		// Add a unique userID to the map with the current time
		h.uniqueIDsMap.Set(userID, time.Now())

		// Check if the limit of unique IDs has been reached
		if len(h.uniqueIDsMap.m) >= h.uniqueIDsLimit {
			// Clear the map after handling the group
			h.uniqueIDsMap = NewSafeMap()
		}

		// Run a goroutine to clean up IDs after a specified time
		go h.CleanupUserAfterTimeout(userID)

		return true
	}

	return false
}

// CleanupUserAfterTimeout is a goroutine to clean up IDs after a specified time.
func (h *UserCacheHandler) CleanupUserAfterTimeout(userID string) {
	select {
	case <-time.After(time.Minute):
		h.uniqueIDsMap.mu.Lock()
		defer h.uniqueIDsMap.mu.Unlock()

		if creationTime, exists := h.uniqueIDsMap.Get(userID); exists && time.Since(creationTime) >= time.Minute {
			h.uniqueIDsMap.Set(userID, time.Now())
			fmt.Printf("User %s removed after timeout.\n", userID)
		}
	}
}

// SaveCacheToFile serializes the map to JSON and saves it to a file.
func (h *UserCacheHandler) SaveCacheToFile(filePath string) error {
	h.uniqueIDsMap.mu.RLock()
	defer h.uniqueIDsMap.mu.RUnlock()

	// Serialize the map to JSON format
	data, err := json.Marshal(h.uniqueIDsMap.m)
	if err != nil {
		return err
	}

	// Save the data to the file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
