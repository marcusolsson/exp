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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Available order statuses.
var (
	StatusCheckoutIncomplete = "checkout_incomplete"
	StatusCheckoutComplete   = "checkout_complete"
	StatusCreated            = "created"
)

// Order ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#resource-properties
type Order struct {
	ID                string             `json:"id,omitempty"`
	MerchantReference *MerchantReference `json:"merchant_reference,omitempty"`
	PurchaseCountry   string             `json:"purchase_country"`
	PurchaseCurrency  string             `json:"purchase_currency"`
	Locale            string             `json:"locale"`
	Status            string             `json:"status,omitempty"`
	Reference         string             `json:"reference,omitempty"`
	Reservation       string             `json:"reservation,omitempty"`
	Recurring         bool               `json:"recurring,omitempty"`
	RecurringToken    string             `json:"recurring_token,omitempty"`
	Cart              Cart               `json:"cart"`
	BillingAddress    *Address           `json:"billing_address,omitempty"`
	ShippingAddress   *Address           `json:"shipping_address,omitempty"`
	Customer          *Customer          `json:"customer,omitempty"`
	GUI               *GUI               `json:"gui,omitempty"`
	Merchant          Merchant           `json:"merchant"`
	StartedAt         string             `json:"started_at,omitempty"`
	CompletedAt       string             `json:"completed_at,omitempty"`
	CreatedAt         string             `json:"created_at,omitempty"`
	LastModifiedAt    string             `json:"last_modified_at,omitempty"`
	ExpiresAt         string             `json:"expires_at,omitempty"`
}

// ErrCheckoutIncomplete ...
var ErrCheckoutIncomplete = errors.New("checkout incomplete")

// Acknowledge ...
func (o *Order) Acknowledge() error {
	if o.Status != StatusCheckoutComplete {
		return ErrCheckoutIncomplete
	}
	o.Status = StatusCreated

	return nil
}

// Render ...
func (o *Order) Render(w io.Writer) {
	if o.GUI != nil {
		w.Write([]byte(o.GUI.Snippet))
	}
}

// MerchantReference ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#merchant_reference-object-properties
type MerchantReference struct {
	OrderID1 string `json:"orderid1"`
	OrderID2 string `json:"orderid2"`
}

// Address ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#address-object-properties
type Address struct {
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	CareOf        string `json:"care_of,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	StreetName    string `json:"street_name,omitempty"`
	StreetNumber  string `json:"street_number,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	City          string `json:"city,omitempty"`
	Country       string `json:"country,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Title         string `json:"title,omitempty"`
}

// Cart ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#cart-object-properties
type Cart struct {
	TotalPriceExclTax int        `json:"total_price_excluding_tax,omitempty"`
	TotalTaxAmount    int        `json:"total_tax_amount,omitempty"`
	TotalPriceInclTax int        `json:"total_price_including_tax,omitempty"`
	Items             []CartItem `json:"items,omitempty"`
}

// CartItem ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#cart-item-object-properties
type CartItem struct {
	Type              string `json:"type,omitempty"`
	Reference         string `json:"reference,omitempty"`
	Name              string `json:"name,omitempty"`
	Qty               int    `json:"quantity,omitempty"`
	EAN               string `json:"ean,omitempty"`
	URI               string `json:"uri,omitempty"`
	ImageURI          string `json:"image_uri,omitempty"`
	UnitPrice         int    `json:"unit_price,omitempty"`
	TotalPriceExclTax int    `json:"total_price_excluding_tax,omitempty"`
	TotalTaxAmount    int    `json:"total_tax_amount,omitempty"`
	TotalPriceInclTax int    `json:"total_price_including_tax,omitempty"`
	DiscountRate      int    `json:"discount_rate,omitempty"`
	TaxRate           int    `json:"tax_rate,omitempty"`
}

// Customer ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#customer-object-properties
type Customer struct {
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	Type        string `json:"type"`
}

// GUI ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#gui-object-properties
type GUI struct {
	Layout  string   `json:"layout,omitempty"`
	Options []string `json:"options,omitempty"`
	Snippet string   `json:"snippet,omitempty"`
}

// Merchant ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#merchant-object-properties
type Merchant struct {
	ID              string `json:"id,omitempty"`
	BackToStoreURI  string `json:"back_to_store_uri,omitempty"`
	TermsURI        string `json:"terms_uri,omitempty"`
	CancelTermsURI  string `json:"cancellation_terms_uri,omitempty"`
	CheckoutURI     string `json:"checkout_uri,omitempty"`
	ConfirmationURI string `json:"confirmation_uri,omitempty"`
	PushURI         string `json:"push_uri,omitempty"`
	ValidationURI   string `json:"validation_uri,omitempty"`
}

// Attachment ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#attachment-object-properties
type Attachment struct {
	Body        string `json:"body"`
	ContentType string `json:"content_type"`
}

// ExternalPaymentMethod ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#external_payment_method-object-properties
type ExternalPaymentMethod struct {
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
	ImageURI    string `json:"image_uri"`
	Fee         int    `json:"fee"`
}

// Options ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#options-object-properties
type Options struct {
	ColorButton                  string   `json:"color_button"`
	ColorButtonText              string   `json:"color_button_text"`
	ColorCheckbox                string   `json:"color_checkbox"`
	ColorCheckboxCheckmark       string   `json:"color_checkbox_checkmark"`
	ColorHeader                  string   `json:"color_header"`
	ColorLink                    string   `json:"color_link"`
	ShippingDetails              string   `json:"shipping_details"`
	PhoneMandatory               bool     `json:"phone_mandatory"`
	AllowSeparateShippingAddress bool     `json:"allow_separate_shipping_address"`
	PackStationEnabled           bool     `json:"packstation_enabled"`
	DateOfBirthMandatory         bool     `json:"date_of_birth_mandatory"`
	AdditionalCheckbox           Checkbox `json:"additional_checkbox"`
}

// Checkbox ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#checkbox-object-properties
type Checkbox struct {
	Text     string `json:"text"`
	Checked  bool   `json:"checked"`
	Required bool   `json:"required"`
}

// MerchantRequested ...
//
// https://developers.klarna.com/en/se/kco-v2/checkout-api#merchant_requested-object-properties
type MerchantRequested struct {
	AdditionalCheckbox bool `json:"additional_checkbox"`
}

// CreateOrder creates a new order.
func (c *Client) CreateOrder(o Order) (string, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return "", err
	}

	// Examples on developer.klarna.com uses unescaped ampersands.
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	req, err := http.NewRequest("POST", c.endpoint+"/checkout/orders", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Accept", "application/vnd.klarna.checkout.aggregated-order-v2+json")
	req.Header.Set("Content-Type", "application/vnd.klarna.checkout.aggregated-order-v2+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status %s", resp.Status)
	}

	u, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return "", err
	}

	// Extract Order ID from Location header.
	s := strings.Split(u.RequestURI(), "/")
	id := s[len(s)-1]

	return id, nil
}

// Order returns an existing Checkout order.
func (c *Client) Order(id string) (*Order, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/checkout/orders/"+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Accept", "application/vnd.klarna.checkout.aggregated-order-v2+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var result Order
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateOrder updates an existing Checkout order.
func (c *Client) UpdateOrder(id string, o *Order) (*Order, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	// Examples on developer.klarna.com uses unescaped ampersands.
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	req, err := http.NewRequest("POST", c.endpoint+"/checkout/orders/"+id, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Accept", "application/vnd.klarna.checkout.aggregated-order-v2+json")
	req.Header.Set("Content-Type", "application/vnd.klarna.checkout.aggregated-order-v2+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Order
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
