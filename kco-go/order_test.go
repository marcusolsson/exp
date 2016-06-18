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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestAggregatedOrder_Create(t *testing.T) {
	orderID := "1A2B3C4D5E6F1A2B3C4D5E6F"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Validate request sent from client.
		freq, err := os.Open("testdata/create_order_request.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer freq.Close()

		want, _ := normalizeJSON(freq)
		got, _ := normalizeJSON(req.Body)

		if !bytes.Equal(want, got) {
			t.Errorf("\nwant = %s;\nhave = %s", want, got)
		}

		// Return canned response.
		w.Header().Set("Content-Type", "application/vnd.klarna.checkout.aggregated-order-v2+json")
		w.Header().Set("Location", "https://checkout.testdrive.klarna.com/checkout/orders/"+orderID)
		w.WriteHeader(http.StatusCreated)
	}))

	client := NewAuthClient("s3cret", srv.URL)

	order := Order{
		MerchantReference: &MerchantReference{
			OrderID1: "123456789",
			OrderID2: "123456789",
		},
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Cart: Cart{
			Items: []CartItem{
				{
					Reference:    "123456789",
					Name:         "Klarna t-shirt",
					Qty:          2,
					EAN:          "1234567890123",
					URI:          "http://example.com/product.php?123456789",
					ImageURI:     "http://example.com/product_image.php?123456789",
					UnitPrice:    12300,
					DiscountRate: 1000,
					TaxRate:      2500,
				},
				{
					Type:      "shipping_fee",
					Reference: "SHIPPING",
					Name:      "Shipping fee",
					Qty:       1,
					UnitPrice: 4900,
					TaxRate:   2500,
				},
			},
		},
		ShippingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		GUI: &GUI{
			Layout: "desktop",
		},
		Merchant: Merchant{
			ID:              "0",
			BackToStoreURI:  "http://example.com",
			TermsURI:        "http://example.com/terms.php",
			CheckoutURI:     "https://example.com/checkout.php",
			ConfirmationURI: "https://example.com/thankyou.php?sid=123&klarna_order={checkout.order.uri}",
			PushURI:         "https://example.com/push.php?sid=123&klarna_order={checkout.order.uri}",
		},
	}

	id, err := client.CreateOrder(order)
	if err != nil {
		t.Fatal(err)
	}

	if id != orderID {
		t.Errorf("id = %s; want = %s", id, orderID)
	}
}

func TestAggregatedOrder_Read(t *testing.T) {
	orderID := "1A2B3C4D5E6F1A2B3C4D5E6F"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f, err := os.Open("testdata/read_order_response.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "application/vnd.klarna.checkout.aggregated-order-v2+json")
		w.WriteHeader(http.StatusOK)

		io.Copy(w, f)
	}))

	client := NewAuthClient("s3cret", srv.URL)

	order, err := client.Order(orderID)
	if err != nil {
		t.Fatal(err)
	}

	if order.ID != orderID {
		t.Errorf("order.ID = %s; want = %s", order.ID, orderID)
	}
}

func TestAggregatedOrder_Update(t *testing.T) {
	orderID := "1A2B3C4D5E6F1A2B3C4D5E6F"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Validate request sent from client.
		freq, err := os.Open("testdata/update_order_request.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer freq.Close()

		want, _ := normalizeJSON(freq)
		got, _ := normalizeJSON(req.Body)

		if !bytes.Equal(want, got) {
			t.Errorf("\nwant = %s;\nhave = %s", want, got)
		}

		// Return canned response.
		fresp, err := os.Open("testdata/update_order_response.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer fresp.Close()

		w.Header().Set("Content-Type", "application/vnd.klarna.checkout.aggregated-order-v2+json")
		w.WriteHeader(http.StatusOK)

		io.Copy(w, fresp)
	}))

	client := NewAuthClient("s3cret", srv.URL)
	client.endpoint = srv.URL

	update := &Order{
		MerchantReference: &MerchantReference{
			OrderID1: "123456789",
			OrderID2: "123456789",
		},
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Cart: Cart{
			Items: []CartItem{
				{
					Reference:    "123456789",
					Name:         "Klarna t-shirt",
					Qty:          4,
					EAN:          "1234567890123",
					URI:          "http://example.com/product.php?123456789",
					ImageURI:     "http://example.com/product_image.php?123456789",
					UnitPrice:    12300,
					DiscountRate: 1000,
					TaxRate:      2500,
				},
				{
					Type:      "shipping_fee",
					Reference: "SHIPPING",
					Name:      "Shipping fee",
					Qty:       1,
					UnitPrice: 4900,
					TaxRate:   2500,
				},
			},
		},
		ShippingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		GUI: &GUI{
			Layout: "desktop",
		},
		Merchant: Merchant{
			ID:              "0",
			TermsURI:        "http://example.com/terms.php",
			CheckoutURI:     "https://example.com/checkout.php",
			ConfirmationURI: "https://example.com/thankyou.php?sid=123&klarna_order={checkout.order.uri}",
			PushURI:         "https://example.com/push.php?sid=123&klarna_order={checkout.order.uri}",
		},
	}

	order, err := client.UpdateOrder(orderID, update)
	if err != nil {
		t.Fatal(err)
	}

	if order.ID != orderID {
		t.Errorf("order.ID = %s; want = %s", order.ID, orderID)
	}
}

