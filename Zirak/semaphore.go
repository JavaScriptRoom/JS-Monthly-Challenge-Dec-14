package main

type SemNone struct {}
type Semaphore chan SemNone

func (sem Semaphore) Acquire () {
    sem <- SemNone{}
}

func (sem Semaphore) Release () {
    <- sem
}
