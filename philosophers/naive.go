package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func main() {
    // Initialize
    messages := make(chan string, 10)
    n := 2
    forks := []*Fork{}
    for i := 0; i < n; i++ {
        forks = append(forks, &Fork{Id: i})
    }
    philosophers := []*Philosopher{}
    for i := 0; i < n; i++ {
        p := Philosopher{Id: i, left: forks[i], output: messages}
        if i == n-1 {
            p.right = forks[0]
        } else {
            p.right = forks[i+1]
        }
        philosophers = append(philosophers, &p)
    }
    fmt.Println(philosophers)

    // Kick off receiver
    go func(messages <-chan string) {
        for m := range messages {
            fmt.Println(m)
        }
    }(messages)

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
    output chan<- string
}

type Fork struct {
    mux sync.Mutex
    Id  int
}

func (p *Philosopher) GetLeft() {
    p.output <- fmt.Sprintf("Worker %d awaiting left fork %d", p.Id, p.left.Id)
    p.left.mux.Lock()
    p.output <- fmt.Sprintf("Worker %d obtained left fork %d", p.Id, p.left.Id)
}

func (p *Philosopher) GetRight() {
    p.output <- fmt.Sprintf("Worker %d awaiting right fork %d", p.Id, p.right.Id)
    p.right.mux.Lock()
    p.output <- fmt.Sprintf("Worker %d obtained right fork %d", p.Id, p.right.Id)
}

func (p *Philosopher) DropLeft() {
    p.output <- fmt.Sprintf("Worker %d releasing left fork %d", p.Id, p.left.Id)
    p.left.mux.Unlock()
    p.output <- fmt.Sprintf("Worker %d released left fork %d", p.Id, p.left.Id)
}

func (p *Philosopher) DropRight() {
    p.output <- fmt.Sprintf("Worker %d releasing right fork %d", p.Id, p.right.Id)
    p.right.mux.Unlock()
    p.output <- fmt.Sprintf("Worker %d released right fork %d", p.Id, p.right.Id)
}

func (p *Philosopher) Wait() {
    p.output <- fmt.Sprintf("Worker %d sleeping", p.Id)
    n := rand.Intn(500)
    t := time.Duration(n) * time.Millisecond
    time.Sleep(t)
    p.output <- fmt.Sprintf("Worker %d awake", p.Id)
}

func (p *Philosopher) Run() {
    // Naive algorithm; deadlocks
    go func() {
        p.output <- fmt.Sprintf("Worker %d start", p.Id)
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
