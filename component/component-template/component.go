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

package componenttemplate

import (
	"fmt"

	"github.com/hootrhino/rhilex/typex"
)

// XComponentMetaInfo provides metadata about the component.
type XComponentMetaInfo struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// CallArgs represents the arguments for a service call.
type CallArgs struct {
	ComponentName string
	ServiceName   string
}

// CallResult represents the result of a service call.
type CallResult struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}

// ServiceSpec defines a service provided by the component.
type ServiceSpec struct {
	CallArgs   CallArgs
	CallResult CallResult
}

// MyComponent is a sample implementation of XComponent.
type MyComponent struct {
	metaInfo XComponentMetaInfo
}

// Init initializes the component with the given configuration.
func (c *MyComponent) Init(cfg map[string]interface{}) error {
	// Initialize your component with the provided configuration.
	fmt.Println("Initializing component with config:", cfg)
	return nil
}

// Start starts the component.
func (c *MyComponent) Start(rhilex typex.Rhilex) error {
	// Start your component, e.g., begin listening for incoming requests.
	fmt.Println("Starting component...")
	return nil
}

// Call handles the incoming call and returns the result.
func (c *MyComponent) Call(args CallArgs) (CallResult, error) {
	// Implement the logic to handle the call based on the provided arguments.
	fmt.Printf("Handling call for service '%s'...\n", args.ServiceName)
	// For demonstration purposes, we just return a dummy result.
	return CallResult{Code: 200, Result: "Success"}, nil
}

// Services returns the list of services provided by the component.
func (c *MyComponent) Services() map[string]ServiceSpec {
	// Define the services provided by this component.
	return map[string]ServiceSpec{
		"ExampleService": {
			CallArgs: CallArgs{
				ComponentName: "MyComponent",
				ServiceName:   "ExampleService",
			},
			CallResult: CallResult{
				Code:   200,
				Result: "Success",
			},
		},
	}
}

// MetaInfo returns the meta information of the component.
func (c *MyComponent) MetaInfo() XComponentMetaInfo {
	return c.metaInfo
}

// Stop stops the component.
func (c *MyComponent) Stop() error {
	// Implement the logic to stop the component gracefully.
	fmt.Println("Stopping component...")
	return nil
}
