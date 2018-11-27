package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func splunkAPI(params url.Values, path string) *http.Response {
	encodedData := params.Encode()
	endpointURL := os.Getenv("SPLUNK_API_URL") + path
	r, _ := http.NewRequest("POST", endpointURL, strings.NewReader(encodedData))
	r.SetBasicAuth("admin", os.Getenv("SPLUNK_PASSWORD"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		splunkClient := &http.Client{
			Transport: tr,
		}
		resp, err := splunkClient.Do(r)
		if err != nil {
			log.Printf("http error: %s\n", err.Error())
			time.Sleep(time.Second * 5)
			continue
		}
		return resp
	}
}

func getParams(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)
	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

func enableHEC() {
	response := splunkAPI(url.Values{}, "/servicesNS/admin/splunk_httpinput/data/inputs/http/http/enable")
	if response.StatusCode != 200 && response.StatusCode != 201 && response.StatusCode != 409 {
		log.Fatalf("error enabling HEC: %d", response.StatusCode)
	}
	defer response.Body.Close()
}

func createHECToken() string {
	rand.Seed(time.Now().UTC().UnixNano())
	response := splunkAPI(url.Values{
		"name":  []string{"mytoken_" + fmt.Sprintf("%d", rand.Int())},
		"index": []string{"main"},
	}, "/servicesNS/admin/splunk_httpinput/data/inputs/http")
	if response.StatusCode != 200 && response.StatusCode != 201 && response.StatusCode != 409 {
		log.Fatalf("error creating HEC token: %d", response.StatusCode)
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("error reading response: %s", err.Error())
	}
	body := string(bodyBytes)
	params := getParams(`name=\"token\"\>(?P<token>.+)\<\/s:key\>`, body)
	log.Println(params)
	token := params["token"]
	if token == "" {
		log.Fatalf("could not create token: %s", body)
	}
	log.Println(token)
	return token
}
