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

package transceivercom

import (
	"time"

	"github.com/hootrhino/rhilex/component/transceivercom"
	"github.com/hootrhino/rhilex/glogger"
)

// Ctrl
func Ctrl(name string, cmd []byte, timeout time.Duration) ([]byte, error) {
	glogger.GLogger.Debugf("transceiver Ctrl:(%s, %s, %d)", name, string(cmd), timeout)
	return DefaultTransceiverCommunicatorManager.Ctrl(name, cmd, timeout)
}

// Unload
func Unload(name string) {
	glogger.GLogger.Infof("transceiver Unload:(%s)", name)
	DefaultTransceiverCommunicatorManager.UnLoad(name)
}

// List
func List() []transceivercom.CommunicatorInfo {
	return DefaultTransceiverCommunicatorManager.List()
}

// List
func GetCommunicator(name string) transceivercom.TransceiverCommunicator {
	return DefaultTransceiverCommunicatorManager.Get(name)
}

// Stop
func Stop() {
	for _, TC := range DefaultTransceiverCommunicatorManager.List() {
		Unload(TC.Name)
	}
}
