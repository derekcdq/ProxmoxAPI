package model

import (
	"context"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/wonderivan/logger"
	"golang.org/x/crypto/ssh"
	"os"
	"strconv"
	"time"
)

func ScpFile(filePath string, fileName string) error {
	clientConfig, _ := auth.PasswordKey(
		prometheusConfig.UserName, prometheusConfig.PassWord, ssh.InsecureIgnoreHostKey())
	timeOut, _ := strconv.Atoi(prometheusConfig.TimeOut)
	client := scp.NewClientWithTimeout(
		prometheusConfig.Host, &clientConfig, time.Duration(timeOut)*time.Second)
	err := client.Connect()
	if err != nil {
		logger.Info(err)
		return err
	}
	f, _ := os.Open(filePath)
	defer client.Close()
	defer f.Close()
	dstPath := prometheusConfig.DstPath + fileName
	err = client.CopyFromFile(context.Background(), *f, dstPath, "0655")
	if err != nil {
		logger.Info(err)
	}
	return err
}
