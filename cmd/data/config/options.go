package config

type Options struct {
	Options *map[string]bool
}

func NewOptions(optionsMap *map[string]bool) *Options {
	return &Options{
		Options: optionsMap,
	}
}

func (c *Config) SetOption(optionKey string, optionValue bool) error {
	(*c.Options.Options)[optionKey] = optionValue
	err := c.DB.SetOptionRepo(optionKey, optionValue)
	return err
}
