// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package tinydhcp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

type Lease struct {
	IP         net.IP
	MAC        net.HardwareAddr
	ExpiryTime time.Time
}

type IPPool struct {
	Subnet      *net.IPNet
	StartIP     net.IP
	EndIP       net.IP
	LeaseTime   time.Duration
	AssignedIPs map[string]Lease
	mu          sync.Mutex
}

type StaticBinding struct {
	MAC net.HardwareAddr
	IP  net.IP
}

type DHCPServer struct {
	IPPools        []*IPPool
	StaticBindings map[string]StaticBinding
	Log            []string
	Gateway        net.IP
}

func nextIP(ip net.IP) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)
	for j := len(result) - 1; j >= 0; j-- {
		result[j]++
		if result[j] != 0 {
			break
		}
	}
	return result
}

func (s *DHCPServer) AddIPPool(subnet string, startIP, endIP string, leaseTime time.Duration) error {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}

	pool := IPPool{
		Subnet:      ipnet,
		StartIP:     net.ParseIP(startIP),
		EndIP:       net.ParseIP(endIP),
		LeaseTime:   leaseTime,
		AssignedIPs: make(map[string]Lease),
	}

	s.IPPools = append(s.IPPools, &pool)
	return nil
}

func (s *DHCPServer) RemoveIPPool(subnet string) error {
	for i, pool := range s.IPPools {
		if pool.Subnet.String() == subnet {
			s.IPPools = append(s.IPPools[:i], s.IPPools[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("subnet not found")
}

func (s *DHCPServer) AssignLease(mac net.HardwareAddr) (net.IP, error) {
	for i, pool := range s.IPPools {
		pool.mu.Lock()
		defer pool.mu.Unlock()

		for ip := pool.StartIP; !ip.Equal(pool.EndIP); ip = nextIP(ip) {
			if _, assigned := pool.AssignedIPs[ip.String()]; !assigned {
				lease := Lease{
					IP:         ip,
					MAC:        mac,
					ExpiryTime: time.Now().Add(pool.LeaseTime),
				}
				pool.AssignedIPs[ip.String()] = lease
				s.IPPools[i] = pool
				return ip, nil
			}
		}
	}
	return nil, fmt.Errorf("no available IPs")
}

func (s *DHCPServer) ReleaseLease(mac net.HardwareAddr) {
	for i, pool := range s.IPPools {
		pool.mu.Lock()
		defer pool.mu.Unlock()

		for ip, lease := range pool.AssignedIPs {
			if lease.MAC.String() == mac.String() {
				delete(pool.AssignedIPs, ip)
				s.IPPools[i] = pool
				return
			}
		}
	}
}

func (s *DHCPServer) AddStaticBinding(mac net.HardwareAddr, ip net.IP) error {
	s.StaticBindings[mac.String()] = StaticBinding{
		MAC: mac,
		IP:  ip,
	}
	return nil
}

func (s *DHCPServer) RemoveStaticBinding(mac net.HardwareAddr) error {
	delete(s.StaticBindings, mac.String())
	return nil
}

func (s *DHCPServer) LogEvent(event string) {
	s.Log = append(s.Log, fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), event))
}

func (s *DHCPServer) GetLeaseHistory() []Lease {
	leases := []Lease{}
	for _, pool := range s.IPPools {
		for _, lease := range pool.AssignedIPs {
			leases = append(leases, lease)
		}
	}
	return leases
}

func (s *DHCPServer) StartServer(interfaceName string, port int) {

	srv, err := server4.NewServer(interfaceName,
		&net.UDPAddr{
			IP: net.ParseIP("0.0.0.0"), Port: port,
		},
		func(conn net.PacketConn, peer net.Addr, pkt *dhcpv4.DHCPv4) {
			mac := pkt.ClientHWAddr
			s.LogEvent(fmt.Sprintf("Request from MAC %s", mac))

			if static, found := s.StaticBindings[mac.String()]; found {
				response, _ := dhcpv4.NewReplyFromRequest(pkt)
				response.YourIPAddr = static.IP
				response.Options.Update(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
				response.Options.Update(dhcpv4.OptRouter(s.Gateway))
				return
			}

			ip, err := s.AssignLease(mac)
			if err != nil {
				s.LogEvent(fmt.Sprintf("No IP available for MAC %s", mac))
				return
			}
			response, _ := dhcpv4.NewReplyFromRequest(pkt)
			response.YourIPAddr = ip
			response.Options.Update(dhcpv4.OptSubnetMask(net.IPv4Mask(255, 255, 255, 0)))
			response.Options.Update(dhcpv4.OptRouter(s.Gateway))
		})
	if err != nil {
		log.Fatalf("Failed to start DHCP server: %v", err)
	}
	log.Println("DHCP Server is running...")
	if err := srv.Serve(); err != nil {
		log.Fatalf("Failed to serve DHCP requests: %v", err)
	}
}
