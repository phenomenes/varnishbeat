// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period    time.Duration `config:"period"`
	Directory string        `config:"directory"`
	Stats     bool          `config:"stats"`
	Log       bool          `config:"log"`
}

var DefaultConfig = Config{
	Period:    1 * time.Second,
	Directory: "",
	Stats:     false,
	Log:       false,
}
