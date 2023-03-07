package model

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/wonderivan/logger"
	"gopkg.in/twindagger/httpsig.v1"
	"net/http"
	"strconv"
	"strings"
)

type SigAuth struct {
	AccessKeyID     string
	AccessKeySecret string
}

type JNode struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	OrgID     string `json:"org_id"`
	Name      string `json:"name"`
	FullValue string `json:"full_value"`
	OrgName   string `json:"org_name"`
}

type Asset struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Platform  string   `json:"platform"`
	Nodes     []string `json:"nodes"`
	Protocols []Protocol
	IsActive  bool `json:"is_active"`
}

type Protocol struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type AssetsList []Asset

var jNodeMap map[string]map[string]string

func (t *SigAuth) Sign(r *http.Request) error {
	headers := []string{"(request-target)", "date"}
	signer, err := httpsig.NewRequestSigner(
		jumpServerConfig.AccessKeyID, jumpServerConfig.AccessKeySecret, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(r, headers, nil)
}

func (t *JNode) GetJNodeMap() map[string]map[string]string {
	method := "/api/v1/assets/nodes/"
	var j JumpServerAPI
	body := j.Get(method)
	var jnodeList []JNode
	err := json.Unmarshal(body, &jnodeList)
	if err != nil {
		logger.Info(err)
	}
	m := make(map[string]map[string]string)
	for _, v := range jnodeList {
		m[v.Value] = make(map[string]string)
		m[v.Value]["id"] = v.ID
	}
	return m
}

func (t *JNode) Create() error {
	method := "/api/v1/assets/nodes/"
	var j JumpServerAPI
	buf, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		logger.Info(err)
	}
	statusCode, body := j.Post(method, string(buf))
	if statusCode == 201 {
		logger.Info("新节点创建成功.")
	} else {
		logger.Info("节点创建失败:" + string(body))
		err = fmt.Errorf("节点创建失败:" + string(body))
	}
	return err
}

func (t Asset) Create() error {
	method := "/api/v1/assets/hosts/"
	var j JumpServerAPI
	assetID, _ := uuid.NewV4()
	t.ID = assetID.String()
	buf, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		logger.Info(err)
	}
	statusCode, body := j.Post(method, string(buf))
	if statusCode == 201 {
		logger.Info("新资产创建成功.")
	} else {
		logger.Info("资产创建失败:" + string(body))
		err = fmt.Errorf("资产创建失败:" + string(body))
	}
	return err
}

func (t Asset) GetFromSinglePool(poolID string) AssetsList {
	poolConfig := new(PoolConfig).Get(poolID)
	var assetsList []Asset
	for _, v := range poolConfig.Data.Members {
		if v.Status != "running" {
			continue
		}
		t.Name = v.Name
		vmConfig := new(VMConfig).Get(v.Node, strconv.Itoa(v.Vmid))
		array := strings.Split(t.Name, "-")
		t.Address = array[1]
		var p Protocol
		switch vmConfig.Data.OSType {
		case "l24", "l26", "solaris":
			t.Platform = jumpServerConfig.LinuxPlatform
			port, _ := strconv.Atoi(jumpServerConfig.SSHPort)
			p = Protocol{
				Name: "ssh",
				Port: port,
			}
		case "wxp", "w2k", "w2k3", "w2k8", "wvista", "win7", "win8", "win10", "win11":
			t.Platform = jumpServerConfig.WinPlatform
			port, _ := strconv.Atoi(jumpServerConfig.RDPPort)
			p = Protocol{
				Name: "rdp",
				Port: port,
			}
		default:
			t.Platform = "2"
			p = Protocol{
				Name: "ssh",
				Port: 22,
			}
		}
		t.Protocols = []Protocol{p}
		t.IsActive = true
		t.Nodes = []string{jNodeMap[poolConfig.Data.Comment]["id"]}
		assetsList = append(assetsList, t)
	}
	return assetsList
}

func InitJNodeMap() {
	jNode := new(JNode)
	jNodeMap = jNode.GetJNodeMap()
}

func SyncAllAssets() {
	if pveConfig.JumpServerEnable != "1" {
		return
	}
	poolList := new(PoolsList).Get()
	for _, v := range poolList.Data {
		jNode := new(JNode)
		jID, _ := uuid.NewV4()
		jNode.ID = jID.String()
		jNode.Value = v.Comment
		jNode.FullValue = "/Default/" + v.Comment
		err := jNode.Create()
		if err != nil {
			logger.Info(err)
		}
	}
	InitJNodeMap()
	for _, v := range poolList.Data {
		assetsList := new(Asset).GetFromSinglePool(v.PoolID)
		for _, v := range assetsList {
			err := v.Create()
			if err != nil {
				logger.Info(err)
			}
		}
	}
}
