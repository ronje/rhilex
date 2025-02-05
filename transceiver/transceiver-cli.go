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

package transceiver

import (
	"time"

	"github.com/hootrhino/rhilex/glogger"
)

// Ctrl
func Ctrl(name string, topic, args []byte, timeout time.Duration) ([]byte, error) {
	glogger.GLogger.Debugf("transceiver Ctrl:(%s, %s, %s, %d)",
		name, string(topic), string(args), timeout)
	return __DefaultTransceiverManager.Ctrl(name, topic, args, timeout)
}

// Unload
func Unload(name string) {
	glogger.GLogger.Infof("transceiver Unload:(%s)", name)
	__DefaultTransceiverManager.UnLoad(name)
}

// List
func List() []CommunicatorInfo {
	return __DefaultTransceiverManager.List()
}

// List
func GetCommunicator(name string) Transceiver {
	return __DefaultTransceiverManager.Get(name)
}

// Stop
func Stop() {
	for _, TC := range __DefaultTransceiverManager.List() {
		Unload(TC.Name)
	}
}
