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
*/

package utils

import (
	"bufio"
	"io"
)

func GetImageFormat(file io.Reader) string {
	tmpFile := bufio.NewReader(file)
	bytes := make([]byte, 4)
	bytes, _ = tmpFile.Peek(4)
	if len(bytes) < 4 {
		return ""
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return ".png"
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return ".jpg"
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return ".gif"
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return ".bmp"
	}
	return ""
}
