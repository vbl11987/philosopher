package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Philosopher struct {
	number          int
	chopsticks      chan bool
	nextPhilosopher *Philosopher
	timesEaten      int
}

func (p Philosopher) waiting() {
	fmt.Printf("%v is waiting\n", p.number)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (p Philosopher) eating() {
	fmt.Printf("starting to eat %v\n", p.number)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (p Philosopher) getChops() {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1e9)
		timeout <- true
	}()

	<-p.chopsticks
	fmt.Printf("%v took his chopsticks\n", p.number)
	select {
	case <-p.nextPhilosopher.chopsticks:
		fmt.Printf("%v got %v's chopsticks\n", p.number, p.nextPhilosopher.number)
		return
	case <-timeout:
		p.chopsticks <- true
		p.waiting()
		p.getChops()
	}
}

func (p *Philosopher) returnChops() {
	// when a philosopher is done with his chopsticks
	// then he returns his and the nextPhilosopher
	fmt.Printf("%v is returning his chopsticks and %v's chopsticks\n", p.number, p.nextPhilosopher.number)
	p.chopsticks <- true
	p.nextPhilosopher.chopsticks <- true
}

func (p *Philosopher) haveDinner(philChan chan *Philosopher) {
	p.waiting()
	p.getChops()
	p.eating()
	p.returnChops()
	if p.timesEaten == 3 {
		philChan <- p
		return
	}
	p.timesEaten++
	p.haveDinner(philChan)
}

func main() {
	philosopers := make([]*Philosopher, 5)
	var p *Philosopher
	//creating the list of Philosophers
	for i := range [5]int{} {
		c := make(chan bool, 1)
		c <- true
		p = &Philosopher{i, c, p, 0}
		philosopers[i] = p
	}
	philosopers[0].nextPhilosopher = p

	philChan := make(chan *Philosopher)
	//created the list of philosopers
	//executing the function to haveDinner

	for _, phil := range philosopers {
		//concurrency
		go phil.haveDinner(philChan)
	}

	for i := 0; i < len(philosopers); i++ {
		p := <-philChan
		fmt.Printf("Philosopher %v is done\n", p.number)
	}
	fmt.Println("closing the application")
}
