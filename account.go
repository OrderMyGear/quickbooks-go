package quickbooks

import (
	"fmt"
	"strconv"
	"strings"
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
	IsActive bool
	Type     string
}

func (a *AccountFilter) Eq() string {
	sqlSelect := "SELECT * FROM Account"
	sqlIsActive := "WHERE Active = " + strconv.FormatBool(a.IsActive)
	sqlType := "AND AccountType = '" + a.Type + "'"
	return strings.Join([]string{sqlSelect, sqlIsActive, sqlType}, " ")
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
