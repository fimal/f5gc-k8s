package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	//"os/signal"
	//"syscall"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"

	//"reflect"
	"context"
	"strconv" // typecast "size" int as string

	"github.com/Jeffail/gabs/v2"
)

func constructQuery(q string, size int) *strings.Reader {
	// Build a query string from string passed to function
	var query = `{"query": {`

	// Concatenate query string with string passed to method call
	query = query + q

	// Use the strconv.Itoa() method to convert int to string
	query = query + `}, "size": ` + strconv.Itoa(size) + `}`
	//log.Println("\nquery:", query)
	isValid := json.Valid([]byte(query)) // returns bool
	// Default query is "{}" if JSON is invalid
	if isValid == false {
		log.Println("constructQuery() ERROR: query string not valid:", query)
		log.Println("Using default match_all query")
		query = "{}"
	} else {
		//log.Println("constructQuery() valid JSON:", isValid)
	}
	// Build a new string from JSON query
	var b strings.Builder
	b.WriteString(query)

	// Instantiate a *strings.Reader object from string
	read := strings.NewReader(b.String())

	// Return a *strings.Reader object
	return read
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func check() (string, error) {
	var totalValueTenMinBefore, totalValueBase float64
	log.SetFlags(0)
	log.Printf("checking\n")
	ctx1 := context.Background()
	ctx2 := context.Background()
	var buf bytes.Buffer

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return "", fmt.Errorf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		return "", fmt.Errorf("Error getting response: %s", err)
	}
	//log.Println(elasticsearch.Version)
	defer res.Body.Close()
	var tenMinBefore = `"bool":{"must":[{"match":{"resp_code":"500"}},{"match":{"waas_tag":"waas-f5gc-ausf"}},{"match":{"req_path":"/nausf-auth/v1/ue-authentications"}},{"range":{"@timestamp":{"gt":"now-3m","lt":"now"}}}]}`
	read1 := constructQuery(tenMinBefore, 1000)
	if err := json.NewEncoder(&buf).Encode(read1); err != nil {
		return "", fmt.Errorf("json.NewEncoder() ERROR: %v", err)
		// Query is a valid JSON object
	}
	// Pass the JSON query to the Golang client's Search() method
	res, err = es.Search(
		es.Search.WithContext(ctx1),
		es.Search.WithIndex("access*"),
		es.Search.WithBody(read1),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return "", fmt.Errorf("Error: Query failed: %v", err)
	}
	jsonParsed, _ := gabs.ParseJSON(streamToByte(res.Body))
	totalValueTenMinBefore, _ = jsonParsed.Path("hits.total.value").Data().(float64)
	//log.Println("This is 10 minutes period Value", totalValueTenMinBefore)

	var baseLine = `"bool":{"must":[{"match":{"resp_code":"500"}},{"match":{"waas_tag":"waas-f5gc-ausf"}},{"match":{"req_path":"/nausf-auth/v1/ue-authentications"}},{"range":{"@timestamp":{"gt":"now-6m","lt":"now-4m"}}}]}`
	read2 := constructQuery(baseLine, 1000)
	if err := json.NewEncoder(&buf).Encode(read2); err != nil {
		return "", fmt.Errorf("json.NewEncoder() ERROR: %v", err)
		// Query is a valid JSON object
	}
	// Pass the JSON query to the Golang client's Search() method
	res2, err := es.Search(
		es.Search.WithContext(ctx2),
		es.Search.WithIndex("access*"),
		es.Search.WithBody(read2),
		es.Search.WithTrackTotalHits(true),
	)
	defer res2.Body.Close()

	if err != nil {
		return "", fmt.Errorf("Error: Query failed %v", err)
	} else {
		jsonParsed, _ := gabs.ParseJSON(streamToByte(res2.Body))
		totalValueBase, _ = jsonParsed.Path("hits.total.value").Data().(float64)
		totalValueBase = totalValueBase / 10.0
		//log.Println("This is BaseLine Total Value", totalValueBase)
	}

	if totalValueTenMinBefore > totalValueBase {
		return "Attack Detected, IMSI  CRACKING Attack! Please check Location nR-CGI = PLMN: 20893 nRCellIdentity: 0x0000000010 tAC: 0x000001", nil
	}
	return "", nil

}

var (
	elasticAddress  string
	elasticPort     string
	elasticUser     string
	elasticPassword string
	esCfg           elasticsearch.Config
	interval        time.Duration
)

func readEnvVars() {
	address, ok := os.LookupEnv("ELASTIC_ADDRESS")
	if !ok {
		address = "10.195.5.182"
	}
	elasticAddress = address

	port, ok := os.LookupEnv("ELASTIC_PORT")
	if !ok {
		port = "31001"
	}
	elasticPort = port

	user, ok := os.LookupEnv("ELASTIC_USER")
	if !ok {
		log.Fatalf("could not find ELASTIC_USER environment variable")
	}
	elasticUser = user

	password, ok := os.LookupEnv("ELASTIC_PASS")
	if !ok {
		log.Fatalf("could not find ELASTIC_PASS environment variable")
	}
	elasticPassword = password

	esCfg = elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("https://%s:%s", elasticAddress, elasticPort),
		},
		Username: elasticUser,
		Password: elasticPassword,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

type trans struct {
	TransID   string
	TransTime string
}

func main() {

	readEnvVars()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)
	truelyDone := make(chan bool)

	go func() {
		sign := <-signals

		log.Printf("received signal = %v\n", sign)
		done <- true
	}()

	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-done:
				log.Println("stopping ticker")
				ticker.Stop()
				log.Println("ticker stopped. exiting for loop")
				truelyDone <- true
				return
			case <-ticker.C:
				result, err := check()
				if err != nil {
					log.Println(err)
					return
				}
				//post result to elastic
				if result != "" {
					err = postResult(result)
					if err != nil {
						log.Println(err)
					} else {
						log.Println("event posted to elastic")
					}
				} else {
					log.Println("no violations found")
				}
			}
		}
	}()

	fmt.Println("ticker running")
	<-truelyDone
	log.Println("exiting main")

}

func postResult(message string) error {
	//curl -k -XPOST -u elastic:es182:31001/forensics/_doc -d@sample.json -H"Content-Type: application/json"

	// create the document to send

	//send an https request to elastic
	templateFile := "event.tmpl"
	basename := filepath.Base(templateFile)
	tmpl, err := template.New(basename).ParseFiles(templateFile)

	if err != nil {
		return err
	}

	var eventBytes bytes.Buffer

	newuid := uuid.New().String()
	//"2020-12-07T12:33:26.457Z"
	event := trans{
		TransID:   newuid,
		TransTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
	}

	err = tmpl.Execute(&eventBytes, event)
	if err != nil {
		log.Println(eventBytes.String())
		return err
	}

	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	postAddress := fmt.Sprintf("https://%s:%s/forensics/_doc", elasticAddress, elasticPort)
	r, err := http.NewRequest("POST", postAddress, &eventBytes)
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(elasticUser, elasticPassword)

	response, err := client.Do(r)
	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("error, status code returned: %v", response.StatusCode)
	}

	return nil
}

func init() {
	flag.DurationVar(&interval, "interval", 10*time.Second, "interval between checks for IMSI cracking")
	flag.Parse()
}
