package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Identity represents a git user identity with name, email, and optional signing key
type Identity struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	SigningKey string `json:"signing_key,omitempty"`
}

// IdentityStore manages the collection of saved identities
type IdentityStore struct {
	Identities []Identity `json:"identities"`
	filePath   string
}

// getStorePath returns the path to the identity store file.
// The store location can be overridden by setting the GITID_STORE env variable,
// which is handy when managing separate stores for work and personal machines.
func getStorePath() (string, error) {
	if override := os.Getenv("GITID_STORE"); override != "" {
		return override, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(homeDir, ".gitid", "identities.json"), nil
}

// LoadStore reads the identity store from disk, creating it if it doesn't exist
func LoadStore() (*IdentityStore, error) {
	path, err := getStorePath()
	if err != nil {
		return nil, err
	}

	store := &IdentityStore{filePath: path}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// Return empty store if file doesn't exist yet
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read identity store: %w", err)
	}

	if err := json.Unmarshal(data, store); err != nil {
		return nil, fmt.Errorf("could not parse identity store: %w", err)
	}

	return store, nil
}

// Save writes the identity store to disk
func (s *IdentityStore) Save() error {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize identities: %w", err)
	}

	// Use 0600 instead of 0644 — identities file may contain sensitive info
	// (e.g. signing key references), so restrict read access to owner only.
	if err := os.WriteFile(s.filePath, data, 0600); err != nil {
		return fmt.Errorf("could not write identity store: %w", err)
	}

	return nil
}

// Add inserts a new identity into the store, returning an error if the ID already exists
func (s *IdentityStore) Add(identity Identity) error {
	for _, existing := range s.Identities {
		if strings.EqualFold(existing.ID, identity.ID) {
			return fmt.Errorf("identity with id %q already exists", identity.ID)
		}
	}
	s.Identities = append(s.Identities, identity)
	return nil
}

// Remove deletes an identity by ID, returning an error if not found
func (s *IdentityStore) Remove(id string) error {
	for i, identity := range s.Identities {
		if strings.EqualFold(identity.ID, id) {
			s.Identities = append(s.Identities[:i], s.Identities[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("identity %q not found", id)
}

// FindByID looks up an identity by its ID
func (s *IdentityStore) FindByID(id string) (*Identity, error) {
	for i, identity := range s.Identities {
		if strings.EqualFold(identity.ID, id) {
			return &s.Identities[i], nil
