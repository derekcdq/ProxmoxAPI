package model

import (
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NodeFile struct {
	Targets []string `json:"targets"`
	Labels  struct {
		Vendor     string `json:"Vendor"`
		Department string `json:"Department"`
	} `json:"labels"`
}

func (t NodeFile) Create(poolID string) (error, string) {
	poolConfig := new(PoolConfig).Get(poolID)
	var targets []string
	var instance string
	if len(poolConfig.Data.Members) == 0 {
		err := fmt.Errorf("no host in the pool:" + poolConfig.Data.Comment)
		return err, "no host in the pool:" + poolConfig.Data.Comment
	}
	for _, v := range poolConfig.Data.Members {
		if v.Status == "running" {
			array := strings.Split(v.Name, "-")
			vmConfig := new(VMConfig).Get(v.Node, strconv.Itoa(v.Vmid))
			if vmConfig.Data.OSType == "l26" {
				instance = array[1] + ":9100"
			} else {
				instance = array[1] + ":9182"
			}
			targets = append(targets, instance)
		}
	}
	t.Targets = targets
	t.Labels.Vendor = prometheusConfig.Vendor
	t.Labels.Department = poolConfig.Data.Comment
	a := []NodeFile{t}
	buf, err := json.MarshalIndent(a, "", "\t")
	//logger.Info(string(buf))
	fileName := poolID + "_node.json"
	exePath, err := os.Executable()
	if err != nil {
		logger.Info(err)
	}
	path, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	filePath := path + "/../prometheusconf/" + fileName
	logger.Info(filePath)
	logger.Info("正在生成json文件,部门：" + poolConfig.Data.Comment)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		logger.Info(err)
	}
	_, err = file.Write(buf)
	if err != nil {
		logger.Info(err)
	}
	_, err = os.Stat(filePath)
	if err != nil {
		logger.Info(err)
	} else {
		logger.Info("文件生成成功,主机数量:" + strconv.Itoa(len(t.Targets)) + "台")
	}
	err = file.Close()
	if err != nil {
		logger.Info(err)
	}
	return err, filePath
}

// CreateAndSyncJsonAll 为所有部门创建prometheus格式的json文件
func CreateAndSyncJsonAll() {
	if pveConfig.PrometheusEnable != "1" {
		return
	}
	poolList := new(PoolsList).Get()
	for _, v := range poolList.Data {
		err, filePath := new(NodeFile).Create(v.PoolID)
		if err != nil {
			logger.Info(err)
		} else {
			fileName := v.PoolID + "_node.json"
			logger.Info("开始传输文件:" + filePath)
			err = ScpFile(filePath, fileName)
			if err != nil {
				logger.Info(err)
			} else {
				logger.Info("文件传输成功.")
			}
		}
	}
}
