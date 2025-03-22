// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package apis

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"gopkg.in/ini.v1"
)

/*
*
  - 上传最新固件, 必须是ZIP包

*
*/
func UploadFirmWare(c *gin.Context, ruleEngine typex.Rhilex) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	saveDir := "./zupgrade/"
	if !utils.PathExists(saveDir) {
		if err := os.MkdirAll(filepath.Dir(saveDir), os.ModePerm); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if err := c.SaveUploadedFile(file, ossupport.FirmwarePath); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 解压、升级
*
 */
func UpgradeFirmWare(c *gin.Context, ruleEngine typex.Rhilex) {
	file, errOpenFile := os.Create(ossupport.UpgradeLogPath)
	if errOpenFile != nil {
		c.JSON(common.HTTP_OK, common.Error(errOpenFile.Error()))
		return
	}
	defer file.Close()
	os.Stdout = file
	os.Stderr = file

	glogger.DefaultOutput("[RHILEX UPGRADE] Current Version: %s", typex.MainVersion)
	uploadPath := "./zupgrade/"        // 固定路径
	tempPath := uploadPath + "temp001" // 固定路径
	errMkdirAll := os.MkdirAll(tempPath, os.ModePerm)
	if errMkdirAll != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] %s", errMkdirAll.Error())))
		return
	}
	// 提前解压文件
	if errUnzip := ossupport.Unzip(ossupport.FirmwarePath, tempPath); errUnzip != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] Unzip error:%s", errUnzip.Error())))
		return
	}
	// 检查 /tmp/temp001/rhilex 的Md5
	md51, errSumMD5 := sumMD5(tempPath + "/" + ossupport.GetExePath())
	if errSumMD5 != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] GetExePath error:%s", errSumMD5.Error())))
		return
	}
	errCheckFileType := CheckFileType(tempPath + "/" + ossupport.GetExePath())
	if errCheckFileType != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] CheckFileType error:%s", errCheckFileType.Error())))
		return
	}
	// 从解压后的目录提取Md5
	readBytes, errReadFile := os.ReadFile(tempPath + "/md5.sum")
	if errReadFile != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] ReadFile md5.sum error:%s", errReadFile.Error())))
		return
	}
	glogger.DefaultOutput("[RHILEX UPGRADE] Compare MD5:[%s]~[%s]", md51, string(readBytes))
	if md51 != string(readBytes) {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] Invalid sum md5!")))
		return
	}
	if errMov := MoveFile(tempPath+"/"+ossupport.GetExePath(),
		ossupport.MainWorkDir+ossupport.GetUpgraderPath()); errMov != nil {

		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] Move Upgrader Failed:%s", errMov)))
		return
	}
	if errChmod := ossupport.ChmodX(ossupport.MainWorkDir + ossupport.GetUpgraderPath()); errChmod != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] ChmodX Upgrader Failed:%s", errChmod)))
		return
	}
	glogger.DefaultOutput("[RHILEX UPGRADE] Remove Temp File:%s", tempPath)
	if errRm := os.RemoveAll(tempPath); errRm != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] Remove Temp File:%s", errRm)))
		return
	}
	glogger.DefaultOutput("[RHILEX UPGRADE] RuleEngine GetConfig:%s", tempPath)
	IniPath := ruleEngine.GetConfig().IniPath
	mainCfg, errLoad := ini.Load(IniPath)
	if errLoad != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] Load Ini File error:%s", errLoad)))
		return
	}
	section := mainCfg.Section("plugin.license_manager")
	license_path, err1 := section.GetKey("license_path")
	if err1 != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] GetKey license_manager.license_path error:%s", err1)))
		return
	}
	key_path, err2 := section.GetKey("key_path")
	if err2 != nil {
		c.JSON(common.HTTP_OK,
			common.Error(glogger.DefaultOutput("[RHILEX UPGRADE] GetKey license_manager.key_path error:%s", err2)))
		return
	}
	glogger.DefaultOutput("[RHILEX UPGRADE] Start Upgrade Process: license_path=%s, key_path=%s",
		license_path.String(), key_path.String())
	go func() {
		time.Sleep(2 * time.Second)
		ossupport.StartUpgradeProcess(IniPath, license_path.String(), key_path.String(), "./rhilex.db")
	}()
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
  - 检查包, 一般包里面会有一个可执行文件和 MD5 SUM 值。要对比一下。
    文件列表:
  - rhilex
  - rhilex.ini
  - md5.sum

*
*/

func sumMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

/*
*
* 移动文件
*
 */
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func CheckFileType(filePath string) error {
	currentArch := runtime.GOARCH
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var magicNumber [4]byte
	_, err = file.Read(magicNumber[:])
	if err != nil {
		return err
	}
	switch {
	case bytes.Equal(magicNumber[:], []byte{0x7F, 'E', 'L', 'F'}):
		// ELF文件，用于Linux
		elfArch, err := checkELFArch(file)
		if err != nil {
			return err
		}
		if elfArch != currentArch {
			return fmt.Errorf("ELF architecture mismatch: %s != %s", elfArch, currentArch)
		}
	case CheckPEFileMagic(magicNumber):
		return nil //fmt.Errorf("not support windows PE Format")
	case CheckDOSHeaderMagic(magicNumber):
		return nil //fmt.Errorf("not support windows DOS Format")
	default:
		return fmt.Errorf("unknown file type")
	}

	return nil
}

// checkELFArch 检查ELF文件的架构
func checkELFArch(file *os.File) (string, error) {
	type elfHeader struct {
		Ident     [16]byte
		Type      uint16
		Machine   uint16
		Version   uint32
		Entry     uint64
		Phoff     uint64
		Shoff     uint64
		Flags     uint32
		Ehsize    uint16
		Phentsize uint16
		Phnum     uint16
		Shentsize uint16
		Shnum     uint16
		Shstrndx  uint16
	}
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	var hdr elfHeader
	err = binary.Read(file, binary.LittleEndian, &hdr)
	if err != nil {
		return "", err
	}
	switch hdr.Machine {
	case 3:
		return "386", nil // x86
	case 62:
		return "amd64", nil // x86_64
	case 40:
		return "arm", nil // ARM
	default:
		return "", fmt.Errorf("unknown ELF architecture")
	}
}

func CheckPEFileMagic(data [4]byte) bool {
	return (uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24) == 0x50450000
}

func CheckDOSHeaderMagic(data [4]byte) bool {
	return (uint32(data[0]) | uint32(data[1])<<8) == 0x5A4D
}
