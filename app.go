package main

import (
	"app/core"
	"fmt"
)

func main(){
	fmt.Println("hello gnproxy :)")

	app.Exec(func(row string) {
		fmt.Println(row)
	})
}


