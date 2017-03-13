# currconv
currency converter in Golang via API calls according to

GET http://your-service/convert?amount={amount}&currency={currency}

EXAMPLE:

Request:

      http://your-service/convert?amount=200&currency=SEK

      Response:

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "amount": 200,
        "currency": "SEK",
        "converted": {
          "AUD": 32.6,
          "BGN": 42.25,
          "CAD": 32.2,
          "CNY": 158.87,
          "CZK": 583.86,
          "DKK": 160.77,
          "EUR": 21.6,
          "GBP": 17.44,
          ...
          "USD": 24.55
        }
      }

