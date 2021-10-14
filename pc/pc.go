package main

import (
	"fmt"
	"github.com/ChrisGora/semaphore"
	"math/rand"
	"sync"
	"time"
)

type buffer struct {
	b                 []int
	size, read, write int
}

func newBuffer(size int) buffer {
	return buffer{
		b:     make([]int, size),
		size:  size,
		read:  0,
		write: 0,
	}
}

func (buffer *buffer) get() int {
	x := buffer.b[buffer.read]
	fmt.Println("Get\t", x, "\t", buffer)
	buffer.read = (buffer.read + 1) % len(buffer.b)
	return x
}

func (buffer *buffer) put(x int) {
	buffer.b[buffer.write] = x
	fmt.Println("Put\t", x, "\t", buffer)
	buffer.write = (buffer.write + 1) % len(buffer.b)
}
//Your solution must not use channels.
//If you could use channels the solution would be t
func producer(buffer *buffer,start, delta int,mutex *sync.Mutex,space semaphore.Semaphore,work semaphore.Semaphore) {
	x:=start
	for{
		space.Wait()
		mutex.Lock()
		buffer.put(x)
		x = x + delta
		work.Post()
		mutex.Unlock()
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}


}



func consumer(buffer *buffer,mutex *sync.Mutex,space semaphore.Semaphore,work semaphore.Semaphore) {
	for {
		work.Wait()
		mutex.Lock()
		_ = buffer.get()
		space.Post()
		mutex.Unlock()
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	buffer := newBuffer(5)

	var mutex = sync.Mutex{}
	work := semaphore.Init(5,0)
	space := semaphore.Init(5,5)

	go producer(&buffer, 1, 1,&mutex,space,work)
	go producer(&buffer, 1000, -1,&mutex,space,work)

	consumer(&buffer,&mutex,space,work)
}


// another version

/*
mutex.lock()
if (!condition){
	mutex.unlock()
 */
