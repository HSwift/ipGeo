package main

import (
	"flag"
	"fmt"
	"github.com/HSwift/ipGeo"
	"os"
)

func lookup(db ipGeo.IPDB, ip string) {
	res, err := db.GetIPLocation(ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Country, res.Area)
}

func openIpv6() *ipGeo.IPv6DB {
	if _, err := os.Stat("./ipv6wry.db"); err != nil {
		fmt.Println("下载IPv6数据库")
		err := ipGeo.DownloadIPv6DB("./ipv6wry.db")
		if err != nil {
			panic(err)
		}
	}
	db, err := ipGeo.OpenIPv6DB("./ipv6wry.db")
	if err != nil {
		panic(err)
	}
	return db
}

func openIpv4() *ipGeo.IPv4DB {
	if _, err := os.Stat("./qqwry.db"); err != nil {
		fmt.Println("下载IPv4数据库")
		err := ipGeo.DownloadIPv4DB("./qqwry.db")
		if err != nil {
			panic(err)
		}
	}
	db, err := ipGeo.OpenIPv4DB("./qqwry.db")
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	var enableIPv6 bool
	var ip string
	var db ipGeo.IPDB
	flag.BoolVar(&enableIPv6, "6", false, "enable ipv6")
	flag.Parse()
	ip = flag.Arg(0)
	if enableIPv6 {
		db = openIpv6()
	} else {
		db = openIpv4()
	}
	lookup(db, ip)
}
