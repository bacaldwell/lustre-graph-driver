// +build linux

package lustre

import (
	"fmt"
	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/go-units"
	"strings"
	"sync"
)

const (
	DefaultDockerBaseImageSize = 10 * 1024 * 1024 * 1024
	DefaultMetaObjectDataSize  = 256
)

type LustreMappingInfo struct {
	Pool     string `json:"pool"`
	Name     string `json:"name"`
	Snapshot string `json:"snap"`
	Device   string `json:"device"`
}

type DevInfo struct {
	Hash        string `json:"hash"`
	Device      string `json:"-"`
	Size        uint64 `json:"size"`
	BaseHash    string `json:"base_hash"` //for delete snapshot
	Initialized bool   `json:"initialized"`

	mountCount int        `json:"-"`
	mountPath  string     `json:"-"`
	lock       sync.Mutex `json:"-"`
}

type MetaData struct {
	Devices     map[string]*DevInfo `json:"Devices"`
	devicesLock sync.Mutex          `json:"-"` // Protects all read/writes to Devices map
}

type LustreSet struct {
	MetaData

	//Options
	baseImageName string
	baseImageSize uint64

	filesystem   string
	mkfsArgs     []string
}
func NewLustreSet(root string, doInit bool, options []string) (*LustreSet, error) {
	devices := &LustreSet{
		MetaData:      MetaData{Devices: make(map[string]*DevInfo)},
		baseImageName: "base_image",
		baseImageSize: DefaultDockerBaseImageSize,
		filesystem:    "ext4",
	}

	for _, option := range options {
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil {
			return nil, err
		}

		key = strings.ToLower(key)

		switch key {
		case "lustre.basesize":
			size, err := units.RAMInBytes(val)
			if err != nil {
				return nil, err
			}
			devices.baseImageSize = uint64(size)
		case "lustre.fs":
			if val != "ext4" && val != "xfs" {
				return nil, fmt.Errorf("Unsupported filesystem %s\n", val)
			}
			devices.filesystem = val
		case "lustre.mkfsarg":
			devices.mkfsArgs = append(devices.mkfsArgs, val)
		default:
			return nil, fmt.Errorf("Unknown option %s\n", key)
		}
	}

	//if err := devices.initLustreSet(doInit); err != nil {
	//	return nil, err
	//}
	return devices, nil
}
