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
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/transceivercom"
	atk01lora "github.com/hootrhino/rhilex/component/transceivercom/atk01-lora"
	ec200a4g "github.com/hootrhino/rhilex/component/transceivercom/ec200a-4g"
	mx01ble "github.com/hootrhino/rhilex/component/transceivercom/mx01-ble"
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

func (TM *TransceiverCommunicatorManager) Load(name string, config transceivercom.TransceiverConfig,
	tc transceivercom.TransceiverCommunicator) error {
	glogger.GLogger.Debugf("transceiver Load:(%s, %v, %s)", name, config, tc.Info().String())
	if _, ok := TM.Transceivers.Load(name); !ok {
		if err := tc.Start(config); err != nil {
			glogger.GLogger.Error(err)
			return err
		}
		TM.Transceivers.Store(name, tc)
		return nil
	}
	return fmt.Errorf("transceiver already loaded: %s", name)
}

func (TM *TransceiverCommunicatorManager) List() []transceivercom.CommunicatorInfo {
	List := []transceivercom.CommunicatorInfo{}
	TM.Transceivers.Range(func(key, value any) bool {
		switch T := value.(type) {
		case transceivercom.TransceiverCommunicator:
			List = append(List, T.Info())
		}
		return true
	})
	return List
}

func (TM *TransceiverCommunicatorManager) Get(name string) transceivercom.TransceiverCommunicator {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case transceivercom.TransceiverCommunicator:
			return T
		}
	}
	return nil
}
func (TM *TransceiverCommunicatorManager) UnLoad(name string) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case transceivercom.TransceiverCommunicator:
			T.Stop()
			TM.Transceivers.Delete(name)
		}
	}
}

func (TM *TransceiverCommunicatorManager) Ctrl(name string, topic, args []byte,
	timeout time.Duration) ([]byte, error) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case transceivercom.TransceiverCommunicator:
			return T.Ctrl(topic, args, timeout)
		}
	}
	return nil, fmt.Errorf("transceiver not exists: %s", name)
}

func (TM *TransceiverCommunicatorManager) Status(name string) (transceivercom.TransceiverStatus, error) {
	if value, ok := TM.Transceivers.Load(name); ok {
		switch T := value.(type) {
		case transceivercom.TransceiverCommunicator:
			return T.Status(), nil
		}
	}
	return transceivercom.TransceiverStatus{}, fmt.Errorf("transceiver not exists: %s", name)
}

/*
*
* Load Default Modules
*
 */
func initDefaultRFModule() {
	env1 := os.Getenv("BLESUPPORT")
	if env1 == "MX01" {
		Config := transceivercom.TransceiverConfig{}
		err1 := utils.INIToStruct(core.GlobalConfig.IniPath, "transceiver.mx01", &Config)
		if err1 != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
		Mx01 := mx01ble.NewMx01BLE(DefaultTransceiverCommunicatorManager.R)
		err := DefaultTransceiverCommunicatorManager.Load(Mx01.Info().Name, Config, Mx01)
		if err != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
	}
	env2 := os.Getenv("4GSUPPORT")
	if env2 == "EC200A" {
		Config := transceivercom.TransceiverConfig{}
		err1 := utils.INIToStruct(core.GlobalConfig.IniPath, "transceiver.ec200a", &Config)
		if err1 != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
		EC200A := ec200a4g.NewEC200ADtu(DefaultTransceiverCommunicatorManager.R)
		err := DefaultTransceiverCommunicatorManager.Load(EC200A.Info().Name, Config, EC200A)
		if err != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
	}
	env3 := os.Getenv("LORASUPPORT")
	if env3 == "ATK01" {
		Config := transceivercom.TransceiverConfig{}
		err1 := utils.INIToStruct(core.GlobalConfig.IniPath, "transceiver.atk01", &Config)
		if err1 != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
		ATK01 := atk01lora.NewATK01Lora(DefaultTransceiverCommunicatorManager.R)
		err := DefaultTransceiverCommunicatorManager.Load(ATK01.Info().Name, Config, ATK01)
		if err != nil {
			glogger.GLogger.Fatal(err1)
			os.Exit(1)
		}
	}
}
