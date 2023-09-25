package cnf_file

import (
	"fmt"
	"os"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
	"gopkg.in/yaml.v2"
)

func readConfigFile(file string) (config []byte, err error) {
	_, err = os.Stat(file)
	if !os.IsNotExist(err) {
		config, err = os.ReadFile(file)
		if err != nil {
			// log.ErrorS("config", "Unable to read config-file '%s'", file)
			return nil, err
		}
	}
	return
}

func readConfig() (config []byte) {
	cwd, err := os.Getwd()
	cwdConfig := "./calamary.yml"
	if err == nil {
		cwdConfig = cwd + "/calamary.yml"
		config, err := readConfigFile(cwdConfig)
		if err == nil && config != nil {
			return config
		}
	}
	config, err = readConfigFile(cnf.ConfigFileAbs)
	if err == nil && config != nil {
		return config
	}
	log.ErrorS("config", fmt.Sprintf(
		"Neither config file could be read: (%s, %s)", cwdConfig, cnf.ConfigFileAbs,
	))
	panic(fmt.Errorf("no valid config file found! (%s, %s)", cwdConfig, cnf.ConfigFileAbs))
}

func Load() {
	log.Info("config", "Loading config from file")
	newConfig := cnf.Config{}
	err := yaml.Unmarshal(readConfig(), &newConfig)
	if err != nil {
		log.ErrorS("config", "Failed to parse config! Check if it is valid!")
		panic(fmt.Errorf("failed to parse config"))
	}
	cnf.C = &newConfig
	newRules := ParseRules(cnf.C.Rules)
	cnf.RULES = &newRules
	log.Debug("config", "Finished loading config")
}
