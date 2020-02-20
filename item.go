// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/validator.v2"
)

// Item represents a QuickBooks Item object (a product type).
type Item struct {
	ID        string `json:"Id,omitempty"`
	SyncToken string `json:"SyncToken,omitempty"`
	//MetaData
	Name        string `json:"Name,omitempty" validate:"max=100"`
	SKU         string `json:"Sku,omitempty"`
	Description string `json:"Description,omitempty" validate:"max=4000"`
	Active      bool   `json:"Active,omitempty"`
	//SubItem
	//ParentRef
	//Level
	//FullyQualifiedName
	Taxable             bool           `json:"Taxable,omitempty"`
	SalesTaxIncluded    bool           `json:"SalesTaxIncluded,omitempty"`
	UnitPrice           json.Number    `json:"UnitPrice,omitempty"`
	Type                string         `json:"Type,omitempty"`
	IncomeAccountRef    *ReferenceType `json:"IncomeAccountRef,omitempty"`
	ExpenseAccountRef   *ReferenceType `json:"ExpenseAccountRef,omitempty"`
	PurchaseDesc        string         `json:"PurchaseDesc,omitempty"`
	PurchaseTaxIncluded bool           `json:"PurchaseTaxIncluded,omitempty"`
	PurchaseCost        json.Number    `json:"PurchaseCost,omitempty"`
	AssetAccountRef     *ReferenceType `json:"AssetAccountRef,omitempty"`
	TrackQtyOnHand      bool           `json:"TrackQtyOnHand,omitempty"`
	//InvStartDate Date
	QtyOnHand          json.Number    `json:"QtyOnHand,omitempty"`
	SalesTaxCodeRef    *ReferenceType `json:"SalesTaxCodeRef,omitempty"`
	PurchaseTaxCodeRef *ReferenceType `json:"PurchaseTaxCodeRef,omitempty"`
}

type ItemFilter struct {
	Name string
}

func (a *ItemFilter) Eq() string {
	sql := fmt.Sprintf("SELECT * FROM Account Name = %s MAXRESULTS %s", a.Name, strconv.Itoa(queryPageSize))
	return sql
}

// FetchItems returns the list of Items in the QuickBooks account. These are
// basically product types, and you need them to create invoices.
func (c *Client) FetchItems(filter *ItemFilter) ([]*Item, error) {
	var response struct {
		QueryResponse struct {
			Items         []*Item `json:"Item"`
			StartPosition int
			MaxResults    int
		}
	}

	sql := filter.Eq()
	if err := c.query(sql, &response); err != nil {
		return nil, err
	}

	return response.QueryResponse.Items, nil
}

// FetchItem returns just one particular Item from QuickBooks, by ID.
func (c *Client) FetchItem(id string) (*Item, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/item/" + id

	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
	}

	var r struct {
		Item Item
		Time Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Item, nil
}

// CreateItem creates the given Item on the QuickBooks server, returning
// the resulting Item object.
func (c *Client) CreateItem(item *Item) (*Item, error) {
	if err := validator.Validate(item); err != nil {
		return nil, err
	}

	u, err := url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/item"

	b, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	r := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", u.String(), r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status: %s", response.Status)
	}

	var i struct {
		Item *Item
	}
	err = json.NewDecoder(response.Body).Decode(&i)
	return i.Item, err
}
