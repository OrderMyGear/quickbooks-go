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

type AccountFilter struct {
	IsActive string // must be either "true" or "false"
	Type     string
}

func (a *AccountFilter) Eq() string {
	sql := "SELECT * FROM Account"
	sql += " WHERE Active = " + a.IsActive
	if a.Type != "" {
		sql += " AND AccountType = '" + a.Type + "'"
	}
	return sql
}

func (c *Client) FetchAccounts(filter *AccountFilter) ([]*Account, error) {
	var response struct {
		QueryResponse AccountQueryResponse `json:"QueryResponse"`
	}

	sql := filter.Eq()
	if err := c.query(sql, &response); err != nil {
		return nil, err
	}

	if len(response.QueryResponse.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts returned for query: %s\n", sql)
	}

	return response.QueryResponse.Accounts, nil
}
