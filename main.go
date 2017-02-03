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
		"NEXT",
		"LAST",
		"UNIT",
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
		fmt.Fprintln(w, strings.Join([]string{
			timer.NextElapse.Local().String(),
			timer.LastTriggered.Local().String(),
			timer.UnitName,
		}, "\t"))
	}

	w.Flush()
}
