package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rgamba/evtwebsocket"
)

// 	"github.com/rgamba/evtwebsocket" this library is used for easy socket connection for fetching and updating the data,
//	this library  has used websocket package provided by go and optimized for easy usage.

type StockCurrency struct {
	Type              string `json:"type"`
	BaseCurrency      string `json:"base_currency"`
	QuoteCurrency     string `json:"quote_currency"`
	Status            string `json:"status"`
	QuantityIncrement string `json:"quantity_increment"`
	TickSize          string `json:"tick_size"`
	TakeRate          string `json:"take_rate"`
	MakeRate          string `json:"make_rate"`
	FeeCurrency       string `json:"fee_currency"`
}

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

type TicketNotification struct {
	Ch   string            `json:"ch"`
	Data map[string]Ticker `json:"data"`
}
type Ticker struct {
	Timestamp          int64  `json:"t"`
	BestAsk            string `json:"a"`
	BestAskQuantity    string `json:"A"`
	BestBid            string `json:"b"`
	BestBidQuantity    string `json:"B"`
	LastPrice          string `json:"c"`
	OpenPrice          string `json:"o"`
	HighPrice          string `json:"h"`
	LowPrice           string `json:"l"`
	BaseAssetVolume    string `json:"v"`
	QuoteAssetVolume   string `json:"q"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	LastTradeID        int64  `json:"L"`
}

type SymbolName struct {
	FullName string `json:"full_name"`
}

var CurrencyData = sync.Map{}
var Symbol = []string{"ETHBTC", "BTCUSDT"}

func Configure() {
	for _, s := range Symbol {
		data := GetData("https://api.hitbtc.com/api/3/public/symbol/", s)
		name := GetPrice("https://api.hitbtc.com/api/3/public/currency/", data.BaseCurrency)
		CurrencyData.Store(s, Currency{
			Id:          data.BaseCurrency,
			FeeCurrency: data.FeeCurrency,
			FullName:    name.FullName,
		})
	}

}

func Fetch() {
	c := evtwebsocket.Conn{
		OnConnected: func(conn *evtwebsocket.Conn) {
			fmt.Println("Connected")
		},
		OnMessage: func(msg []byte, conn *evtwebsocket.Conn) {
			fmt.Printf("Received message: %s\n", msg)
			data := TicketNotification{}
			err := json.Unmarshal(msg, &data)
			if err == nil {
				for _, sym := range Symbol {
					if d, ok := data.Data[sym]; ok {
						if v, ok := CurrencyData.Load(sym); ok {
							v1 := v.(Currency)
							v1.Ask = d.BestAsk
							v1.Last = d.LastPrice
							v1.High = d.HighPrice
							v1.Low = d.LowPrice
							v1.Bid = d.BestBid
							v1.Open = d.OpenPrice
							CurrencyData.Store(sym, v1)
						}
					}
				}
			}
		},
		OnError: func(err error) {
			fmt.Printf("** ERROR **\n%s\n", err.Error())
		},
		PingIntervalSecs: 10,
	}
	// Connect
	err := c.Dial("wss://api.hitbtc.com/api/3/ws/public", "")
	if err != nil {
		fmt.Println(err)
	}

	msg := evtwebsocket.Msg{
		Body: []byte(`{
			"method": "subscribe",
			"ch": "ticker/1s",
			"params": {
				"symbols": ["ETHBTC","BTCUSDT"]
			},
			"id": 123
		}`),
	}
	c.Send(msg)
}

func GetData(url, symbol string) StockCurrency {
	data := MakeRequest(url + symbol)
	s := StockCurrency{}
	json.Unmarshal(data, &s)
	return s
}

func GetPrice(url, symbol string) SymbolName {
	data := MakeRequest(url + symbol)
	s := SymbolName{}
	json.Unmarshal(data, &s)
	return s
}

func MakeRequest(URL string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Add("Accept-Encoding", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error::", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	return resBody
}
