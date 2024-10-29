package apis

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

func InitBackupRoute() {
	backupApi := server.RouteGroup(server.ContextUrl("/backup"))
	{
		backupApi.GET(("/download"), server.AddRoute(DownloadSqlite))
		backupApi.POST(("/upload"), server.AddRoute(UploadSqlite))
		backupApi.GET(("/snapshot"), server.AddRoute(SnapshotDump))
		backupApi.GET(("/runningLog"), server.AddRoute(GetRunningLog))
	}
}

/*
*
* 备份Sqlite文件
*
 */
func DownloadSqlite(c *gin.Context, ruleEngine typex.Rhilex) {
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	files := []string{"./rhilex_datacenter.db", "./rhilex.db"}
	zipFilename := "./backup.zip"
	if err := utils.Zip(zipFilename, files); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	dir := wd
	c.Writer.WriteHeader(http.StatusOK)
	c.FileAttachment(fmt.Sprintf("%s/%s", dir, zipFilename),
		fmt.Sprintf("rhilex_backup_%d.zip", time.Now().UnixNano()))
}

/*
*
* 上传恢复
*
 */
const (
	zipHeaderBytes = 4
	zipHeader      = "PK\x03\x04"
)

// IsValidZip 检查给定的文件路径是否指向一个有效的ZIP文件。
// 它通过检查文件的头部签名来确定文件是否为ZIP文件。
func IsValidZip(filePath string) (bool, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()
	header := make([]byte, zipHeaderBytes)
	_, err = io.ReadFull(file, header)
	if err != nil {
		return false, fmt.Errorf("Error reading file header: %w", err)
	}
	isZip := bytes.Equal(header, []byte(zipHeader))
	return isZip, nil
}

func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*
*
* 上传zip文件, 必须包含数据中心和配置中心两个文件
*
 */
func UploadSqlite(c *gin.Context, ruleEngine typex.Rhilex) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	fileName := "recovery.zip"
	if err := os.MkdirAll(filepath.Dir(ossupport.RecoverBackupPath), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	if err := c.SaveUploadedFile(file, ossupport.RecoverBackupPath+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	if ok, err := IsValidZip(ossupport.RecoverBackupPath + fileName); !ok {
		os.Remove(ossupport.RecoverBackupPath + fileName)
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	// 解压
	if err := utils.Unzip(ossupport.RecoverBackupPath+fileName, ossupport.RecoverBackupPath); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	if ok, err := FileExists(ossupport.RecoveryDbPath); !ok {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	if ok, err := FileExists(ossupport.RecoveryDataCenterPath); !ok {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	if _, err := ReadSQLiteFileMagicNumber(ossupport.RecoveryDbPath); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	if _, err := ReadSQLiteFileMagicNumber(ossupport.RecoveryDataCenterPath); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData("Upload Db Backup Success"))
	ossupport.StartRecoverProcess()

}

// https://www.sqlite.org/fileformat.html
func ReadSQLiteFileMagicNumber(filePath string) ([16]byte, error) {
	MagicNumber := [16]byte{}
	file, err := os.Open(filePath)
	if err != nil {
		return MagicNumber, err
	}
	defer file.Close()
	binary.Read(file, binary.BigEndian, &MagicNumber)
	if string(MagicNumber[:]) == "SQLite format 3\x00" {
		return MagicNumber, nil
	}
	return MagicNumber, fmt.Errorf("Invalid Sqlite Db ,MagicNumber:%v error", MagicNumber)
}
