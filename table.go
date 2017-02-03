package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/systemd-timers/systemd"
	"github.com/reconquest/loreley"
)

const (
	oneSecond = 1000 * 1000
	oneMinute = 60 * oneSecond
)

var (
	headers = []string{
		"UNIT",
		"LAST",
		"RESULT",
		"EXECUTION TIME",
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
			formatExecutionTime(timer.LastExecutionTime),
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

func formatExecutionTime(executionTime uint64) string {
	if executionTime == 0 {
		return "n/a"
	}

	if executionTime < oneSecond {
		return "Less than a second"
	}

	if executionTime < 2*oneSecond {
		return "1 second"
	}

	if executionTime < 60*oneSecond {
		return fmt.Sprintf("%s seconds", strconv.Itoa(int(executionTime/oneSecond)))
	}

	if executionTime < 61*oneSecond {
		return "1 minute"
	}

	if executionTime < 62*oneSecond {
		return "1 minute 1 second"
	}

	if executionTime < 2*oneMinute {
		return fmt.Sprintf("1 minute %s seconds", strconv.Itoa(int((executionTime-oneMinute)/oneSecond)))
	}

	if (executionTime-oneMinute)/oneSecond%60 < 1 {
		return fmt.Sprintf("%s minutes", strconv.Itoa(int(executionTime/oneSecond/60)))
	}

	if (executionTime-oneMinute)/oneSecond%60 < 2 {
		return fmt.Sprintf("%s minutes 1 second", strconv.Itoa(int(executionTime/oneSecond/60)))
	}

	return fmt.Sprintf("%s minutes %s seconds", strconv.Itoa(int(executionTime/oneSecond/60)), strconv.Itoa(int((executionTime-oneMinute)/oneSecond%60)))
}
