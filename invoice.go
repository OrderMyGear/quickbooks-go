// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Invoice struct {
	Lines       []Line        `json:"Line"`
	CustomerRef ReferenceType `json:"CustomerRef"`
}

type Line struct {
	Amount              float64             `json:"Amount"`
	Description         string              `json:"Description"`
	DetailType          string              `json:"DetailType"`
	SalesItemLineDetail SalesItemLineDetail `json:"SalesItemLineDetail"`
}

// SalesItemLineDetail ...
type SalesItemLineDetail struct {
	ItemRef   ReferenceType `json:"ItemRef"`
}

// CreateInvoice creates the given Invoice on the QuickBooks server, returning
// the resulting Invoice object.
func (c *Client) CreateInvoice(inv *Invoice) (*Invoice, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/invoice"
	var j []byte
	j, err = json.Marshal(inv)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
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
		Invoice Invoice
		Time    time.Time
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Invoice, err
}
