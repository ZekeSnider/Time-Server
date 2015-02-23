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
	config "command/config"
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"
)

//This struct stores information passed to templates
type Page struct {
	Time         string
	UserName     string
	ErrorMessage string
	UTCTime      string
}

//Struct used to store JSON data
type userLogin struct {
	Name string
	UUID string
}

//Declaring global variables used by different functions
var UserMap map[string]string
var templatePath string
var absoluteETCPath string

//Creating a mutex to synchronize state of the map
var mutex = &sync.Mutex{}

//These functions lock the mutex, modify the map, then unlock the mutex
func setMap(inputUserName string, inputUUID string) {
	mutex.Lock()
	UserMap[inputUUID] = inputUserName
	mutex.Unlock()
}

func getMap(inputUUID string) string {
	mutex.Lock()
	userName := UserMap[inputUUID]
	mutex.Unlock()
	return userName
}

//imports user list from a JSON file. Accepts a path to the JSON, and a bool which tells it to modify the master map or not.
//returns a copy of the imported map for checking integrity
func importUserList(JSONFilePath string, ModifyMainMap bool) map[string]string {
	UserMapCheck := make(map[string]string)

	//opening the file
	fileContent, err := ioutil.ReadFile(JSONFilePath)
	if err != nil {
		log.Criticalf("File error: %v\n", err)
	} else {
		//Creating array of structs to hold jsondata
		var userLogins []userLogin

		//parsing the json
		err = json.Unmarshal(fileContent, &userLogins)
		if err != nil {
			log.Errorf("Unmarshal error %d", err)
		} else {
			//looping through every element and importing.
			for index, _ := range userLogins {
				if ModifyMainMap {
					setMap(userLogins[index].Name, userLogins[index].UUID)
				}
				UserMapCheck[userLogins[index].UUID] = userLogins[index].Name
			}
		}
	}
	return UserMapCheck

	//http://golang.org/pkg/encoding/json/#pkg-examples used for reference
}

//Exports the user data to a JSON file. Returns a bool of integrity check of whether or not
//the exported data was identacle to the original data or not
func exportUserList(JSONFilePath string) bool {

	//copying the original map
	mutex.Lock()
	userMapCopy := UserMap
	mutex.Unlock()

	//Creating array of structs to hold json data
	var allLogins []userLogin

	//looping over every element in map
	for k := range userMapCopy {
		//appending element to the array
		aLogin := userLogin{
			UUID: k,
			Name: userMapCopy[k],
		}
		allLogins = append(allLogins, aLogin)
	}
	//parsing into JSON
	jsonData, err := json.Marshal(allLogins)

	if err != nil {
		log.Errorf(err.Error())
	}

	//Creating the JSON file
	dumpJSONFile, err := os.Create(JSONFilePath)

	if err != nil {
		log.Errorf(err.Error())
	}

	//Writing to tit
	defer dumpJSONFile.Close()
	dumpJSONFile.Write(jsonData)
	dumpJSONFile.Close()

	//Checks if the newly created data is correct by importing it and comparing it with the original data
	if reflect.DeepEqual(importUserList(JSONFilePath, false), userMapCopy) {
		return true
	} else {
		return false
	}
	//https://www.socketloop.com/tutorials/golang-convert-csv-data-to-json-format-and-save-to-file
	//used for reference
}

//returns username from a UUID via HTTP request
func getHandler(w http.ResponseWriter, r *http.Request) {
	//attempting to get the name
	userUUID := r.FormValue("cookie")
	userName := getMap(userUUID)

	//if there was no name, send a header error
	if userName == "" {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
	}

	//otherwise, print the name
	fmt.Fprintf(w, "%s", userName)

}

func setHandler(w http.ResponseWriter, r *http.Request) {
	//getting values from url
	userUUID := r.FormValue("cookie")
	userName := r.FormValue("name")

	//set the map
	setMap(userName, userUUID)

	//if the data is bad, throw an error
	if userName == "" || userUUID == "" {
		w.WriteHeader(400)
	} else {

		w.WriteHeader(200)
	}

	//print the data
	fmt.Fprintf(w, "%s %s", userName, userUUID)
}

//pages made to use during testing that displays all elements in map
func displayAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello:\n")
	for k := range UserMap {
		fmt.Fprintf(w, "%s %s", k, UserMap[k])
	}
}

//Complete proccess to create a new JSON dump file
func exportRoutine() {
	log.Debugf("Starting JSON backup.")

	//renaming the current file with .bak
	err := os.Rename(config.DumpFile, config.DumpFile+".bak")
	if err != nil {
		log.Errorf("There was no json file to backup.")
	}

	//Exporting the current user list to the dump file location
	exportResult := exportUserList(config.DumpFile)

	//if successful, delete the backup
	if exportResult {
		err := os.Remove(config.DumpFile + ".bak")
		if err != nil {
			log.Errorf("There was no json backup to delete")
		}
	} else { //otherwise, remove the new file, and rename the old backup.
		err := os.Remove(config.DumpFile)
		if err != nil {
			log.Errorf("Couldn't delete bad json file!")
		}
		err = os.Rename(config.DumpFile+".bak", config.DumpFile)
		if err != nil {
			log.Errorf("couldn't restore backup file!")
		}
	}
	log.Debugf("JSON backup finished!")

}

//repeats backup on a timeframe specified by the CheckPointInterval flag
func repeatBackup(t time.Duration) {
	for _ = range time.Tick(t) {
		exportRoutine()
	}
}
func main() {
	//creating map
	UserMap = make(map[string]string)

	//loading logger
	Logger, err := log.LoggerFromConfigAsFile(config.AuthLog)
	if err != nil {
		log.Errorf("Error loading the seelog config", err)
	}

	log.ReplaceLogger(Logger)

	defer log.Flush()

	//starting another thread to run the json backups
	go repeatBackup(time.Millisecond * time.Duration(config.CheckPointInterval))

	log.Infof("The log file path is %s", config.ServerLog)

	//importing the user list
	_ = importUserList(config.DumpFile, true)

	//starting the server
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/all", displayAll)
	err = http.ListenAndServe(config.GetAuthPort(), nil)
	log.Errorf("Server err:%v", err)
}
