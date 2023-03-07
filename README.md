# 1.ProxmoxAPI项目介绍
本项目以ProxmoxVE为iaaS基础环境，将ProxmoxVE平台下的虚拟机信息自动同步到Prometheus平台以及Jumpserver平台（仅支持Jumpserver3.x版本），实现了VM虚拟机新增后的数据同步自动化

# 2.初始化配置
初次使用需要进行配置文件的配置，本项目所有配置文件都在configs下面

# 3.配置文件说明
## pve.json
{ <br>
  "pveConfig": {<br>
    "apiUrl": "https://pvehost:8006/api2/json/",                       //ProxmoxVE服务器的API接口地址，具体可参照Proxmox官网。<br>
    "token": "PVEAPIToken=root@pam!pvetoken=1232132132131231321",      //需要在数据中心->权限->API令牌中设置<br>
    "prometheusEnable": "1",                                           //是否开启往prometheus同步配置文件，1为启用，0为禁用<br>
    "jumpserverEnable": "1"                                            //是否开启往jumpserver同步主机信息，1为启动，0为禁用<br>
  }<br>
}<br>

## prometheus.json
{
  "prometheusConfig": {
    "vendor": "公司名",
    "host": "serverIP:ssh端口",
    "userName": "root",
    "passWord": "password",
    "dstPath": "/opt/prometheus/conf.d/",
    "timeOut": "600"
  }
}

## log.json
{
  "TimeFormat":"2006-01-02 15:04:05",
  "Console": {
    "level": "TRAC",
    "color": true
  },
  "File": {
    "filename": "../logs/promoxapi.log",
    "level": "TRAC",
    "daily": true,
    "maxlines": 1000000,
    "maxsize": 1,
    "maxdays": -1,
    "append": true,
    "permit": "0660"
  },
  "Conn": {
    "net":"tcp",
    "addr":"10.1.55.10:1024",
    "level": "Warn",
    "reconnect":true,
    "reconnectOnMsg":false
  }
}

## jumpserver.json
{
  "jumpServerConfig": {
    "apiUrl": "http://jumpserver.org",
    "accessKeyID": "xxx",
    "accessKeySecret": "xxx",
    "sshPort": "22",
    "rdpPort": "3389"
  }
}
