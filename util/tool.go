package util

import (
	"crypto/md5"
	"fmt"
	"github.com/ccqstark/gdufsclub/middleware"
	"os"
)

//加盐md5
func Md5SaltCrypt(str string) string {

	Cryptstr := "salt" + str + "shit"
	Cryptstr = fmt.Sprintf("%x", md5.Sum([]byte(Cryptstr)))
	return Cryptstr
}

// 判断所给路径文件/文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		middleware.Log.Error(err.Error())
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

