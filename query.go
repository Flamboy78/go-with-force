package soql

import (
	"encoding/json"
	"fmt"
	"klnef/go-with-force/internal/credentials"
	"net/http"
	"net/url"
)

type Resource struct {
	session credentials.ResourceCalls
}

type Error struct {
	ErrorCode string
	Message   string
	Fields    []string
}

type QueryResponse struct {
	Done           bool
	TotalSize      int
	NextRecordsURL string
	Records        []map[string]interface{}
}

func NewResource(session credentials.ResourceCalls) (Resource, error) {
	if session == nil {
		return Resource{}, fmt.Errorf("error : session can not be nil")
	}

	return Resource{
		session: session,
	}, nil
}

func (r *Resource) Query(query string) (QueryResponse, error) {
	if query == "" {
		return QueryResponse{}, fmt.Errorf("soql resource query: query can not be nil")
	}

	request, err := r.queryRequest(query)
	if err != nil {
		return QueryResponse{}, err
	}

	response, err := r.queryResponse(request)
	if err != nil {
		return QueryResponse{}, err
	}
	return response, nil
}

func (r *Resource) queryRequest(query string) (*http.Request, error) {
	const endpoint string = "/query"
	queryURL := r.session.ServiceURL() + endpoint + "/"

	form := url.Values{}
	form.Add("q", query)
	queryURL += "?" + form.Encode()

	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	r.session.AuthorizationHeader(request)
	return request, nil
}

func (r *Resource) queryResponse(request *http.Request) (QueryResponse, error) {
	response, err := r.session.Client().Do(request)

	if err != nil {
		return QueryResponse{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var queryErrs []Error
		err = decoder.Decode(&queryErrs)
		var errMsg error
		if err == nil {
			for _, queryErr := range queryErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", queryErr.ErrorCode, queryErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}

		return QueryResponse{}, errMsg
	}

	var resp QueryResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return QueryResponse{}, err
	}

	return resp, nil
}
