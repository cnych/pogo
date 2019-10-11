package files

import (
	"os"
	"path"
	"strings"
)

/**
判断某个文件是否存在且有可执行权限
 */
func XOk(filepath string) (result bool) {
	fileInfo, err := os.Stat(filepath)
	if err != nil && os.IsNotExist(err) {
		return
	}
	mode := fileInfo.Mode().String()
	if strings.Index(mode, "x") > 0 {
		result = true
	}
	return
}

/**
获取某个命令的完整路径
 */
func Which(name string) (filePath string) {
	for _, value := range os.Environ() {
		if strings.Index(value, "PATH=") != 0 {
			continue
		}
		paths := strings.Split(value[5:], ":")
		for _, ph := range paths {
			fp := path.Join(ph, name)
			if XOk(fp) {
				filePath = fp
				break
			}
		}
		if len(filePath) > 0 {
			break
		}
	}
	return
}
