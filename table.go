package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/systemd-timers/systemd"
	"github.com/reconquest/loreley"
)

var (
	headers = []string{
		"UNIT",
		"LAST",
		"RESULT",
		"NEXT",
		"SCHEDULE",
	}
)

func generateTable(timers []*systemd.Timer) (string, error) {
	buf := &bytes.Buffer{}

	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, timer := range timers {
		var lastTriggered, result, nextElapse string

		if timer.LastTriggered.IsZero() {
			lastTriggered = "n/a"
			result = "n/a"
		} else {
			lastTriggered = timer.LastTriggered.Local().String()

			if timer.Result == "success" {
				result = "<fg 2>success<reset>"
			} else {
				result = "<fg 1>failed<reset>"
			}
		}

		if timer.NextElapse.IsZero() {
			nextElapse = "n/a"
		} else {
			nextElapse = timer.NextElapse.Local().String()
		}

		fmt.Fprintln(w, strings.Join([]string{
			timer.Name,
			lastTriggered,
			result,
			nextElapse,
			timer.Schedule,
		}, "\t"))
	}

	w.Flush()

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	table, err := loreley.CompileAndExecuteToString(
		buf.String(),
		nil,
		nil,
	)
	if err != nil {
		return "", err
	}

	return table, nil
}
