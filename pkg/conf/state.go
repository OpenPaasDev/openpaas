package conf

import "time"

type State struct {
	IsInitialised  bool      `yaml:"is_initialised"`
	LastNodeUpdate time.Time `yaml:"last_node_update"`
	LastNodeIDs    []string  `yaml:"last_node_ids"`
}

func (f *State) Save(baseDir string) error {
	return nil
}

func LoadState(baseDir string) (*State, error) {
	return nil, nil
}
