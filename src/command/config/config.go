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

package config

import (
	"flag"
	log "github.com/cihub/seelog"
	"path/filepath"
)

var AbsoluteETCPath string
var TemplatePath string
var ServerPort string

var Logger log.LoggerInterface

func init() {
	//declaring command line flags for the time server
	//All flag functionality is detailed in the included README.
	portPointer := flag.String("port", "8080", "Server port number")
	versionBoolPointer := flag.Bool("v", false, "Display server version bool")
	templatePathPointer := flag.String("template", "templates", "Path to templates")
	logPathPointer := flag.String("log", "seelog.xml", "Name of log file")

	//parsing the flags
	flag.Parse()

	//setting up the logging library to load configuration from the specified file
	AbsoluteETCPath, _ = filepath.Abs("etc/")

	Logger, err := log.LoggerFromConfigAsFile(AbsoluteETCPath + "/" + *logPathPointer)
	if err != nil {
		log.Errorf("Error loading the seelog config", err)
	}
	log.ReplaceLogger(Logger)

	defer log.Flush()


	//logging startup and replacing the default logger
	log.Infof("The log file path is %s", AbsoluteETCPath+"/"+*logPathPointer)


	//Outputting server version number if it is requested in command line flags
	if *versionBoolPointer == true {
		log.Infof("Personalized Time Server version 1.3")
	}


	//adding a ":" to the port number to match the format requested by http.ListenAndServe
	ServerPort := ":" + *portPointer
	_ = ServerPort

	TemplatePath = *templatePathPointer
}
