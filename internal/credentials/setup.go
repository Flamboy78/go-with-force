package credentials

import (
	"net/http"
)

type Setup struct {
	PasswordCredentials PasswordCredentials
	Protocol            *http.Client
}
