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

package crontask

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CronExpr struct {
	Second     string
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

func (expr CronExpr) String() string {
	bytes, _ := json.Marshal(expr)
	return string(bytes)
}

// ParseCronExpr parses a cron expression and returns a CronExpr struct.
func ParseCronExpr(expr string) (CronExpr, error) {
	// Split the cron expression by spaces
	parts := strings.Fields(expr)
	// A cron expression can have 5 or 6 parts (if seconds are included)
	if len(parts) < 5 || len(parts) > 6 {
		return CronExpr{}, errors.New("invalid cron expression: expected 5 or 6 fields")
	}
	// If only 5 parts, prepend "0" to assume that seconds are set to 0
	if len(parts) == 5 {
		parts = append([]string{"0"}, parts...)
	}

	// Validate each part
	if err := validateCronPart(parts[0], 0, 59, "second"); err != nil {
		return CronExpr{}, err
	}
	if err := validateCronPart(parts[1], 0, 59, "minute"); err != nil {
		return CronExpr{}, err
	}
	if err := validateCronPart(parts[2], 0, 23, "hour"); err != nil {
		return CronExpr{}, err
	}
	if err := validateCronPart(parts[3], 1, 31, "day of month"); err != nil {
		return CronExpr{}, err
	}
	if err := validateCronPart(parts[4], 1, 12, "month"); err != nil {
		return CronExpr{}, err
	}
	if err := validateCronPart(parts[5], 0, 6, "day of week"); err != nil {
		return CronExpr{}, err
	}

	// Return the parsed CronExpr
	return CronExpr{
		Second:     parts[0],
		Minute:     parts[1],
		Hour:       parts[2],
		DayOfMonth: parts[3],
		Month:      parts[4],
		DayOfWeek:  parts[5],
	}, nil
}

// validateCronPart validates that the cron expression part is within the allowed range.
func validateCronPart(part string, min, max int, name string) error {
	if part == "*" {
		return nil // '*' means any value, so it's always valid
	}

	// Check for step expression (e.g. */5)
	if strings.HasPrefix(part, "*/") {
		step, err := strconv.Atoi(part[2:])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid %s: %s", name, part)
		}
		return nil
	}

	// Check for list or range
	for _, field := range strings.Split(part, ",") {
		if strings.Contains(field, "-") {
			rangeParts := strings.Split(field, "-")
			if len(rangeParts) != 2 {
				return fmt.Errorf("invalid range in %s: %s", name, field)
			}
			start, err1 := strconv.Atoi(rangeParts[0])
			end, err2 := strconv.Atoi(rangeParts[1])
			if err1 != nil || err2 != nil || start < min || end > max || start > end {
				return fmt.Errorf("invalid range in %s: %s", name, field)
			}
		} else {
			// Single value
			val, err := strconv.Atoi(field)
			if err != nil || val < min || val > max {
				return fmt.Errorf("invalid %s: %s", name, part)
			}
		}
	}
	return nil
}
