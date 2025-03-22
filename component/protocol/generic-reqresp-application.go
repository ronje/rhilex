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

package protocol

type GenericAppLayer struct {
	errTxCount int32 // 错误包计数器
	errRxCount int32 // 错误包计数器
	transport  *Transport
}

func NewGenericAppLayerAppLayer(config ExchangeConfig) *GenericAppLayer {
	return &GenericAppLayer{
		errTxCount: 0,
		errRxCount: 0,
		transport:  NewTransport(config),
	}
}

// Request 发送请求并获取响应
// 先调用 Write 方法发送请求帧，若发送成功则调用 Read 方法读取响应帧
func (app *GenericAppLayer) Request(appFrame *ApplicationFrame) (*ApplicationFrame, error) {
	// 发送请求帧
	if err := app.Write(appFrame); err != nil {
		return nil, err
	}
	// 读取响应帧
	return app.Read()
}

// Write 编码并发送应用帧
// 先对应用帧进行编码，编码成功则通过 transport 发送编码后的数据
func (app *GenericAppLayer) Write(appFrame *ApplicationFrame) error {
	appBytes, err := appFrame.Encode()
	if err != nil {
		app.incrementErrTxCount()
		return err
	}
	return app.transport.Write(appBytes)
}

// Read 读取并解码应用帧
// 先通过 transport 读取数据，读取成功则对数据进行解码得到应用帧
func (app *GenericAppLayer) Read() (*ApplicationFrame, error) {
	bytes, err := app.transport.Read()
	if err != nil {
		app.incrementErrRxCount()
		return nil, err
	}
	frame, err := DecodeApplicationFrame(bytes)
	if err != nil {
		app.incrementErrRxCount()
		return nil, err
	}
	return frame, nil
}

// Status 获取传输层状态
func (app *GenericAppLayer) Status() error {
	return app.transport.Status()
}

// Close 关闭传输层连接
func (app *GenericAppLayer) Close() error {
	return app.transport.Close()
}

// incrementErrTxCount 增加发送错误计数
func (app *GenericAppLayer) incrementErrTxCount() {
	app.errTxCount++
}

// incrementErrRxCount 增加接收错误计数
func (app *GenericAppLayer) incrementErrRxCount() {
	app.errRxCount++
}
