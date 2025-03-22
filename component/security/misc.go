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

var privateKeyPath = "./.encrypt.priv"
var publicKeyPath = "./.encrypt.pub"

// InitSecurityLicense 初始化安全证书
func InitSecurityLicense() error {

	privateKey, publicKey, err := GenLocalSecurityLicense()
	if err != nil {
		return err
	}
	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		return fmt.Errorf("failed to write private key file: %v", err)
	}
	err = os.WriteFile(publicKeyPath, []byte(publicKey), 0644)
	if err != nil {
		return fmt.Errorf("failed to write public key file: %v", err)
	}
	return nil
}

// GenLocalSecurityLicense 生成本地安全证书
func GenLocalSecurityLicense() (privateKey, publicKey string, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key pair: %v", err)
	}

	// 编码私钥
	derStream := x509.MarshalPKCS1PrivateKey(priKey)
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	privateKey = string(pem.EncodeToMemory(privateKeyBlock))

	// 编码公钥
	puKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(puKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %v", err)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	publicKey = string(pem.EncodeToMemory(publicKeyBlock))

	return privateKey, publicKey, nil
}

// readKeyFromFile 从文件中读取密钥
func readKeyFromFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}
	return data, nil
}

// parsePublicKey 解析公钥
func parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return publicKey, nil
}

// parsePrivateKey 解析私钥
func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	return privateKey, nil
}

// RSA_Encrypt RSA加密
func RSA_Encrypt(plainText []byte, publicKeyPath string) ([]byte, error) {
	publicKeyData, err := readKeyFromFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := parsePublicKey(publicKeyData)
	if err != nil {
		return nil, err
	}
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return cipherText, nil
}

// RSA_Decrypt RSA解密
func RSA_Decrypt(cipherText []byte, privateKeyPath string) ([]byte, error) {
	privateKeyData, err := readKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := parsePrivateKey(privateKeyData)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return plainText, nil
}

// CheckLicense 检查本地证书文件是否存在
func CheckLicense() bool {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// ReadLocalKeys 读取本地的私钥和公钥文件，并返回它们的字符串表示
func ReadLocalKeys() (string, string, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read private key file: %v", err)
	}
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read public key file: %v", err)
	}
	return string(privateKeyPEM), string(publicKeyPEM), nil
}
