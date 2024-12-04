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

package tinyrouter

import (
	"fmt"
	"log"
	"net"

	"github.com/vishvananda/netlink"
)

// RouterManager 提供路由表的管理功能
type RouterManager struct{}

// AddRoute 添加一个新的静态路由
func (rm *RouterManager) AddRoute(destination, gateway, iface string) error {
	destIP, destNet, err := net.ParseCIDR(destination)
	if err != nil {
		return fmt.Errorf("invalid destination: %w", err)
	}
	gwIP := net.ParseIP(gateway)
	if gwIP == nil {
		return fmt.Errorf("invalid gateway IP")
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("failed to find interface %s: %w", iface, err)
	}

	route := &netlink.Route{
		Dst:       &net.IPNet{IP: destIP, Mask: destNet.Mask},
		Gw:        gwIP,
		LinkIndex: link.Attrs().Index,
	}

	if err := netlink.RouteAdd(route); err != nil {
		return fmt.Errorf("failed to add route: %w", err)
	}
	log.Printf("Added route to %s via %s on interface %s", destination, gateway, iface)
	return nil
}

// DeleteRoute 删除指定的静态路由
func (rm *RouterManager) DeleteRoute(destination, gateway string) error {
	destIP, destNet, err := net.ParseCIDR(destination)
	if err != nil {
		return fmt.Errorf("invalid destination: %w", err)
	}
	gwIP := net.ParseIP(gateway)
	if gwIP == nil {
		return fmt.Errorf("invalid gateway IP")
	}

	route := &netlink.Route{
		Dst: &net.IPNet{IP: destIP, Mask: destNet.Mask},
		Gw:  gwIP,
	}

	if err := netlink.RouteDel(route); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	log.Printf("Deleted route to %s via %s", destination, gateway)
	return nil
}

// ListRoutes 列出当前系统的所有路由条目
func (rm *RouterManager) ListRoutes() ([]netlink.Route, error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	for _, route := range routes {
		log.Printf("Route: %+v", route)
	}
	return routes, nil
}

// FindRoute 查找指定目标的路由条目
func (rm *RouterManager) FindRoute(destination string) (*netlink.Route, error) {
	destIP, destNet, err := net.ParseCIDR(destination)
	if err != nil {
		return nil, fmt.Errorf("invalid destination: %w", err)
	}

	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	for _, route := range routes {
		if route.Dst != nil && route.Dst.IP.Equal(destIP) && route.Dst.Mask.String() == destNet.Mask.String() {
			log.Printf("Found route to %s: %+v", destination, route)
			return &route, nil
		}
	}

	return nil, fmt.Errorf("route to %s not found", destination)
}

// UpdateRoute 修改指定的路由条目
func (rm *RouterManager) UpdateRoute(destination, gateway, iface string) error {
	// 删除旧路由条目
	route, err := rm.FindRoute(destination)
	if err != nil {
		return fmt.Errorf("failed to find route: %w", err)
	}

	if err := netlink.RouteDel(route); err != nil {
		return fmt.Errorf("failed to delete old route: %w", err)
	}

	// 添加新的路由条目
	if err := rm.AddRoute(destination, gateway, iface); err != nil {
		return fmt.Errorf("failed to add new route: %w", err)
	}
	log.Printf("Updated route to %s via %s on interface %s", destination, gateway, iface)
	return nil
}

func TestR1() {
	rm := RouterManager{}

	// 示例：添加一个路由条目
	if err := rm.AddRoute("192.168.1.0/24", "192.168.1.1", "eth1"); err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}

	// 示例：列出所有路由
	routes, err := rm.ListRoutes()
	if err != nil {
		log.Fatalf("Failed to list routes: %v", err)
	}
	fmt.Println("All routes:", routes)

	// 示例：查找特定路由
	if route, err := rm.FindRoute("192.168.1.0/24"); err == nil {
		fmt.Println("Found route:", route)
	} else {
		log.Println("Route not found:", err)
	}

	// 示例：删除路由条目
	if err := rm.DeleteRoute("192.168.1.0/24", "192.168.1.1"); err != nil {
		log.Fatalf("Failed to delete route: %v", err)
	}
}
