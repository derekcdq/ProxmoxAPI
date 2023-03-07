package model

import (
	"encoding/json"
	"github.com/wonderivan/logger"
)

type PoolsList struct {
	Data []struct {
		PoolID  string `json:"poolid"`
		Comment string `json:"comment"`
	} `json:"data"`
}

type PoolConfig struct {
	Data struct {
		Comment string   `json:"comment"`
		Members []VMInfo `json:"members"`
	} `json:"data"`
}

type VMInfo struct {
	Node     string  `json:"node"`
	Vmid     int     `json:"vmid"`
	Disk     int     `json:"disk"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Mem      int     `json:"mem"`
	ID       string  `json:"id"`
	Template int     `json:"template"`
	CPU      float64 `json:"cpu"`
	OSType   string  `json:"ostype"`
	Uptime   int     `json:"uptime"`
}

type VMConfig struct {
	Data struct {
		OSType string
		Cpu    string
		Cores  int
		Memory int
	} `json:"data"`
}

func (t *PoolsList) Get() *PoolsList {
	method := "pools/"
	var p ProxmoxAPI
	body := p.Get(method)
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	return t
}

func (t *PoolConfig) Get(poolID string) *PoolConfig {
	method := "pools/" + poolID
	var p ProxmoxAPI
	body := p.Get(method)
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	return t
}

func (t *VMConfig) Get(node string, vmID string) *VMConfig {
	method := "nodes/" + node + "/qemu/" + vmID + "/config/"
	var p ProxmoxAPI
	body := p.Get(method)
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	return t
}
