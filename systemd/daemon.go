package systemd

import (
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
)

const (
	serviceSuffix   = ".service"
	serviceUnitType = "Service"
	timerSuffix     = ".timer"
	timerUnitType   = "Timer"
)

// Client represents systemd D-Bus API client.
type Client struct {
	conn *dbus.Conn
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

		timer := NewTimer(strings.TrimSuffix(unit.Name, timerSuffix))

		serviceProps, err := c.conn.GetUnitTypeProperties(timer.ServiceName(), serviceUnitType)
		if err != nil {
			return []*Timer{}, err
		}

		var execMainStartTimestamp, execMainExitTimestamp uint64

		if v, ok := serviceProps["ExecMainStartTimestamp"]; ok {
			if v2, ok2 := v.(uint64); ok2 {
				execMainStartTimestamp = v2
			}
		}

		if v, ok := serviceProps["ExecMainExitTimestamp"]; ok {
			if v2, ok2 := v.(uint64); ok2 {
				execMainExitTimestamp = v2
			}
		}

		if execMainStartTimestamp < execMainExitTimestamp {
			timer.LastExecutionTime = execMainExitTimestamp - execMainStartTimestamp
		} else {
			timer.LastExecutionTime = 0
		}

		if v, ok := serviceProps["Result"]; ok {
			if result, ok2 := v.(string); ok2 {
				timer.Result = result
			}
		}

		timerProps, err := c.conn.GetUnitTypeProperties(timer.TimerName(), timerUnitType)
		if err != nil {
			return []*Timer{}, err
		}

		if v, ok := timerProps["TimersCalendar"]; ok {
			if s, ok2 := v.([][]interface{}); ok2 {
				// []interface {}{"OnCalendar", "*-*-* 06,18:00:00", 0x5475da471b800}
				if len(s) > 0 && len(s[0]) > 1 {
					if schedule, ok3 := s[0][1].(string); ok3 {
						timer.Schedule = schedule
					}
				}
			}
		}

		if v, ok := timerProps["LastTriggerUSec"]; ok {
			if lastTriggerUSec, ok2 := v.(uint64); ok2 {
				if lastTriggerUSec == 0 {
					timer.LastTriggered = time.Time{}
				} else {
					timer.LastTriggered = time.Unix(int64(lastTriggerUSec)/1000/1000, 0)
				}
			}
		}

		if v, ok := timerProps["NextElapseUSecRealtime"]; ok {
			if nextElapseUSecRealtime, ok2 := v.(uint64); ok2 {
				if nextElapseUSecRealtime == 0 {
					timer.NextElapse = time.Time{}
				} else {
					timer.NextElapse = time.Unix(int64(nextElapseUSecRealtime)/1000/1000, 0)
				}
			}
		}

		timers = append(timers, timer)
	}

	return timers, nil
}

// Timer represents systemd timer
type Timer struct {
	Name              string    `json:"name"`
	Schedule          string    `json:"schedule"`
	LastTriggered     time.Time `json:"last_triggered"`
	NextElapse        time.Time `json:"next_elapse"`
	Result            string    `json:"result"`
	Active            bool      `json:"active"`
	LastExecutionTime uint64    `json:"last_execution_time"`
}

// NewTimer creates new Timer object
func NewTimer(name string) *Timer {
	return &Timer{
		Name: name,
	}
}

// ServiceName returns service unit name
func (t *Timer) ServiceName() string {
	return t.Name + serviceSuffix
}

// TimerName returns timer unit name
func (t *Timer) TimerName() string {
	return t.Name + timerSuffix
}
