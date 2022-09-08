package bpf_map

import (
	"fmt"
	"os"

	"github.com/cilium/ebpf"
)

// err = m.Put(EndpointKey{IP: 6}, EndpointInfo{
// 	IfIndex: 2,
// 	LxcID:   3,
// 	MAC:     4,
// 	NodeMAC: 5,
// })

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetMapByPinned(pinPath string, opts ...*ebpf.LoadPinOptions) *ebpf.Map {
	var options *ebpf.LoadPinOptions
	if len(opts) == 0 {
		options = &ebpf.LoadPinOptions{}
	} else {
		options = opts[0]
	}
	m, err := ebpf.LoadPinnedMap(pinPath, options)
	if err != nil {
		fmt.Println("GetMapByPinned failed: ", err.Error())
	}
	return m
}

func createMap(
	name string,
	_type ebpf.MapType,
	keySize uint32,
	valueSize uint32,
	maxEntries uint32,
	flags uint32,
) (*ebpf.Map, error) {
	spec := ebpf.MapSpec{
		Name:       name,
		Type:       ebpf.Hash,
		KeySize:    keySize,
		ValueSize:  valueSize,
		MaxEntries: maxEntries,
		Flags:      flags,
	}
	m, err := ebpf.NewMap(&spec)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// 该方法在同一节点上调用多次但是只会创建一个同名的 map
func CreateOnceMapWithPin(
	pinPath string,
	name string,
	_type ebpf.MapType,
	keySize uint32,
	valueSize uint32,
	maxEntries uint32,
	flags uint32,
) (*ebpf.Map, error) {
	if PathExists(pinPath) {
		return GetMapByPinned(pinPath), nil
	}
	m, err := createMap(
		name,
		_type,
		keySize,
		valueSize,
		maxEntries,
		flags,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}
