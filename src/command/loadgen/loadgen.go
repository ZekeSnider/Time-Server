//Zeke Snider
//CSS 490 Assignment 3

/*
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"fmt"
	//log "github.com/cihub/seelog"
	counter "command/counter"
	"flag"
	"net/http"
	"time"
)

var Rate int
var Burst int
var Timeout int
var Runtime int
var Url string

//repeats backup on a timeframe specified by the CheckPointInterval flag
func load() {
	duration := time.Duration((Burst * 1000000 / Rate)) * time.Microsecond
	for _ = range time.Tick(duration) {
		for i := 1; i <= Burst; i++ {
			go loadRoutine()
		}
	}
}

func loadRoutine() {
	//loading the page with a timeout specified by command line flag.
	loadClient := http.Client{
		Timeout: time.Duration(Timeout) * time.Second,
	}
	resp, err := loadClient.Get(Url)

	counter.IncrementValue("Total")

	if err != nil {
		counter.IncrementValue("Errors")
	} else {
		resp.Body.Close()
		//getting statuscode, then incrementing related map.
		status := resp.StatusCode
		if status == 503 {
			counter.IncrementValue("Errors")
		} else if status >= 100 && status < 200 {
			counter.IncrementValue("100s")
		} else if status >= 200 && status < 300 {
			counter.IncrementValue("200s")
		} else if status >= 300 && status < 400 {
			counter.IncrementValue("300s")
		} else if status >= 400 && status < 500 {
			counter.IncrementValue("400s")
		} else if status >= 500 && status < 600 {
			counter.IncrementValue("500s")
		}

	}

	//reading the body of the request to check for errors

}

//Setting all counter map values to 0 on initialization
func init() {
	counter.ResetMapValue("Total")
	counter.ResetMapValue("100s")
	counter.ResetMapValue("200s")
	counter.ResetMapValue("300s")
	counter.ResetMapValue("400s")
	counter.ResetMapValue("500s")
	counter.ResetMapValue("Errors")
}

func main() {
	//getting flags, parsing them, then storing to values

	ratePointer := flag.Int("rate", 200, "average rate of requests (per second)")
	burstPointer := flag.Int("burst", 20, "number of concurrent requests to issue")
	timeoutPointer := flag.Int("timeout", 5, "max time to wait for a response")
	runtimePointer := flag.Int("runtime", 10, "number of seconds to proccess")
	urlPointer := flag.String("url", "http://localhost:8080/", "URL to sample")

	flag.Parse()

	Rate = *ratePointer
	Burst = *burstPointer
	Timeout = *timeoutPointer
	Runtime = *runtimePointer
	Url = *urlPointer

	//run the tests on another go runtime
	go load()

	//going to sleep for specified time to allow the test to run
	fmt.Printf("Sleeping for %d.\n", time.Duration(Runtime)*time.Second)
	time.Sleep(time.Duration(Runtime) * time.Second)
	fmt.Printf("Done sleeping.\n")

	//Getting a copy of the results map
	countMap := counter.GetMapCopy()

	//Printing the results to the log
	fmt.Printf("%v: %v\n", "Total", countMap["Total"])
	fmt.Printf("%v: %v\n", "100s", countMap["100s"])
	fmt.Printf("%v: %v\n", "200s", countMap["200s"])
	fmt.Printf("%v: %v\n", "300s", countMap["300s"])
	fmt.Printf("%v: %v\n", "400s", countMap["400s"])
	fmt.Printf("%v: %v\n", "500s", countMap["500s"])
	fmt.Printf("%v: %v\n", "Errors", countMap["Errors"])
}
