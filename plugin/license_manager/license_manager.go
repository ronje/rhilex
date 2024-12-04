package licensemanager

/*
*
* 证书管理器
 */

import (
	"encoding/base64"
	"os"
	"runtime"
	"time"

	"github.com/hootrhino/rhilex-common-misc/misc"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"gopkg.in/ini.v1"
)

// LicenseManager 证书管理
type LicenseManager struct {
}

/*
*
* 开源版本默认是给一个授权证书
*
 */
func NewLicenseManager(r typex.Rhilex) *LicenseManager {
	return &LicenseManager{}
}
func validateLicense(key_path, license_path string) error {
	licBytesB64, err := os.ReadFile(license_path)
	if err != nil {
		utils.CLog("[LOAD LICENSE] load license file failed")
		os.Exit(0)
	}
	keyBytes, err := os.ReadFile(key_path)
	if err != nil {
		utils.CLog("[LOAD LICENSE] load key file failed")
		os.Exit(0)
	}
	licBytes, err2 := base64.StdEncoding.DecodeString(string(licBytesB64))
	if err2 != nil {
		utils.CLog("[LOAD LICENSE] decode key file failed")
		os.Exit(0)
	}
	privateKey, errParse := misc.ParsePrivateKey(keyBytes)
	if errParse != nil {
		utils.CLog("[LOAD LICENSE] parse key file failed")
		os.Exit(0)
	}
	adminSalt, err := misc.RSADecrypt(licBytes, privateKey)
	if err != nil {
		utils.CLog("[LOAD LICENSE] decrypt key file failed")
		os.Exit(0)
	}
	LocalLicense, err := utils.ParseAuthInfo(string(adminSalt))
	if err != nil {
		utils.CLog("[LOAD LICENSE] parse auth info failed")
		os.Exit(0)
	}
	LocalLicense.License = string(licBytesB64)
	T1 := time.UnixMilli(LocalLicense.BeginAuthorize)
	T2 := time.UnixMilli(LocalLicense.EndAuthorize)
	T1s := T1.Format("2006-01-02 15:04:05")
	T2s := T2.Format("2006-01-02 15:04:05")
	//
	if !LocalLicense.ValidateTime() {
		utils.CLog("[LOAD LICENSE] License has expired, Valid from %s to %s\n", T1s, T2s)
		os.Exit(0)
	}
	// validate local mac
	localMac := ""
	var err3 error
	if runtime.GOOS == "windows" {
		localMac, err3 = ossupport.GetWindowsFirstMacAddress()
	}
	if runtime.GOOS == "linux" {
		localMac, err3 = ossupport.GetLinuxMacAddr(LocalLicense.Iface)
	}
	if err3 != nil {
		utils.CLog("[LOAD LICENSE] fetch local mac address failed")
		os.Exit(0)
	}
	if localMac != LocalLicense.MAC {
		utils.CLog("[LOAD LICENSE] Local Mac:%s; certificate Mac:%s", localMac, LocalLicense.MAC)
		os.Exit(0)
	}
	typex.License = LocalLicense
	utils.CLog("[LOAD LICENSE] license load success")
	return nil
}
func (dm *LicenseManager) Init(section *ini.Section) error {
	license_path, err1 := section.GetKey("license_path")
	if err1 != nil {
		utils.CLog("[LOAD LICENSE] load license file failed")
		os.Exit(0)
	}
	key_path, err2 := section.GetKey("key_path")
	if err2 != nil {
		utils.CLog("[LOAD LICENSE] load key file failed")
		os.Exit(0)
	}
	return validateLicense(key_path.String(), license_path.String())
}

// Start 未实现
func (dm *LicenseManager) Start(typex.Rhilex) error {
	return nil
}
func (dm *LicenseManager) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

// Stop 未实现
func (dm *LicenseManager) Stop() error {
	return nil
}

func (dm *LicenseManager) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "LicenseManager",
		Name:        "LicenseManager",
		Version:     "v0.0.1",
		Description: "Rhilex License Manager",
	}
}
