package nmmodules

type WifiConn struct {
	ID        string   `json:"id"`
	SSID      string   `json:"name"`
	Strength  int      `json:"strength"`
	Mode      string   `json:"mode"`
	Frequency string   `json:"frequency"`
	Security  []string `json:"security"`
	Active    bool     `json:"active"`
}

type Ip struct {
	Addr string
	Mask int
}
type IPConfiguration struct {
	Ips      []Ip
	Gw       []Ip
	DnsAddrs []Ip
}

type WifiManager interface {
	List() ([]WifiConn, error)
	Save(ssid, pass string) error
}
