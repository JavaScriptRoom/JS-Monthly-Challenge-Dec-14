package main

import "os"
import "fmt"
import "strconv"
import "sort"
import "math"
import "github.com/PuerkitoBio/goquery"

//import "golang.org/x/net/html"

var urlRoot string
var starCountForUser map[string]int

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected an integer room number")
		return
	}
	roomNum := os.Args[1]
	urlRoot = "http://chat.stackoverflow.com/rooms/info/" + roomNum

	starCountForUser = make(map[string]int)

	processPage(urlRoot + "?tab=stars")

	sortData()
}

func processPage(url string) {
	//fmt.Println("Processing: " + url)
	document, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Errorf("Failed to get/parse document: %s", err)
		return
	}

	document.Find(".monologue").Each(func(i int, message *goquery.Selection) {
		processMessage(message)
	})

	nextPage, hasNextPage := document.Find("a[rel=next]").Attr("href")
	if hasNextPage {
		processPage(urlRoot + nextPage)
	}

}

func processMessage(message *goquery.Selection) {
	usernameNode := message.Find(".username").Nodes[0]
	//How deep we have to look varies based on whether or not its an anonymous user
	for usernameNode.FirstChild != nil {
		usernameNode = usernameNode.FirstChild
	}
	username := usernameNode.Data
	starCount := 0
	//Unlikely that a single .monologue will contain multiple star counts... but just in case
	message.Find(".stars .times").Each(func(i int, times *goquery.Selection) {
		var count int
		if times.Nodes[0].FirstChild != nil {
			countStr := times.Nodes[0].FirstChild.Data
			count, _ = strconv.Atoi(countStr)
		} else {
			count = 1
		}
		//fmt.Printf("%d stars for %v\n", count, username)
		starCount += count
	})
	starCountForUser[username] += starCount
}

type StarData struct {
	Username  string
	StarCount int
}

type ByStarCount []StarData

func (a ByStarCount) Len() int           { return len(a) }
func (a ByStarCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStarCount) Less(i, j int) bool { return a[i].StarCount > a[j].StarCount }

func sortData() {
	sortData := make(ByStarCount, len(starCountForUser))
	var i = 0
	for user, count := range starCountForUser {
		data := StarData{user, count}
		sortData[i] = data
		i += 1
	}

	sort.Sort(sortData)

	width := int(math.Log10(float64(len(sortData)))) + 1

	for i, userData := range sortData {
		fmt.Printf("#%"+strconv.Itoa(width)+"v - %v has %v stars \n", i+1, userData.Username, userData.StarCount)
	}
}
