package main

import (
	"sync"
	"fmt"
	"time"

	"./webservice"
	"github.com/kardianos/service"
)

// -- install (not support VerifyPin)
// $ sc.exe create "GoLinkage" binpath= "c:\GoLinkage.exe" DisplayName= "GoLinkage" start= auto
// $ sc.exe description "GoLinkage" "Web service for smartcard connection and LinkageCenter connection"

// -- uninstall
// $ sc.exe delete "GoLinkage"

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync    sync.Mutex
)

const serviceName = "GoLinkage Service"
const serviceDescription = "Web service for smartcard connection and LinkageCenter connection"

// ----------------------- Service Program -----------------------------
type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Printf("[Service %v] Started\n", s.String())
	writingSync.Lock()
   	serviceIsRunning = true
   	writingSync.Unlock()
	go p.run()
	return nil
}
 
func (p program) Stop(s service.Service) error {
	writingSync.Lock()
	serviceIsRunning = false
	writingSync.Unlock()
	for programIsRunning {
	   fmt.Printf("[Service %v] Stopping...\n", s.String())
	   time.Sleep(1 * time.Second)
	}
	fmt.Printf("[Service %v] Stopped\n", s.String())
	return nil
}

func (p program) run() {
	webservice.Run()
}

// -------------------------------------------------------------

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	 }
	 prg := &program{}
	 s, err := service.New(prg, serviceConfig)
	 if err != nil {
		fmt.Printf("[Service] Cannot create the service: %v\n", err.Error())
	 }
	 err = s.Run()
	 if err != nil {
		fmt.Printf("[Service %v] Cannot start the service: %v\n", s.String(), err.Error())
	 }
}
