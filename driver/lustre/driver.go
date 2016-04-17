// +build linux

package lustre

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/mount"
	"github.com/bacaldwell/lustre-graph-driver/driver"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type LustreDriver struct {
	home string
	*LustreSet
}

func init() {
	graphdriver.Register("lustre", Init)
}

func Init(home string, options []string) (graphdriver.Driver, error) {
	if err := os.MkdirAll(home, 0700); err != nil && !os.IsExist(err) {
		log.Errorf("Lustre create home dir %s failed: %v", err)
		return nil, err
	}

	lustreSet, err := NewLustreSet(home, true, options)
	if err != nil {
		return nil, err
	}

	if err := mount.MakePrivate(home); err != nil {
		return nil, err
	}

	d := &LustreDriver{
		LustreSet: lustreSet,
		home:   home,
	}
	return d, nil
}

func (d *LustreDriver) String() string {
	return "lustre"
}

func (d *LustreDriver) Status() [][2]string {
	status := [][2]string{
		{"Pool Objects", ""},
	}
	return status
}

func (d *LustreDriver) GetMetadata(id string) (map[string]string, error) {
	info := d.LustreSet.Devices[id]

	metadata := make(map[string]string)
	metadata["BaseHash"] = info.BaseHash
	metadata["DeviceSize"] = strconv.FormatUint(info.Size, 10)
	metadata["DeviceName"] = info.Device
	return metadata, nil
}

func (d *LustreDriver) Cleanup() error {
	err := d.LustreSet.Shutdown()

	if err2 := mount.Unmount(d.home); err2 == nil {
		err = err2
	}

	return err
}

func (d *LustreDriver) Create(id, parent string) error {
	if err := d.LustreSet.AddDevice(id, parent); err != nil {
		return err
	}
	return nil
}

func (d *LustreDriver) Remove(id string) error {
	if !d.LustreSet.HasDevice(id) {
		return nil
	}

	if err := d.LustreSet.DeleteDevice(id); err != nil {
		return err
	}

	mountPoint := path.Join(d.home, "mnt", id)
	if err := os.RemoveAll(mountPoint); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (d *LustreDriver) Get(id, mountLabel string) (string, error) {
	mp := path.Join(d.home, "mnt", id)

	if err := os.MkdirAll(mp, 0755); err != nil && !os.IsExist(err) {
		return "", err
	}

	if err := d.LustreSet.MountDevice(id, mp, mountLabel); err != nil {
		return "", err
	}

	rootFs := path.Join(mp, "rootfs")
	if err := os.MkdirAll(rootFs, 0755); err != nil && !os.IsExist(err) {
		d.LustreSet.UnmountDevice(id)
		return "", err
	}

	idFile := path.Join(mp, "id")
	if _, err := os.Stat(idFile); err != nil && os.IsNotExist(err) {
		// Create an "id" file with the container/image id in it to help reconscruct this in case
		// of later problems
		if err := ioutil.WriteFile(idFile, []byte(id), 0600); err != nil {
			d.LustreSet.UnmountDevice(id)
			return "", err
		}
	}

	return rootFs, nil
}

func (d *LustreDriver) Put(id string) error {
	if err := d.LustreSet.UnmountDevice(id); err != nil {
		log.Errorf("Warning: error unmounting device %s: %s", id, err)
		return err
	}
	return nil
}

func (d *LustreDriver) Exists(id string) bool {
	return d.LustreSet.HasDevice(id)
}
