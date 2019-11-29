package main

import "io/ioutil"

func main(){
	ioutil.WriteFile("abc",[]byte("asdfad"),0777)
}