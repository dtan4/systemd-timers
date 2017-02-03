package systemd

import (
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
)

const (
	timerSuffix   = ".timer"
	timerUnitType = "Timer"
)

// Client represents systemd D-Bus API client.
type Client struct {
	conn *dbus.Conn
}

// Timer represents systemd timer
type Timer struct {
	UnitName      string    `json:"unit_name"`
	Schedule      string    `json:"schedule"`
	LastTriggered time.Time `json:"last_triggered"`
	NextElapse    time.Time `json:"next_elapse"`
	Result        string    `json:"result"`
	Active        bool      `json:"active"`
}

// NewClient creates new Client object
func NewClient(conn *dbus.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// NewConn establishes a new connection to D-Bus
func NewConn() (*dbus.Conn, error) {
	return dbus.New()
}

// ListTimers returns installed systemd timers
func (c *Client) ListTimers() ([]*Timer, error) {
	units, err := c.conn.ListUnits()
	if err != nil {
		return []*Timer{}, err
	}

	timers := []*Timer{}

	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, timerSuffix) {
			continue
		}

		timer := &Timer{
			UnitName: unit.Name,
			Active:   unit.ActiveState == "active",
		}

		props, err := c.conn.GetUnitTypeProperties(unit.Name, timerUnitType)
		if err != nil {
			return []*Timer{}, err
		}

		if v, ok := props["TimersCalendar"]; ok {
			if s, ok2 := v.([][]interface{}); ok2 {
				// []interface {}{"OnCalendar", "*-*-* 06,18:00:00", 0x5475da471b800}
				if len(s) > 0 && len(s[0]) > 1 {
					if schedule, ok3 := s[0][1].(string); ok3 {
						timer.Schedule = schedule
					}
				}
			}
		}

		if v, ok := props["LastTriggerUSec"]; ok {
			if lastTriggerUSec, ok2 := v.(uint64); ok2 {
				if lastTriggerUSec == 0 {
					timer.LastTriggered = time.Time{}
				} else {
					timer.LastTriggered = time.Unix(int64(lastTriggerUSec)/1000/1000, 0)
				}
			}
		}

		if v, ok := props["NextElapseUSecRealtime"]; ok {
			if nextElapseUSecRealtime, ok2 := v.(uint64); ok2 {
				if nextElapseUSecRealtime == 0 {
					timer.NextElapse = time.Time{}
				} else {
					timer.NextElapse = time.Unix(int64(nextElapseUSecRealtime)/1000/1000, 0)
				}
			}
		}

		if v, ok := props["Result"]; ok {
			if result, ok2 := v.(string); ok2 {
				timer.Result = result
			}
		}

		timers = append(timers, timer)
	}

	return timers, nil
}
