package errwrapper

import (
	"errors"
	"fmt"
	"strings"
)

const separator = ":"

// Wrap will append errs to existing err and return new error with all previous errors
func Wrap(existingErr error, newErrs ...error) error {
	var errMessages []string
	for _, e := range newErrs {
		if e == nil {
			continue
		}
		errMessages = append(errMessages, e.Error())
	}

	errsFormatted := strings.Join(errMessages, separator)

	if len(errsFormatted) == 0 {
		return nil
	}

	if existingErr == nil {
		existingErr = errors.New(errsFormatted)
	} else {
		existingErr = fmt.Errorf("%s%s%s", existingErr, separator, errsFormatted)
	}

	return existingErr
}
