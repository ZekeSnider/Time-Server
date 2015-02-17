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

package authserver

import (
	log "github.com/cihub/seelog"
	config "command/config"
	"html"
	"net/http"
	"strings"
	"os/exec"
	"time"
	"sync"
)

//This struct stores information passed to templates
type Page struct {
	Time         string
	UserName     string
	ErrorMessage string
	UTCTime      string
}

//Declaring global variables used by different functions
var userMap map[string]string
var templatePath string
var absoluteETCPath string

//Creating a mutex to synchronize state of the map
var mutex = &sync.Mutex{}

func CheckUserCookie(r *http.Request) string {
	//This function checks if the user has a cookie
	cookie, _ := r.Cookie("TimeServerSession")

	if cookie != nil {
		//if they do have a cookie their information
		//is retrieved from the map and returned
		value := cookie.Value

		mutex.Lock()
		name := userMap[value]
		mutex.Unlock()

		return name
	} else {
		//otherwise an empty string is returned
		return ""
	}
}

func SetUserCookie(w http.ResponseWriter, r *http.Request, inputName string) {
	//running the name through EscapeString to sanitize the data
	inputName = html.EscapeString(inputName)

	//Getting a UUID from the unix command
	userUUIDByte, err := exec.Command("uuidgen").Output()

	//If there was an error generating the UUID, log it
	if err != nil {
		log.Errorf("Error generating the UUID. %d", err)
	}

	//Converting the UUID to a string
	userUUIDString := strings.Replace(string(userUUIDByte), "\n", "", -1)

	//Generating and setting cookie that stores UUID and expires in 180 days from now
	userCookie := &http.Cookie{Name: "TimeServerSession", Value: userUUIDString, Expires: time.Now().Add(180 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(w, userCookie)

	log.Debugf("User logging in. UUID: %s Name: %s", userUUIDString, inputName)

	//Storing the association between the user's UUID and name in the internal map
	mutex.Lock()
	userMap[userUUIDString] = inputName
	mutex.Unlock()
}

func LogoutCookie(w http.ResponseWriter, r *http.Request) {
	//Checking for a cookie
	cookie, _ := r.Cookie("TimeServerSession")

	//If there is a cookie, the cookie is replaced with one that expires now, so it is removed from the browser
	//by the client. The user is then reidrected in 10 seconds to the homepage after the good bye message is displayed.
	if cookie != nil {
		value := cookie.Value
		deleteCookie := &http.Cookie{Name: "TimeServerSession", Value: value, Expires: time.Now(), HttpOnly: true}
		http.SetCookie(w, deleteCookie)
		log.Debugf("User with UUID logged out. UUID: %s", value)
	}
}

func init( ) {
	log.ReplaceLogger(config.Logger)
}