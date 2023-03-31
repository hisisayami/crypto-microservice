package main

import (
	"crypto_service/data"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Currency struct {
	Id          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}

// 	"github.com/labstack/echo/v4" this is echo server for starting the server. this makes http request
//	handling easy and quite better than other libraries. This library used basic http package of golang to create its own implementation.

func main() {
	e := echo.New()
	e.GET("/currency/all", func(c echo.Context) error {
		output := map[string]interface{}{}
		data.CurrencyData.Range(func(key, value any) bool {
			output[key.(string)] = value
			return true
		})
		return c.JSON(http.StatusOK, output)
	})

	e.GET("/currency/:symbol", func(c echo.Context) error {
		param := c.Param("symbol")

		if data, ok := data.CurrencyData.Load(param); ok {
			return c.JSON(http.StatusOK, data)
		}
		return c.JSON(http.StatusNotFound, "symbol not found")

	})
	data.Configure()
	data.Fetch()
	e.Logger.Fatal(e.Start(":8088"))
}
