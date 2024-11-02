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

package activation

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/**
 * 申请证书
 *
 */
func GetLicense(Host, Sn, Iface, Mac, U, P string) (string, string, string, error) {
	conn, err := grpc.NewClient(Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", "", "", fmt.Errorf("Dail to Server %s Error", Host)
	}
	defer conn.Close()
	client := NewDeviceActivationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req := &ActivationRequest{
		Sn:       Sn,
		Iface:    Iface,
		Mac:      Mac,
		Username: U,
		Password: P,
	}
	resp, err1 := client.ActivateDevice(ctx, req)
	if err1 != nil {
		return "", "", "", fmt.Errorf("Activate Device error, check your network status")
	}
	if !resp.Success {
		return "", "", "", fmt.Errorf("Activate Device failed, maybe server is panic")
	}
	return resp.Certificate, resp.Privatekey, resp.License, nil
}
