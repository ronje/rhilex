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

package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/ossupport"
)

// FetchLoadLicense rhilex active -H 127.0.0.1 -U admin -P 123456
func FetchLoadLicense(host, SN, username, password, Iface, macAddr string) error {
	activeParams := fmt.Sprintf(`%s&%s&%s&%s&%s&0&0`,
		SN, username, password, Iface, macAddr)
	CLog("\n*>> BEGIN LICENCE ACTIVE\n"+
		"*# Vendor Admin: (%s, %s).\n"+
		"*# Local Iface: (%s).\n"+
		"*# Local Mac Address: (%s).\n"+
		"*# Try to request license from server:(%s).",
		username, password, Iface, macAddr, host)
	filePath := fmt.Sprintf("license_%v.zip", time.Now().UnixMilli())
	err := Download(host, activeParams, filePath)
	if err != nil {
		return fmt.Errorf("Request failed")
	}
	fmt.Println("*# License fetch success, save as: " + filePath)
	fmt.Println("*<< END LICENCE ACTIVE")
	return nil
}

func RSADecrypt(License, Key []byte) ([]byte, error) {
	block, _ := pem.Decode(Key)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, License)
}

/*
*
* MD5
*
 */
func SumMd5(inputString string) string {
	hasher := md5.New()
	io.WriteString(hasher, inputString)
	hashBytes := hasher.Sum(nil)
	md5String := fmt.Sprintf("%x", hashBytes)
	return md5String
}

type LocalLicense struct {
	DeviceID          string `json:"device_id"`
	AuthorizeAdmin    string `json:"authorize_admin"`
	AuthorizePassword string `json:"authorize_password"`
	BeginAuthorize    int64  `json:"begin_authorize"`
	EndAuthorize      int64  `json:"end_authorize"`
	Iface             string `json:"iface"`
	MAC               string `json:"mac"`
	License           string `json:"license"`
}

func (ll *LocalLicense) ToString() string {
	beginTime := time.UnixMilli(ll.BeginAuthorize).Format(time.RFC3339)
	endTime := time.UnixMilli(ll.EndAuthorize).Format(time.RFC3339)
	return fmt.Sprintf(`
** Device SN          : %s
** Authorize Admin    : %s
** Authorize Password : %s
** Begin Authorize    : %s
** End Authorize      : %s
** Authorized MAC     : %s
`,
		ll.DeviceID,
		ll.AuthorizeAdmin,
		ll.AuthorizePassword,
		beginTime,
		endTime,
		ll.MAC)
}

func (ll LocalLicense) ValidateTime() bool {
	Now := time.Now().UnixMilli()
	V := ll.EndAuthorize - Now
	if (ll.BeginAuthorize > Now) && (V <= 0) {
		return false
	}
	return true
}

// 00001 & rhino & hoot & eth0 & FF:FF:FF:FF:FF:FF & 0 & 0
func ParseAuthInfo(info string) (LocalLicense, error) {
	var ll LocalLicense
	ss := strings.Split(info, "&")
	if len(ss) != 7 {
		return ll, fmt.Errorf("failed to parse: %s", info)
	}

	beginAuthorize, err1 := strconv.ParseInt(ss[5], 10, 64)
	if err1 != nil {
		return ll, fmt.Errorf("failed to parse BeginAuthorize: %w", err1)
	}
	endAuthorize, err2 := strconv.ParseInt(ss[6], 10, 64)
	if err2 != nil {
		return ll, fmt.Errorf("failed to parse EndAuthorize: %w", err2)
	}

	ll.DeviceID = ss[0]
	ll.AuthorizeAdmin = ss[1]
	ll.AuthorizePassword = ss[2]
	ll.Iface = ss[3]
	ll.MAC = ss[4]
	ll.BeginAuthorize = beginAuthorize
	ll.EndAuthorize = endAuthorize
	return ll, nil
}
func ValidateLicense(key_path, license_path string) (LocalLicense, error) {
	LocalLicense := LocalLicense{}
	licBytesB64, err := os.ReadFile(license_path)
	if err != nil {
		return LocalLicense, err
	}
	keyBytes, err := os.ReadFile(key_path)
	if err != nil {
		return LocalLicense, err
	}
	licBytes, err := base64.StdEncoding.DecodeString(string(licBytesB64))
	if err != nil {
		return LocalLicense, err
	}
	adminSalt, err := RSADecrypt(licBytes, keyBytes)
	if err != nil {
		return LocalLicense, err
	}

	LocalLicense, err = ParseAuthInfo(string(adminSalt))
	if err != nil {
		return LocalLicense, err
	}
	//
	if !LocalLicense.ValidateTime() {
		return LocalLicense, fmt.Errorf("Invalid Auth Time!")
	}
	localMac := ""
	var err3 error
	if runtime.GOOS == "windows" {
		localMac, err3 = ossupport.GetWindowsFirstMacAddress()
	}
	if runtime.GOOS == "linux" {
		localMac, err3 = ossupport.GetLinuxMacAddr(LocalLicense.Iface)
	}
	if err3 != nil {
		return LocalLicense, err3
	}
	if localMac != LocalLicense.MAC {
		return LocalLicense, fmt.Errorf("Local mac is not matched!")
	}
	return LocalLicense, nil
}
