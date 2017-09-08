package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"log"
	"net/url"
	"strings"
	"time"
)

var mangaUrl = ""

func main() {
	u, err := url.Parse(mangaUrl)
	if err != nil {
		log.Fatalf("Failed to parse url:%v", err)
	}
	baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	fmt.Println(mangaUrl)
	if err := page.Navigate(mangaUrl); err != nil {
		log.Fatalf("Failed to navigate:%v", err)
	}
	time.Sleep(time.Second * 5)
	html, _ := page.HTML()
	reader := strings.NewReader(string(html))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatalf("Failed to parse html:%v", err)
	}
	urls := []string{}
	doc.Find("#leftside table.listing tr td a").Each(func(i int, link *goquery.Selection) {
		url, exists := link.Attr("href")
		if exists {
			urls = append(urls, url)
		}
	})
	for i := len(urls) - 1; i >= 0; i-- {
		fmt.Printf("\"%s%s\",\n", baseUrl, urls[i])
	}
	fmt.Println(len(urls))
}
