package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/zserge/lorca"
)

var siteURL string

type cookie struct {
	Domain  string  `json:"domain"`
	Name    string  `json:"name"`
	Value   string  `json:"value"`
	Expires float64 `json:"expires"`
}

func main() {
	flag.StringVar(&siteURL, "url", "", "Site URL")
	flag.Parse()

	if siteURL == "" {
		log.Fatal("url is not provided")
	}
	fmt.Printf("Site URL: %s\n", siteURL)

	startURL := "data:text/html,<html><head><title>Connecting to site: " + siteURL + "</title><head></html>"
	ui, err := lorca.New(startURL, "", 480, 430)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Load(siteURL)

	go func() {
		currentURL := ""
		for strings.ToLower(currentURL) != strings.ToLower(siteURL) {
			newURL := ui.Eval("window.location.href").String()
			if currentURL != newURL {
				fmt.Printf("%s\n", newURL)
				currentURL = newURL
			}
			time.Sleep(500 * time.Microsecond)
		}
		cookies := ui.Send("Network.getCookies", nil)
		if cookies.Err() != nil {
			fmt.Printf("%s\n", cookies.Err())
			return
		}
		for _, c := range cookies.Object()["cookies"].Array() {
			cc := &cookie{}
			c.To(&cc)
			fmt.Printf("%#v\n", cc)

			sec, dec := math.Modf(cc.Expires)
			expireTime := time.Unix(int64(sec), int64(dec*(1e9)))
			fmt.Printf("%s\n", expireTime)
		}
		ui.Close()
	}()

	<-ui.Done()
}
