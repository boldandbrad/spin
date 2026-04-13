package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zalando/go-keyring"
)

const keyringService = "spin"

type Profile struct {
	Username   string `json:"username"`
	SessionKey string `json:"-"`
}

type ProfileStore struct {
	Profiles []Profile `json:"profiles"`
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home, _ = os.Getwd()
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "spin")
	case "linux":
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg != "" {
			return filepath.Join(xdg, "spin")
		}
		return filepath.Join(home, ".local", "share", "spin")
	default:
		return filepath.Join(home, ".spin")
	}
}

func ensureAppDataDir() error {
	return os.MkdirAll(appDataDir(), 0700)
}

func profilesFile() string {
	return filepath.Join(appDataDir(), "profiles.json")
}

func activeProfileFile() string {
	return filepath.Join(appDataDir(), "active_profile")
}

func setCredential(username, sessionKey string) error {
	return keyring.Set(keyringService, username, sessionKey)
}

func getCredential(username string) (*Profile, error) {
	item, err := keyring.Get(keyringService, username)
	if err != nil {
		return nil, fmt.Errorf("credential not found for user %s: %w", username, err)
	}
	return &Profile{Username: username, SessionKey: item}, nil
}

func deleteCredential(username string) error {
	return keyring.Delete(keyringService, username)
}

func LoadProfiles() (*ProfileStore, error) {
	if err := ensureAppDataDir(); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	dataFile := profilesFile()

	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &ProfileStore{}, nil
		}
		return nil, fmt.Errorf("failed to read profiles: %w", err)
	}

	var store ProfileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse profiles: %w", err)
	}

	return &store, nil
}

func SaveProfiles(store *ProfileStore) error {
	dataFile := profilesFile()

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	if err := os.WriteFile(dataFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write profiles: %w", err)
	}

	return nil
}

func AddProfile(username string, sessionKey string) error {
	store, err := LoadProfiles()
	if err != nil {
		return err
	}

	for _, p := range store.Profiles {
		if p.Username == username {
			return fmt.Errorf("profile %s already exists", username)
		}
	}

	if err := setCredential(username, sessionKey); err != nil {
		return fmt.Errorf("failed to store credential: %w", err)
	}

	store.Profiles = append(store.Profiles, Profile{
		Username:   username,
		SessionKey: sessionKey,
	})

	if err := SaveProfiles(store); err != nil {
		return err
	}

	if len(store.Profiles) == 1 {
		if err := SetActiveProfile(username); err != nil {
			return err
		}
	}

	return nil
}

func ListProfiles() ([]Profile, error) {
	store, err := LoadProfiles()
	if err != nil {
		return nil, err
	}
	return store.Profiles, nil
}

func GetActiveProfile() (string, error) {
	dataFile := activeProfileFile()

	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			profiles, err := ListProfiles()
			if err != nil {
				return "", err
			}
			if len(profiles) == 0 {
				return "", fmt.Errorf("no active profile and no profiles found")
			}
			return profiles[0].Username, nil
		}
		return "", fmt.Errorf("failed to read active profile: %w", err)
	}

	return string(data), nil
}

func SetActiveProfile(username string) error {
	store, err := LoadProfiles()
	if err != nil {
		return err
	}

	found := false
	for _, p := range store.Profiles {
		if p.Username == username {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("profile %s not found", username)
	}

	dataFile := activeProfileFile()
	if err := os.WriteFile(dataFile, []byte(username), 0600); err != nil {
		return fmt.Errorf("failed to write active profile: %w", err)
	}

	return nil
}

func DeleteProfile(username string) error {
	store, err := LoadProfiles()
	if err != nil {
		return err
	}

	found := -1
	for i, p := range store.Profiles {
		if p.Username == username {
			found = i
			break
		}
	}

	if found == -1 {
		return fmt.Errorf("profile %s not found", username)
	}

	store.Profiles = append(store.Profiles[:found], store.Profiles[found+1:]...)

	if err := deleteCredential(username); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to delete credential from keyring: %v\n", err)
	}

	if err := SaveProfiles(store); err != nil {
		return err
	}

	active, err := GetActiveProfile()
	if err == nil && active == username {
		if len(store.Profiles) > 0 {
			return SetActiveProfile(store.Profiles[0].Username)
		}
		os.Remove(activeProfileFile())
	}

	return nil
}

func ProfileExists(username string) bool {
	profiles, err := ListProfiles()
	if err != nil {
		return false
	}
	for _, p := range profiles {
		if p.Username == username {
			return true
		}
	}
	return false
}

func ResolveProfile(profileFlag string) (string, error) {
	if profileFlag != "" {
		return profileFlag, nil
	}
	return GetActiveProfile()
}

func GetCredential(username string) (*Profile, error) {
	return getCredential(username)
}
