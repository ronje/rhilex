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
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

func Random(length int, level int) string {
	var characterSet string
	switch level {
	case 1:
		characterSet = "abcdefghijklmnopqrstuvwxyz"
	case 2:
		characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case 3:
		characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case 4:
		characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>/?"
	default:
		characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	var result []rune
	for i := 0; i < length; i++ {
		index := rand.Intn(len(characterSet))
		result = append(result, rune(characterSet[index]))
	}

	return string(result)
}

type HmacType = string

var (
	HmacTypeSha256 HmacType = "hmacsha256"
	HmacTypeSha1   HmacType = "hmacsha1"
	HmacTypeMd5    HmacType = "hmacmd5"
)

func HmacSha256(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha1(data string, secret []byte) string {
	h := hmac.New(sha1.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