func TestUnmarshal_CreateOrderRequest(t *testing.T) {
	f, err := os.Open("testdata/create_order_request.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var order Order
	if err := json.NewDecoder(f).Decode(&order); err != nil {
		t.Fatal(err)
	}

	want := Order{
		MerchantReference: &MerchantReference{
			OrderID1: "123456789",
			OrderID2: "123456789",
		},
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Cart: Cart{
			Items: []CartItem{
				{
					Reference:    "123456789",
					Name:         "Klarna t-shirt",
					Qty:          2,
					EAN:          "1234567890123",
					URI:          "http://example.com/product.php?123456789",
					ImageURI:     "http://example.com/product_image.php?123456789",
					UnitPrice:    12300,
					DiscountRate: 1000,
					TaxRate:      2500,
				},
				{
					Type:      "shipping_fee",
					Reference: "SHIPPING",
					Name:      "Shipping fee",
					Qty:       1,
					UnitPrice: 4900,
					TaxRate:   2500,
				},
			},
		},
		ShippingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		GUI: &GUI{
			Layout: "desktop",
		},
		Merchant: Merchant{
			ID:              "0",
			BackToStoreURI:  "http://example.com",
			TermsURI:        "http://example.com/terms.php",
			CheckoutURI:     "https://example.com/checkout.php",
			ConfirmationURI: "https://example.com/thankyou.php?sid=123&klarna_order={checkout.order.uri}",
			PushURI:         "https://example.com/push.php?sid=123&klarna_order={checkout.order.uri}",
		},
	}

	if !reflect.DeepEqual(order, want) {
		t.Errorf("\nwant = %+v;\nhave = %+v", want, order)
	}
}

func TestUnmarshal_ReadOrderResponse(t *testing.T) {
	f, err := os.Open("testdata/read_order_response.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var order Order
	if err := json.NewDecoder(f).Decode(&order); err != nil {
		t.Fatal(err)
	}

	want := Order{
		ID: "1A2B3C4D5E6F1A2B3C4D5E6F",
		MerchantReference: &MerchantReference{
			OrderID1: "123456789",
			OrderID2: "123456789",
		},
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Status:           "created",
		Reference:        "Q1W2E3",
		Reservation:      "123456789",
		StartedAt:        "2012-01-18T11:45:00+01:00",
		CompletedAt:      "2012-01-18T11:51:00+01:00",
		CreatedAt:        "2012-01-18T11:52:00+01:00",
		LastModifiedAt:   "2012-01-18T11:52:00+01:00",
		ExpiresAt:        "2012-02-01T11:52:00+01:00",
		ShippingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			CareOf:        "Testperson",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		BillingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			CareOf:        "Testperson",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		Cart: Cart{
			TotalPriceExclTax: 20280,
			TotalTaxAmount:    6760,
			TotalPriceInclTax: 27040,
			Items: []CartItem{
				{
					Reference:         "123456789",
					Name:              "Klarna t-shirt",
					Qty:               2,
					UnitPrice:         12300,
					DiscountRate:      1000,
					TaxRate:           2500,
					TotalPriceExclTax: 16605,
					TotalTaxAmount:    5535,
					TotalPriceInclTax: 22140,
				},
				{
					Type:              "shipping_fee",
					Reference:         "SHIPPING",
					Name:              "Shipping fee",
					Qty:               1,
					UnitPrice:         4900,
					TaxRate:           2500,
					TotalPriceExclTax: 3675,
					TotalTaxAmount:    1225,
					TotalPriceInclTax: 4900,
				},
			},
		},
		Customer: &Customer{
			Type: "person",
		},
		GUI: &GUI{
			Layout:  "desktop",
			Snippet: "...",
		},
		Merchant: Merchant{
			ID:              "0",
			TermsURI:        "http://example.com/terms.php",
			CheckoutURI:     "https://example.com/checkout.php",
			ConfirmationURI: "https://example.com/thankyou.php?sid=123&klarna_order={checkout.order.uri}",
			PushURI:         "https://example.com/push.php?sid=123&klarna_order={checkout.order.uri}",
		},
	}

	if !reflect.DeepEqual(order, want) {
		t.Errorf("\nwant = %+v;\nhave = %+v", want, order)
	}
}

func TestUnmarshal_UpdateOrderRequest(t *testing.T) {
	f, err := os.Open("testdata/update_order_request.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var order Order
	if err := json.NewDecoder(f).Decode(&order); err != nil {
		t.Fatal(err)
	}

	want := Order{
		MerchantReference: &MerchantReference{
			OrderID1: "123456789",
			OrderID2: "123456789",
		},
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Cart: Cart{
			Items: []CartItem{
				{
					Reference:    "123456789",
					Name:         "Klarna t-shirt",
					Qty:          4,
					EAN:          "1234567890123",
					URI:          "http://example.com/product.php?123456789",
					ImageURI:     "http://example.com/product_image.php?123456789",
					UnitPrice:    12300,
					DiscountRate: 1000,
					TaxRate:      2500,
				},
				{
					Type:      "shipping_fee",
					Reference: "SHIPPING",
					Name:      "Shipping fee",
					Qty:       1,
					UnitPrice: 4900,
					TaxRate:   2500,
				},
			},
		},
		ShippingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		GUI: &GUI{
			Layout: "desktop",
		},
		Merchant: Merchant{
			ID:              "0",
			TermsURI:        "http://example.com/terms.php",
			CheckoutURI:     "https://example.com/checkout.php",
			ConfirmationURI: "https://example.com/thankyou.php?sid=123&klarna_order={checkout.order.uri}",
			PushURI:         "https://example.com/push.php?sid=123&klarna_order={checkout.order.uri}",
		},
	}

	if !reflect.DeepEqual(order, want) {
		t.Errorf("\nwant = %+v;\nhave = %+v", want, order)
	}
}

func TestUnmarshal_UpdateOrderResponse(t *testing.T) {
	f, err := os.Open("testdata/update_order_response.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var order Order
	if err := json.NewDecoder(f).Decode(&order); err != nil {
		t.Fatal(err)
	}

	want := Order{
		ID: "1A2B3C4D5E6F1A2B3C4D5E6F",
		MerchantReference: &MerchantReference{
			OrderID1: "123456789",
			OrderID2: "123456789",
		},
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Status:           "checkout_incomplete",
		Reference:        "Q1W2E3",
		Reservation:      "123456789",
		StartedAt:        "2012-01-18T11:45:00+01:00",
		CompletedAt:      "2012-01-18T11:51:00+01:00",
		CreatedAt:        "2012-01-18T11:52:00+01:00",
		LastModifiedAt:   "2012-01-18T11:52:00+01:00",
		ExpiresAt:        "2012-02-01T11:52:00+01:00",
		ShippingAddress: &Address{
			GivenName:     "Testperson-se",
			FamilyName:    "Approved",
			CareOf:        "Testperson",
			StreetAddress: "Stårgatan 1",
			PostalCode:    "12345",
			City:          "Ankeborg",
			Country:       "se",
			Email:         "checkout@testdrive.klarna.com",
			Phone:         "0765260000",
		},
		Cart: Cart{
			TotalPriceExclTax: 36885,
			TotalTaxAmount:    12295,
			TotalPriceInclTax: 49180,
			Items: []CartItem{
				{
					Reference:         "123456789",
					Name:              "Klarna t-shirt",
					Qty:               4,
					UnitPrice:         12300,
					DiscountRate:      1000,
					TaxRate:           2500,
					TotalPriceExclTax: 33210,
					TotalTaxAmount:    11070,
					TotalPriceInclTax: 44280,
				},
				{
					Type:              "shipping_fee",
					Reference:         "SHIPPING",
					Name:              "Shipping fee",
					Qty:               1,
					UnitPrice:         4900,
					TaxRate:           2500,
					TotalPriceExclTax: 3675,
					TotalTaxAmount:    1225,
					TotalPriceInclTax: 4900,
				},
			},
		},
		Customer: &Customer{
			Type: "person",
		},
		GUI: &GUI{
			Layout:  "desktop",
			Snippet: "...",
		},
		Merchant: Merchant{
			ID:              "0",
			TermsURI:        "http://example.com/terms.php",
			CheckoutURI:     "https://example.com/checkout.php",
			ConfirmationURI: "https://example.com/thankyou.php?sid=123&klarna_order={checkout.order.uri}",
			PushURI:         "https://example.com/push.php?sid=123&klarna_order={checkout.order.uri}",
		},
	}

	if !reflect.DeepEqual(order, want) {
		t.Errorf("\nwant = %+v;\nhave = %+v", want, order)
	}
}

func normalizeJSON(r io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	json.Indent(&buf, b, "", "  ")

	return buf.Bytes(), nil
}
