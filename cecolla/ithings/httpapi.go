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

package ithings

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ApiResponse struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data SchemaResponse `json:"data"`
}

func (O ApiResponse) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type SchemaResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data SchemaSimple `json:"data"`
}

func (O SchemaResponse) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

/**
 * 获取物模型
 *
 */
func FetchIthingsSchema(host, productID, deviceName, mqttUser, mqttPwd string) (SchemaSimple, error) {
	// http://%s/api/v1/things/device/edge/send/{handle}/{type}/{productID}/{deviceName}
	url := "%s/api/v1/gateway/device/edge/send/gateway/thing"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{"method": "getSchema","msgToken": "%s"}`, uuid.NewString()))
	client := &http.Client{}
	fmt.Println(fmt.Sprintf(url, host), " Authorization:", BasicAuthEncode(mqttUser, mqttPwd))
	req, err := http.NewRequest(method, fmt.Sprintf(url, host), payload)
	if err != nil {
		return SchemaSimple{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", BasicAuthEncode(mqttUser, mqttPwd))
	// req.Header.Add("productID", productID)
	// req.Header.Add("deviceName", deviceName)
	res, err := client.Do(req)
	if err != nil {
		return SchemaSimple{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return SchemaSimple{}, err
	}
	fmt.Println(string(body))
	apiResponse := ApiResponse{}
	if errUnmarshal := json.Unmarshal(body, &apiResponse); errUnmarshal != nil {
		return SchemaSimple{}, errUnmarshal
	}
	return apiResponse.Data.Data, nil
}

func BasicAuthEncode(username, password string) string {
	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encodedAuth
}
