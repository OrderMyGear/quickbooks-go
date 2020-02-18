// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

/*
Package quickbooks provides access to Intuit's QuickBooks Online API.

NOTE: This library is very incomplete. I just implemented the minimum for my
use case. Pull requests welcome :)

 // Do this after you go through the normal OAuth process.
 var client = oauth2.NewClient(ctx, tokenSource)

 // Initialize the client handle.
 var qb = quickbooks.Client{
	 Client: client,
	 Endpoint: quickbooks.SandboxEndpoint,
	 RealmID: "some company account ID"'
 }

 // Make a request!
 var companyInfo, err = qb.FetchCompanyInfo()
*/
package quickbooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Client is your handle to the QuickBooks API.
type Client struct {
	// Get this from oauth2.NewClient().
	Client *http.Client
	// Set to ProductionEndpoint or SandboxEndpoint.
	Endpoint EndpointURL
	// The account ID you're connecting to.
	RealmID string
}

// FetchCompanyInfo returns the QuickBooks CompanyInfo object. This is a good
// test to check whether you're connected.
func (c *Client) FetchCompanyInfo() (*CompanyInfo, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/companyinfo/" + c.RealmID
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
		CompanyInfo CompanyInfo
		Time        Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.CompanyInfo, err
}

// query makes the specified QBO `query` and unmarshals the result into `out`
func (c *Client) query(query string, out interface{}) error {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return err
	}
	u.Path = "/v3/company/" + c.RealmID + "/query"

	var v = url.Values{}
	v.Add("query", query)
	u.RawQuery = v.Encode()
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// getByID will build and send the request to get the given resource by the given ID
// ex: c.getByID("account", "3", &accountResponse)
func (c *Client) getByID(resource, id string, response interface{}) error {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("/v3/company/%s/%s/%s", c.RealmID, resource, id)

	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
	}

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return err
	}
	return nil
}
