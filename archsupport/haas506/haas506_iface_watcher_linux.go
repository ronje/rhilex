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

package haas506

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hootrhino/rhilex/ossupport"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func InitNetworkIfaceWatcher() {
	NetworkIfaceWatcher := NewNetworkIfaceWatcher(logrus.New())
	NetworkIfaceWatcher.SetCallback(func(e IfaceEvent) {
		if e.Status == "up" {
			up, err := ossupport.IsInterfaceUp(e.Iface)
			if err != nil {
				log.Println("Get Interface State error:", err)
				return
			}
			if up {
				log.Println("Interface change state, try to DHCP:", e.String())
				if err1 := ossupport.AllocateIPAddressUsingUdhcpc(e.Iface); err1 != nil {
					log.Println("DHCP error:", e.Iface, err1)
				}
			}
		}
	})
	go NetworkIfaceWatcher.StartWatch()
	time.Sleep(1 * time.Second)
	InitDevTree()
}

type IfaceEvent struct {
	Iface        string
	Status       string //"up"|"down"
	HardwareAddr string
}

// String returns a formatted string representation of the NetlinkEvent.
func (e IfaceEvent) String() string {

	return fmt.Sprintf(`Netlink Event == Interface: %s Status: %s HardwareAddr: %s`,
		e.Iface,
		e.Status,
		e.HardwareAddr,
	)
}

type NetworkIfaceWatcher struct {
	logger   *logrus.Logger
	done     chan struct{}
	callback func(e IfaceEvent)
}

func NewNetworkIfaceWatcher(logger *logrus.Logger) *NetworkIfaceWatcher {
	return &NetworkIfaceWatcher{logger: logger, done: make(chan struct{})}
}
func (w *NetworkIfaceWatcher) SetCallback(callback func(e IfaceEvent)) {
	w.callback = callback
}
func (w *NetworkIfaceWatcher) StartWatch() error {
	ch := make(chan netlink.LinkUpdate)
	if err := netlink.LinkSubscribe(ch, w.done); err != nil {
		return err
	}
	for {
		select {
		case <-context.Background().Done():
			return nil
		case <-w.done:
			return nil
		case state := <-ch:
			attrs := state.Link.Attrs()
			if attrs != nil {
				event := IfaceEvent{
					Iface:        attrs.Name,
					Status:       attrs.OperState.String(),
					HardwareAddr: attrs.HardwareAddr.String(),
				}
				w.logger.Debug(event.String())
				if w.callback != nil {
					w.callback(event)
				}
			}
		default:
		}
	}
}
func (w *NetworkIfaceWatcher) StopWatch() {
	w.done <- struct{}{}
}
