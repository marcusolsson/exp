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

package main

import (
	"fmt"
	"log"

	"github.com/marcusolsson/kco-go"
)

func main() {
	var (
		merchantID   = "0"
		sharedSecret = "sharedSecret"
	)

	client := kco.NewAuthClient(sharedSecret, kco.TestEnvironmentURL)

	merchant := kco.Merchant{
		ID:              merchantID,
		TermsURI:        "http://example.com/terms.html",
		CheckoutURI:     "http://example.com/checkout",
		ConfirmationURI: "http://example.com/thank-you?klarna_order_id={checkout.order.id}",
		PushURI:         "http://example.com/push?klarna_order_id={checkout.order.id}",
	}

	cart := kco.Cart{
		Items: []kco.CartItem{
			{
				Qty:          1,
				Reference:    "123456789",
				Name:         "Klarna t-shirt",
				UnitPrice:    12300,
				DiscountRate: 1000,
				TaxRate:      2500,
			},
			{
				Qty:       1,
				Type:      "shipping_fee",
				Reference: "SHIPPING",
				Name:      "Shipping fee",
				UnitPrice: 4900,
				TaxRate:   2500,
			}},
	}

	newOrder := kco.Order{
		PurchaseCountry:  "se",
		PurchaseCurrency: "sek",
		Locale:           "sv-se",
		Cart:             cart,
		Merchant:         merchant,
	}

	orderID, err := client.CreateOrder(newOrder)
	if err != nil {
		log.Fatal(err)
	}

	order, err := client.Order(orderID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully created order", order.ID)
}
