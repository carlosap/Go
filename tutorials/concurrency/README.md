# Concurrency

* creating multiple processes that execute independently (not simultaneusly)
* goroutine are schedule by the go runtime.
* wait smaller than threads
* Go Manages goroutines
* Less switching
* Faster start-up times
* Safe Communication
* Go uses the 'Actor Model' communicating sequential processes (CSP)
* GOROUTINE   ---->>>Channel<<<<------GOROUTINE (channel are like pipes)
* Concurrency vs. Parallelism (not the same - but related)- parallelism is simultaneous execution

### Channels

* buffered:      myChannel := make(chan int, 5) given a size (buffers *5).Gorutines drop the data and buffers (queue). Async
* unbuffered:    myChannel := make(chan int) Unbuffered channels wait for another gorutine to be ready to collect. Sync


```javascript
package main
import (
	"fmt"
	"time"
	"sync"
	"runtime" 
)

func main(){

	//enables parallelism with the runtime lib
	runtime.GOMAXPROCS(2) 
	//Sync
	var waitGrp sync.WaitGroup
	waitGrp.Add(2)

	//Example of blocking functions
	func(){
		time.Sleep(5 * time.Second)
		fmt.Println("Hello")
	}()


	func(){
		fmt.Println("Hello2")
	}()

	//Adding Goroutines. Adding the key word "go"
	//creates a concurrency program
	go func(){
		//this function will be block but it will "switch" to 
		//the one below
		defer waitGrp.Done()
		time.Sleep(5 * time.Second)
		fmt.Println("goroutine1")
	}()

	//this will print first
	go func(){
		defer waitGrp.Done()
		fmt.Println("goroutine2")
	}()

	//ensures all concurrencies are completed 
	waitGrp.Wait()
}

```