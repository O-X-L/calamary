package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/cnf/cnf_file"
	"github.com/superstes/calamary/log"
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
	var modeValidate bool
	flag.BoolVar(&modeValidate, "v", false, "Only validate config")
	flag.StringVar(&cnf.ConfigFileAbs, "f", cnf.ConfigFileAbs, "Path to the config file")
	flag.Parse()

	cnf.C = &cnf.Config{}

	if modeValidate {
		cnf.C.Service.Debug = true
		cnf_file.Load(true, true)
		log.Info("config", "Config validated successfully")
		os.Exit(0)
	}

	welcome()
	cnf_file.Load(false, true)
	service := &service{}
	go startPrometheusExporter()
	service.start()
	service.signalHandler()
}
