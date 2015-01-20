//Zeke Snider
//CSS 490 Homework 1
//Time Server version 1.0

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
	"fmt"
	"time"
	"flag"
	"log"
	"strings"
    //"net/http/cookiejar"
    "os/exec"
)

//Declaring map with string indexs to store UDID ints
//Used to store user logins.
var userMap map[string]string

func homeHandler (w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		pageError(w,r)
		return
	}

	cookie,_ := r.Cookie("TimeServer")

	if cookie != nil {
		value := cookie.Value
		name := userMap[value]
		log.Println(name)
		log.Println(value)
		fmt.Fprint(w, "<html><body><p>hello, ", name, ". </p></body></html>")
	} else {
		http.Redirect(w, r, "/login", 302)
	}

}

func loginHandler (w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	if r.Method == "GET" {
		fmt.Fprint(w, "<html><body><form method=\"post\" action=\"login\">What is your name, Earthling?<input type=\"text\" name=\"name\" size=\"50\"><input type=\"submit\"></form></p></body></html>")
	} else if r.Method == "POST" {
		userName:= r.FormValue("name")
		userUUIDByte,err := exec.Command("uuidgen").Output()

		userUUIDString := string(userUUIDByte)

		userUUIDString = strings.Replace(userUUIDString,"\n", "",-1)
		if err != nil {
			log.Fatal(err)
		}

		userCookie := &http.Cookie{Name: "TimeServer", Value: userUUIDString, Expires:time.Now().Add(180*24*time.Hour), HttpOnly:true}
		
		http.SetCookie(w, userCookie)

		log.Println(userName)
		userMap[userUUIDString] = userName

		http.Redirect(w, r, "/", 302)
	}
}


func logoutHandler (w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "<html><head><META http-equiv=\"refresh\" content=\"10;URL=/\"><body><p>Good-bye.</p></body></html>")

	
}

func timeHandler (w http.ResponseWriter, r *http.Request) {
	//getting the current time
	currentTime := time.Now()

	//serving the current time page

	//Formating the current time in the proper format and storing it to a variable
	pageTime := currentTime.Format("03:04:05PM")
	//printing the head css styles to the page
	fmt.Fprint(w, "<html><head><style>p {font-size: xx-large} span.time {color: red}</style></head>")
	//printing the body of the page
	fmt.Fprint(w, "<body><p>The time is now <span class=\"time\">", pageTime, "</span>.</p></body></html>")
}
func pageError(w http.ResponseWriter, r *http.Request){
	//writing 404 not found error to html header
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//Displaying custom 404 text to the error page
	fmt.Fprintf(w, "<p>These are not the URLs you're looking for.</p>")

}
func main() {

	//Creating the map top store the UUID and name values
	userMap = make(map[string]string)

	//declaring command line flags for the time server

	//Port number (optional): declares what port the server should launch on.
	//defaults to 8080
	portPointer := flag.String("port", "8080", "Server port number")

	//Version output (optional): if true, the version number will be 
	//output to the console.
	versionBoolPointer := flag.Bool("v", false, "Display server version bool")

	//parsing the flags
	flag.Parse()

	//Outputting server version number if it is requested in command line flags
	if *versionBoolPointer == true {
		fmt.Print("Personalized time server version 1.0.2")
	}

	//adding a ":" to the port number to match the format requested by http.ListenAndServe
	portNumber := ":"+*portPointer



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