// Copyright 2016 Marcus Olsson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kco

import "net/http"

// Client ...
type Client struct {
	httpClient *http.Client
	endpoint   string
}

// Entrypoints
var (
	LiveEnvironmentURL = "https://checkout.klarna.com"
	TestEnvironmentURL = "https://checkout.testdrive.klarna.com"
)

// NewCustomClient returns a client with a custom HTTP client.
func NewCustomClient(u string, client *http.Client) *Client {

	return &Client{
		httpClient: client,
		endpoint:   u,
	}
}

// NewAuthClient returns an authenticated client.
func NewAuthClient(secret string, u string) *Client {
	tt := newBearerTokenTransport(secret)
	client := &http.Client{
		Transport: tt,
	}
	return NewCustomClient(u, client)
}

// NewClient returns a client for making Klarna API calls.
func NewClient() *Client {
	return NewCustomClient(LiveEnvironmentURL, http.DefaultClient)
}
