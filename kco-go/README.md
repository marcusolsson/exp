# Klarna Checkout Go SDK

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/marcusolsson/kco-go)
[![License MIT](https://img.shields.io/badge/license-apache-lightgrey.svg?style=flat)](LICENSE)

This is an __unofficial__ client package for accessing the Klarna Checkout API.

__WARNING:__ This package is under heavy development and currently __NOT__ production-ready. You would be better off using any of the official SDKs available at this point.

## Example

```go
func main() {
	var (
		merchantID   = "your_merchant_id"
		sharedSecret = "your_shared_secret"
	)

	client := kco.NewAuthClient(sharedSecret)

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
```

## Documentation

Additional documentation can be found at [https://developers.klarna.com](https://developers.klarna.com).

## License

Klarna Checkout Go SDK is licensed under [Apache License, Version 2.0](LICENSE)
