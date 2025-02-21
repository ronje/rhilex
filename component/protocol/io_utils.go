// Copyright (C) 2025 wwhai
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

package protocol

import (
	"context"
	"io"
	"time"
)

/*
*
* 在timeout内读取N个字节
*
 */
func ReadInWill(ctx context.Context, rwc io.ReadWriteCloser, timeout time.Duration) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var responseData []byte
	totalRead := 0
	tempBuffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return totalRead, responseData, ctx.Err()
			}
			return totalRead, responseData, ctx.Err()
		default:
			n, err := rwc.Read(tempBuffer)
			if err != nil {
				if err == io.EOF {
					return totalRead, responseData, nil
				}
				if ctx.Err() == context.DeadlineExceeded {
					return totalRead, responseData, ctx.Err()
				}
				return totalRead, responseData, err
			}
			if n > 0 {
				responseData = append(responseData, tempBuffer[:n]...)
				totalRead += n
			}
		}
	}
}
