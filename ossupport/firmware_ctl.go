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

package ossupport

import (
	"archive/zip"
	"debug/elf"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/*
*
* 给Linux下 ELF 文件增加可执行权限
*
 */
func ChmodX(filePath string) error {
	file, err := elf.Open(filePath)
	if err != nil {
		return err
	}
	elfHeader := file.FileHeader
	file.Close()
	if elfHeader.Type == elf.ET_EXEC {
		if err := os.Chmod(filePath, 0755); err != nil {
			return err
		}
		return nil
	}
	return nil
}

/*
*
* Stop RHILEX
*
 */
func StopRhilex() error {
	bytes, err1 := os.ReadFile(MainExePidPath)
	if err1 != nil {
		return fmt.Errorf("ReadFile error: %s", err1)
	}
	pid, err2 := strconv.Atoi(string(bytes))
	if err2 != nil {
		return fmt.Errorf("strconv error: %s", err2)
	}
	process, err3 := os.FindProcess(pid)
	if err3 != nil {
		return fmt.Errorf("FindProcess error: %s", err3)
	}
	err4 := process.Signal(syscall.SIGINT)
	if err4 != nil {
		err4 = process.Signal(syscall.SIGTERM)
	}
	if err4 != nil {
		return fmt.Errorf("Signal error: %d", err4)
	}
	time.Sleep(1 * time.Second)
	if process != nil {
		err5 := process.Kill()
		if err5 != nil {
			return fmt.Errorf("Kill error: %d", err5)
		}
	}

	return nil
}

/*
*
* 重启
*
 */
func RestartRhilex() error {
	return StopRhilex()
}

/*
*
* 恢复上传的DB
1 停止RHILEX
2 删除老DB
3 复制新DB到路径
3 删除PID,停止守护进程
4 重启(脚本会新建PID)
- path: /usr/local/rhilex, args: recover=true
*
*/
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

/**
 * 兼容Linux、Windows的文件名
 *
 */
func GetExePath() string {
	if runtime.GOOS == "linux" {
		return "rhilex"
	}
	if runtime.GOOS == "windows" {
		return "rhilex.exe"
	}
	return ""
}
func GetUpgraderPath() string {
	if runtime.GOOS == "linux" {
		return "rhilex-upgrader"
	}
	if runtime.GOOS == "windows" {
		return "rhilex-upgrader.exe"
	}
	return ""
}

/*
*
* 数据备份
*
 */
func StartRecoverProcess() {
	cmd := exec.Command("./rhilex", "recover", "-recover=true")
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		log.Println("Start Recover Process Failed:", err)
		return
	}
	os.Exit(0)
}

/*
*
* 启用升级进程
*
 */
func StartUpgradeProcess(s1, s2, s3, s4 string) {
	cmd := exec.Command(MainWorkDir+GetUpgraderPath(), "upgrade",
		"-upgrade=true",
		fmt.Sprintf("-inipath=%s", s1),
		fmt.Sprintf("-licpath=%s", s2),
		fmt.Sprintf("-keypath=%s", s3),
		fmt.Sprintf("-rundbpath=%s", s4),
	)
	cmd.SysProcAttr = NewSysProcAttr()
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		log.Println("Start Upgrade Process Failed:", err)
		return
	}
	os.Exit(0)
}

func Reboot() error {
	return RebootLocal()
}

/**
 * 解压文件
 *
 */
func UnzipFirmware(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		// 检查文件路径是否安全，防止 zip slip 攻击
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// 创建目录
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		// 创建文件
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outFile.Close()
		// 解压文件内容
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	return nil
}

/*
*
  - Unzip 指令包装,不通用

unzip -o -d /usr/local ./zupload/firmware.zip
*/
func Unzip(zipFile, destDir string) error {
	cmd := exec.Command("unzip", "-o", "-d", destDir, zipFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unzip file: %s, %s", err.Error(), string(out))
	}
	return nil
}

/**
 * 备份老版本
 *
 */

func BackupOldVersion(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destDir := filepath.Dir(dest)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return err
		}
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	return err
}

/*
*
* 移动文件
*
 */
func MoveFile(sourcePath, destPath string) error {

	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}
