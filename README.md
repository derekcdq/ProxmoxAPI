# 1.ProxmoxAPI项目介绍
本项目以ProxmoxVE为IaaS基础环境，将ProxmoxVE平台下的虚拟机信息自动同步到Prometheus平台以及Jumpserver平台（仅支持Jumpserver3.x版本），实现了VM虚拟机新增后的数据同步自动化
由于ProxmoxVE的接口无法直接取到虚拟机实际的IP，所以本项目的实现依赖于VM的标准化命名以及Proxmox里的资源池设置。

# 2.ProxmoxVE的资源池及虚拟机名设置规则
资源池命名规则：在数据中心->权限->资源池里设置好以部门为单位的资源池，例如，资源池名称为IT,备注为 企业IT部
虚拟机命名规则：IT-192.168.10.101

# 2.初始化配置
初次使用需要进行配置文件的配置，本项目所有配置文件都在configs下面

# 3.配置文件说明
## pve.json
```
{ 
  "pveConfig": {
  "apiUrl": "https://pvehost:8006/api2/json/",                    //ProxmoxVE服务器的API接口地址，具体可参照Proxmox官网
    "token": "PVEAPIToken=root@pam!pvetoken=1232132132131231321", //需要在数据中心->权限->API令牌中设置
    "prometheusEnable": "1",                                      //是否开启往prometheus同步配置文件，1为启用，0为禁用
    "jumpserverEnable": "1"                                       //是否开启往jumpserver同步主机信息，1为启动，0为禁用
  }
}
```

## prometheus.json
```
{
  "prometheusConfig": {
    "vendor": "公司名",                        //会在prometheus生成vendor=XX的标签，一般填写公司名或者是云厂商名称
    "host": "serverIP:ssh端口",                //prometheus服务器的IP以及ssh监听端口
    "userName": "root",                       //具有scp权限的ssh账号
    "passWord": "password",                   //ssh账号对应的密码
    "dstPath": "/opt/prometheus/conf.d/",     //scp超时时间
    "timeOut": "600"
  }
}
```

## log.json
```
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
```

## jumpserver.json
```
{
  "jumpServerConfig": {
    "apiUrl": "http://jumpserver.org",         //Jumpserver的hostname或IP地址
    "accessKeyID": "xxx",                      //Jumpserver的accessKeyID
    "accessKeySecret": "xxx",                  //Jumpserver的Secret
    "sshPort": "22",
    "rdpPort": "3389",
    "linuxPlatform": "29",                     //当判断虚拟机为linux操作系统的时候，使用jumpserver编号为29的platform配置
    "winPlatform": "5"                         //当判断虚拟机为windows操作系统的时候，使用jumpserver编号为5的platform配置
  }
}
```

# 4.运行和编译
## 运行
```
git clone https://github.com/derekcdq/ProxmoxAPI.git
cd /opt/PromoxAPI/
go run main

## 编译
  cd /opt/ProxmoxAPI
  go build -o bin/proxmoxApi
  ./bin/proxmoxApi
  
  
