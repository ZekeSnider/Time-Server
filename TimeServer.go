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
import 
(
	"net/http"
	"html/template"
	"html"
	"fmt"
	"time"
	"flag"
	"log"
	"strings"
    "os/exec"
    "sync"
)

type TimePage struct {
    Time string
    UserString string
}

type LoginPage struct {
    ErrorMessage string
}

type HomePage struct {
	UserName string
}


//Declaring map with string indexs to store UDID ints
//Used to store user logins.
var userMap map[string]string


var templatePath string
//Creating a mutex to synchronize state of the map
var mutex = &sync.Mutex{}

func homeHandler (w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	//Since the homeHandler handles all pages, it must forward to a page error
	//if the client is not requesting the homepage (index.html or /)
	if (r.URL.Path != "/") && (r.URL.Path != "/index.html") {
		pageError(w,r)
		return
	}

	//Checking for a client cookie
	cookie,_ := r.Cookie("TimeServerSession")

	//If there is a cookie, the name will be retrieved from the UUID and internal map, and
	//the page will display hello, name.
	if cookie != nil {
		value := cookie.Value

		mutex.Lock()
		name := userMap[value]
		mutex.Unlock()


		tmpl := template.New("Home page")
		tmpl, err := tmpl.ParseFiles(templatePath + "/menu.tmpl", templatePath + "/home.tmpl",)

		if err != nil {
			fmt.Printf("parsing template: %s\n", err)
			return
		}

		theHome := HomePage{
			UserName: name,
		}

		tmpl.ExecuteTemplate(w,"page", theHome)


	} else {	
		//if there is no cookie, the client will be redirected to the login page
		http.Redirect(w, r, "/login", 302)
	}

}


func loginHandler (w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

	tmpl := template.New("Login page")
	tmpl, err := tmpl.ParseFiles(templatePath + "/menu.tmpl", templatePath + "/login.tmpl",)

	if err != nil {
		fmt.Printf("parsing template: %s\n", err)
		return
	}


	

	
	//If the method is GET (regular page load), the login form will be displayed
	if r.Method == "GET" {

		theLogin := LoginPage{
			ErrorMessage: "",
		}

		tmpl.ExecuteTemplate(w,"page", theLogin)

	//if the method is POST (form submit) the form data will be parsed and handled
	} else if r.Method == "POST" {
		//getting the name from the form
		userName:= r.FormValue("name")

		//the data will only be processed if the name is not empty
		if userName != "" {

			//running the name through EscapeString to sanitize the data
			userName = html.EscapeString(userName)

			//Getting a UUID from the unix command
			userUUIDByte,err := exec.Command("uuidgen").Output()

			//If there was an error generating the UUID, log it
			if err != nil {
				log.Fatal(err)
			}

			//Converting the UUIS to a string
			userUUIDString := string(userUUIDByte)

			//Removing the newline for the string to be stored in the cookie
			userUUIDString = strings.Replace(userUUIDString,"\n", "",-1)

			//Generating cookie that stores UUID and expires in 180 days from now
			userCookie := &http.Cookie{Name: "TimeServerSession", Value: userUUIDString, Expires:time.Now().Add(180*24*time.Hour), HttpOnly:true}
			
			//Setting the cookie for the client
			http.SetCookie(w, userCookie)


			log.Printf("User logging in. UUID: %s Name: %s", userUUIDString, userName)

			//Storing the association between the user's UUID and name in the internal map
			mutex.Lock()
			userMap[userUUIDString] = userName
			mutex.Unlock()

			//Redirecting the user to the homepage
			http.Redirect(w, r, "/", 302)

		//if the name was empty no data is processed and 
		//a copy of the login page with the text "C'mon, I need a name" is displayed to the user
		} else {
			theLogin := LoginPage {
				ErrorMessage: "C'mon, I need a name.",
			}

			tmpl.ExecuteTemplate(w,"page", theLogin)
		}
	}
}


func logoutHandler (w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	//Checking for a cookie
	cookie,_ := r.Cookie("TimeServerSession")

	//If there is a cookie, the cookie is replaced with one that expires now, so it is removed from the browser
	//by the client. The user is then reidrected in 10 seconds to the homepage after the good bye message is displayed.
	if cookie != nil {
		value := cookie.Value
		deleteCookie := &http.Cookie{Name: "TimeServerSession", Value: value, Expires:time.Now(), HttpOnly:true}
		http.SetCookie(w, deleteCookie)
		log.Printf("User with UUID %s logged out.", value)
	}

	fmt.Fprint(w, "<html><head><META http-equiv=\"refresh\" content=\"10;URL=/\"><body><p>Good-bye.</p></body></html>")

	
}


func timeHandler (w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

	//Formating the current time in the proper format and storing it to a variable
	pageTime := time.Now().Format("03:04:05PM")

	tmpl := template.New("Time page")
	tmpl, err := tmpl.ParseFiles(templatePath + "/menu.tmpl", templatePath + "/time.tmpl",)

	if err != nil {
		fmt.Printf("parsing template: %s\n", err)
		return
	}

	//checking for a client's cookie
	cookie,_ := r.Cookie("TimeServerSession")

	var user string
	//if a cookie exists (the user is logged in), the UUID from the cookie is used
	//to retrieve the userName from the internal map. The username is then printed
	if cookie != nil {
		value := cookie.Value

		mutex.Lock()
		name := userMap[value]
		mutex.Unlock()

		user = ", " + name
	} else {
		user = ""
	}


	theTime := TimePage {
		Time: pageTime,
		UserString: user,
	}

	tmpl.ExecuteTemplate(w,"page", theTime)

}
func pageError(w http.ResponseWriter, r *http.Request){
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	//writing 404 not found error to html header
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//Displaying custom 404 text to the error page
	fmt.Fprintf(w, "<p>These are not the URLs you're looking for.</p>")

}
func main() {

	//Creating the map top store the UUID and name values
	userMap = make(map[string]string)



	//Whenever the map is access, the mutex is locked beforehand and locked after
	//to ensure exclusive access to the usermap

	//declaring command line flags for the time server

	//Port number (optional): declares what port the server should launch on.
	//defaults to 8080
	portPointer := flag.String("port", "8080", "Server port number")

	//Version output (optional): if true, the version number will be 
	//output to the console.
	versionBoolPointer := flag.Bool("v", false, "Display server version bool")

	//Template path (optional): if specified the templates will be loaded 
	//from this path
	templatePathPointer := flag.String("template", "templates", "Path to templates")

	//parsing the flags
	flag.Parse()

	//Outputting server version number if it is requested in command line flags
	if *versionBoolPointer == true {
		fmt.Print("Personalized time server version 1.1.3")
	}

	//adding a ":" to the port number to match the format requested by http.ListenAndServe
	portNumber := ":"+*portPointer

	templatePath = *templatePathPointer


    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

    http.HandleFunc("/login", loginHandler)
    //http.HandleFunc("/login/:name", loginActionHandler)

    http.HandleFunc("/logout", logoutHandler)

	//If the /time page is requested, the time will be displayed
    http.HandleFunc("/time", timeHandler)

    //If any other page is requested, a 404 page will be displayed
    http.HandleFunc("/", homeHandler)





    //attempting to start the server on the requested port.
    //if there are any errors they will be stored to the err variable
    err := http.ListenAndServe(portNumber, nil)
    
    //if there was any errors in starting the server, they will be displayed
    //to the console and the program will exit.
    if err!= nil {
    	log.Fatal(err)
    }


}