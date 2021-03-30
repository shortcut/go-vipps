# THIS IS A [FORK](https://github.com/torfjor/go-vipps)
Please keep this in mind if something seems wonky.

# Go Vipps
Community maintained Go client library for the [Vipps](https://vipps.no) E-commerce and Recurring payments APIs. Please see Vipps' own documentation on their [Developer page](https://vipps.no/developer/).

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/torfjor/go-vipps)

## Installation

Install go-vipps with:

```sh
go get -u github.com/shortcut/go-vipps
```

Then, import it using:

``` go
import (
    "github.com/shortcut/go-vipps"
)
```

## Usage

```go
package main

import (
	"context"
	"log"
	"os"
	"time"
	"github.com/shortcut/go-vipps"
	"github.com/shortcut/go-vipps/auth"
	"github.com/shortcut/go-vipps/ecom"
)

func main() {
	credentials := vipps.Credentials{
		ClientID:           os.Getenv("CLIENT_ID"),
		ClientSecret:       os.Getenv("CLIENT_SECRET"),
		APISubscriptionKey: os.Getenv("API_KEY"),
	}
	env := vipps.EnvironmentTesting
	authClient := auth.NewClient(env, credentials)
	client := ecom.NewClient(vipps.ClientConfig{
		HTTPClient: authClient,
		Environment: env,
	})
	
	mobileNumber := 97777776
	amount := 1000
	orderID := "8b84-0ad5258beb0f"
	transactionText := "A transaction"
	
	cmd := ecom.InitiatePaymentCommand{
		MerchantInfo: ecom.MerchantInfo{
			MerchantSerialNumber: "CHANGETHIS",
			CallbackURL:          "https://some.endpoint.no/callbacks",
			RedirectURL:          "https://some.endpoint.no/redirect",
			ConsentRemovalURL:    "https://some.endpoint.no/consentremoval",
			IsApp:                false,
			PaymentType:          ecom.PaymentTypeRegular,
		},
		CustomerInfo: ecom.CustomerInfo{
			MobileNumber: mobileNumber,
		},
		Transaction:  ecom.Transaction{
			Amount: amount,
			OrderID: orderID,
			TransactionText: transactionText,
		},
	}
	
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	p, err := client.InitiatePayment(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}
	// Do something with p
}

```
