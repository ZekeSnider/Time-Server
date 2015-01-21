# Assignment 2: Personalized Time Server
Zeke Snider  
CSS 490C  
1/11/15  

## Overview
This program extends the Assignment 1: Time Server to include more functionality. 

The time server uses the following golang libraries:
* net/http
* html"
* fmt
* time
* flag
* log
* strings
* os/exec
* sync

## Build Instructions
The program requires no special build instructions. Simply CD to the directory containing TimeServer.go then run "go build TimeServer.go", then "./TimeServer" and any flags if required. The server will then launch and display http request logs to the console.  

Alternatively, the included Makefile can be used to build and run the server.

## Design Notes
The sync library is used to lock and unlock the internal UUID name map so that the map does not go out of sync between multiple requests at the same time.  

When the server stops, the server's internal map is deleted but the client's cookie still remains unless they logout before the server stops. This is because the cookies are set to expire in 180 days. Thus, if a user logs in, the server restarts, then the user tries to load the homepage it will display "Hello, ." because it cannot pull the name from the map anymore. This could be solved in a future iteration by implementing permanent user map storage.  

If any errors occur when starting the server, the errors are output to the console and the program is terminated. For example, this can happen if the starting port you provide is already in use by another service.


## Page List 

###Homepage ("/", "index.html")
The homepage displays a simple html page that says "hello, name" if the user is logged in via a cookie. If no cookie is present, the user will be redirected to the /login page  

###Login ("/login", "/login?name=name")
The login page displays a form requesting the user's name, with a submit button when first loaded. When the submit button is pressed, the page reloads with the form /login?name=name using POST method. The user is then authenticated via a cookie and a UUID. the UUID is mapped to the username by an internal map.  

Because the login is done by a POST method, it is not possible to modify the url to login without using the form. The name is sanitized with EscapeString to ensure that no html injection is possible. If the user supplies an empty name and presses submit, the page reloads the form displaying the text "C'mon, I need a name."

###Logout ("/logout")
The logout page resets the client's cookie with a new expired cookie so that the client browser will delete the cookie and deauthenticate. A message is displayed for 10 seconds then the user is redirected to "/", which redirects to "/login"

###Time ("/time")
The time page functions the same as the "/time" page in the Simple Time Server, except the time is addresssed to the user's name if a cookie is present.  

## Flag List

### -v
If "v" is supplied on the command line arguments as true, the server will print the current version number of the server to the command line after launching. The default value is false.

### -port
If "port" is supplied on the command line, it will be used as the starting port for server. By default the server will start on port 8080 if no port is supplied.

