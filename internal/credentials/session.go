package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Session struct {
	PassCred *SessionResponse
	Setup    Setup
}

type SessionResponse struct {
	Access_Token string
	Instance_URL string
	Id           string
	Token_Type   string
	Issued_At    string
	Signature    string
}

type ResourceCalls interface {
	ServiceURL() string
	AuthorizationHeader(*http.Request)
	Client() *http.Client
}

const oauthEndpoint = "/services/oauth2/token"

func Auth(setup Setup) (*Session, error) {
	if err := IsValidSetup(setup); err != nil {
		return nil, err
	}

	request, err := SessionPasswordRequest(setup.PasswordCredentials)

	if err != nil {
		return nil, err
	}

	response, err := SessionPasswordResponse(request, setup.Protocol)

	if err != nil {
		return nil, err
	}

	return &Session{
		PassCred: response,
		Setup:    setup,
	}, nil

}

func SessionPasswordRequest(passCred PasswordCredentials) (*http.Request, error) {
	oauthUrl := passCred.Url + oauthEndpoint

	body, err := UrlEncoder(passCred)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, oauthUrl, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")
	return request, nil
}

func SessionPasswordResponse(request *http.Request, client *http.Client) (*SessionResponse, error) {
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(" Error response : %d %s %+v", response.StatusCode, response.Status, response.Body)
	}

	// body, err := io.ReadAll(response.Body)
	// fmt.Printf("Decoded response : %s", string(body))
	decoded := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var sessionResponse SessionResponse
	if err := decoded.Decode(&sessionResponse); err != nil {
		return nil, err
	}

	fmt.Printf("Session response : %v", &sessionResponse)

	return &sessionResponse, nil

}

func (session *Session) ServiceURL() string {
	return fmt.Sprintf("%s/services/data/v52.0", session.PassCred.Instance_URL)
}

func (session *Session) AuthorizationHeader(request *http.Request) {
	auth := fmt.Sprintf("%s %s", session.PassCred.Token_Type, session.PassCred.Access_Token)
	request.Header.Add("Authorization", auth)
}

func (session *Session) Client() *http.Client {
	return session.Setup.Protocol
}

func UrlEncoder(passCred PasswordCredentials) (io.Reader, error) {
	encodedUrl := url.Values{}
	encodedUrl.Add("grant_type", "password")
	encodedUrl.Add("username", passCred.Username)
	encodedUrl.Add("password", passCred.Password)
	encodedUrl.Add("client_id", passCred.ClientId)
	encodedUrl.Add("client_secret", passCred.ClientSecret)

	return strings.NewReader(encodedUrl.Encode()), nil
}

func IsValidSetup(setup Setup) error {
	if setup.PasswordCredentials == (PasswordCredentials{}) {
		return errors.New("credentials cannot be empty")
	}
	return nil
}
