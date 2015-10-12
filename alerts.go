package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func (alert *Alert) ApplyFunction(values []float64) float64 {
	var applied_function float64
	if alert.Function == "average" {
		applied_function = 0
		for _, i := range values {
			applied_function += float64(i)
		}
		applied_function = applied_function / float64(alert.Limit)
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

	values := query(fmt.Sprintf("%s limit %d", alert.Query, alert.Limit))

	applied_function := alert.ApplyFunction(values)

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
		message := fmt.Sprintf("*[!] %s triggered!* Value: %.2f | Trigger: %s %d",
			alert.Name, applied_function, alert.Trigger.Operator, alert.Trigger.Value)
		color.Red(message)

		for _, n := range alert.Notifiers {
			n.Run(message)
		}

	} else {
		color.Green(fmt.Sprintf("[+] %s passed.", alert.Name))
	}

}
