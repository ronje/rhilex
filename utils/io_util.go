package utils

import (
	"context"
	"errors"
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
* 通过一个定时时间片读取
*
 */
func SliceReceive(ctx context.Context, io io.ReadWriteCloser, resultBuffer []byte,
	showError bool, timeout time.Duration) (int, error) {
	N, B, E := ReadInWill(ctx, io, timeout)
	if E != nil {
		if showError {
			if strings.Contains(E.Error(), "timeout") {
				return N, nil
			}
			return N, E
		}
		return N, E
	}
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
* 在timeout内读取N个字节
*
 */
func ReadInWill(ctx context.Context, rwc io.ReadWriteCloser, timeout time.Duration) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var responseData []byte
	totalRead := 0
	tempBuffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return totalRead, responseData, ctx.Err()
			}
			return totalRead, responseData, ctx.Err()
		default:
			n, err := rwc.Read(tempBuffer)
			if err != nil {
				if err == io.EOF {
					return totalRead, responseData, nil
				}
				if ctx.Err() == context.DeadlineExceeded {
					return totalRead, responseData, ctx.Err()
				}
				return totalRead, responseData, err
			}
			if n > 0 {
				responseData = append(responseData, tempBuffer[:n]...)
				totalRead += n
			}
		}
	}
}
