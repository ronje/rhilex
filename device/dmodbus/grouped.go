package dmodbus

import "sort"

type GroupedRegister struct {
	SlaverId byte   `json:"slaverId" validate:"required" title:"从机ID"`
	Function int    `json:"function" validate:"required" title:"Modbus功能"`
	Address  uint16 `json:"address" validate:"required" title:"地址"`
	Quantity uint16 `json:"quantity" validate:"required" title:"数量"`
}

// GroupModbusRegister groups modbus registers for optimal batch reading
func GroupModbusRegister(registers []GroupedRegister) [][]GroupedRegister {
	if len(registers) == 0 {
		return [][]GroupedRegister{}
	}

	// Sort registers by SlaverId, Function, and Address
	sort.Slice(registers, func(i, j int) bool {
		if registers[i].SlaverId != registers[j].SlaverId {
			return registers[i].SlaverId < registers[j].SlaverId
		}
		if registers[i].Function != registers[j].Function {
			return registers[i].Function < registers[j].Function
		}
		return registers[i].Address < registers[j].Address
	})

	var result [][]GroupedRegister
	var currentGroup []GroupedRegister

	for i, reg := range registers {
		if i == 0 {
			currentGroup = append(currentGroup, reg)
			continue
		}

		prev := currentGroup[len(currentGroup)-1]
		// Check continuity for optimal grouping and ensure no overlap
		if reg.SlaverId == prev.SlaverId && reg.Function == prev.Function &&
			reg.Address == prev.Address+prev.Quantity {
			currentGroup = append(currentGroup, reg)
		} else {
			result = append(result, currentGroup)
			currentGroup = []GroupedRegister{reg}
		}
	}

	if len(currentGroup) > 0 {
		result = append(result, currentGroup)
	}

	return result
}
