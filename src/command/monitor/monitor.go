//Zeke Snider
//CSS 490 Assignment 6

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
	"flag"
	"strings"
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Sample struct{
    time string
    value float64
}
var Results map[string] map[string][]Sample
var interval int
var runTime int
var targetList []string

func monitor() {
	//runs a request every interval seconds
	for _ = range time.Tick(time.Duration(interval) * time.Second) {
		monitorRoutine()
	}
}
func monitorRoutine() {
	//looping over url list
	for index, _ := range targetList {

		//getting the monitor page
		monitorURL := targetList[index] + "/monitor"
		resp, err := http.Get(monitorURL)

		//if the submap for this url has not been created yet, make it
		if _, ok := Results[monitorURL]; !ok {
			Results[monitorURL] = make(map[string][]Sample, 10)
		}

		//carch error in page get
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			//reading the body
			pageBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Printf(err.Error())
			}

			//creating a map to store the new entries
			var newEntries map[string]interface{}

			//parsing the json from the page body
			json.Unmarshal(pageBody, &newEntries)

			//getting current time
			const layout = "1/2/2006 3:04:05pm (MST)"
			t := time.Now()

			///looping over all json entries
			for index, _ := range newEntries {
				
				//converting value to float
				valueFloat := newEntries[index].(float64)

				//pairing it with the time in sample struct
				newSample := Sample{
					time: string(t.Format(layout)),
					value: valueFloat,
				}	

				//adding this to the map's map's sample array
				Results[monitorURL][index] = append(Results[monitorURL][index], newSample)
			}
		}
	}
}

func removeComma(inputString string) string {
	//removes the comma and newline character from the end of the string
	inputString = strings.TrimSuffix(inputString, ",\n")
	//replaces the newline character
	inputString += fmt.Sprintf("\n")
	//returning new string
	return inputString
}

func main() {
	//parsing the flags
	targetsPointer := flag.String("targets", "http://localhost:8080", "URLs of targets seperated by commas")
	intervalPointer := flag.Int("sampleIntervalSec", 4, "Monitor interval time")
	runTimePointer := flag.Int("runTimeSec", 10, "Monitor run time")

	flag.Parse()

	targets := *targetsPointer
	interval = *intervalPointer
	_ = interval
	runTime = *runTimePointer


	//The results are stored in a map, the key is the json and the values are a struct of the url and time of request
	Results = make(map[string] map[string] []Sample)
	
	//URLs are seperated by comma in the flag
	targetList = strings.Split(targets, ",")

	go monitor()

	//sleeping for the runtime while the monitors run
	time.Sleep(time.Duration(runTime) * time.Second)

	var jsonString string
	//creating a json output string

	jsonString += "\n{\n"
	for i, _ := range Results { //looping through all url targets
		jsonString += fmt.Sprintf("  \"%v\": {\n", i)
		for j:= range Results[i] { //looping through all counters
			jsonString += fmt.Sprintf("    \"%v\": [\n", j)
			for k := range Results[i][j] { //looping through all time/value pairs
				jsonString += fmt.Sprintf("      {\n")
				jsonString += fmt.Sprintf("       \"Time\": \"%v\",\n", Results[i][j][k].time)
				jsonString += fmt.Sprintf("       \"Value\": %v\n", Results[i][j][k].value)
				jsonString += fmt.Sprintf("      },\n")
			}
			//if that was the last element displayed, remove the trailing comma
			jsonString = removeComma(jsonString) 
			jsonString += fmt.Sprintf("    ],\n") 
		}
		//if that was the last element displayed, remove the trailing comma
		jsonString = removeComma(jsonString)

		jsonString += fmt.Sprintf("  },\n")
	}
	//if that was the last element displayed, remove the trailing comma
	jsonString = removeComma(jsonString)

	jsonString += "}\n"

	//done, output the final json
	fmt.Printf("%v", jsonString)
}
