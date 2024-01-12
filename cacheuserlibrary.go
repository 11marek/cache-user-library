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
	// Add fields for the database configuration, e.g., address, port, password, etc .
}

// UserCacheHandler handles the caching logic.
type UserCacheHandler struct {
	uniqueIDsMap   map[string]time.Time
	uniqueIDsLimit int
	mu             sync.RWMutex
	DatabaseConfig *DatabaseConfig
}

// NewUserCacheHandler initializes a new caching handler.
func NewUserCacheHandler(uniqueIDsLimit int, dbConfig *DatabaseConfig) *UserCacheHandler {
	return &UserCacheHandler{
		uniqueIDsMap:   make(map[string]time.Time),
		uniqueIDsLimit: uniqueIDsLimit,
		DatabaseConfig: dbConfig,
	}
}

// IsUserCached checks if the user is cached.
func (h *UserCacheHandler) IsUserCached(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, exists := h.uniqueIDsMap[userID]
	return exists
}

// HandleCaching is the core caching logic.
func (h *UserCacheHandler) HandleCaching(userID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if the userID is already in the group
	if _, exists := h.uniqueIDsMap[userID]; !exists {
		// Add a unique userID to the map with the current time
		h.uniqueIDsMap[userID] = time.Now()

		// Check if the limit of unique IDs has been reached
		if len(h.uniqueIDsMap) >= h.uniqueIDsLimit {
			// Clear the map after handling the group
			h.uniqueIDsMap = make(map[string]time.Time)
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
		h.mu.Lock()
		defer h.mu.Unlock()

		if creationTime, exists := h.uniqueIDsMap[userID]; exists && time.Since(creationTime) >= time.Minute {
			delete(h.uniqueIDsMap, userID)
			fmt.Printf("User %s removed after timeout.\n", userID)
		}
	}
}

// SaveCacheToFile serializes the map to JSON and saves it to a file.
func (h *UserCacheHandler) SaveCacheToFile(filePath string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Serialize the map to JSON format
	data, err := json.Marshal(h.uniqueIDsMap)
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
