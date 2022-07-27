package utils

import (
	"bufio"
	"io"
	"os"
)

// 逐行读取, 一行是一个[]byte, 多行就是[][]byte
func ReadByLine(filename string) (lines [][]byte, err error) {
	fp, err := os.Open(filename) // 获取文件指针
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	bufReader := bufio.NewReader(fp)
	for {
		line, err := bufReader.ReadBytes('\n') // 按行读
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
		} else {
			lines = append(lines, line)
		}
	}
	return
}

// * 一块一块地读取, 即给一个缓冲, 分多次读到缓冲中
//按字节读取，将整个文件读取到缓冲区buffer
func ReadByBlock(filename string) (content []byte, err error) {
	fp, err := os.Open(filename) // 获取文件指针
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	const bufferSize = 64 // 缓冲大小, 每次读取64个字节
	buffer := make([]byte, bufferSize)
	for {
		// 注意这里要取bytesRead, 否则有问题
		bytesRead, err := fp.Read(buffer) // 文件内容读取到buffer中
		content = append(content, buffer[:bytesRead]...)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				return nil, err
			}
		}
	}
	return
}
