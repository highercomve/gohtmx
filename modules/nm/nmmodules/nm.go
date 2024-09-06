package nmmodules

type WifiConn struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Strength  int32    `json:"strength"`
	Frequency int      `json:"frequency"`
	BitRate   int      `json:"bitrate"`
	Security  []string `json:"security"`
	Saved     bool     `json:"saved"`
}

type Ip struct {
	Addr string
	Mask int32
}
type IPConfiguration struct {
	Ips      []Ip
	Gw       []Ip
	DnsAddrs []Ip
}

type WifiManager interface {
	List() ([]WifiConn, error)
	Save(config *WifiConn) error
}
