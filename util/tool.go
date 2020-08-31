package util

import (
	"crypto/md5"
	"fmt"
)

//加盐md5
func Md5SaltCrypt(str string) string {

	Cryptstr := "salt" + str + "shit"
	Cryptstr = fmt.Sprintf("%x", md5.Sum([]byte(Cryptstr)))
	return Cryptstr
}
