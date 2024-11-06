package utils

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/glogger"
)

// CreateZip 压缩指定的文件列表到一个ZIP文件中。
func Zip(zipFilename string, filenames []string) error {
	// 创建ZIP文件
	zipFile, err := os.Create(zipFilename)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建一个新的ZIP压缩器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历要压缩的文件列表
	for _, filename := range filenames {
		// 打开要压缩的文件
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// 获取文件信息
		info, err := file.Stat()
		if err != nil {
			return err
		}

		// 创建ZIP文件中的文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置压缩文件的名称（在ZIP文件中），不包含路径
		header.Name = filepath.Base(filename)
		header.Method = zip.Deflate

		// 将文件头写入ZIP文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 将文件内容写入ZIP文件
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}

	return nil
}

// 解压缩文件到指定目录
func Unzip(zipPath string, targetDir string) error {
	// 打开要解压的 zip 文件
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 遍历 zip 文件中的文件并解压
	for _, file := range zipFile.File {
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destPath := targetDir + "/" + file.Name
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
*
* GetPwd
*
 */
func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	return dir
}

/*
*
* DEBUG使用
*
 */
func TraceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	var info [7]float64
	info[0] = float64(ms.HeapObjects)
	info[1] = BtoMB(ms.HeapAlloc)
	info[2] = BtoMB(ms.TotalAlloc)
	info[3] = BtoMB(ms.HeapSys)
	info[4] = BtoMB(ms.HeapIdle)
	info[5] = BtoMB(ms.HeapReleased)
	info[6] = BtoMB(ms.HeapIdle - ms.HeapReleased)

	for _, v := range info {
		fmt.Printf("%v,\t", v)
	}
	fmt.Println()
}
func BtoMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

/*
*
* Byte to Mbyte
*
 */
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

/*
*
* 获取操作系统发行版版本
runtime.GOARCH:

	386: 32-bit Intel/AMD x86 architecture
	amd64: 64-bit Intel/AMD x86 architecture
	arm: ARM architecture (32-bit)
	arm64: ARM architecture (64-bit)
	ppc64: 64-bit PowerPC architecture
	ppc64le: 64-bit little-endian PowerPC architecture
	mips: MIPS architecture (32-bit)
	mips64: MIPS architecture (64-bit)
	s390x: IBM System z architecture (64-bit)
	wasm: WebAssembly architecture

runtime.GOOS:

	darwin: macOS
	freebsd: FreeBSD
	linux: Linux
	windows: Windows
	netbsd: NetBSD
	openbsd: OpenBSD
	plan9: Plan 9
	dragonfly: DragonFly BSD

*
*/
func GetOSDistribution() (string, error) {
	if runtime.GOOS == "windows" {
		return runtime.GOOS, nil
	}
	// Linux 有很多发行版, 目前特别要识别一下Openwrt
	if runtime.GOOS == "linux" {
		if PathExists("/etc/os-release") {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "cat", "/etc/os-release")
			output, err := cmd.Output()
			if err != nil {
				return runtime.GOOS, err
			}
			osIssue := strings.ToLower(string(output))
			if strings.Contains((osIssue), "openwrt") {
				return "openwrt", nil
			}
			if strings.Contains((osIssue), "ubuntu") {
				return "ubuntu", nil
			}
			if strings.Contains((osIssue), "debian") {
				return "debian", nil
			}
			if strings.Contains((osIssue), "armbian") {
				return "armbian", nil
			}
			if strings.Contains((osIssue), "deepin") {
				return "deepin", nil
			}
		}
	}
	return runtime.GOOS, nil
}
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

/*
*
* 获取Ubuntu的版本
*
 */
func GetUbuntuVersion() (string, error) {
	// lsb_release -ds -> Ubuntu 22.04.1 LTS
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "lsb_release", "-ds")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	info := strings.ToLower(strings.TrimSpace(string(output)))
	if strings.Contains(info, "ubuntu") {
		if strings.Contains(info, "16.") {
			return "ubuntu16", nil
		}
		if strings.Contains(info, "18.") {
			return "ubuntu18", nil
		}
		if strings.Contains(info, "20.") {
			return "ubuntu20", nil
		}
		if strings.Contains(info, "22.") {
			return "ubuntu22", nil
		}
		if strings.Contains(info, "24.") {
			return "ubuntu24", nil
		}
	}
	return "", fmt.Errorf("unsupported OS:%s", info)
}

/*
*
* 检查命令是否存在
*
 */

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
