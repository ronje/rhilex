package apis

import "github.com/hootrhino/rhilex/component/apiserver/server"

// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

func LoadSystemSettingsAPI() {
	//
	// 系统设置
	//

	iFacesApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		iFacesApi.GET(("/ctrlTree"), server.AddRoute(GetDeviceCtrlTree))
		iFacesApi.GET(("/netStatus"), server.AddRoute(GetNetworkStatus))
	}
	// ethernet
	ethApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		ethApi.POST("/eth", server.AddRoute(SetEthNetwork))
		ethApi.GET("/eth", server.AddRoute(GetEthNetwork))
		ethApi.GET("/eths", server.AddRoute(AllEthNetwork))
	}
	// wifi
	wifiApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		wifiApi.GET("/wifi", server.AddRoute(GetWifi))
		wifiApi.POST("/wifi", server.AddRoute(SetWifi))
		wifiApi.GET("/wifis", server.AddRoute(AllWlanNetwork))
		wifiApi.GET("/wifi/scanSignal", server.AddRoute(ScanWIFIList))
	}
	// time
	timesApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		// time
		timesApi.GET("/time", server.AddRoute(GetSystemTime))
		timesApi.POST("/time", server.AddRoute(SetSystemTime))
		timesApi.PUT("/ntp", server.AddRoute(UpdateTimeByNtp))
	}
	// 4g module
	settings4GApi := server.RouteGroup(server.ContextUrl("/mn4g"))
	{
		settings4GApi.POST(("/turnon"), server.AddRoute(Turnon4g))
		settings4GApi.POST(("/turnoff"), server.AddRoute(Turnoff4g))
		settings4GApi.GET(("/info"), server.AddRoute(Get4GBaseInfo))
		settings4GApi.GET("/apn", server.AddRoute(GetAPN))
		settings4GApi.POST("/apn", server.AddRoute(SetAPN))
	}
	// 固件
	settingsFirmware := server.RouteGroup(server.ContextUrl("/firmware"))
	{
		settingsFirmware.POST("/reboot", server.AddRoute(Reboot))
		settingsFirmware.POST("/recoverNew", server.AddRoute(RecoverNew))
		settingsFirmware.POST("/restartRhilex", server.AddRoute(ReStartRhilex))
		settingsFirmware.POST("/upload", server.AddRoute(UploadFirmWare))
		settingsFirmware.POST("/upgrade", server.AddRoute(UpgradeFirmWare))
		settingsFirmware.GET("/upgradeLog", server.AddRoute(GetUpGradeLog))
		settingsFirmware.GET("/vendorKey", server.AddRoute(GetVendorKey))
	}

}
