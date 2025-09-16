package state

import (
	"encoding/json"
	"errors"
	"os"
)

type State struct {
	Version   int             `json:"version"`
	Resources []ResourceState `json:"resources"`
}

type ResourceState struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

func LoadState(path string) (*State, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var st State
	if err := json.Unmarshal(b, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func SaveState(path string, st *State) error {
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}

	dir := ".mini-terra"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, b, 0644)
}

func PrettyJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func FindResource(st *State, typ, name string) (*ResourceState, int) {
	if st == nil {
		return nil, -1
	}
	for i := range st.Resources {
		r := &st.Resources[i]
		if r.Type == typ && r.Name == name {
			return r, i
		}
	}
	return nil, -1
}

func NewEmptyState() *State {
	return &State{Version: 1, Resources: []ResourceState{}}
}

func EnsureState(st **State) {
	if *st == nil {
		*st = NewEmptyState()
	}
}

var ErrNotFound = errors.New("resource not found")
