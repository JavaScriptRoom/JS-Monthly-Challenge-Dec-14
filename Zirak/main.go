package main

import (
    "fmt"

    "os"
    "flag"

    "strconv"

    "github.com/olekukonko/tablewriter"
)

func main() {
    var roomid int
    flag.IntVar(&roomid, "roomid", 17, "Room id to grab stars from")

    flag.Parse()

    fmt.Println("Fetching page count...")
    pagesCount := CountStarPages(roomid)
    fmt.Println("Count:", pagesCount)

    humans := make(map[uint]*Humanoid)

    // Make sure we only make 4 requests at the same time, so we won't spam SO
    sem := make(Semaphore, 4)

    for p := 0; p < pagesCount; p += 1 {
        sem.Acquire()

        url := FormatStarUrl(roomid, p)
        go CollectHumanoids(sem, url, humans)
    }

    fmt.Println("=====================")
    fmt.Println("Done!")

    sortedHumanoids := SortHumanoids(humans)
    table := tablewriter.NewWriter(os.Stdout)

    table.SetHeader([]string{"uid", "Name", "Star count"})
    table.SetAlignment(tablewriter.ALIGN_LEFT)

    for _, humanoid := range sortedHumanoids {
        // fmt.Println(humanoid)
        table.Append([]string{
            strconv.FormatUint(uint64(humanoid.id), 10),
            humanoid.name,
            strconv.FormatUint(uint64(humanoid.starCount), 10),
        })
    }

    table.Render()
}

func CollectHumanoids (sem Semaphore, url string, humans map[uint]*Humanoid) {
    c := make(chan *Humanoid)

    fmt.Fprintln(os.Stderr, "Collecting humanoids from", url)
    go ExtractHumanoidsFromUrl(url, c)

    for humanoid := range c {
        // fmt.Println("Humanoid!", humanoid)

        existingHumanoid, existed := humans[humanoid.id]

        if existed {
            existingHumanoid.starCount += humanoid.starCount
        } else {
            humans[humanoid.id] = humanoid
        }
    }

    fmt.Fprintln(os.Stderr, "Finished", url)
    sem.Release()
}
