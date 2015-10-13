package main

import (
	"fmt"
	"github.com/fatih/color"
	"math"
	"os"
)

func (alert *Alert) ApplyFunction(values []float64) float64 {
	var applied_function float64 = 0
	if alert.Function == "average" {
		applied_function = 0
		for _, i := range values {
			applied_function += float64(i)
		}
		applied_function = applied_function / float64(alert.Limit)
	} else if alert.Function == "max" {
		applied_function = 0
		for _, i := range values {
			applied_function = math.Max(applied_function, i)
		}
	} else if alert.Function == "min" {
		applied_function = 0
		for _, i := range values {
			applied_function = math.Min(applied_function, i)
		}
	}
	return applied_function
}

func (alert *Alert) Run() {
	for _, n := range alert.NotifiersRaw {
		alert.Notifiers = append(alert.Notifiers, Notifier{Name: n})
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Println("Query: ", fmt.Sprintf("%s limit %d", alert.Query, alert.Limit))
	}

	values := query(fmt.Sprintf("%s where time > now() - %s limit %d", alert.Query, alert.Timeshift, alert.Limit))

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
			n.Run(message)
		}

	} else {
		color.Green(fmt.Sprintf("[+] %s passed. (%.2f)", alert.Name, applied_function))
	}

}
