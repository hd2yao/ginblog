package utils

import (
    "fmt"
    "gopkg.in/ini.v1"
)

var (
    AppMode  string
    HttpPort string
    JwtKey   string

    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string

    Zone        int
    AccessKey   string
    SecretKey   string
    Bucket      string
    QiniuServer string
)

func init() {
    file, err := ini.Load("config/config.ini")
    if err != nil {
        fmt.Println("配置文件读取错误，请检查文件路径:", err)
    }
    LoadServer(file)
    LoadData(file)
    LoadQiniu(file)
}

func LoadServer(file *ini.File) {
    AppMode = file.Section("server").Key("AppMode").MustString("debug")
    HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
    JwtKey = file.Section("server").Key("JwtKey").MustString("")
}

func LoadData(file *ini.File) {
    DBHost = file.Section("database").Key("DbHost").MustString("localhost")
    DBPort = file.Section("database").Key("DbPort").MustString("3306")
    DBUser = file.Section("database").Key("DbUser").MustString("ginblog")
    DBPassword = file.Section("database").Key("DbPassWord").String()
    DBName = file.Section("database").Key("DbName").MustString("ginblog")
}

func LoadQiniu(file *ini.File) {
    Zone = file.Section("qiniu").Key("Zone").MustInt(1)
    AccessKey = file.Section("qiniu").Key("AccessKey").String()
    SecretKey = file.Section("qiniu").Key("SecretKey").String()
    Bucket = file.Section("qiniu").Key("Bucket").String()
    QiniuServer = file.Section("qiniu").Key("QiniuServer").String()
}
