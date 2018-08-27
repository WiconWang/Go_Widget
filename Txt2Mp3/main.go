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
	"log"
	"encoding/base64"
	"crypto/md5"
	"net/http"
	"strconv"
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
	const NEW_EXTENSION = "mp3"

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

	var myIP = getIp()

	//检索txt文件夹
	txtList,err = getFileList(path + SOURCE_PATH + checkSystemPath(), EXTENSION)
	if len(txtList) == 0 {
		fmt.Printf("未检索到 %s 文件", EXTENSION)
		fmt.Println(".")
		goto ERR
	}

	for i := 0; i < len(txtList); i++ {
		getTxtContent,err = ReadAll(txtList[i]); if err == nil {
			fmt.Printf("文件  %s 转换启动",filepath.Base(txtList[i]))
			fmt.Println(" --> ")
			contents := string(getTxtContent[:])
			body := sendPost(URL,AUE,APP_ID,API_KEY,myIP,contents)
			//var data []byte = []byte(contents)
			writeMp3(path + RESOURCE_PATH + checkSystemPath(),
				strings.Replace(
					filepath.Base(txtList[i]),
					filepath.Ext(filepath.Base(txtList[i])),
					"."+NEW_EXTENSION,
					-1),
				body)
			fmt.Printf("转换完成")
			fmt.Println("------")
		}
	}
	fmt.Println("-- FINISH --")

ERR:
	// 结束
	end()
	return
}


func sendPost(url,aue,appId,appKey,myIP,content string) (body []byte){
	var (
		curTime string
		param string
		checkSum string
		//jsonStr string
	)

	// 预生成一些校验数据
	curTime = strconv.FormatInt(time.Now().Unix(),10)
	param = "{\"aue\":\"" + aue + "\",\"auf\":\"audio/L16;rate=16000\",\"voice_name\":\"xiaoyan\",\"engine_type\":\"intp65\"}"
	paramBase64 := base64.StdEncoding.EncodeToString([]byte(param))
	checkSum = fmt.Sprintf("%x", md5.Sum([]byte(appKey + curTime + paramBase64))) //将[]byte转成16进制

	// 清理无关符号
	content = strings.Replace(content,"\"","",-1)
	content = strings.Replace(content,"'","",-1)
	content = strings.Replace(content,"`","",-1)
	content = strings.Replace(content,"\n\t","",-1)
	content = strings.Replace(content,"\n","",-1)
	//jsonStr := []byte(`{"text": "`+ content +`"}`)
	//fmt.Println(jsonStr)
	// 启动POST请求

	//data := make(url.Values)
	//data["text"] = []string{content}

	//v := url.Values{}
	//v.Set("huifu", "hello world")
	//params := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码

	params :=  "text=" + content

	//jsonStr = "{'text': "+ content +"}"
	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	//req.Body
	//req.Body.Set("asdfsdf","13123");
	req.Header.Set("X-CurTime", curTime)
	req.Header.Set("X-Param", paramBase64)
	req.Header.Set("X-Appid", appId)
	req.Header.Set("X-CheckSum", checkSum)
	req.Header.Set("X-Real-Ip", myIP)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	//接收请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("-->")
	//fmt.Println("reqonse Headers:", req.Header)
	//fmt.Println("reqonse Body:", req.Body)
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return
}


func getIp() (ip string) {
	ip = "127.0.0.1"
	//addrs, err := net.InterfaceAddrs()
	//if err != nil {
	//	os.Stderr.WriteString("Oops:" + err.Error())
	//	os.Exit(1)
	//}
	//for _, a := range addrs {
	//	if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	//		if ipnet.IP.To4() != nil {
	//			os.Stdout.WriteString(ipnet.IP.String() + "\n")
	//		}
	//	}
	//}
	//os.Exit(0)
	return
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

//使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func writeMp3(path string, fileName string,data []byte)  (res bool, err error) {
	res = true
	tmpfn := filepath.Join(path, fileName)
	if err := ioutil.WriteFile(tmpfn, data, 0666); err != nil {
		res = false
		log.Fatal(err)
	}
	return
}

//遍历所有文件并取出txt
func getFileList(filePath string,fileExtension string) (listFile []string, err error) {
	//var (
	//	osType = os.Getenv("GOOS") // 获取系统类型
	//)
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
