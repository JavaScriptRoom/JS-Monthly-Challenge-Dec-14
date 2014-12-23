package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

func getNumPages() (int, error){
	doc, err := goquery.NewDocument("http://chat.stackoverflow.com/rooms/info/17/javascript/?tab=stars")
	if err != nil {
		return 0, err
	}
	text := doc.Find(".page-numbers.dots ~ * span.page-numbers").Text()
	numStr := strings.Split(text, " ")[0]
	numPagesStr := strings.TrimSpace(numStr)
	return strconv.Atoi(numPagesStr)
}

func scrapePage(stars map[string]int, j int){
	doc, err := goquery.NewDocument("http://chat.stackoverflow.com/rooms/info/17/javascript/?tab=stars&page=" + strconv.Itoa(j))
	if err != nil {
		fmt.Println("Error in page", j, ": ", err)
		return
	}

	doc.Find(".monologue").Each(func(i int, s *goquery.Selection){
		user := strings.TrimSpace(s.Find(".username").Text())
		starsText := strings.TrimSpace(s.Find(".times").Text())
		starsForPost, err :=  strconv.Atoi(starsText) // yolo
		if err != nil { // one star
			starsForPost = 1
		}
		stars[user] += starsForPost
	})

}

func main() {

	pagesNum, err := getNumPages()
	if err != nil {
		fmt.Println("Failed to get page numbers: ", err)
		return
	}
	stars := make(map[string]int)
	fmt.Println("Going to scrape ", pagesNum, " pages")
	for i := 0; i < pagesNum ; i++ {
		fmt.Println("Scraping page ", i)
		scrapePage(stars, i)
	}

	fmt.Println("Got here!")
	for key, value := range stars {
    	fmt.Println("Key:", key, "Value:", value)
	}
	
}
