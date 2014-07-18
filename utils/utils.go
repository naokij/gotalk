/*
Copyright 2014 Jiang Le

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

部分代码来自 https://github.com/beego/wetalk/blob/master/modules/utils/tools.go

*/

package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
)

var (
	HZRange = [15][2]int{
		{11904, 11929},
		{11931, 12019},
		{12032, 12245},
		{12293, 12293},
		{12295, 12295},
		{12321, 12329},
		{12344, 12347},
		{13312, 19893},
		{19968, 40908},
		{63744, 64109},
		{64112, 64217},
		{131072, 173782},
		{173824, 177972},
		{177984, 178205},
		{194560, 195101},
	}
)

//按照一个汉字算两个英文的方式计算字符串长度
func HZStringLength(s string) (length int) {
	r := []rune(s)
	for _, v := range r {
		var isHz bool
		for _, subRange := range HZRange {
			if int(v) >= subRange[0] && int(v) <= subRange[1] {
				isHz = true
				break
			}
		}
		if isHz {
			length += 2
		} else {
			length++
		}
	}
	return length
}

// Random generate string
func GetRandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// Encode string to md5 hex value
func EncodeMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// convert any numeric value to int64
func ToInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}
