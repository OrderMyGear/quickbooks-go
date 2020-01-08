package quickbooks

import (
	"fmt"
)

type Account struct {
	ID      string `json:"Id"`
	Name    string `json:"Name"`
	Active  bool   `json:"Active"`
	SubType string `json:"AccountSubType"`
}

func (c *Client) FetchAccounts() ([]*Account, error) {
	var r struct {
		QueryResponse struct {
			Accounts []*Account
		}
	}

	err := c.query("SELECT * FROM Account", &r)
	if err != nil {
		return nil, err
	}

	if len(r.QueryResponse.Accounts) == 0 {
		return nil, fmt.Errorf("chart of accounts is empty")
	}

	return r.QueryResponse.Accounts, nil
}
