package config

type MachineConfig struct {
	Host      string `mapstructure:"host"`
	Available bool   `mapstructure:"available"`
}

type MachineList []MachineConfig

type LBConfig struct {
	LB map[string]MachineList `mapstructure:"lb"`
}

func (ml MachineList) GetAvailableMachine() []MachineConfig {
	res := make([]MachineConfig, 0)
	for _, machine := range ml {
		res = append(res, machine)
	}

	return res
}
