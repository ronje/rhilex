#! /bin/bash

# Function to check command execution
check_command() {
    if [ $? -ne 0 ]; then
        echo "Error: $1 failed"
        exit 1
    fi
}

# Set environment paths
setup_environment() {
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export GOBIN=$GOPATH/bin
    export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
}

# Main function to generate proto files
generate_proto() {
    echo -e "\033[42;33m>>>\033[0m [BEGIN] >>>"
# RhilexRpc
    echo ">>> Generating RhilexRpc Proto..."
    protoc -I ./component/rhilexrpc --go_out=./component/rhilexrpc --go_opt paths=source_relative \
        --go-grpc_out=./component/rhilexrpc --go-grpc_opt paths=source_relative \
        ./component/rhilexrpc/rhilexrpc.proto
    check_command "RhilexRpc Proto generation"
    echo ">>> Generate RhilexRpc Proto OK"

# Stream
    echo ">>> Generating XStream Proto..."
    protoc -I ./component/xstream --go_out ./component/xstream --go_opt paths=source_relative \
        --go-grpc_out=./component/xstream --go-grpc_opt paths=source_relative \
        ./component/xstream/xstream.proto
    check_command "XStream Proto generation"
    echo ">>> Generate XStream Proto OK."

# AI Base
    echo ">>> Generating Aibase Proto..."
    protoc -I ./component/aibase/grpc \
        --go_out ./component/aibase/grpc \
        --go_opt paths=source_relative \
        --go-grpc_out=./component/aibase/grpc \
        --go-grpc_opt paths=source_relative \
        ./component/aibase/grpc/aibase.proto
    check_command "Aibase Proto generation"
    echo ">>> Generate AIBase Proto OK."

# Activation
    echo ">>> Generating Activation Proto..."
    protoc -I ./component/activation \
        --go_out ./component/activation \
        --go_opt paths=source_relative \
        --go-grpc_out=./component/activation \
        --go-grpc_opt paths=source_relative \
        ./component/activation/activation.proto
    check_command "Activation Proto generation"
    echo ">>> Generate Activation Proto OK."

# Stream
    echo ">>> Generating Tunnel Proto... "
    protoc -I ./plugin/tunnel --go_out ./plugin/tunnel --go_opt paths=source_relative \
        --go-grpc_out=./plugin/tunnel --go-grpc_opt paths=source_relative \
        ./plugin/tunnel/tunnel.proto
    check_command "Tunnel Proto generation"
    echo ">>> Generate Tunnel Proto OK."

    echo -e "\033[42;33m>>>\033[0m [FINISHED] >>>"
}

# Main execution
setup_environment
generate_proto
