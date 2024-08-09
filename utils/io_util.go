package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// 直接把io包下面的同名函数抄过来了，加上了Context支持，主要解决读取超时问题
var errShortBuffer = errors.New("short buffer")
var errEOF = errors.New("EOF")
var errUnexpectedEOF = errors.New("unexpected EOF")
var errTimeout = errors.New("read Timeout")

/*
* 读取字节，核心原理是一个一个读，这样就不会出问题.
*
 */

func ReadAtLeast(ctx context.Context, r io.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		n = 0
		err = errShortBuffer
		return
	}
	for n < min && err == nil {
		select {
		case <-ctx.Done():
			err = errTimeout
			return
		default:
		}
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == errEOF {
		err = errUnexpectedEOF
	}
	return
}

/*
*
* 时间片读写请求(该函数已经被优化为ReadInLeastTimeout，除了历史遗留以后不要用，后期全面迁移)
*
 */
func SliceRequest(ctx context.Context,
	iio io.ReadWriteCloser, writeBytes []byte,
	resultBuffer []byte,
	showError bool,
	td time.Duration) (int, error) {
	_, errW := iio.Write(writeBytes)
	if errW != nil {
		return 0, errW
	}
	return SliceReceive(ctx, iio, resultBuffer, showError, td)
}

/*
*
* 读取数据的时候，如果出现错误就返回
*
 */
func SliceReceiveWithError(ctx context.Context,
	iio io.ReadWriteCloser, resultBuffer []byte, td time.Duration) (int, error) {
	return SliceReceive(ctx, iio, resultBuffer, true, td)
}

/*
*
* 读取数据的时候，如果出现错误则判断是否是串口引起的超时，如果是就忽略
*
 */
func SliceReceiveWithoutError(ctx context.Context,
	iio io.ReadWriteCloser, resultBuffer []byte, td time.Duration) (int, error) {
	return SliceReceive(ctx, iio, resultBuffer, false, td)
}

/*
*
* 通过一个定时时间片读取
*
 */
func SliceReceive(ctx context.Context, io io.ReadWriteCloser, resultBuffer []byte,
	showError bool, timeout time.Duration) (int, error) {
	N, B := ReadInLeastTimeout(ctx, io, timeout)
	copy(resultBuffer, B[:N])
	return N, nil
}

/*
*
* Slice分页计算器
*
 */
func Paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}

/*
*
* 自定义日志
*
 */
func CLog(format string, v ...interface{}) {
	timestamp := time.Now().UTC().Format("2006/01/02 15:04:05")
	logMsg := fmt.Sprintf(format, v...)
	logLine := fmt.Sprintf("[%s] %s", timestamp, logMsg)
	fmt.Print(logLine)
	fmt.Println()
}

/*
*
* 在timeout内读取N个字节
*
 */
func ReadInLeastTimeout(ctx context.Context,
	io io.ReadWriteCloser, timeout time.Duration) (int, []byte) {
	var responseData [256]byte
	CtxR, Cancel := context.WithTimeout(context.Background(), timeout)
	acc := 0
	defer Cancel()
	for {
		select {
		case <-ctx.Done():
			return acc, responseData[:acc]
		case <-CtxR.Done():
			return acc, responseData[:acc]
		default:
		}
		N, errRead := io.Read(responseData[acc:])
		if errRead != nil {
			if strings.Contains(errRead.Error(), "timeout") {
				if N > 0 {
					acc += N
				}
				continue
			}
		}
		if N > 0 {
			acc += N
		}
	}
}

/*
*
* 十六进制打印字节
*
 */

func BeautifulHex(b []byte) string {
	var sb strings.Builder
	for i, byteValue := range b {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%02X", byteValue))
	}
	return sb.String()
}
