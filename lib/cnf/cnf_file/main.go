package cnf_file

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/log"
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

func Load(validationMode bool, fail bool) {
	log.Info("config", "Loading config from file")
	newConfig := cnf.Config{}
	err := yaml.Unmarshal(readConfig(), &newConfig)
	if err != nil {
		log.ErrorS("config", "Failed to parse config! Check if it's schema is valid!")
		if !fail {
			return
		}
		panic(fmt.Errorf("failed to parse config: %v", err))
	}
	if !validateConfig(newConfig, fail) {
		if !fail {
			log.ErrorS("config", "Failed to validate config!")
			return
		}
		panic(fmt.Errorf("failed to vaidate config!"))
	}
	if !validationMode {
		cnf.C = &newConfig
	}
	newRules := ParseRules(cnf.C.Rules)
	cnf.RULES = &newRules
	log.Debug("config", "Finished loading config")
	log.Debug("config", fmt.Sprintf("CONFIG DUMP: %+v", cnf.C))
	log.Debug("config", fmt.Sprintf("PARSED RULES: %+v", cnf.RULES))
}
