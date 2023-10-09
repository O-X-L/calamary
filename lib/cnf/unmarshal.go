package cnf

import (
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

// allow single string to be supplied
type YamlStringArray []string

func (a *YamlStringArray) UnmarshalYAML(value *yaml.Node) error {
	var multi []string
	err := value.Decode(&multi)
	if err != nil {
		var single string
		err := value.Decode(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

// apply defaults from tags on unmarshal
func (s *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	type plain Config
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func (s *ServiceConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	type plain ServiceConfig
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func (s *ServiceListener) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	type plain ServiceListener
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func (s *ServiceTimeout) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	type plain ServiceTimeout
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func (s *ServiceOutput) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	type plain ServiceOutput
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func (s *ServiceMetrics) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(s); err != nil {
		return err
	}

	type plain ServiceMetrics
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}
