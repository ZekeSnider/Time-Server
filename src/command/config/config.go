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
	"path/filepath"
)

var AbsoluteETCPath string
var TemplatePath string
var ServerPort string
var AuthPort string
var AuthHost string
var DumpFile string
var ResponseTime int
var DeviationTime int
var ServerLog string
var AuthLog string
var DisplayVersionBool bool
var MaxConnections int
var CheckPointInterval int

func GetAuthPort() string {
	return AuthPort
}
func init() {
	//declaring command line flags for the time server
	//All flag functionality is detailed in the included README.
	portPointer := flag.String("port", "8080", "Server port number")
	authPortPointer := flag.String("authport", "8090", "Auth Server port number")
	authHostPointer := flag.String("authhost", "localhost", "Auth Server hostname")
	versionBoolPointer := flag.Bool("v", false, "Display server version bool")
	templatePathPointer := flag.String("template", "templates", "Path to templates")
	logPathPointer := flag.String("log", "seelog.xml", "Name of log file")
	dumpFilePointer := flag.String("dumpfile", "dumpfile.json", "Name of dump file")
	checkPointIntervalPointer := flag.Int("checkpointinterval", 30000, "Interval to save logins")
	responseTimePointer := flag.Int("response", 30, "Average simulated response time")
	deviationTimePointer := flag.Int("deviation", 30, "Average simulated deviation time")
	maxConnectionsPointer := flag.Int("maxinflight", 0, "maximum number of inflight requests handled")
	authLogPointer := flag.String("authlog", "authlog.xml", "Name of auth log")

	//parsing the flags
	flag.Parse()

	//setting up the logging library to load configuration from the specified file
	AbsoluteETCPath, _ = filepath.Abs("etc/")

	DisplayVersionBool =  *versionBoolPointer
	TemplatePath = *templatePathPointer
	ResponseTime = *responseTimePointer
	DeviationTime = *deviationTimePointer
	MaxConnections = *maxConnectionsPointer
	DumpFile = AbsoluteETCPath + "/" + *dumpFilePointer
	CheckPointInterval = *checkPointIntervalPointer
	AuthHost = "http://" + *authHostPointer
	ServerPort = ":" + *portPointer
	AuthPort = ":" + *authPortPointer
	ServerLog = AbsoluteETCPath + "/" + *logPathPointer
	AuthLog = AbsoluteETCPath + "/" + *authLogPointer	


}
