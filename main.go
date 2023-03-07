package main

import (
	"net/http"
	"proxmoxApi/model"
	"time"
)

func init() {
	model.InitLogConfig()
	model.InitConfig()
}

func main() {
	//sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//go func() {
	//	sig := <-sigs
	//	logger.Info(sig)
	//	done <- true
	//}()
	//
	//if os.Getppid() != 1 {
	//	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	//	cmd.Stdin = os.Stdin
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	//	err := cmd.Start()
	//	if err != nil {
	//		logger.Info(err)
	//	}
	//	os.Exit(0)
	//}

	go func() {
		for {
			model.CreateAndSyncJsonAll()
			time.Sleep(time.Duration(60) * time.Second)
		}
	}()

	go func() {
		for {
			model.SyncAllAssets()
			time.Sleep(time.Duration(60) * time.Second)
		}
	}()

	for {
		http.ListenAndServe("0.0.0.0:8888", nil)
	}

}
