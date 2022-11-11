package tools

import (
	"errors"
	"fmt"
	"os"
)

//判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

//创建文件夹
func CreateDir(path string) error {
	_exist, _err := HasDir(path)
	if _err != nil {
		return errors.New(fmt.Sprintf("获取文件夹异常 -> %v\n", _err))
	}
	if _exist {
		return nil
	}
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("获取文件夹异常 -> %v\n", _err))
	}
	return nil

}
