package app

import (
	"bytes"
	"fmt"
	"os/exec"
)

type StdHandler func(row string)

func Exec(handler StdHandler) {
	//执行子进程
	cmd := exec.Command("sh", "abc.sh")
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	cmd.Run()

	for {
		line, err := stdOut.ReadString(byte('\n'))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(line)
	}

}
