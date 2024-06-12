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

package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"github.com/vapourismo/knx-go/knx/util"
)

// go test -timeout 30s -run ^Test_knx_group_tunnel github.com/hootrhino/rhilex/test -v -count=1
func Test_knx_group_tunnel(t *testing.T) {
	// Setup logger for auxiliary logging. This enables us to see log messages from internal
	// routines.
	util.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Connect to the gateway.
	client, err := knx.NewGroupTunnel("127.0.0.1:3671", knx.TunnelConfig{
		ResendInterval:    500 * time.Millisecond,
		HeartbeatInterval: 10 * time.Second,
		ResponseTimeout:   10 * time.Second,
		SendLocalAddress:  true,
		UseTCP:            false,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Close upon exiting. Even if the gateway closes the connection, we still have to clean up.
	defer client.Close()

	// Send 20.5°C to group 1/2/3.
	err = client.Send(knx.GroupEvent{
		Command:     knx.GroupWrite,
		Destination: cemi.NewGroupAddr3(1, 2, 3),
		Data:        dpt.DPT_9001(20.5).Pack(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages from the gateway. The inbound channel is closed with the connection.
	for msg := range client.Inbound() {
		var temp dpt.DPT_9001

		err := temp.Unpack(msg.Data)
		if err != nil {
			continue
		}

		util.Logger.Printf("%+v: %v", msg, temp)
	}
}
