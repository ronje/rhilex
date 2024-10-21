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
	"fmt"
	"os"
	"sync"
	"time"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/utils"

	"github.com/hootrhino/rhilex/glogger"

	"github.com/hootrhino/rhilex/typex"
)

var DefaultTransceiverCommunicatorManager *TransceiverCommunicatorManager

type TransceiverCommunicatorManager struct {
	Transceivers sync.Map
	R            typex.Rhilex
}

func InitTransceiverCommunicatorManager(R typex.Rhilex) {
	DefaultTransceiverCommunicatorManager = &TransceiverCommunicatorManager{
		R:            R,
		Transceivers: sync.Map{},
	}
	initDefaultRFModule()
}

func (TM *TransceiverCommunicatorManager) Load(name string, config TransceiverConfig,
	tc TransceiverCommunicator) error {
	glogger.GLogger.Debugf("Transceiver Communicator Load:(%s, %v, %s)",
		name, config, tc.Info().String())
	if _, ok := TM.Transceivers.Load(name); !ok {
		if err := tc.Start(config); err != nil {
			glogger.GLogger.Error(err)
			return err
		}
		TM.Transceivers.Store(name, tc)
		return nil
	}
	return fmt.Errorf("Transceiver already loaded: %s", name)
}

func (TM *TransceiverCommunicatorManager) List() []CommunicatorInfo {
	List := []CommunicatorInfo{}
	TM.Transceivers.Range(func(key, value any) bool {
		switch T := value.(type) {
		case TransceiverCommunicator:
			List = append(List, T.Info())
		}
		return true
	})
	return List
}

func (TM *TransceiverCommunicatorManager) Get(name string) TransceiverCommunicator {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case TransceiverCommunicator:
			return T
		}
	}
	return nil
}
func (TM *TransceiverCommunicatorManager) UnLoad(name string) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case TransceiverCommunicator:
			T.Stop()
			TM.Transceivers.Delete(name)
		}
	}
}

func (TM *TransceiverCommunicatorManager) Ctrl(name string, topic, args []byte,
	timeout time.Duration) ([]byte, error) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case TransceiverCommunicator:
			return T.Ctrl(topic, args, timeout)
		}
	}
	return nil, fmt.Errorf("Transceiver not exists: %s", name)
}

func (TM *TransceiverCommunicatorManager) Status(name string) (TransceiverStatus, error) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case TransceiverCommunicator:
			return T.Status(), nil
		}
	}
	return TransceiverStatus{}, fmt.Errorf("Transceiver not exists: %s", name)
}

/*
*
* Load Default Modules
*
 */
func initDefaultRFModule() {
	env := os.Getenv("TRANSCEIVER")
	if env == "default_transceiver" {
		Config := TransceiverConfig{}
		err1 := utils.INIToStruct(core.GlobalConfig.IniPath, fmt.Sprintf("transceiver.%s", env), &Config)
		if err1 != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
		Config.Name = env
		Transceiver := NewTransceiver(DefaultTransceiverCommunicatorManager.R)
		err2 := DefaultTransceiverCommunicatorManager.Load(env, Config, Transceiver)
		if err2 != nil {
			glogger.GLogger.Fatal(err2)
			os.Exit(1)
		}
	}
}
