package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/influxdb/influxdb/client"
	flag "github.com/ogier/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	MyHost        = "localhost"
	MyPort        = 8086
	MyDB          = "statsd"
	MyMeasurement = "shapes"
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

func query(query string) []int64 {
	ret := []int64{}
	res, err := queryDB(query)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range res[0].Series[0].Values {
		_, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		val, _ := row[1].(json.Number).Int64()
		ret = append(ret, val)
		//log.Printf("[%2d] %s: %d\n", i, t.Format(time.Stamp), val)
	}
	return ret
}

var err error
var con *client.Client

type Trigger struct {
	Operator string
	Value    int64
}

type Alert struct {
	Name     string
	Type     string
	Function string
	Limit    int
	Query    string
	Trigger  Trigger
}

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

func main() {
	var file *string = flag.StringP("config", "c", "", "Config file to use")
	flag.Parse()

	setupInflux()

	alerts := []Alert{}

	data, _ := ioutil.ReadFile(*file)
	err := yaml.Unmarshal(data, &alerts)
	if err != nil {
		panic(err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("%+v\n", alerts)
	}

	for _, alert := range alerts {
		if os.Getenv("DEBUG") == "true" {
			fmt.Println("Query: ", fmt.Sprintf("%s limit %d", alert.Query, alert.Limit))
		}
		values := query(fmt.Sprintf("%s limit %d", alert.Query, alert.Limit))
		var applied_function float64
		if alert.Function == "average" {
			applied_function = 0
			for _, i := range values {
				applied_function += float64(i)
			}
			applied_function = applied_function / float64(alert.Limit)
		}
		if os.Getenv("DEBUG") == "true" {
			fmt.Println("Applied Func: ", applied_function)
		}
		alert_triggered := false
		switch alert.Trigger.Operator {
		case ">":
			alert_triggered = float64(alert.Trigger.Value) > applied_function
		case "<":
			alert_triggered = applied_function < float64(alert.Trigger.Value)
		}
		if alert_triggered {
			color.Red(fmt.Sprintf("[!] %s triggered!", alert.Name))
		} else {
			color.Green(fmt.Sprintf("[+] %s passed.", alert.Name))
		}

	}
}
