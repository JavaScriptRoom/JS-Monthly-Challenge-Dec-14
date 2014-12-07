package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

func main() {

	var sitep = flag.String("site", "stackoverflow", "chat site's base url")
	var roomp = flag.Int("roomid", -1, "id of targeted chatroom")

	flag.Parse()

	var site = *sitep
	var room = *roomp

	log.Printf("Targeting Site : %s \n", site)

	if room < 0 {

		var roomselectfmt string = "http://chat.%s.com/rooms/?tab=all&sort=people"

		target := fmt.Sprintf(roomselectfmt, site)

		log.Printf("Attempting to resolve target room id\n")
		log.Printf("Target Endpoint : %s\n", target)

		findPage, err := goquery.NewDocument(target)

		if err != nil {
			log.Fatal(err)
		}

		var roomidquery = "[id^=\"room-\"]"

		var node = findPage.Find(roomidquery).First()

		if node == nil {
			log.Fatal("no matching rooms found at the requested page")
		}

		var roomidattr, exists = node.Attr("id")

		if !exists {
			log.Fatal("no matching ref")
		}

		var idsplit = strings.Split(roomidattr, "-")

		room, err = strconv.Atoi(idsplit[1])

		if err != nil {
			log.Fatal(err)
		}

	}

	log.Printf("Targeting Room : %d \n", room)

	baseurl := "http://chat.%s.com/rooms/info/%d/%s"

	currurl := fmt.Sprintf(baseurl, site, room, "?tab=stars")

	log.Printf("Targeting Endpoint :%s", currurl)

	currpage, err := goquery.NewDocument(currurl)

	if err != nil {
		log.Fatal(err)
	}

	storagep := MakePersonMap()

	var lasttext = currpage.Find("div.pager span.page-numbers.dots + a span.page-numbers").First().Text()
	var last, errlast = strconv.Atoi(lasttext)

	if errlast != nil {
		log.Fatalf("Unable to determine last index")
	}
	var index = 1

	log.Printf("Parsing %d Pages\n", last)

	for {

		msgquery := "#content > .monologue"
		starquery := " .messages .flash > .stars > .times"
		userquery := ".signature > .tiny-signature > .username > a"

		currpage.Find(msgquery).Each(func(i int, s *goquery.Selection) {
			starvalue := 1

			star := s.Find(starquery).Text()

			if star != "" {
				starvalue, err = strconv.Atoi(star)

				if err != nil {
					log.Fatal(err)
				}
			}

			username := s.Find(userquery).Text()

			_, ok := storagep[username]

			if !ok {
				storagep[username] = &Person{name: username, stars: 0}
			}

			storagep[username].stars += starvalue
		})

		index++

		if index > last {
			log.Printf("finished traversing lists")
			break
		}

		nextlnk := fmt.Sprintf("?tab=stars&page=%d", index)
		currurl = fmt.Sprintf(baseurl, site, room, nextlnk)

		currpage, err = goquery.NewDocument(currurl)

		if err != nil {
			log.Fatal(err)
			break
		}

	}

	var values []*Person

	for key := range storagep {
		values = append(values, storagep[key])
	}

	stars := func(person1, person2 *Person) bool {
		
		if person1.stars > person2.stars {
			return true
		}
		
		if person1.stars == person2.stars && person1.name < person2.name {
			return true
		} 
		
		return false
	}

	By(stars).Sort(values)

	var jsonppl []string

	for _, person := range values {
		jsonppl = append(jsonppl, fmt.Sprintf("{ \"username\" : \"%s\" , \"stars\" : %d }", person.name, person.stars))
	}

	jsonresult := strings.Join(jsonppl, ",\n\t")

	log.Print("[\n\t", jsonresult, "\n]")
}
