package machine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// MachineID represents a unique identifier for a machine in the system.
type MachineID string

type MachineStatus string

const (
	Idle          MachineStatus = "idle"
	Starting      MachineStatus = "starting"
	Running       MachineStatus = "running"
	Paused        MachineStatus = "paused"
	Stopped       MachineStatus = "stopped"
	Fault         MachineStatus = "fault"
	Maintenance   MachineStatus = "maintenance"
	EmergencyStop MachineStatus = "emergency_stop"
)

type Machine struct {
	MachineID   MachineID     `json:"machine_id"`
	Name        string        `json:"name"`
	Status      MachineStatus `json:"status"`
	Controllers []*Controller `json:"controllers"`
}

type ControllerType string

const (
	PLC        ControllerType = "plc"
	CNC        ControllerType = "cnc"
	SensorNode ControllerType = "sensor_node"
	Gateway    ControllerType = "gateway"
	EdgeDevice ControllerType = "edge_device"
	RemoteIO   ControllerType = "remote_io"
)

type Controller struct {
	ControllerID string          `json:"controller_id"`
	Type         ControllerType  `json:"type"`
	Vendor       string          `json:"vendor"`
	Model        string          `json:"model"`
	Protocol     string          `json:"protocol"`
	Driver       string          `json:"driver"`
	Address      string          `json:"address"`
	Points       []*Point        `json:"points"`
	DriverConfig json.RawMessage `json:"driver_config"`
	Options      map[string]any  `json:"options"`

	// driver Driver `json:"-"`
}

type DataType string

const (
	BOOL   DataType = "bool"
	INT    DataType = "int"
	FLOAT  DataType = "float"
	STRING DataType = "string"
)

type AccessMode string

const (
	ReadOnly  AccessMode = "read_only"
	WriteOnly AccessMode = "write_only"
	ReadWrite AccessMode = "read_write"
)

type Point struct {
	Name         string          `json:"name"`
	Display      string          `json:"display"`
	Type         DataType        `json:"type"`
	Access       AccessMode      `json:"access"`
	Unit         string          `json:"unit"`
	DriverConfig json.RawMessage `json:"driver_config"`
	Options      map[string]any  `json:"options"`

	value *Value `json:"-"`
}

func (p *Point) Value() *Value {
	if p.value == nil {
		return nil
	}

	return p.value
}

type Value struct {
	Type  DataType  `json:"type"`
	Value any       `json:"value"`
	Time  time.Time `json:"time"`
}

func (v *Value) SetValueWithTime(value, time time.Time) error {
	err := v.SetValue(value)
	if err != nil {
		return err
	}

	v.Time = time
	return nil
}

func (v *Value) SetValue(value any) error {
	switch val := value.(type) {
	case bool:
		v.Type = BOOL
		v.Value = val

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		s := fmt.Sprintf("%v", val)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		v.Type = INT
		v.Value = i

	case float32, float64:
		s := fmt.Sprintf("%v", val)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		v.Type = FLOAT
		v.Value = f

	case string:
		v.Type = STRING
		v.Value = val

	default:
		return errors.New("unsupported value type")
	}

	v.Time = time.Now()

	return nil
}
