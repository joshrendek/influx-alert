package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
	"time"
)

func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

func query(query string) []float64 {
	ret := []float64{}
	res, err := queryDB(query)
	if err != nil {
		log.Fatal(err)
	}
	for i, row := range res[0].Series[0].Values {
		t, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		val, _ := row[1].(json.Number).Float64()
		ret = append(ret, val)
		if os.Getenv("DEBUG") == "true" {
			log.Printf("[%2d] %s: %d\n", i, t.Format(time.Stamp), val)
		}
	}
	return ret
}

var con *client.Client

func setupInflux() {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", MyHost, MyPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: "root",
		Password: "root",
	}

	con, err = client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("DEBUG") == "true" {
		log.Printf("Connected in %v | Version: %s", dur, ver)
	}
}
