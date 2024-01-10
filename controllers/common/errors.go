package common

import "errors"

type InvalidConfiguration error

func InvalidConfigurationError(msg string) InvalidConfiguration {
	return errors.New(msg)
}
