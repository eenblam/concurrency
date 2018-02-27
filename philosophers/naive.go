package main

import (
    "fmt"
    "log"
    "math/rand"
    "sync"
    "time"
)

func main() {
    // Initialize
    n := 10
    forks := []*Fork{}
    for i := 0; i < n; i++ {
        forks = append(forks, &Fork{Id: i})
    }
    philosophers := []*Philosopher{}
    for i := 0; i < n; i++ {
        p := Philosopher{Id: i, left: forks[i]}
        if i == n-1 {
            p.right = forks[0]
        } else {
            p.right = forks[i+1]
        }
        philosophers = append(philosophers, &p)
    }
    fmt.Println(philosophers)

    // Kick off workers
    for _,p := range philosophers {
        p.Run()
    }

    // Hang program while workers run
    var input string
    fmt.Scanln(&input)
}

type Philosopher struct {
    Id    int
    left  *Fork
    right *Fork
}

type Fork struct {
    mux sync.Mutex
    Id  int
}

func (p *Philosopher) GetLeft() {
    //p.output <- fmt.Sprintf("Worker %d awaiting left fork %d", p.Id, p.left.Id)
    log.Printf("Worker %d awaiting left fork %d", p.Id, p.left.Id)
    p.left.mux.Lock()
    log.Printf("Worker %d obtained left fork %d", p.Id, p.left.Id)
}

func (p *Philosopher) GetRight() {
    log.Printf("Worker %d awaiting right fork %d", p.Id, p.right.Id)
    p.right.mux.Lock()
    log.Printf("Worker %d obtained right fork %d", p.Id, p.right.Id)
}

func (p *Philosopher) DropLeft() {
    log.Printf("Worker %d releasing left fork %d", p.Id, p.left.Id)
    p.left.mux.Unlock()
    log.Printf("Worker %d released left fork %d", p.Id, p.left.Id)
}

func (p *Philosopher) DropRight() {
    log.Printf("Worker %d releasing right fork %d", p.Id, p.right.Id)
    p.right.mux.Unlock()
    log.Printf("Worker %d released right fork %d", p.Id, p.right.Id)
}

func (p *Philosopher) Wait() {
    log.Printf("Worker %d sleeping", p.Id)
    n := rand.Intn(500)
    t := time.Duration(n) * time.Millisecond
    time.Sleep(t)
    log.Printf("Worker %d awake", p.Id)
}

func (p *Philosopher) Run() {
    // Naive algorithm; deadlocks
    go func() {
        log.Printf("Worker %d start", p.Id)
        for {
            p.Wait()
            p.GetLeft()
            p.GetRight()
            p.Wait()
            p.DropRight()
            p.DropLeft()
        }
    }()
}
