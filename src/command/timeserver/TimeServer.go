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
	counter "command/counter"
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"html"
	"html/template"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var Logger log.LoggerInterface
var currentConnections int
var concurrentMutex = &sync.Mutex{}

//This struct stores information passed to templates
type Page struct {
	Time         string
	UserName     string
	ErrorMessage string
	UTCTime      string
}

//Attempts to get name from a cookie
func CheckUserCookie(r *http.Request) string {
	//This function checks if the user has a cookie
	cookie, _ := r.Cookie("TimeServerSession")

	if cookie != nil {
		//if they do have a cookie their UUID is retrieved
		value := cookie.Value

		//formatting request to auth server
		authPath := config.AuthHost + config.AuthPort + "/get?cookie=" + value
		resp, err := http.Get(authPath)
		if err != nil {
			return ""
		} else {
			//getting the result and returning it
			pageBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Errorf(err.Error())
			}
			return string(pageBody)
		}

	} else {
		//otherwise an empty string is returned
		return ""
	}
}

//shortcut functions to lock/unlock and increment/decrement the concurrent counter
func incrementCurrent() {
	concurrentMutex.Lock()
	currentConnections++
	concurrentMutex.Unlock()
}
func decrementCurrent() {
	concurrentMutex.Lock()
	currentConnections--
	concurrentMutex.Unlock()
}

//if the MaxConnections flag is 0, the server will always allow any connection.
//otherwise, it will only allow if it currently handling less than the MaxConnections
func checkCurrent() bool {
	if config.MaxConnections == 0 {
		incrementCurrent()
		return true
	} else if currentConnections < config.MaxConnections {
		incrementCurrent()
		return true
	} else {
		return false
	}
}

//Calculates a random number from a normal distribution based on the command line flags
func getServerWait() float64 {
	waitTime := math.Abs(rand.NormFloat64()*float64(config.ResponseTime) + float64(config.DeviationTime))
	return waitTime
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

	//formatting request to the auth server
	authPath := config.AuthHost + config.AuthPort + "/set?cookie=" + userUUIDString + "&name=" + inputName

	//sending the request
	resp, err := http.Get(authPath)
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		pageBody, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Criticalf(err.Error())
		}
		_ = pageBody
	}
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if checkCurrent() {
		log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

		//Since the homeHandler handles all pages, it must forward to a page error
		//if the client is not requesting the homepage (index.html or /)
		if (r.URL.Path != "/") && (r.URL.Path != "/index.html") {
			pageError(w, r)
			log.Warnf("404 Page Not Found: %s %s %s", r.RemoteAddr, r.Method, r.URL)
			return
		}

		//Checking for a client cookie
		user := CheckUserCookie(r)

		//if there is a cookie, the homepage is displayed with the user's username
		if user != "" {
			loadPage(w, "home", &Page{UserName: user})

		} else {
			//if there is no cookie, the client will be redirected to the login page
			http.Redirect(w, r, "/login", 302)
		}
		decrementCurrent()
	} else {
		w.WriteHeader(500)
	}

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if checkCurrent() {
		counter.IncrementValue("login")

		log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

		//If the method is GET (regular page load), the login form will be displayed
		if r.Method == "GET" {
			loadPage(w, "login", &Page{})
			//if the method is POST (form submit) the form data will be parsed and handled
		} else if r.Method == "POST" {
			//getting the name from the form
			userName := r.FormValue("name")
			//the data will only be processed if the name is not empty

			if userName != "" {
				SetUserCookie(w, r, userName)

				//Redirecting the user to the homepage
				http.Redirect(w, r, "/", 302)

			} else {
				//if the name was empty no data is processed and
				//a copy of the login page with the text "C'mon, I need a name" is displayed to the user
				loadPage(w, "login", &Page{ErrorMessage: "C'mon, I need a name."})
			}
		}
		decrementCurrent()
	} else {
		w.WriteHeader(500)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if checkCurrent() {
		log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)
		LogoutCookie(w, r)

		loadPage(w, "logout", &Page{})
		decrementCurrent()
	} else {
		w.WriteHeader(500)
	}

}

func loadPage(w http.ResponseWriter, inputTemplate string, p *Page) {
	//simulating a delatyed response time using a normal distribution
	time.Sleep(time.Duration(getServerWait()) * time.Millisecond)

	//Loads the specified template content inside of the menu template.
	tmpl := template.New("Page")
	tmpl, err := tmpl.ParseFiles(config.AbsoluteETCPath+"/"+config.TemplatePath+"/menu.tmpl", config.AbsoluteETCPath+"/"+config.TemplatePath+"/"+inputTemplate+".tmpl")

	if err != nil {
		log.Errorf("Execute template error: %d", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	counter.IncrementValue("200s")
	//executing the template
	err = tmpl.ExecuteTemplate(w, "page", p)

	if err != nil {
		log.Errorf("Execute template error: %d", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	if checkCurrent() {
		log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

		//Formating the current time in the proper format and storing it to a variable
		pageTime := time.Now().Format("03:04:05PM")
		pageUTCTime := time.Now().UTC().Format("3:04:05PM")

		//If the user is logged in, their name will be displayed on the page
		user := CheckUserCookie(r)
		if user != "" {
			counter.IncrementValue("time-user")
			user = ", " + user
		} else {
			counter.IncrementValue("time-anon")
		}

		loadPage(w, "time", &Page{Time: pageTime, UserName: user, UTCTime: pageUTCTime})
		decrementCurrent()
	} else {
		w.WriteHeader(500)
	}

}

func monitorHandler(w http.ResponseWriter, r *http.Request) {
	counter.IncrementValue("200s")
	Mapcopy := counter.GetMapCopy()

	jsonDump, _ := json.Marshal(Mapcopy)
	fmt.Fprintf(w, string(jsonDump))

}

func pageError(w http.ResponseWriter, r *http.Request) {
	if checkCurrent() {
		counter.IncrementValue("404s")
		log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

		//writing 404 not found error to html header
		w.WriteHeader(http.StatusNotFound)

		//Displaying custom 404 text to the error page
		loadPage(w, "error", &Page{})
		decrementCurrent()
	} else {
		w.WriteHeader(500)
	}

}

func main() {
	counter.ResetMapValue("login")
	counter.ResetMapValue("time-user")
	counter.ResetMapValue("time-anon")
	counter.ResetMapValue("200s")
	counter.ResetMapValue("404s")

	currentConnections = 0

	//loading the logger
	Logger, err := log.LoggerFromConfigAsFile(config.ServerLog)
	if err != nil {
		log.Errorf("Error loading the seelog config", err)
	}
	log.ReplaceLogger(Logger)

	defer log.Flush()

	log.Infof("The log file path is %s", config.ServerLog)

	//Outputting server version number if it is requested in command line flags
	if config.DisplayVersionBool == true {
		log.Infof("Personalized Time Server version 1.3")
	}

	//the /css/ directory is served as a file directory so the html header can access the stylesheet
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("etc/css"))))

	//Calling handlers for other pages
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/monitor", monitorHandler)
	http.HandleFunc("/", homeHandler)

	//attempting to start the server on the requested port.
	//if there are any errors they will be displayed
	err = http.ListenAndServe(config.ServerPort, nil)

	fmt.Printf("Server err:%v", err)
}
