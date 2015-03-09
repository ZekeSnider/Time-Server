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

package counter

import (
	"sync"
)

var countMap map[string]int
var mutex = &sync.Mutex{}


//Increments the value of a key in the map
func IncrementValue(key string) {
	mutex.Lock()

	//checking if the map exists 
	_, existsBool := countMap[key]

	//if the value exists in the map, increment it
	if existsBool {		
		countMap[key]++
	} else { 
		//otherwise, set it to 1. If you try to increment an empty map element go will crash.
		countMap[key] = 1
	}

	mutex.Unlock()

}

//Resets a map element to 0. 
func ResetMapValue(key string) {
	mutex.Lock()
	countMap[key] = 0
	mutex.Unlock()
}

//Returns a copy of the map
func GetMapCopy() map[string] int {
	mutex.Lock()
	mapCopy := countMap
	mutex.Unlock()

	return mapCopy
}

//Creating the map on initialization
func init() {	
	countMap = make(map[string]int)
}
