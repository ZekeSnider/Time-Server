# Assignment 4: AuthServer
Zeke Snider  
CSS 490C  
2/20/15  

## Overview
This program has added an authentication server to "Assignment 3: Refactoring". The time server still sets cookies, but requests to check the usermap are now sent via http requests the authserver. There is also a config package which is used to load the command line flags. More features are added as seen in the flags section below.


## Build Instructions
The included makefile should be used to build and the run the project. "make build", should be used to build, then "./bin/TimeServer" can be used to run the time server, and "./bin/AuthServer" to run the authserver. I moved all the flag parsing to the config.go package, but for some reason I wasn't able to get both the auth server and time server to load the same flags. Thus, you must supply the appropriate flags to both executables on launch. 

## Design Notes
The seelog library has been used to replace the old logging system. The log output file is stored at "/etc/timeserver.log". Errors are logged using seelog error logging. Regular requests or other traces are logged using debug logging. General info are logged using info log.

The following template files are required for proper functionality of the server:  

* error.tmpl
* home.tmpl
* login.tmpl
* logout.tmpl
* menu.tmpl
* time.tmpl  

The menu.tmpl is the main template file, and the other files are used as subtemplates.  

All cookie functionality has been refactored into checkUserCookie, setUserCookie logoutCookie, which manage login/logout and session management.

The sync library is used to lock and unlock the internal UUID name map so that the map does not go out of sync between multiple requests at the same time.  

When the server stops, the server's internal map is deleted but the client's cookie still remains unless they logout before the server stops. 


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

### -log
If "log" is supplied on the command line, it will be used as the path to the seelog configuration. For example, if -log="logconfig.xml" is suppliod, the log will be loaded from /etc/logconfig.xml. The default value is seelog.xml

### -authlog
Same as -log, except used for the auth server. Default value is authlog.xml

### -authport, -authhost
Specifies where the auth server should start, and where the time server should look for it. Default value is http://localhost:8090

### -response, -deviation
Each request to the time server is now throttled to simulate a delay in generating the time. The throttle delay is generated using a random number from a normal distribution of response and deviation (in ms). 

### -maxinflight
Configures the maximum number of requests that the time server should handle at once. Whenever a request is handled, a protected counter is checked to see if under the maxinflight number. If it is under the limit, the counter is incremented, the request is handled, then the counter is decremented. If the server cannot handle the request, an error is thrown.

### -dumpfile, -checkpointinterval
Specifies where the auth server should backup the userlist to (json file), and how often it should do so. The auth server starts another go thread which will backup the userlist every checkpointinterval ms while the server is running. This exported file is checked for integrity. When the server starts the next time the user list is loaded back from the file.

### -template
If "template" is supplied on the command line, it will be used as the path to the templates directory. For example, if -log="sampletemplates" is suppliod, the templates will be loaded from /etc/samepletemplates. The default value is templates.
