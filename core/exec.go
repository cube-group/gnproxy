package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type StdHandler func(row string)

func Exec(handler StdHandler) {
	cmd := exec.Command("nginx", "-g", "'daemon off;'")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error=>", err.Error())
		return
	}
	cmd.Start() // Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		handler(line)
	}

	cmd.Wait()

}
