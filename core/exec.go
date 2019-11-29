package core

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type StdHandler func(row string)

//执行nginx
func Exec(handler StdHandler) {
	cmd := exec.Command("nginx")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic("StdoutPipe Err " + err.Error())
	}
	if err := cmd.Start(); err != nil {
		panic("nginx Start Err " + err.Error())
	}

	fmt.Println("pid:", cmd.Process.Pid)
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			fmt.Println("ReadString Err " + err2.Error())
			break
		}
		handler(line)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait Err " + err.Error())
	}
	fmt.Println("nginx over shit")
}

//nginx reload
func Reload() error {
	return exec.Command("nginx", "-s", "reload").Run()
}

//nginx test config
func Test() error {
	return exec.Command("nginx", "-t").Run()
}
