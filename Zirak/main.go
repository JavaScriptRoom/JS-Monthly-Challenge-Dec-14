package main

import (
    "fmt"
)

func main() {
    roomid := 17

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
    for _, humanoid := range sortedHumanoids {
        fmt.Println(humanoid)
    }
}

func CollectHumanoids (sem Semaphore, url string, humans map[uint]*Humanoid) {
    c := make(chan *Humanoid)

    fmt.Println("Collecting humanoids from", url)
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

    fmt.Println("Finished", url)
    sem.Release()
}
