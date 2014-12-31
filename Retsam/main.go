package main

import "os"
import "fmt"
import "strconv"
import "sort"
import "math"
import "strings"
import "github.com/PuerkitoBio/goquery"

//import "golang.org/x/net/html"

var urlRoot string
var starCountForUserId map[string]int
var usernameForUserId map[string]string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected an integer room number")
		return
	}
	roomNum := os.Args[1]
	urlRoot = "http://chat.stackoverflow.com/rooms/info/" + roomNum

	starCountForUserId = make(map[string]int)
	usernameForUserId = make(map[string]string)

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
	userId := getUserId(message)
	username := getUserName(message)
	starCount := getStarCount(message)

	if _, hasName := usernameForUserId[userId]; !hasName {
		usernameForUserId[userId] = username
	}

	starCountForUserId[userId] += starCount
}

func getUserId(message *goquery.Selection) string {
	var userId string
	for _, attr := range message.Nodes[0].Attr {
		if attr.Key == "class" {
			for _, class := range strings.Split(attr.Val, " ") {
				if strings.HasPrefix(class, "user-") {
					userId = class
					break
				}
			}
			break
		}
	}

	if userId == "" {
		panic("Couldn't find an user id")
	}

	return userId
}

func getUserName(message *goquery.Selection) string {
	usernameNode := message.Find(".username").Nodes[0]

	//How deep we have to look varies based on whether or not its an anonymous user
	for usernameNode.FirstChild != nil {
		usernameNode = usernameNode.FirstChild
	}
	return usernameNode.Data
}

func getStarCount(message *goquery.Selection) int {
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
	return starCount
}

type StarData struct {
	UserId    string
	StarCount int
}

type ByStarCount []StarData

func (a ByStarCount) Len() int           { return len(a) }
func (a ByStarCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStarCount) Less(i, j int) bool { return a[i].StarCount > a[j].StarCount }

func sortData() {
	sortData := make(ByStarCount, len(starCountForUserId))
	var i = 0
	for userId, count := range starCountForUserId {
		data := StarData{userId, count}
		sortData[i] = data
		i += 1
	}

	sort.Sort(sortData)

	width := int(math.Log10(float64(len(sortData)))) + 1

	for i, userData := range sortData {
		username := usernameForUserId[userData.UserId]
		fmt.Printf("#%"+strconv.Itoa(width)+"v - %v has %v stars \n", i+1, username, userData.StarCount)
	}
}
