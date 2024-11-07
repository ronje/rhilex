package server

import (
	"context"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/internotify"
	"github.com/hootrhino/rhilex/component/supervisor"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 南向资源监控器 5秒检查一下状态
*
 */
func StartInSupervisor(InCtx context.Context, in *typex.InEnd, ruleEngine typex.Rhilex) {
	UUID := in.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(in.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For InEnd:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)
	for {
		select {
		case <-typex.GCTX.Done():
			{
				glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-InCtx.Done():
			{
				glogger.GLogger.Debugf("Source Context cancel:%v, supervisor exit", UUID)
				return
			}
		default:
			{
			}
		}
		// 被删除后就直接退出监督进程
		currentIn := ruleEngine.GetInEnd(UUID)
		if currentIn == nil {
			glogger.GLogger.Debugf("Source:%v Deleted, supervisor exit", UUID)
			return
		}
		// STOP 设计特殊状态,标记被彻底删除的资源
		// 资源可能不会及时DOWN
		if currentIn.Source.Status() == typex.SOURCE_DOWN {
			ErrMsg := ""
			Slot := intercache.GetSlot("__DefaultRuleEngine")
			if Slot != nil {
				CacheValue, ok := Slot[currentIn.UUID]
				if ok {
					ErrMsg = CacheValue.ErrMsg
				}
			}
			info := fmt.Sprintf("Source:(%s,%s) DOWN, supervisor try to Restart, error message: %s",
				UUID, currentIn.Name, ErrMsg)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  `WARNING`,
				Event: "event.inend.down",
				Ts:    uint64(time.Now().UnixMilli()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestInEnd(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

/*
*
* 北向资源监控器 5秒检查一下状态
*
 */
func StartOutSupervisor(OutCtx context.Context, out *typex.OutEnd, ruleEngine typex.Rhilex) {
	UUID := out.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(out.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For OutEnd:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)

	for {
		select {
		case <-typex.GCTX.Done():
			glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
			return
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-OutCtx.Done():
			glogger.GLogger.Debugf("OutEnd Context cancel:%v, supervisor exit", UUID)
			return
		default:
			{
			}
		}
		// 被删除后就直接退出监督进程
		currentOut := ruleEngine.GetOutEnd(UUID)
		if currentOut == nil {
			glogger.GLogger.Debugf("OutEnd:%v Deleted, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentOut.Target.Status() == typex.SOURCE_DOWN {
			ErrMsg := ""
			Slot := intercache.GetSlot("__DefaultRuleEngine")
			if Slot != nil {
				CacheValue, ok := Slot[currentOut.UUID]
				if ok {
					ErrMsg = CacheValue.ErrMsg
				}
			}
			info := fmt.Sprintf("OutEnd:(%s,%s) DOWN, supervisor try to Restart, error message: %s",
				UUID, currentOut.Name, ErrMsg)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  `WARNING`,
				Event: "event.outend.down",
				Ts:    uint64(time.Now().UnixMilli()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestOutEnd(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

/*
*
* 设备监控器 5秒检查一下状态
*
 */
func StartDeviceSupervisor(DeviceCtx context.Context, device *typex.Device, ruleEngine typex.Rhilex) {
	UUID := device.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(device.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For Device:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)

	for {
		select {
		case <-typex.GCTX.Done():
			{
				glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-DeviceCtx.Done():
			{
				glogger.GLogger.Debugf("Device Context cancel:%v, supervisor exit", UUID)
				return
			}
		default:
		}
		// 被删除后就直接退出监督进程
		currentDevice := ruleEngine.GetDevice(UUID)
		if currentDevice == nil {
			glogger.GLogger.Debugf("Device:%v Deleted, supervisor exit", UUID)
			return
		}

		// 资源可能不会及时DOWN
		currentDeviceStatus := currentDevice.Device.Status()
		if currentDeviceStatus == typex.DEV_DOWN {
			ErrMsg := ""
			Slot := intercache.GetSlot("__DefaultRuleEngine")
			if Slot != nil {
				CacheValue, ok := Slot[currentDevice.UUID]
				if ok {
					ErrMsg = CacheValue.ErrMsg
				}
			}
			info := fmt.Sprintf("Device:(%s,%s) DOWN, supervisor try to Restart, error message: %s",
				UUID, currentDevice.Name, ErrMsg)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  `WARNING`,
				Event: "event.device.down",
				Ts:    uint64(time.Now().UnixMilli()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestDevice(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

/*
*
* 设备监控器 5秒检查一下状态
*
 */
func StartCecollaSupervisor(CecollaCtx context.Context, cecolla *typex.Cecolla, ruleEngine typex.Rhilex) {
	UUID := cecolla.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(cecolla.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For Cecolla:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)

	for {
		select {
		case <-typex.GCTX.Done():
			{
				glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-CecollaCtx.Done():
			{
				glogger.GLogger.Debugf("Cecolla Context cancel:%v, supervisor exit", UUID)
				return
			}
		default:
		}
		// 被删除后就直接退出监督进程
		currentCecolla := ruleEngine.GetCecolla(UUID)
		if currentCecolla == nil {
			glogger.GLogger.Debugf("Cecolla:%v Deleted, supervisor exit", UUID)
			return
		}

		// 资源可能不会及时DOWN
		currentCecollaStatus := currentCecolla.Cecolla.Status()
		if currentCecollaStatus == typex.CEC_DOWN {
			ErrMsg := ""
			Slot := intercache.GetSlot("__DefaultRuleEngine")
			if Slot != nil {
				CacheValue, ok := Slot[currentCecolla.UUID]
				if ok {
					ErrMsg = CacheValue.ErrMsg
				}
			}
			info := fmt.Sprintf("Cecolla:(%s,%s) DOWN, supervisor try to Restart, error message: %s",
				UUID, currentCecolla.Name, ErrMsg)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  `WARNING`,
				Event: "event.cecolla.down",
				Ts:    uint64(time.Now().UnixMilli()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestCecolla(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}
