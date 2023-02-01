package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
)

type Shop struct {
	ShopName string `csv:"shop_name"`
	URL      string `csv:"url"`
}

func Scrape() []Shop {
	// 読み込むページ
	var Shops []Shop
	for i := 1; i <= 26; i++ {
		fmt.Println("")
		fmt.Printf("page %d start", i)
		fmt.Println("")
		res, err := http.Get(fmt.Sprintf("https://monobito.com/all_shop/all/1/?p=%d&&show=s", i))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Bodyを読み込む
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(".cf").Each(func(i int, s *goquery.Selection) {
			var shop Shop
			shopName := s.Find(".info").Find("p").Find("a").Text()
			shopUrl, isExist := s.Find(".img").Find("a").Attr("href")
			if !isExist {
				shopUrl = ""
			}
			shop = Shop{
				ShopName: shopName,
				URL:      fmt.Sprintf("https://monobito.com%s", shopUrl),
			}
			Shops = append(Shops, shop)
		})
	}
	return Shops
}

func main() {
	//contents
	fmt.Println("start!")
	c := Scrape()
	f, err := os.Create("shop3.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gocsv.MarshalFile(c, f)
	fmt.Println("successfully ended!")
}
