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

package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	PUB_KEY  = "RSA PUBLIC KEY"
	PRIV_KEY = "RSA PRIVATE KEY"
)

func InitSecurityLicense() {
	GenLocalSecurityLicense()
}

// bits 生成的公私钥对的位数，一般为 1024 或 2048
// privateKey 生成的私钥
// publicKey 生成的公钥
func GenLocalSecurityLicense() (privateKey, publicKey string) {
	priKey, err2 := rsa.GenerateKey(rand.Reader, 2048)
	if err2 != nil {
		panic(err2)
	}

	derStream := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	prvKey := pem.EncodeToMemory(block)
	puKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(puKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey := pem.EncodeToMemory(block)
	privateKey = string(prvKey)
	publicKey = string(pubKey)
	os.WriteFile("./.encrypt.priv", []byte(privateKey), 0755)
	os.WriteFile("./.encrypt.pub", []byte(publicKey), 0755)
	return
}

// RSA加密
// plainText 要加密的数据
// path 公钥匙文件地址
func RSA_Encrypt(plainText []byte, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	block, _ := pem.Decode(buf)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// RSA解密
func RSA_Decrypt(cipherText []byte, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err1 := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err1 != nil {
		return nil, err1
	}
	return plainText, nil
}

// CheckLicense 检查本地证书文件是否存在
func CheckLicense() bool {
	_, err := os.Stat("./.encrypt.priv")
	if os.IsNotExist(err) {
		//fmt.Printf("Private key file '%s' does not exist.\n", ".encrypt.priv")
		return false
	}
	_, err = os.Stat("./.encrypt.pub")
	if os.IsNotExist(err) {
		//fmt.Printf("Public key file '%s' does not exist.\n", ".encrypt.pub")
		return false
	}
	return true
}

// ReadLocalKeys 读取本地的私钥和公钥文件，并返回它们的字符串表示
func ReadLocalKeys() (string, string, error) {
	privateKeyPEM, err := os.ReadFile("./.encrypt.priv")
	if err != nil {
		return "", "", fmt.Errorf("failed to read private key file: %v", err)
	}
	publicKeyPEM, err := os.ReadFile("./.encrypt.pub")
	if err != nil {
		return "", "", fmt.Errorf("failed to read public key file: %v", err)
	}
	return string(privateKeyPEM), string(publicKeyPEM), nil
}
