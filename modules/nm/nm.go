package nm

import (
	"github.com/highercomve/gohtmx/modules/nm/nmmodules"
)

type NetworkManager struct{}

func Init() nmmodules.WifiManager {
	return &NetworkManager{}
}

func (nm *NetworkManager) List() (conns []nmmodules.WifiConn, err error) {
	return conns, err
}

func (nm *NetworkManager) Save(conn *nmmodules.WifiConn) error {
	return nil
}
