package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fatih/color"
)

func (alert *Alert) ApplyFunction(values []float64) float64 {
	var appliedFunction float64

	if len(values) > 0 {
		appliedFunction = values[0]
	}

	if alert.Function == "average" {
		for _, i := range values {
			appliedFunction += float64(i)
		}
		appliedFunction = appliedFunction / float64(len(values))
	} else if alert.Function == "max" {
		for _, i := range values {
			appliedFunction = math.Max(appliedFunction, i)
		}
	} else if alert.Function == "min" {
		for _, i := range values {
			appliedFunction = math.Min(appliedFunction, i)
		}
	}
	return appliedFunction
}

func (alert *Alert) Setup() {
	hash := md5.Sum([]byte(alert.Name))
	alert.Hash = hex.EncodeToString(hash[:])
	for _, n := range alert.NotifiersRaw {
		alert.Notifiers = append(alert.Notifiers, Notifier{Name: n})
	}

}

func (alert *Alert) Run() {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("Query: ", fmt.Sprintf("%s limit %d", alert.Query, alert.Limit))
	}

	groupByQuery := ""
	if len(alert.GroupBy) > 0 {
		groupByQuery = fmt.Sprintf("GROUP BY time(%s)", alert.GroupBy)
	}

	values := query(fmt.Sprintf("%s where time > now() - %s %s limit %d",
		alert.Query, alert.Timeshift, groupByQuery, alert.Limit))

	applied_function := alert.ApplyFunction(values)

	if os.Getenv("DEBUG") == "true" {
		fmt.Println("Applied Func: ", applied_function)
	}

	alert_triggered := false
	switch alert.Trigger.Operator {
	case "gt":
		alert_triggered = applied_function > float64(alert.Trigger.Value)
	case "lt":
		alert_triggered = applied_function < float64(alert.Trigger.Value)
	}

	if alert_triggered {
		message := fmt.Sprintf("*[!] %s triggered!* Value: %.2f | Trigger: %s %d",
			alert.Name, applied_function, alert.Trigger.Operator, alert.Trigger.Value)
		color.Red(message)
		alertAlreadyTriggered := false
		tMutex.Lock()
		if v, ok := triggeredAlerts[alert.Hash]; ok {
			color.Yellow(fmt.Sprintf("[already triggered at %s] %s", v.TriggeredAt, message))
			alertAlreadyTriggered = true
		} else {
			triggeredAlerts[alert.Hash] = TriggeredAlert{Hash: alert.Hash, TriggeredAt: time.Now()}
		}
		tMutex.Unlock()
		if !alertAlreadyTriggered {
			for _, n := range alert.Notifiers {
				n.Run(message, true)
			}
		}

	} else {
		tMutex.Lock()
		if _, ok := triggeredAlerts[alert.Hash]; ok {
			delete(triggeredAlerts, alert.Hash)
			message := fmt.Sprintf("*[+] %s resolved * Value: %.2f | Trigger: %s %d",
				alert.Name, applied_function, alert.Trigger.Operator, alert.Trigger.Value)
			for _, n := range alert.Notifiers {
				n.Run(message, false)
			}
			color.Green("[+] %s - Alert resolved.", alert.Name)
		}
		tMutex.Unlock()
		color.Green(fmt.Sprintf("[+] %s passed. (%.2f)", alert.Name, applied_function))
	}

}
