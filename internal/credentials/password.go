package credentials

import (
	"errors"
)

type PasswordCredentials struct {
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
	Url          string
}

func GeneratePasswordCredentials(passCred PasswordCredentials) (*PasswordCredentials, error) {
	if err := IsValidPasswordCredentials(passCred); err != nil {
		return nil, err
	}

	return &PasswordCredentials{
		Username:     passCred.Username,
		Password:     passCred.Password,
		ClientId:     passCred.ClientId,
		ClientSecret: passCred.ClientSecret,
		Url:          passCred.Url,
	}, nil
}

func IsValidPasswordCredentials(passCred PasswordCredentials) error {
	switch {
	case len(passCred.Username) == 0:
		return errors.New("username cannot be empty")
	case len(passCred.Password) == 0:
		return errors.New("password cannot be empty")
	case len(passCred.ClientId) == 0:
		return errors.New("clientId cannot be empty")
	case len(passCred.ClientSecret) == 0:
		return errors.New("clientSecret cannot be empty")
	case len(passCred.Url) == 0:
		return errors.New("url cannot be empty")
	}

	return nil
}
