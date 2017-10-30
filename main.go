/*
    This rewrite of Oberon is a command line utility for indexing courses and
    course-related data in Golang.
*/
package main

import "github.com/ahermida/Goberon/cmd"

//start up cli
func main() {
		cli := cmd.CLI{}
		cli.Run()
}
