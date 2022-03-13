# ipGeo
用于查询IP地理信息的GO库，同时支持IPv4和IPv6，使用纯真IPv4数据库和zxinc数据库，提供查询和下载功能，不产生多余依赖和输出。

## 使用

### 引入ipGeo
```go
import "github.com/HSwift/ipGeo"
```

### 下载数据库
```go
ipGeo.DownloadIPv4DB("./ipv4.db") # ipv4数据库
ipGeo.DownloadIPv6DB("./ipv6.db") # IPv6数据库
```

### 查询IP地址
IPv4
```go
ipdb,_ := ipGeo.OpenIPv4DB("./ipv4.db")
res,_ := ipdb.GetIPLocation("1.1.1.1")
println(res.Country,res.Area)
```
IPv6
```go
ipdb,_ := ipGeo.OpenIPv6DB("./ipv6.db")
res,_ := ipdb.GetIPLocation("2400::1")
println(res.Country,res.Area)
```

### 完整实例

参考 [main.go](/cmd/ipGeo/main.go)

## 注意

数据库格式错误会产生未定义的错误。

## 感谢

[纯真IP数据库](https://www.cz88.net/) 提供的IPv4数据

[zxinc数据库](http://ip.zxinc.org/) 提供的IPv6数据

https://github.com/zu1k/nali

https://github.com/freshcn/qqwry

https://gist.github.com/rmb122/ec7f305679ae9921a79b571d56390a74
