package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/db/cache"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Request struct {
	Name         string
	Url          string
	Method       string
	Payload      string
	Header       http.Header
}
type Requests []Request

type RequestMethods struct {
	POST string
	GET string
}

func (r Requests) Execute() []string {
	//start := time.Now()
	//secs := time.Since(start).Seconds()
	//fmt.Printf("Number of CPU %d and parallel instances: %d\n", cpus, parallel)
	var errorLock sync.Mutex
	var updateLock sync.Mutex
	errors := make([]string, 0)
	var wg sync.WaitGroup
	semaphore := make(chan int, parallel)
	c := &cache.Cache{}
	//ch := make(chan string)
	for _, request := range r {
		wg.Add(1)
		go func(r Request) {
			defer wg.Done()
			semaphore <- 1
			c.Delete(r.Name)
			//fmt.Printf("request url: %s\n", r.Url)
			payload := strings.NewReader(r.Payload)
			client := &http.Client{}
			req, err := http.NewRequest(r.Method, r.Url, payload)
			if err != nil {
				errorLock.Lock()
				defer errorLock.Unlock()
				errors = append(errors,
					fmt.Sprintf("[1]%s error: %s", r.Url, err))
			}

			//need header
			if r.Header != nil {
				req.Header = r.Header
			}

			res, err := client.Do(req)
			if err != nil {
				errorLock.Lock()
				defer errorLock.Unlock()
				errors = append(errors,
					fmt.Sprintf("[2]%s error: %s", r.Url, err))
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				errorLock.Lock()
				defer errorLock.Unlock()
				errors = append(errors,
					fmt.Sprintf("[3]%s error: %s", r.Url, err))
			} else {
				//fmt.Printf("Found Body: %s\n", string(body))
				updateLock.Lock()
				defer updateLock.Unlock()
				c.Set(r.Name, string(body))
			}
			<-semaphore
		}(request)
	}
	wg.Wait()

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errors))
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
	}
	//fmt.Printf("%.2f elapsed with response length:\n", secs)
	return errors
}

func (r Request) GetResponse() []byte {
	c := &cache.Cache{}
	body := c.Get(r.Name)
	return []byte(body)
}

//func (r Request) GetValue() interface{} {
//	c := &cache.Cache{}
//	body := c.Get(r.Name)
//	err := json.Unmarshal([]byte(body), &r.Value)
//	if err != nil {
//		fmt.Println("unmarshal body response: ", err)
//	}
//	return r.Value
//}


func (r Request) Execute() []string {
	var requests = Requests{}
	requests = append(requests, r)
	errors := requests.Execute()
	return errors
}

//func MakeRequest(request Request, ch chan<-string) {
//	start := time.Now()
//	secs := time.Since(start).Seconds()
//	fmt.Printf("request url: %s\n", request.Url)
//
//	payload := strings.NewReader(request.Payload)
//
//	client := &http.Client{}
//	req, err := http.NewRequest(request.Method, request.Url, payload)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	//need header
//	if request.Header != nil {
//		req.Header = request.Header
//	}
//
//
//	res, err := client.Do(req)
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	iType := getInteerfaceType(request.ValueType)
//
//	if iType == "AccessToken" {
//		accessToken := &AccessToken{}
//		err = json.Unmarshal(body, accessToken)
//		if err != nil {
//			fmt.Println("unmarshal body response: ", err)
//		}
//	}
//
//	//ch <- fmt.Sprintf("%.2f elapsed with response length:", secs)
//	ch <- fmt.Sprintf("%.2f elapsed with response length:", secs)
//}


