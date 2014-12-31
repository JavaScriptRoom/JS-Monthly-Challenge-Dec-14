package main

import (
    "fmt"
    "strconv"
    "regexp"

    "github.com/PuerkitoBio/goquery"
)

func ExtractHumanoidsFromUrl (url string, c chan *Humanoid) {
    fmt.Println("Navigating to", url)
    doc, err := goquery.NewDocument(url)
    fmt.Println("Done!")

    if err != nil {
        fmt.Println("Oh noes, error!", err)
        return
    }

    ExtractHumanoidsFromDoc(doc, c)
}

func ExtractHumanoidsFromDoc (doc *goquery.Document, c chan *Humanoid) {
    fmt.Println("Extracting stars from document")
    messages := doc.Find(".monologue")

    for i := range messages.Nodes {
        c <- ExtractHumanoidFromMessage(messages.Eq(i))
    }

    // XXX is this the correct place to close the channel?
    close(c)
}

func ExtractHumanoidFromMessage (message *goquery.Selection) (*Humanoid) {
    useridRegexp := regexp.MustCompile(`(\d+)`)

    // XXX yeah...handle error
    starCount, _ := strconv.ParseUint(message.Find(".times").Text(), 10, 64)
    if starCount == 0 {
        starCount = 1
    }
    // fmt.Println("Star count:", starCount)

    userContainer := message.Find(".username a").Eq(0)
    href, _ := userContainer.Attr("href")

    username := userContainer.Text()
    userid, _ := strconv.ParseUint(useridRegexp.FindString(href), 10, 64)

    // fmt.Println("Username:", username)
    // fmt.Println("Userid:", userid)

    return &Humanoid {uint(userid), username, uint(starCount)}
}

func CountStarPages (roomid int) int {
    doc,   _ := goquery.NewDocument(FormatStarUrl(roomid, 1))
    pages, _ := strconv.ParseInt(doc.Find(".page-numbers.dots").Eq(1).Next().Text(), 10, 64)

    return int(pages)
}

func FormatStarUrl (roomid int, page int) string {
    base := "http://chat.stackoverflow.com/rooms/info/%d/javascript/?tab=stars&page=%d"

    return fmt.Sprintf(base, roomid, page)
}
