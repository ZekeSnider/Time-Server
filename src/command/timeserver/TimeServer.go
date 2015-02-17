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

package timeserver

import (
	auth "command/authserver"
	config "command/config"
	log "github.com/cihub/seelog"
	"html/template"
	"net/http"
	"time"
)

//This struct stores information passed to templates
type Page struct {
	Time         string
	UserName     string
	ErrorMessage string
	UTCTime      string
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
	user := auth.CheckUserCookie(r)

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
			auth.SetUserCookie(w, r, userName)

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

	auth.LogoutCookie(w, r)

	loadPage(w, "logout", &Page{})

}

func loadPage(w http.ResponseWriter, inputTemplate string, p *Page) {

	//Loads the specified template content inside of the menu template.
	tmpl := template.New("Page")
	tmpl, err := tmpl.ParseFiles(config.AbsoluteETCPath+"/"+config.TemplatePath+"/menu.tmpl", config.AbsoluteETCPath+"/"+config.TemplatePath+"/"+inputTemplate+".tmpl")

	if err != nil {
		log.Errorf("Execute template error: %d", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//executing the template
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
	user := auth.CheckUserCookie(r)
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
	log.ReplaceLogger(config.Logger)
	//Whenever map is access, the mutex is locked beforehand and locked after
	//to ensure exclusive access to the usermap

	//declaring command line flags for the time server
	//All flag functionality is detailed in the included README.

	//the /css/ directory is served as a file directory so the html header can access the stylesheet
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("etc/css"))))

	//Calling handlers for other pages
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", homeHandler)

	//attempting to start the server on the requested port.
	//if there are any errors they will be displayed
	err := http.ListenAndServe(config.ServerPort, nil)

	config.Logger.Infof("hello world")
	log.Errorf("Server err:%v", err)

}
