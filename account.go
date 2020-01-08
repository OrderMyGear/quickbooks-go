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
	SubType string `json:"AccountSubType"`
}

func (c *Client) FetchChartOfAccounts() ([]*Account, error) {
	var response struct {
		QueryResponse AccountQueryResponse `json:"QueryResponse"`
	}

	err := c.query("SELECT * FROM Account", &response)
	if err != nil {
		return nil, err
	}

	if len(response.QueryResponse.Accounts) == 0 {
		return nil, fmt.Errorf("chart of accounts is empty")
	}

	return response.QueryResponse.Accounts, nil
}
