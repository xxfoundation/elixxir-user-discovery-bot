////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	"bytes"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/primitives/utils"
	"strings"
	"sync"
)

type BlackList struct {
	list         map[string]bool // Contains the list of keys
	file         string          // Absolute URL for the blacklist file
	sync.RWMutex                 // Only allows one writer at a time
}

// Initialises a map for the blacklist with the specified keys. The
// updateFinished channel receives a value when the map finises updating.
func InitBlackList(filePath string) *BlackList {
	bl := BlackList{
		list: make(map[string]bool),
		file: filePath,
	}
	// Update the blacklist from the specified file
	bl.UpdateBlacklist()

	return &bl
}

// Initialises a map for the blacklist with the specified keys.
func (bl *BlackList) UpdateBlacklist() {
	// Get list of strings from the blacklist file
	list, err := BlacklistFileParse(bl.file)
	// If the file was read successfully, update the list
	if err == nil {
		// Enable write lock while writing to the map
		bl.Lock()

		// Reset the blacklist to empty
		bl.list = make(map[string]bool)

		// Add all the keys to the map
		for _, key := range list {
			bl.list[key] = true
		}

		// Disable write lock when writing is done
		bl.Unlock()
	}
}

// Initialises a map for the blacklist with the specified keys.
func (bl *BlackList) UpdateBlackList() {
	// Get list of strings from the blacklist file
	list, err := BlacklistFileParse(bl.file)

	// If the file was read successfully, update the list
	if err == nil {
		// Enable write lock while writing to the map
		bl.Lock()

		// Reset the blacklist to empty
		bl.list = make(map[string]bool)

		// Add all the keys to the map
		for _, key := range list {
			bl.list[key] = true
		}

		// Disable write lock when writing is done
		bl.Unlock()
	}
}

// Parses the given file and stores each value in a slice. Returns the slice and
// an error. The file is expected to have value separated by new lines.
func BlacklistFileParse(filePath string) ([]string, error) {
	// Load file contents into memory
	data, err := utils.ReadFile(filePath)
	if err != nil {
		globals.Log.ERROR.Printf("Failed to read file: %v", err)
		return []string{}, err
	}

	// Convert the data to string, triim whitespace, and normalize new lines
	dataStr := strings.TrimSpace(string(normalizeNewlines(data)))

	// Return empty slice if the file is empty or only contains whitespace
	if dataStr == "" {
		return []string{}, nil
	}

	// Split the data at new lines and place in slice
	return strings.Split(dataStr, "\n"), nil
}

// Searches if the specified key exists in the blacklist. Returns true if it
// exists and false otherwise.
func (wl *BlackList) Exists(key string) bool {
	// Enable read lock while reading from the map
	wl.RLock()

	// Check if the key exists in the map
	_, ok := wl.list[key]
	// Disable read lock when reading is done
	wl.RUnlock()

	return ok
}

// Normalizes \r\n (Windows) and \r (Mac) into \n (UNIX).
func normalizeNewlines(d []byte) []byte {
	// Replace CR LF \r\n (Windows) with LF \n (UNIX)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)

	// Replace CF \r (Mac) with LF \n (UNIX)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)

	return d
}
