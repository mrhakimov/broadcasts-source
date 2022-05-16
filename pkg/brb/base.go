package brb

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	sourceHost = "http://localhost:8080"
	//byzantineFactor = 3
)

var (
	hosts     []string
	sentEcho  = false
	sentReady = false
	delivered = false
	echos     map[string]string
	readys    map[string]string
	n         int
	f         int
)

func clear() {
	sentEcho = false
	sentReady = false
	delivered = false
	echos = make(map[string]string)
	readys = make(map[string]string)
}

func init() {
	file, err := os.OpenFile("/Users/mukkhakimov/Documents/itmo/thesis/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Lmicroseconds | log.Ldate)
	log.SetOutput(file)

	data, err := ioutil.ReadFile("/Users/mukkhakimov/Documents/itmo/thesis/hosts.txt")
	if err != nil {
		log.Fatalln("unable to read hosts: ", err)
	}

	allData := strings.Split(string(data), "\n")
	hosts = allData[1:]
	f1, _ := strconv.ParseInt(allData[0], 10, 32)
	f = int(f1)
	//log.Println("hosts: ", hosts)

	n = len(hosts)
	//f = n / byzantineFactor

	clear()
}

func makeClearRequest(host string) {
	_, err := http.Get(fmt.Sprintf("%s/brb/clear", host))
	if err != nil {
		//log.Fatalf("unable to reach host '%s'", host)
	}
}

func makeRequest(rType, message, from, host string) {
	//log.Println(fmt.Sprintf("%s -> %s - %s %s", from, host, rType, message))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/brb/deliver/%s", host, rType), nil)
	if err != nil {
		//log.Fatalln("unable to reach host: ", err)
	}

	q := req.URL.Query()
	q.Add("message", message)
	q.Add("from", from)
	req.URL.RawQuery = q.Encode()

	//log.Println(req.URL.String())

	_, err = http.Get(req.URL.String())
	if err != nil {
		//log.Fatalf("unable to reach host '%s'", host)
	}
}

func Broadcast(w http.ResponseWriter, r *http.Request) {
	log.Println("broadcast started")

	message := r.URL.Query().Get("message")
	if message == "" {
		log.Fatalln("empty messages are not allowed")
	}

	for _, host := range hosts {
		makeClearRequest(host)
	}

	for _, host := range hosts {
		makeRequest("send", message, sourceHost, host)
	}

	w.WriteHeader(http.StatusOK)
}

func Send(_ http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	from := r.URL.Query().Get("from")

	if from == sourceHost && !sentEcho {
		sentEcho = true
		for _, host := range hosts {
			makeRequest("echo", message, sourceHost, host)
		}
	}

	checkReady1()
	checkReady2()
	checkDeliver()
}

func Echo(_ http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	from := r.URL.Query().Get("from")

	if _, ok := echos[from]; !ok {
		echos[from] = message
	}

	checkReady1()
	checkReady2()
	checkDeliver()
}

func Ready(_ http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	from := r.URL.Query().Get("from")

	if _, ok := readys[from]; !ok {
		readys[from] = message
	}

	checkReady1()
	checkReady2()
	checkDeliver()
}

func Clear(_ http.ResponseWriter, _ *http.Request) {
	clear()
}

func checkReady1() {
	messagesCnt := make(map[string]int)

	for _, host := range hosts {
		messagesCnt[echos[host]]++
	}

	for message, cnt := range messagesCnt {
		if message != "" && cnt > (n+f)/2 && !sentReady {
			sentReady = true
			for _, host := range hosts {
				makeRequest("ready", message, sourceHost, host)
			}
		}
	}
}

func checkReady2() {
	messagesCnt := make(map[string]int)

	for _, host := range hosts {
		messagesCnt[readys[host]]++
	}

	for message, cnt := range messagesCnt {
		if message != "" && cnt > f && !sentReady {
			sentReady = true
			for _, host := range hosts {
				makeRequest("ready", message, sourceHost, host)
			}
		}
	}
}

func checkDeliver() {
	messagesCnt := make(map[string]int)

	for _, host := range hosts {
		messagesCnt[readys[host]]++
	}

	for message, cnt := range messagesCnt {
		if message != "" && cnt > 2*f && !delivered {
			delivered = true
			log.Printf("%s delivered message %s", sourceHost, "0") // replace 0 with message

			//clear()
		}
	}
}
