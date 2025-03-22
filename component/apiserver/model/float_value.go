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
package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
)

type Decimal float64

// NewDecimal creates a new Decimal from a float64 value and returns a pointer to it.
func NewDecimal(f float64) *Decimal {
	d := Decimal(f)
	return &d
}

func (d Decimal) Value() (driver.Value, error) {
	return fmt.Sprintf("%.3f", d), nil
}

// Scan implements the sql.Scanner interface for the Decimal type.
// It can now handle string, *float64, and *float32 inputs.
func (d *Decimal) Scan(input any) error {
	switch v := input.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*d = Decimal(f)
	case float64:
		*d = Decimal(v)
	case *float64:
		if v == nil {
			return errors.New("nil pointer passed to Decimal.Scan for *float64")
		}
		*d = Decimal(*v)
	case float32:
		*d = Decimal(v)
	case *float32:
		if v == nil {
			return errors.New("nil pointer passed to Decimal.Scan for *float32")
		}
		*d = Decimal(float64(*v))
	default:
		return errors.New("invalid input type for Decimal")
	}
	return nil
}

// ToFloat64 converts Decimal to *float64.
func (d Decimal) ToFloat64() *float64 {
	f := float64(d)
	return &f
}
