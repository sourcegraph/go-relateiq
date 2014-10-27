package relateiq

import (
	"encoding/json"
	"net/http"
	"time"
)

// An Account represents companies (or other entities). Accounts can
// have any kind of relationship with your company -- they could be
// leads, clients, former clients, or partners of your Organization.
type Account struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ModifiedDate Time   `json:"modifiedDate"`

	// TODO(sqs): add fieldValues
}

// Time is a time.Time whose JSON encoding is the number of
// milliseconds since the Unix epoch.
type Time struct{ time.Time }

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UnixNano() / int64(time.Millisecond))
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(b []byte) error {
	var msec int64
	if err := json.Unmarshal(b, &msec); err != nil {
		return err
	}
	*t = Time{time.Unix(0, 1000000*msec).In(time.UTC)}
	return nil
}

// AccountsService communicates with the account-related endpoints of
// the RelateIQ API.
type AccountsService struct {
	c *Client
}

// AccountsListOptions specifies options for listing accounts.
type AccountsListOptions struct {
	IDs []string `url:"_ids,omitempty,comma"` // a list of account IDs
	ListOptions
}

// List returns all accounts in your organization.
func (s *AccountsService) List(opt AccountsListOptions) ([]*Account, *http.Response, error) {
	req, err := s.c.NewRequest("GET", "accounts", opt, nil)
	if err != nil {
		return nil, nil, err
	}

	var list *struct {
		Objects []*Account `json:"objects"`
	}
	resp, err := s.c.Do(req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list.Objects, resp, err
}
