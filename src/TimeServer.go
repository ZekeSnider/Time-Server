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
	"flag"
	log "github.com/cihub/seelog"
	"html"
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
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

//Declaring global variables used by different functions
var userMap map[string]string
var Logger log.LoggerInterface
var templatePath string
var absoluteETCPath string

//Creating a mutex to synchronize state of the map
var mutex = &sync.Mutex{}

func checkUserCookie(r *http.Request) string {
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

	//Since the homeHandler handles all pages, it must forward to a page error
	//if the client is not requesting the homepage (index.html or /)
	if (r.URL.Path != "/") && (r.URL.Path != "/index.html") {
		pageError(w, r)
		log.Warnf("404 Page Not Found: %s %s %s", r.RemoteAddr, r.Method, r.URL)
		return
	}

	//Checking for a client cookie
	user := checkUserCookie(r)

	//if there is a cookie, the homepage is displayed with the user's username
	if user != "" {
		loadPage(w, "home", &Page{UserName: user})

	} else {
		//if there is no cookie, the client will be redirected to the login page
		http.Redirect(w, r, "/login", 302)
	}

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

			//running the name through EscapeString to sanitize the data
			userName = html.EscapeString(userName)

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

			log.Debugf("User logging in. UUID: %s Name: %s", userUUIDString, userName)

			//Storing the association between the user's UUID and name in the internal map
			mutex.Lock()
			userMap[userUUIDString] = userName
			mutex.Unlock()

			//Redirecting the user to the homepage
			http.Redirect(w, r, "/", 302)

		} else {
			//if the name was empty no data is processed and
			//a copy of the login page with the text "C'mon, I need a name" is displayed to the user
			loadPage(w, "login", &Page{ErrorMessage: "C'mon, I need a name."})
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)
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

	loadPage(w, "logout", &Page{})

}

func loadPage(w http.ResponseWriter, inputTemplate string, p *Page) {

	//Loads the specified template content inside of the menu template.
	tmpl := template.New("Page")
	tmpl, err := tmpl.ParseFiles(absoluteETCPath+"/"+templatePath+"/menu.tmpl", absoluteETCPath+"/"+templatePath+"/"+inputTemplate+".tmpl")

	if err != nil {
		log.Errorf("Execute template error: %d", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "page", p)

	if err != nil {
		log.Errorf("Execute template error: %d", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

	//Formating the current time in the proper format and storing it to a variable
	pageTime := time.Now().Format("03:04:05PM")
	pageUTCTime := time.Now().UTC().Format("3:04:05PM")

	//If the user is logged in, their name will be displayed on the page
	user := checkUserCookie(r)
	if user != "" {
		user = ", " + user
	}

	loadPage(w, "time", &Page{Time: pageTime, UserName: user, UTCTime: pageUTCTime})

}

func pageError(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Page Loaded. RemoteAddress:%s Method:%s URL:%s", r.RemoteAddr, r.Method, r.URL)

	//writing 404 not found error to html header
	w.WriteHeader(http.StatusNotFound)

	//Displaying custom 404 text to the error page
	loadPage(w, "error", &Page{})

}

func main() {

	//Creating the map top store the UUID and name values
	userMap = make(map[string]string)

	//Whenever the map is access, the mutex is locked beforehand and locked after
	//to ensure exclusive access to the usermap

	//declaring command line flags for the time server
	//All flag functionality is detailed in the included README.
	portPointer := flag.String("port", "8080", "Server port number")
	versionBoolPointer := flag.Bool("v", false, "Display server version bool")
	templatePathPointer := flag.String("template", "templates", "Path to templates")
	logPathPointer := flag.String("log", "seelog.xml", "Name of log file")

	//parsing the flags
	flag.Parse()


	//setting up the logging library to load configuration from the specified file
	absoluteETCPath, _ = filepath.Abs("etc/")
	defer log.Flush()
	Logger, err := log.LoggerFromConfigAsFile(absoluteETCPath + "/" + *logPathPointer)
	if err != nil {
		log.Errorf("Error loading the seelog config", err)
	}

	//logging startup and replacing the default logger
	log.Infof("The log file path is %s", absoluteETCPath+"/"+*logPathPointer)
	log.ReplaceLogger(Logger)

	//Outputting server version number if it is requested in command line flags
	if *versionBoolPointer == true {
		log.Infof("Personalized Time Server version 1.2")
	}

	//adding a ":" to the port number to match the format requested by http.ListenAndServe
	portNumber := ":" + *portPointer

	templatePath = *templatePathPointer

	//the /css/ directory is served as a file directory so the html header can access the stylesheet
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("etc/css"))))

	//Calling handlers for other pages
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", homeHandler)

	//attempting to start the server on the requested port.
	//if there are any errors they will be displayed
	err = http.ListenAndServe(portNumber, nil)

	log.Errorf("Server err:%v", err)

}
