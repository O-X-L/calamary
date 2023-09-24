package main

import (
	"fmt"

	"github.com/superstes/calamary/cnf/cnf_file"
)

func welcome() {
	fmt.Printf("\n   ______      __                                \n")
	fmt.Println("  / ____/___ _/ /___ _____ ___  ____ ________  __")
	fmt.Println(" / /   / __ `/ / __ `/ __ `__ \\/ __ `/ ___/ / / /")
	fmt.Println("/ /___/ /_/ / / /_/ / / / / / / /_/ / /  / /_/ / ")
	fmt.Println("\\____/\\__,_/_/\\__,_/_/ /_/ /_/\\__,_/_/   \\__, /  ")
	fmt.Println("                                        /____/   ")
	fmt.Printf("by Superstes\n\n")
}

func main() {
	welcome()
	cnf_file.Load()
	service := &service{}
	_, cancel := service.start()
	defer service.shutdown(cancel)
	service.signalHandler()
}
