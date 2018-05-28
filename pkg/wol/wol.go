package wol

import (
	"fmt"
	"net"

	"log"

	"github.com/mritd/wol/pkg/mp"
	"github.com/mritd/wol/pkg/utils"
)

type Machine struct {
	Mac                string
	BroadcastInterface string
	BroadcastIP        string
	Port               int
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

func (m *Machine) Wake() {
	localAddr, err := m.ipFromInterface(m.BroadcastInterface)
	utils.CheckAndExit(err)
	broadcastAddr := fmt.Sprintf("%s:%d", m.BroadcastIP, m.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	utils.CheckAndExit(err)
	p, err := mp.New(m.Mac)
	utils.CheckAndExit(err)
	bs, err := p.Marshal()
	utils.CheckAndExit(err)
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	utils.CheckAndExit(err)
	defer conn.Close()
	n, err := conn.Write(bs)
	utils.CheckAndExit(err)
	if n != 102 {
		log.Printf("Magic packet sent was %d bytes (expected 102 bytes sent)", n)
	} else {
		log.Printf("Magic packet sent successfully to %s\n", m.Mac)
	}
}
