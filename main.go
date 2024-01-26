package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	apiAddress = flag.String("address", "", "NetBox API address")
	apiToken   = flag.String("token", "", "NetBox API token")
	siteID     = flag.String("site", "", "Site ID")
	newPrefix  = flag.String("prefix", "", "Prefix which is going to be added to NetBox")
	prefixTag  = flag.String("tag", "", "Tag which is going to be added to the added prefix")
)

type PrefixPayload struct {
	Prefix string `json:"prefix"`
	Site   string `json:"site"`
	Status string `json:"status"`
	Tags   []int  `json:"tags"`
}

func main() {
	flag.Parse()

	if *apiAddress == "" {
		*apiAddress = os.Getenv("NETBOX_API_ADDRESS")
		if *apiAddress == "" {
			log.Errorln(`Please set the NetBox API address.
			API address can be set with address flag or
			by setting NETBOX_API_ADDRESS environment variable.`)
			os.Exit(1)
		}
	}

	if *apiToken == "" {
		*apiToken = os.Getenv("NETBOX_API_TOKEN")
		if *apiToken == "" {
			log.Errorln(`Please set the NetBox API Token.
			API token can be set with token flag or
			by setting NETBOX_API_TOKEN environment variable.`)
			os.Exit(1)
		}
	}

	if *siteID == "" {
		log.Errorln(`Please set Site ID with site flag.`)
		os.Exit(1)
	}

	if *newPrefix == "" {
		log.Errorln(`Please set the prefix which you want to add
			with prefix flag.`)
		os.Exit(1)
	}

	addPrefixToNetbox()
}

func addPrefixToNetbox() {
	netboxApiToken := "Token " + *apiToken

	tagID := 0
	if *prefixTag != "" {
		tagID = atoi(*prefixTag)
	}

	data := PrefixPayload{
		Prefix: *newPrefix,
		Site:   *siteID,
		Status: "reserved",
		Tags:   []int{tagID},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to Marshal JSON:", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", *apiAddress, body)
	if err != nil {
		log.Println("Failed to build a request:", err)
	}
	req.Header.Set("Authorization", netboxApiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Failed to send a request:", err)
	}
	defer resp.Body.Close()

	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(answer)))
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error converting %s to integer: %v", s, err)
	}
	return i
}
