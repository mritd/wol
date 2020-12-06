package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"text/template"

	"github.com/mritd/logger"
	"gopkg.in/yaml.v2"

	"fmt"

	"strings"
)

type Machine struct {
	Name               string `yaml:"name"`
	Mac                string `yaml:"mac"`
	BroadcastInterface string `yaml:"broadcast_interface,omitempty"`
	BroadcastIP        string `yaml:"broadcast_ip,omitempty"`
	Port               int    `yaml:"port,omitempty"`
}

func (m *Machine) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawMachine Machine
	raw := rawMachine{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.BroadcastIP == "" {
		raw.BroadcastIP = "255.255.255.255"
	}
	if raw.Port == 0 {
		raw.Port = 7
	}

	*m = Machine(raw)
	return nil
}

func (m *Machine) MarshalYAML() (interface{}, error) {
	if m.BroadcastIP == "" {
		m.BroadcastIP = "255.255.255.255"
	}
	if m.Port == 0 {
		m.Port = 7
	}
	m.Mac = strings.ToUpper(m.Mac)
	return m, nil
}

// Copy from https://github.com/sabhiram/go-wol/blob/4fd002b5515afaf46b3fe9a9b24ef8c245944f36/cmd/wol/wol.go#L39
func (m *Machine) ipFromInterface(iface string) (*net.UDPAddr, error) {
	ief, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	addrs, err := ief.Addrs()
	if err == nil && len(addrs) <= 0 {
		err = fmt.Errorf("no address associated with interface %s", iface)
	}
	if err != nil {
		return nil, err
	}

	// Validate that one of the addr's is a valid network IP address.
	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			// Verify that the DefaultMask for the address we want to use exists.
			if ip.IP.DefaultMask() != nil {
				return &net.UDPAddr{
					IP: ip.IP,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no address associated with interface %s", iface)
}

func (m *Machine) Wake() error {
	var localAddr *net.UDPAddr
	var err error
	if strings.TrimSpace(m.BroadcastInterface) != "" {
		localAddr, err = m.ipFromInterface(m.BroadcastInterface)
		if err != nil {
			return err
		}
	}

	broadcastAddr := fmt.Sprintf("%s:%d", m.BroadcastIP, m.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return err
	}
	p, err := New(m.Mac)
	if err != nil {
		return err
	}
	bs, err := p.Marshal()
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	n, err := conn.Write(bs)
	if err != nil {
		return err
	}
	if n != 102 {
		logger.Warnf("Magic packet sent was %d bytes (expected 102 bytes sent)", n)
	} else {
		logger.Infof("Magic packet sent successfully to %s", m.Mac)
	}
	return nil
}

type WolConfig struct {
	configPath string
	Machines   []*Machine `yaml:"machines"`
}

func (cfg *WolConfig) SetConfigPath(configPath string) {
	cfg.configPath = configPath
}

func (cfg *WolConfig) Write() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfg.configPath, out, 0644)
}

func (cfg *WolConfig) WriteTo(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.Write()
}

func (cfg *WolConfig) Load() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	buf, err := ioutil.ReadFile(cfg.configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, cfg)
}

func (cfg *WolConfig) LoadFrom(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.Load()
}

func (cfg *WolConfig) FindMachine(m *Machine) (int, *Machine) {
	for i, ma := range cfg.Machines {
		if m.Name == ma.Name {
			return i, ma
		}
	}

	for i, ma := range cfg.Machines {
		if m.Mac == ma.Mac {
			return i, ma
		}
	}
	return 0, nil
}

func (cfg *WolConfig) AddMachine(m *Machine) error {
	_, fm := cfg.FindMachine(m)
	if fm != nil {
		return fmt.Errorf("machine [%v] already exist", m)
	}
	cfg.Machines = append(cfg.Machines, m)
	return cfg.Write()
}

func (cfg *WolConfig) DelMachine(m *Machine) error {
	idx, fm := cfg.FindMachine(m)
	if fm == nil {
		return fmt.Errorf("not found machine [%v]", m)
	}
	cfg.Machines[idx] = cfg.Machines[len(cfg.Machines)-1]
	cfg.Machines[len(cfg.Machines)-1] = nil
	cfg.Machines = cfg.Machines[:len(cfg.Machines)-1]
	return cfg.Write()
}

func (cfg *WolConfig) Print() error {
	tpl := `Name            Mac
---------------------------------
{{range . }}{{ .Name | ListLayout }}{{ .Mac | ToUpper }}
{{end}}`
	t, err := template.New("").Funcs(map[string]interface{}{
		"ListLayout": ListLayout,
		"ToUpper":    strings.ToUpper,
	}).Parse(tpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, cfg.Machines)
	if err != nil {
		return err
	}
	fmt.Println(buf.String())
	return nil
}

func ListLayout(name string) string {
	if len(name) < 16 {
		return fmt.Sprintf("%-16s", name)
	} else {
		return name[:16]
	}
}

func ExampleConfig() string {
	out, _ := yaml.Marshal(WolConfig{
		Machines: []*Machine{
			{
				Name: "iMac",
				Mac:  "e0:d5:5e:6e:30:c9",
			},
		},
	})
	return string(out)
}
