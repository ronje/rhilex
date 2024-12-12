<!--
 Copyright (C) 2024 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# MyComponent Template README
## Overview
This README provides instructions on how to use the `MyComponent` template, which is a starting point for creating a component in a Go project. The template includes a basic structure for initializing, starting, stopping, and handling service calls within your component.
## Features
- Initialization and configuration of the component.
- Starting and stopping the component.
- Handling service calls with arguments and returning results.
- Providing metadata about the component.
- Simulated main function to demonstrate the component lifecycle.
## Getting Started
### Prerequisites
- Go programming language installed (version 1.16 or later recommended).
- A Go workspace set up on your machine.
### Installation
1. Clone the repository or copy the template files into your project directory.
2. Import the template into your Go module by replacing `yourmodule` with your actual module path in the import statements.
### Usage
#### Initializing the Component
To initialize the component, call the `Init` method with a configuration map:
```go
cfg := map[string]interface{}{
    // Your configuration options here.
}
err := component.Init(cfg)
if err != nil {
    log.Fatalf("Failed to initialize component: %v", err)
}
```
#### Starting the Component
To start the component, call the `Start` method:
```go
err := component.Start(rhilexInstance)
if err != nil {
    log.Fatalf("Failed to start component: %v", err)
}
```
Replace `rhilexInstance` with the actual instance of the `Rhilex` type that your component interacts with.
#### Handling Service Calls
To handle a service call, use the `Call` method:
```go
callArgs := CallArgs{
    ComponentName: "MyComponent",
    ServiceName:   "ExampleService",
}
callResult, err := component.Call(callArgs)
if err != nil {
    log.Fatalf("Failed to call service: %v", err)
}
```
#### Getting Meta Information
To retrieve the meta information of the component, use the `MetaInfo` method:
```go
metaInfo := component.MetaInfo()
fmt.Printf("Component MetaInfo: %+v\n", metaInfo)
```
#### Stopping the Component
To stop the component, call the `Stop` method:
```go
err := component.Stop()
if err != nil {
    log.Fatalf("Failed to stop component: %v", err)
}
```
## Example
The `main.go` file provides an example of how to use the `MyComponent` template. You can run the example directly to see the component lifecycle in action:
```bash
go run main.go
```
## Contributing
Contributions to the `MyComponent` template are welcome. Please follow the standard Go contribution guidelines and ensure that your code adheres to the project's coding standards.
## License
The `MyComponent` template is licensed under [LICENSE FILE] - please check the LICENSE file for details.
## Contact
For any questions or suggestions, please open an issue on the project's repository or reach out to [YOUR CONTACT INFORMATION].

