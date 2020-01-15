package quickbooks

import (
	"fmt"
)

type AccountQueryResponse struct {
	Accounts []*Account `json:"Account"`
}

type Account struct {
	ID      string `json:"Id"`
	Name    string `json:"Name"`
	Active  bool   `json:"Active"`
	Type    string `json:"AccountType"`
}

func (c *Client) FetchAccounts(sql string) ([]*Account, error) {
	var response struct {
		QueryResponse AccountQueryResponse `json:"QueryResponse"`
	}

	err := c.query(sql, &response)
	if err != nil {
		return nil, err
	}

	if len(response.QueryResponse.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts returned for query: %s\n", sql)
	}

	return response.QueryResponse.Accounts, nil
}
