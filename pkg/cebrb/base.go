package cebrb

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
	hosts       []string
	sentInit    = false
	sentWitness = false
	delivered   = false
	witnesses   map[string]string
	n           int
	f           int
)

func clear() {
	sentInit = false
	sentWitness = false
	delivered = false
	witnesses = make(map[string]string)
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
	_, err := http.Get(fmt.Sprintf("%s/cebrb/clear", host))
	if err != nil {
		//log.Fatalf("unable to reach host '%s'", host)
	}
}

func makeRequest(rType, message, from, host string) {
	//log.Println(fmt.Sprintf("%s -> %s - %s %s", from, host, rType, message))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cebrb/deliver/%s", host, rType), nil)
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
		makeRequest("init", message, sourceHost, host)
	}

	w.WriteHeader(http.StatusOK)
}

func Clear(_ http.ResponseWriter, _ *http.Request) {
	clear()
}

func Init(_ http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	from := r.URL.Query().Get("from")

	if from == sourceHost && !sentInit {
		sentInit = true
		for _, host := range hosts {
			makeRequest("witness", message, sourceHost, host)
		}
	}
}

func Witness(_ http.ResponseWriter, r *http.Request) {
	initMessage := r.URL.Query().Get("message")
	from := r.URL.Query().Get("from")

	if _, ok := witnesses[from]; !ok {
		witnesses[from] = initMessage
	}

	messagesCnt := make(map[string]int)

	for _, host := range hosts {
		messagesCnt[witnesses[host]]++
	}

	for message, cnt := range messagesCnt {
		if message != "" && cnt >= n-2*f && !sentWitness {
			sentWitness = true

			for _, host := range hosts {
				makeRequest("witness", message, sourceHost, host)
			}
		}

		if message != "" && cnt >= n-f && !delivered {
			delivered = true
			log.Printf("%s delivered message %s", sourceHost, "0") // replace 0 with message
		}
	}
}
