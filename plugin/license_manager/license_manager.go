package licensemanager

/*
*
* 证书管理器
 */

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/hootrhino/rhilex/glogger"
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
	errMsg := "License loading failed."
	licBytesB64, err := os.ReadFile(license_path)
	if err != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	keyBytes, err := os.ReadFile(key_path)
	if err != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	licBytes, err := base64.StdEncoding.DecodeString(string(licBytesB64))
	if err != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	adminSalt, err := RSADecrypt(licBytes, keyBytes)
	if err != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	LocalLicense, err := utils.ParseAuthInfo(string(adminSalt))
	if err != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	LocalLicense.License = string(licBytesB64)
	T1 := time.UnixMilli(LocalLicense.BeginAuthorize)
	T2 := time.UnixMilli(LocalLicense.EndAuthorize)
	T1s := T1.Format("2006-01-02 15:04:05")
	T2s := T2.Format("2006-01-02 15:04:05")
	//
	if !LocalLicense.ValidateTime() {
		fmt.Printf("License has expired, Valid from %s to %s\n", T1s, T2s)
		glogger.GLogger.Fatalf("License has expired, Valid from %s to %s", T1s, T2s)
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
		fmt.Println(errMsg)
		fmt.Println(err3)
		glogger.GLogger.Fatal(err3)
		os.Exit(0)
	}
	if localMac != LocalLicense.MAC {
		fmt.Println(errMsg)
		glogger.GLogger.Debugf("Local Mac:%s; certificate Mac:%s", localMac, LocalLicense.MAC)
		glogger.GLogger.Fatal("Local certificate and hardware information do not match.")
		os.Exit(0)
	}
	typex.License = LocalLicense
	fmt.Println("[∫∫] -----------------------------------")
	fmt.Println("|>>| * Type     *", LocalLicense.Type)
	fmt.Println("|>>| * DeviceID *", LocalLicense.DeviceID)
	fmt.Println("|>>| * Admin    *", LocalLicense.AuthorizeAdmin)
	fmt.Println("|>>| * MacAddr  *", LocalLicense.MAC)
	fmt.Println("|>>| * Begin    *", T1s)
	fmt.Println("|>>| * End      *", T2s)
	fmt.Println("[∫∫] -----------------------------------")
	return nil
}
func (dm *LicenseManager) Init(section *ini.Section) error {
	errMsg := "License loading failed."
	license_path, err1 := section.GetKey("license_path")
	if err1 != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
		os.Exit(0)
	}
	key_path, err2 := section.GetKey("key_path")
	if err2 != nil {
		fmt.Println(errMsg)
		glogger.GLogger.Fatal(errMsg)
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
