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
	Type              string `json:"type"` // FREETRIAL | COMMERCIAL
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
** License Type       : %s
** Device SN          : %s
** Authorize Admin    : %s
** Authorize Password : %s
** Begin Authorize    : %s
** End Authorize      : %s
** Authorized MAC     : %s
`,
		ll.Type,
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

// type & 00001 & rhino & hoot & eth0 & FF:FF:FF:FF:FF:FF & 0 & 0
func ParseAuthInfo(info string) (LocalLicense, error) {
	var ll LocalLicense

	ss := strings.Split(info, "&")
	if len(ss) != 8 {
		return ll, fmt.Errorf("failed to parse: %s", info)
	}

	beginAuthorize, err1 := strconv.ParseInt(ss[6], 10, 64)
	if err1 != nil {
		return ll, fmt.Errorf("failed to parse BeginAuthorize: %w", err1)
	}
	endAuthorize, err2 := strconv.ParseInt(ss[7], 10, 64)
	if err2 != nil {
		return ll, fmt.Errorf("failed to parse EndAuthorize: %w", err2)
	}
	if ss[0] == "" {
		ll.Type = "FREETRIAL"
	} else {
		ll.Type = ss[0]
	}
	ll.DeviceID = ss[1]
	ll.AuthorizeAdmin = ss[2]
	ll.AuthorizePassword = ss[3]
	ll.Iface = ss[4]
	ll.MAC = ss[5]
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
