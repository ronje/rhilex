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
	"strings"
	"sync"
	"time"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/utils"
	"gopkg.in/ini.v1"

	"github.com/hootrhino/rhilex/glogger"

	"github.com/hootrhino/rhilex/typex"
)

var __DefaultTransceiverManager *TransceiverManager

type TransceiverManager struct {
	Transceivers sync.Map
	R            typex.Rhilex
}

func InitTransceiverManager(R typex.Rhilex) {
	__DefaultTransceiverManager = &TransceiverManager{
		R:            R,
		Transceivers: sync.Map{},
	}
	LoadTransceiverModules()
}

func (TM *TransceiverManager) Load(name string, config TransceiverConfig,
	tc Transceiver) error {
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

func (TM *TransceiverManager) List() []CommunicatorInfo {
	List := []CommunicatorInfo{}
	TM.Transceivers.Range(func(key, value any) bool {
		switch T := value.(type) {
		case Transceiver:
			List = append(List, T.Info())
		}
		return true
	})
	return List
}

func (TM *TransceiverManager) Get(name string) Transceiver {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case Transceiver:
			return T
		}
	}
	return nil
}
func (TM *TransceiverManager) UnLoad(name string) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case Transceiver:
			T.Stop()
			TM.Transceivers.Delete(name)
		}
	}
}

func (TM *TransceiverManager) Ctrl(name string, topic, args []byte,
	timeout time.Duration) ([]byte, error) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case Transceiver:
			return T.Ctrl(topic, args, timeout)
		}
	}
	return nil, fmt.Errorf("Transceiver not exists: %s", name)
}

func (TM *TransceiverManager) Status(name string) (TransceiverStatus, error) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case Transceiver:
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
func LoadTransceiverModules() {
	iniConfigFile, _ := ini.ShadowLoad(core.GlobalConfig.IniPath)
	sections := iniConfigFile.ChildSections("transceiver")
	for _, section := range sections {
		name := strings.TrimPrefix(section.Name(), "transceiver.")
		enable, errGetKey := section.GetKey("enable")
		if errGetKey != nil {
			continue
		}
		if !enable.MustBool(false) {
			glogger.GLogger.Warnf("transceiver is disable:%s", name)
			continue
		}
		config := TransceiverConfig{Name: name}
		errMap := utils.InIMapToStruct(section, &config)
		if errMap != nil {
			glogger.GLogger.Fatal(errMap)
			os.Exit(1)
		}
		Transceiver := NewTransceiver(__DefaultTransceiverManager.R)
		errLoad := __DefaultTransceiverManager.Load(config.Name, config, Transceiver)
		if errLoad != nil {
			glogger.GLogger.Fatal(errLoad)
			os.Exit(1)
		}
	}
}
