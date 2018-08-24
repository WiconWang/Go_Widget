package main

import (
	"fmt"
	"time"
)

func main()  {

	fmt.Println("hello world")
	res := false
	res = time.Now().Unix() < 1533714306 + 30 * 24 * 60 * 60
	fmt.Println(res)
}

func getValid() (bool res){
	res = time.Now().Unix() < 1533714306 + 30 * 24 * 60 * 60
	return
}
