package common

import (
	"crypto/sha1"
	"fmt"
	"crypto/md5"
)

func SHA1string(s string)(string){
	digestbyte := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", digestbyte)
}

func GenSessionToken(useruuid, sessionuuid, password string)string {
	return SHA1string(useruuid + sessionuuid + password)
}

func MD5string(s string)(string){
	digestbyte := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", digestbyte)
}
