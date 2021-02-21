package config

import "sync"

var once sync.Once

type Configuration struct {
	BasePath string
	Langs    map[string]string `json:"langs"`
}

var (
	instance *Configuration
)

func NewConfig() *Configuration {
	once.Do(func() { // <-- atomic, does not allow repeating
		instance = &Configuration{}
	})
	return instance
}
