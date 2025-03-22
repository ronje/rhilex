package dmodbus

import "testing"

func Test_GroupModbusRegister(t *testing.T) {
	registers := []GroupedRegister{
		// 连续点位，1个字节
		{3, 1, 1, 1},
		{3, 1, 2, 1},
		{3, 1, 3, 1},

		// 离散点位
		{3, 2, 5, 1},
		{3, 2, 8, 1},

		// 交叉地址点位
		{3, 3, 10, 2},
		{3, 3, 11, 2},

		// 连续点位，4个字节
		{3, 4, 20, 4},
		{3, 4, 24, 4},

		// 2个字节点位
		{3, 5, 30, 2},
		{3, 5, 32, 2},

		// 不同从站
		{3, 6, 40, 1},
		{3, 7, 50, 1},

		// 乱序点位
		{1, 1, 15, 1},
		{2, 1, 14, 1},
		{3, 1, 13, 1},
		{1, 1, 15, 1},
		{2, 1, 14, 1},
		{3, 1, 13, 1},
		{1, 1, 15, 1},
		{2, 1, 14, 1},
		{3, 1, 13, 1},
	}

	grouped := GroupModbusRegister(registers)
	for i, group := range grouped {
		t.Logf("[ Group %d ]\n", i+1)
		for _, reg := range group {
			t.Logf("  SlaverId: %d, Function: %d, Address: %d, Quantity: %d\n",
				reg.SlaverId, reg.Function, reg.Address, reg.Quantity)
		}
	}
}
