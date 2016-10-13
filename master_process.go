package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	i "puppet_monitoring/impl"
	"puppet_monitoring/rpc"
	"runtime"
	"strconv"
	"syscall"
)

// global variable (load once)
var settings = i.Settings{}.LoadSettings()

// runs master process
func run_master_process() {

	log.Printf("PID:%v\n", os.Getpid())

	// limit CPU usage
	runtime.GOMAXPROCS(1)

	// check for existing pid file
	if check_pid_file() {
		fmt.Println(settings.PidFile + " already exists! (other instance run?)")
		os.Exit(1)
	}

	create_pid_file()
	defer kill_pid()

	// creating server socket
	var laddr, err = net.ResolveTCPAddr("tcp", settings.Ip+":"+strconv.Itoa(settings.Port))
	ln, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		panic(err)
	}
	defer ln.Close()

	log.Println("listening on", ln.Addr())

	// creating service to handle server tcp socket
	service := i.Service{}.NewService()

	// creating puppet environment collection
	envs := i.EnvironmentCollection{}.NewEnvironmentCollection()

	envs.Conf = &settings

	// setup service params
	service.SetEnvCollection(&envs)

	// and run as go routine
	go service.HandleListener(ln)

	// creating and starting the rpc server to handle commands from outside
	rpcsrv := rpc.RPCServer{Envs: &envs}
	rpcsrv.CreateServer(settings)

	// handle SIGINT and SIGTERM
	ch := make(chan os.Signal)
	var sig os.Signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// awaiting signals
	select {
	case sig = <-ch:
		log.Println(sig)
	}

	// tell service to stop
	service.Stop()
}

// Check if pid file exists
func check_pid_file() bool {
	var _, err = os.Stat(settings.PidFile)
	return err == nil
}

// Create pid file
func create_pid_file() {

	// create file
	var fd, err = os.Create(settings.PidFile)
	fd.Close()
	if err != nil {
		fmt.Println("Error creating pid file!")
		panic(err)
	}
	// set rw-r--r--
	os.Chmod(settings.PidFile, 0644)
	fd, err = os.OpenFile(settings.PidFile, os.O_RDWR, 0644)
	defer fd.Close()
	// writing current process id
	var _, werr = fd.WriteString(strconv.Itoa(os.Getpid()))
	if werr != nil {
		fmt.Println("Error write pid file!")
		panic(werr)
	}
	fd.Sync()

}

// Remove pid file
func kill_pid() {
	err := os.Remove(settings.PidFile)
	if err != nil {
		log.Println(err)
	}
}
