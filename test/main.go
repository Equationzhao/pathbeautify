package main

import (
	"fmt"

	"github.com/Equationzhao/pathbeautify"
)

func main() {
	fmt.Println(pathbeautify.Beautify("~/Downloads"))
	// output: /home/equationzhao/Downloads

	fmt.Println(pathbeautify.Beautify("../.../..../.."))
	// output: ../../../../../../..
}
