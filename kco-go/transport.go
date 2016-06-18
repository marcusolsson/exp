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

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type bearerRoundTripper struct {
	Base   http.RoundTripper
	Secret string
}

// RoundTrip generates the bearer token.
//
// Every API request has to be authenticated by using the Authentication HTTP
// request header. Only a proprietary Klarna authentication scheme is
// supported, in the format of Klarna authorization header.
//
// The authorization header is calculated for each request using these steps of
// this formula: base64(hex(sha256 (request_payload + shared_secret)))
func (rt *bearerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r := cloneRequest(req)

	var buf bytes.Buffer

	if req.Body != nil {
		io.Copy(&buf, req.Body)
		r.Body = ioutil.NopCloser(&buf)
	}

	data := append(buf.Bytes(), []byte(rt.Secret)...)

	var sum string
	sum = fmt.Sprintf("%x", sha256.Sum256(data))
	sum = base64.URLEncoding.EncodeToString([]byte(sum))

	r.Header.Set("Authorization", "Klarna "+sum)

	return rt.Base.RoundTrip(r)
}

func newBearerTokenTransport(secret string) *bearerRoundTripper {
	return &bearerRoundTripper{
		Base:   http.DefaultTransport,
		Secret: secret,
	}
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
