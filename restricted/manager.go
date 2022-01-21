///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package restricted

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/user-discovery-bot/validation"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// Manager contains the list of restricted usernames and regex in memory and
// handles the checking if registered usernames are restricted. Also handles
// the dynamic updating of the lists in memory when their source file changes.
type Manager struct {
	// File path to the restricted username list
	usernamePath string

	// File path to the restricted regex list
	regexPath string

	// List of restricted usernames that are not allowed to be registered
	usernames map[string]struct{}

	// List of regular expressions that registered usernames cannot match
	regexes []*regexp.Regexp

	mux sync.RWMutex
}

// NewManager initialises the restricted Manager with the contents of the
// username and regex lists. Also starts a thread that dynamically updates the
// lists on file change.
func NewManager(usernamePath, regexPath string, quit chan struct{}) (*Manager, error) {
	// Construct a map of restricted usernames
	usernames, err := usernameListParser(usernamePath)
	if err != nil {
		return nil, err
	}

	// Construct list of restricted regex
	regexes, err := regexListParser(regexPath)
	if err != nil {
		return nil, err
	}

	// Construct the restricted username manager
	m := &Manager{
		usernamePath: usernamePath,
		regexPath:    regexPath,
		usernames:    usernames,
		regexes:      regexes,
	}

	// Start the file watcher
	err = m.fileWatch(quit)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// usernameListParser reads the list of restricted usernames from the given
// filepath and returns a map of them. Usernames must be seperated by new lines.
// Usernames are canonicalized before added to the map.
func usernameListParser(path string) (map[string]struct{}, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Scan the file line by line and add each username to the list
	scanner := bufio.NewScanner(f)
	usernames := make(map[string]struct{})
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			usernames[validation.Canonicalize(line)] = struct{}{}
		}
	}

	if err = scanner.Err(); err != nil {
		_ = f.Close() // Ignore error; scanner error takes precedence
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return usernames, nil
}

// regexListParser reads the list of restricted regular expressions from the
// given filepath and returns a list of them compiled. Regular expressions must
// be seperated by new lines.
func regexListParser(path string) ([]*regexp.Regexp, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Scan the file line by line and add each regex to the list
	scanner := bufio.NewScanner(f)
	var regexes []*regexp.Regexp
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			regex, err := regexp.Compile(line)
			if err != nil {
				_ = f.Close() // Ignore error; regex error takes precedence
				return nil, err
			}
			regexes = append(regexes, regex)
		}
	}

	if err = scanner.Err(); err != nil {
		_ = f.Close() // Ignore error; scanner error takes precedence
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return regexes, nil
}

func (m *Manager) updateLists() error {
	// Get new username list
	usernames, err := usernameListParser(m.usernamePath)
	if err != nil {
		return err
	}

	// Get new regex list
	regexes, err := regexListParser(m.regexPath)
	if err != nil {
		return err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	m.usernames = usernames
	m.regexes = regexes

	return nil
}

func (m *Manager) updateUsernamesFromFile() error {
	// Get new username list
	usernames, err := usernameListParser(m.usernamePath)
	if err != nil {
		return err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	m.usernames = usernames

	return nil
}

func (m *Manager) updateRegexesFromFile() error {
	// Get new regex list
	regexes, err := regexListParser(m.regexPath)
	if err != nil {
		return err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	m.regexes = regexes

	return nil
}

// IsRestricted checks if the username is matches any restricted usernames or
// restricted regular expressions. Usernames must be canonicalized before they
// are passed in.
func (m *Manager) IsRestricted(username string) bool {
	m.mux.RLock()
	defer m.mux.RUnlock()

	_, exists := m.usernames[username]
	if exists {
		return exists
	}

	return m.matchRestrictedRegex(username)
}

// matchRestrictedRegex checks if the username matches any of the restricted
// regular expressions.
func (m *Manager) matchRestrictedRegex(username string) bool {
	for _, regex := range m.regexes {
		if regex.MatchString(username) {
			return true
		}
	}

	return false
}

// fileWatch watches for changes to restricted username list files and
// dynamically updates the list in memory when the files change.
func (m *Manager) fileWatch(quit chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Errorf("failed to initialize watcher for restricted username lists: %+v", err)
	}

	go func() {
		jww.INFO.Print("Starting restricted username file watcher.")
		defer func() {
			err = watcher.Close()
			if err != nil {
				jww.ERROR.Printf("Failed to close restricted username file watcher: %+v", err)
			}
		}()

		for {
			select {
			case <-quit:
				jww.INFO.Print("Quitting restricted username file watcher.")
				return
			case event := <-watcher.Events:
				jww.DEBUG.Printf("Restricted username file watcher: file %q op %s", event.Name, event.Op)
				if event.Op == fsnotify.Write || event.Op == fsnotify.Create {
					if strings.Contains(m.usernamePath, event.Name) {
						jww.INFO.Printf("Updating restricted usernames from file %q.", event.Name)
						err = m.updateUsernamesFromFile()
						if err != nil {
							jww.ERROR.Printf("Failed to update restricted username list: %+v", err)
						}
					} else if strings.Contains(m.regexPath, event.Name) {
						jww.INFO.Printf("Updating restricted regex from file %q.", event.Name)
						err = m.updateRegexesFromFile()
						if err != nil {
							jww.ERROR.Printf("Failed to update restricted regex list: %+v", err)
						}
					}
				} else if event.Op == fsnotify.Remove {
					jww.ERROR.Printf("Restricted username file watcher: %q was deleted", event.Name)
				}
			case err := <-watcher.Errors:
				jww.ERROR.Printf("Restricted username file watcher encountered an error: %+v", err)
			}
		}
	}()

	err = watcher.Add(m.usernamePath)
	if err != nil {
		return errors.Errorf("could not add %q to restricted username file watcher: %+v", m.usernamePath, err)
	}

	err = watcher.Add(m.regexPath)
	if err != nil {
		return errors.Errorf("could not add %q to restricted username file watcher: %+v", m.regexPath, err)
	}
	return nil
}

// NewManagerForTesting creates a new Manager without a file backend to only be
// used for testing.
func NewManagerForTesting(usernames map[string]struct{},
	regexes []*regexp.Regexp, x interface{}) *Manager {
	switch x.(type) {
	case *testing.T, *testing.M, *testing.B, *testing.PB:
		break
	default:
		jww.FATAL.Panicf("NewManagerForTesting can only be used for testing.")
	}

	return &Manager{
		usernamePath: "",
		regexPath:    "",
		usernames:    usernames,
		regexes:      regexes,
	}
}
