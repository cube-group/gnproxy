package main

import (
	"fmt"
	"strings"
)

func main(){
	fmt.Println(strings.Index("abc","ab"))
	fmt.Println("abc"[len("abc"):])
}
