package main

import (
	"fmt"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
	"errors"
	"time"
	"io/ioutil"
)

func main() {
	const URL = "http://api.xfyun.cn/v1/service/v1/tts"
	const AUE = "lame"
	const APP_ID = "5b77e2bf"
	const API_KEY = "bc0736a9a9a6ca740c1311bf43821c4c"
	//# 硬盘路径(原视频存放路径)
	const SOURCE_PATH = "txt"
	////# 切割后的视频存放路径
	const RESOURCE_PATH = "audio"
	const EXTENSION = "txt"

	var (
		err error
		path string
		txtList []string
		getTxtContent []byte
	)


	if !getValid() {
		fmt.Println("过期了，找隔壁的老王再要一份吧。。")
		return
	}

	if path, err = getCurrentPath(); err != nil {
		fmt.Println("无法检测到目标文件夹")
		return
	}

	if !mkDir(path + SOURCE_PATH + checkSystemPath()) {
		fmt.Printf("无法新建 %s 文件夹，请手动创建", SOURCE_PATH)
		fmt.Println(".")
		return
	}

	if !mkDir(path + RESOURCE_PATH + checkSystemPath()) {
		fmt.Printf("无法新建 %s 文件夹，请手动创建", RESOURCE_PATH)
		fmt.Println(".")
		return
	}


	//检索txt文件夹
	txtList,err = getFileList(path + SOURCE_PATH + checkSystemPath(), EXTENSION)
	if len(txtList) == 0 {
		fmt.Printf("未检索到 %s 文件", EXTENSION)
		fmt.Println(".")
		goto ERR
	}

	for i := 0; i < len(txtList); i++ {
		getTxtContent,err = ReadAll(txtList[i])
		fmt.Println(txtList[i])
		fmt.Println(getTxtContent)
		//var str string = string(data[:])
		//fmt.Println("------.-")
	}
	fmt.Println("-- error --")
	fmt.Println(err)
ERR:
	// 结束
	end()
	return
}

func ReadAll(filePth string) ([]byte, error) {
	fmt.Println(filePth)
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err	}

	return ioutil.ReadAll(f)
}


//遍历所有文件并取出txt
func getFileList(filePath string,fileExtension string) (listFile []string, err error) {
	//var (
	//	osType = os.Getenv("GOOS") // 获取系统类型
	//)
	fmt.Println(filePath)
	fmt.Println("---===---")
	err = filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		//var strRet string
		//strRet, _ = os.Getwd()
		//if osType == "windows" {
		//	strRet += "\\"
		//} else if osType == "linux" {
		//	strRet += "/"
		//}
		//if f == nil {
		//	return err
		//}
		//if f.IsDir() {
		//	return nil
		//}
		//fmt.Println(strRet)
		//fmt.Println(path)
		//fmt.Println("-<<==---")
		//strRet += path //+ "\r\n"
		//用strings.HasSuffix(src, suffix)//判断src中是否包含 suffix结尾 EXTENSION
		ok := strings.HasSuffix(path, fileExtension)
		if ok {
			listFile = append(listFile, path) //将目录push到listFile []string中
		}
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	return
}



//取当前路径
func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

//检测有效期
func getValid() (res bool) {
	return time.Now().Unix() < 1533714306+30*24*60*60
}

//检测系统类型
func checkSystemPath() (path string) {
	//前边的判断是否是系统的分隔符
	if os.IsPathSeparator('\\') {
		path = "\\"
	} else {
		path = "/"
	}
	return
}

//检测文件夹是否存在，不存在则新建
func mkDir(path string) (res bool) {
	var err error
	res = false
	//文件是否存在，不存在则新建
	_, err = os.Stat(path);
	if err != nil {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("创建目录失败",path)
			return
		}
		res = true
		fmt.Printf("创建目录成功",path)
		return
	}
	res = true
	//fmt.Println("已存在目录" + path + "成功")
	return
}

//结束
func end() {
	var name string
	fmt.Println("------")
	fmt.Println("程序结束。")
	fmt.Scanf("%s", &name)
}
