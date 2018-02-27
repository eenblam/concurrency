package main

import (
    "fmt"
    "log"
    //"math/rand"
    "time"
)

func main() {
    /*TODO
    * This is wrong. I based this on the Dining Philosopher's wikipedia article:
    * https://en.wikipedia.org/wiki/Dining_philosophers_problem
    * before reading the paper https://www.cs.utexas.edu/users/misra/scannedPdf.dir/DrinkingPhil.pdf
    *
    * The former omits the latter's token passing strategy entirely.
    * It provides an unclear specification on how to proceed when initial requests
    * are received, since if I secure someone else's fork mine is dirty.
    * This behavior hasn't been patched in my implementation either.
    *
    * However, I think that "fixing" the algorithm to allow for initial forks
    * to be clean (and nonetheless passed from the initial state)
    * still allows for a deadlock if everyone finishes thinking at once
    * while forks are still evenly distributed.
    */

    n := 4
    // Create n philosophers
    ps := []*Philosopher{}
    for i := 0; i < n; i++ {
        p := &Philosopher{Id: i}
        ps = append(ps, p)
    }
    // Link them together *after* creating them all
    for i, p := range ps {
        prevIdx := i - 1
        if i == 0 {
            prevIdx = n - 1
        }
        nextIdx := i + 1
        if i == n - 1 {
            nextIdx = 0
        }
        prev := ps[prevIdx]
        next := ps[nextIdx]

        // Let neighbors request forks
        leftRequested := make(chan int)
        p.LeftRequested = leftRequested
        prev.RequestFromRight = leftRequested
        rightRequested := make(chan int)
        p.RightRequested = rightRequested
        next.RequestFromLeft = rightRequested

        // Send forks to neighbors
        passLeftChan := make(chan int)
        p.PassLeftChan = passLeftChan
        prev.ReceiveRightChan = passLeftChan
        passRightChan := make(chan int)
        p.PassRightChan = passRightChan
        next.ReceiveLeftChan = passRightChan

        // The other cases will be set by neighbors :)
    }

    // Run them all as goroutines
    for _,p := range ps {
        go p.Start()
    }

    var input string
    fmt.Scanln(&input)
}

// Can really drop the forks entirely with this model...
// ...but they do allow some better insight into state consistency.
type Fork int

const (
    Clean Fork = iota
    Dirty
    Empty
)

type Philosopher struct {
    Id int
    LeftRequested <-chan int
    RightRequested <-chan int
    PassLeftChan chan<- int
    PassRightChan chan<- int

    RequestFromLeft chan<- int
    RequestFromRight chan<- int
    ReceiveLeftChan <-chan int
    ReceiveRightChan <-chan int

    Left Fork
    Right Fork
}

func (p *Philosopher) Start() {
    for {
        p.Think()
        p.Hungry()
        p.Eat()
    }
}

func (p *Philosopher) Think() {
    log.Printf("Worker %d thinking", p.Id)
    //t := time.Duration(rand.Intn(10)) * time.Second
    t := time.Duration(p.Id) * time.Second
    timeout := time.After(t)
    for {
        select {
        case <-p.LeftRequested:
            log.Printf("Worker %d received request from left", p.Id)
            p.PassLeft()
        case <-p.RightRequested:
            log.Printf("Worker %d received request from right", p.Id)
            p.PassRight()
        case <-timeout:
            log.Printf("Worker %d done thinking", p.Id)
            break
        }
    }
    p.GetLeft()
    p.GetRight()
}

func (p *Philosopher) GetLeft() {
    log.Printf("Worker %d requesting left fork", p.Id)
    p.RequestFromLeft <- 0
}

func (p *Philosopher) GetRight() {
    log.Printf("Worker %d requesting right fork", p.Id)
    p.RequestFromRight <- 0
}

func (p *Philosopher) PassLeft() {
    log.Printf("Worker %d passing fork left", p.Id)
    if p.Left != Dirty {
        log.Printf("Error: left fork requested from worker %d but fork is %d", p.Id, p.Left)
    }
    p.Left = Empty
    p.PassLeftChan <- 0
}

func (p *Philosopher) PassRight() {
    log.Printf("Worker %d passing fork right", p.Id)
    if p.Right != Dirty {
        log.Printf("Error: right fork requested from worker %d but fork is %d", p.Id, p.Right)
    }
    p.Right = Empty
    p.PassRightChan <- 0
}

func (p *Philosopher) Eat() {
    if p.Left != Clean {
        log.Printf("Error: worker %d eating with missing or dirty left fork", p.Id)
    }
    if p.Right != Clean {
        log.Printf("Error: worker %d eating with missing or dirty right fork", p.Id)
    }
    p.Left = Dirty
    p.Right = Dirty
}

// Hungry is a state that traps the philosopher until both forks are clean
func (p *Philosopher) Hungry() {
    select {
    case <-p.ReceiveLeftChan:
        <-p.ReceiveRightChan
    case <-p.ReceiveRightChan:
        <-p.ReceiveLeftChan
    }
    p.Left = Clean
    p.Right = Clean
}
