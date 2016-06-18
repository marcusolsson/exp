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
		sharedSecret = "sharedSecret"
		orderID      = "ABC123"
	)

	client := kco.NewAuthClient(sharedSecret, kco.TestEnvironmentURL)

	order, err := client.Order(orderID)
	if err != nil {
		log.Fatal(err)
	}

	order.Cart = kco.Cart{
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

	updatedOrder, err := client.UpdateOrder(orderID, order)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully updated order", updatedOrder.ID)
}
