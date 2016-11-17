package main

import (
	"fmt"
	"math"
	"os"

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

		for _, n := range alert.Notifiers {
			fmt.Printf("<-> Alert sending: %+v\n", alert)
			n.Run(message)
		}

	} else {
		color.Green(fmt.Sprintf("[+] %s passed. (%.2f)", alert.Name, applied_function))
	}

}
