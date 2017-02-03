package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/systemd-timers/systemd"
)

var (
	headers = []string{
		"UNIT",
		"LAST",
		"RESULT",
		"NEXT",
	}
)

func main() {
	conn, err := systemd.NewConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := systemd.NewClient(conn)

	timers, err := client.ListTimers()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, timer := range timers {
		var lastTriggered, result, nextElapse string

		if timer.LastTriggered.IsZero() {
			lastTriggered = "n/a"
			result = "n/a"
		} else {
			lastTriggered = timer.LastTriggered.Local().String()
			result = timer.Result
		}

		if timer.NextElapse.IsZero() {
			nextElapse = "n/a"
		} else {
			nextElapse = timer.NextElapse.Local().String()
		}

		fmt.Fprintln(w, strings.Join([]string{
			timer.UnitName,
			lastTriggered,
			result,
			nextElapse,
		}, "\t"))
	}

	w.Flush()
}
