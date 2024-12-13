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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package supervisor

import (
	"context"
	"fmt"
	"sync"

	"github.com/hootrhino/rhilex/typex"
)

var __DefaultSuperVisorAdmin *SuperVisorAdmin

// SuperVisorAdmin 管理所有的 Supervisor 实例
type SuperVisorAdmin struct {
	Ctx         context.Context
	Locker      sync.Locker
	SuperVisors map[string]*SuperVisor
	rhilex      typex.Rhilex
}

// SuperVisor 表示一个 Supervisor 实例
type SuperVisor struct {
	SlaverId string
	Ctx      context.Context
	Cancel   context.CancelFunc
}

// NewSuperVisorAdmin 创建一个新的 SuperVisorAdmin 实例
func NewSuperVisorAdmin(ctx context.Context, rhilex typex.Rhilex) *SuperVisorAdmin {
	return &SuperVisorAdmin{
		Ctx:         ctx,
		Locker:      &sync.Mutex{},
		SuperVisors: make(map[string]*SuperVisor),
		rhilex:      rhilex,
	}
}

// RegisterSuperVisor 注册一个新的 Supervisor 实例
func (s *SuperVisorAdmin) RegisterSuperVisor(slaverId string) (*SuperVisor, error) {
	s.Locker.Lock()
	defer s.Locker.Unlock()

	if old, ok := s.SuperVisors[slaverId]; ok {
		old.Cancel()
		delete(s.SuperVisors, slaverId)
	}

	ctx, cancel := context.WithCancel(s.Ctx)
	supervisor := &SuperVisor{SlaverId: slaverId, Ctx: ctx, Cancel: cancel}
	s.SuperVisors[slaverId] = supervisor

	return supervisor, nil
}

// UnRegisterSuperVisor 注销一个 Supervisor 实例
func (s *SuperVisorAdmin) UnRegisterSuperVisor(UUID string) error {
	s.Locker.Lock()
	defer s.Locker.Unlock()

	if sv, ok := s.SuperVisors[UUID]; ok {
		sv.Cancel()
		delete(s.SuperVisors, UUID)
		return nil
	}

	return fmt.Errorf("supervisor with UUID %s not found", UUID)
}

// StopSuperVisor 停止一个 Supervisor 实例
func (s *SuperVisorAdmin) StopSuperVisor(UUID string) error {
	s.Locker.Lock()
	defer s.Locker.Unlock()

	if sv, ok := s.SuperVisors[UUID]; ok {
		sv.Cancel()
		return nil
	}

	return fmt.Errorf("supervisor with UUID %s not found", UUID)
}
func (s *SuperVisorAdmin) StopSupervisorAdmin() {
	s.Locker.Lock()
	defer s.Locker.Unlock()
	for _, sv := range s.SuperVisors {
		sv.Cancel()
	}
}

/*
*
* 初始化超级Admin
*
 */
func InitResourceSuperVisorAdmin(rhilex typex.Rhilex) {
	__DefaultSuperVisorAdmin = &SuperVisorAdmin{
		Ctx:         context.Background(),
		Locker:      &sync.Mutex{},
		SuperVisors: map[string]*SuperVisor{},
		rhilex:      rhilex,
	}
}

/*
*
* 启动Supervisor的时候注册
*
 */
func RegisterSuperVisor(SlaverId string) *SuperVisor {
	SuperVisor, err := __DefaultSuperVisorAdmin.RegisterSuperVisor(SlaverId)
	if err != nil {
		return nil
	}
	return SuperVisor
}

/*
*
* Supervisor进程退出的时候执行
*
 */
func UnRegisterSuperVisor(UUID string) {
	__DefaultSuperVisorAdmin.UnRegisterSuperVisor(UUID)
}

/*
*
* 停止一个Supervisor
*
 */
func StopSuperVisor(UUID string) {
	__DefaultSuperVisorAdmin.StopSuperVisor(UUID)

}
func StopSupervisorAdmin() {
	__DefaultSuperVisorAdmin.StopSupervisorAdmin()
}
