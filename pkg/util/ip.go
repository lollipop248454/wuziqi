package util

import (
	"io/ioutil"
	"log"
	"net/http"
)

func AskIpAddr(ip string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://43.142.132.19:5000?ip="+ip, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ip)
	log.Printf("%s\n", bodyText)
}
