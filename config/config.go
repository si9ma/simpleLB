package config

type HostConfig struct {
	Host      string `yaml:"host"`
	Available bool   `yaml:"available"`
}

type LBConfig struct {
	LB map[string][]HostConfig `yaml:"lb"`
}
