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

type Device struct {
	Name               string `yaml:"name"`
	Mac                string `yaml:"mac"`
	BroadcastInterface string `yaml:"broadcast_interface,omitempty"`
	BroadcastIP        string `yaml:"broadcast_ip,omitempty"`
	Port               int    `yaml:"port,omitempty"`
}

func (d *Device) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawDevice Device
	raw := rawDevice{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.BroadcastIP == "" {
		raw.BroadcastIP = "255.255.255.255"
	}
	if raw.Port == 0 {
		raw.Port = 7
	}

	*d = Device(raw)
	return nil
}

// Copy from https://github.com/sabhiram/go-wol/blob/4fd002b5515afaf46b3fe9a9b24ef8c245944f36/cmd/wol/wol.go#L39
func (d *Device) ipFromInterface(iface string) (*net.UDPAddr, error) {
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

func (d *Device) Wake() error {
	var localAddr *net.UDPAddr
	var err error
	if strings.TrimSpace(d.BroadcastInterface) != "" {
		localAddr, err = d.ipFromInterface(d.BroadcastInterface)
		if err != nil {
			return err
		}
	}

	broadcastAddr := fmt.Sprintf("%s:%d", d.BroadcastIP, d.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return err
	}
	p, err := New(d.Mac)
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
		logger.Infof("Magic packet sent successfully to %s", d.Mac)
	}
	return nil
}

type WolConfig struct {
	configPath string
	Devices    []*Device `yaml:"devices"`
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

func (cfg *WolConfig) FindDevice(d *Device) (int, *Device) {
	for i, dev := range cfg.Devices {
		if d.Name == dev.Name {
			return i, dev
		}
	}

	for i, dev := range cfg.Devices {
		if d.Mac == dev.Mac {
			return i, dev
		}
	}
	return 0, nil
}

func (cfg *WolConfig) AddDevice(d *Device) error {
	_, fm := cfg.FindDevice(d)
	if fm != nil {
		return fmt.Errorf("device [%v] already exist", d)
	}
	cfg.Devices = append(cfg.Devices, d)
	return cfg.Write()
}

func (cfg *WolConfig) DelDevice(d *Device) error {
	idx, fm := cfg.FindDevice(d)
	if fm == nil {
		return fmt.Errorf("not found device [%v]", d)
	}
	cfg.Devices[idx] = cfg.Devices[len(cfg.Devices)-1]
	cfg.Devices[len(cfg.Devices)-1] = nil
	cfg.Devices = cfg.Devices[:len(cfg.Devices)-1]
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
	err = t.Execute(&buf, cfg.Devices)
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
		Devices: []*Device{
			{
				Name: "iMac",
				Mac:  "e0:d5:5e:6e:30:c9",
			},
		},
	})
	return string(out)
}
