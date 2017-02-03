package main

import (
	"fmt"
	"os"

	"github.com/dtan4/systemd-timers/systemd"
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

	table, err := generateTable(timers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(table)
}
