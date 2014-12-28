package main

import (
    "fmt"
    "strconv"
    "regexp"

    "sort"

    // used with some sorrow, because x/net/html is just not usable atm
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
    starCount, _ := strconv.ParseInt(message.Find(".times").Text(), 10, 64)
    if starCount == 0 {
        starCount = 1
    }
    // fmt.Println("Star count:", starCount)

    userContainer := message.Find(".username a").Eq(0)
    href, _ := userContainer.Attr("href")

    username := userContainer.Text()
    userid, _ := strconv.ParseInt(useridRegexp.FindString(href), 10, 64)

    // fmt.Println("Username:", username)
    // fmt.Println("Userid:", userid)

    return &Humanoid {userid, username, starCount}
}

func main() {
    c := make(chan *Humanoid)

    go ExtractHumanoidsFromUrl("http://chat.stackoverflow.com/rooms/info/17/javascript/?tab=stars", c)

    humans := make(map[int64]*Humanoid)

    for humanoid := range c {
        fmt.Println("Humanoid!", humanoid)

        existingHumanoid, existed := humans[humanoid.id]

        if existed {
            existingHumanoid.starCount += humanoid.starCount
        } else {
            humans[humanoid.id] = humanoid
        }
    }

    fmt.Println("=====================")
    fmt.Println("Done!")

    sortedHumanoids := sortHumanoids(humans)
    for _, humanoid := range sortedHumanoids {
        fmt.Println(humanoid)
    }
}
