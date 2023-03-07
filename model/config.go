package model

import (
	"github.com/akkuman/parseConfig"
	"github.com/wonderivan/logger"
	"os"
	"path/filepath"
)

type ConfigFile struct {
	FileName string
	FileDir  string
	FilePath string
}

type PveConfig struct {
	ApiUrl           string
	Token            string
	PrometheusEnable string
	JumpServerEnable string
}

type PrometheusConfig struct {
	Vendor    string
	LocalPath string
	Host      string
	UserName  string
	PassWord  string
	DstPath   string
	TimeOut   string
}

type JumpServerConfig struct {
	ApiUrl          string
	AccessKeyID     string
	AccessKeySecret string
	SSHPort         string
	RDPPort         string
	LinuxPlatform   string
	WinPlatform     string
}

var pveConfig *PveConfig
var prometheusConfig *PrometheusConfig
var jumpServerConfig *JumpServerConfig

func (t *ConfigFile) Init(fileName string, fileDir string) *ConfigFile {
	exePath, err := os.Executable()
	if err != nil {
		logger.Info(err)
	}
	path, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	filePath := path + "/../" + fileDir + "/" + fileName
	t.FileName = fileName
	t.FileDir = fileDir
	t.FilePath = filePath
	return t
}

func (t *PveConfig) Init() interface{} {
	configFile := new(ConfigFile)
	configFile.Init("pve.json", "configs")
	config := parseConfig.New(configFile.FilePath)
	t.ApiUrl = config.Get("pveConfig > apiUrl").(string)
	t.Token = config.Get("pveConfig > token").(string)
	t.PrometheusEnable = config.Get("pveConfig > prometheusEnable").(string)
	t.JumpServerEnable = config.Get("pveConfig > jumpserverEnable").(string)
	return t
}
func (t *PrometheusConfig) Init() interface{} {
	configFile := new(ConfigFile)
	configFile.Init("prometheus.json", "configs")
	config := parseConfig.New(configFile.FilePath)
	t.Vendor = config.Get("prometheusConfig > vendor").(string)
	t.Host = config.Get("prometheusConfig > host").(string)
	t.UserName = config.Get("prometheusConfig > userName").(string)
	t.PassWord = config.Get("prometheusConfig > passWord").(string)
	t.DstPath = config.Get("prometheusConfig > dstPath").(string)
	t.TimeOut = config.Get("prometheusConfig > timeOut").(string)
	return t
}

func (t *JumpServerConfig) Init() interface{} {
	configFile := new(ConfigFile)
	configFile.Init("jumpserver.json", "configs")
	config := parseConfig.New(configFile.FilePath)
	t.ApiUrl = config.Get("jumpServerConfig > apiUrl").(string)
	t.AccessKeyID = config.Get("jumpServerConfig > accessKeyID").(string)
	t.AccessKeySecret = config.Get("jumpServerConfig > accessKeySecret").(string)
	t.SSHPort = config.Get("jumpServerConfig > sshPort").(string)
	t.RDPPort = config.Get("jumpServerConfig > rdpPort").(string)
	t.LinuxPlatform = config.Get("jumpServerConfig > linuxPlatform").(string)
	t.WinPlatform = config.Get("jumpServerConfig > winPlatform").(string)
	return t
}

func InitConfig() {
	pveConfig = new(PveConfig)
	pveConfig.Init()
	prometheusConfig = new(PrometheusConfig)
	prometheusConfig.Init()
	jumpServerConfig = new(JumpServerConfig)
	jumpServerConfig.Init()
}
