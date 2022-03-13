package ipGeo

import (
	"testing"
)

func TestDownloadIPv4DB(t *testing.T) {
	if err := DownloadIPv4DB("./qqwry.db"); err != nil {
		t.Error(err)
	}
}

func TestDownloadIPv6DB(t *testing.T) {
	if err := DownloadIPv6DB("./ipv6wry.db"); err != nil {
		t.Error(err)
	}
}

func TestIPv4(t *testing.T) {
	db, _ := OpenIPv4DB("./qqwry.db")
	res, _ := db.GetIPLocation("1.1.1.1")

	println(res.Country)
	println(res.Area)
}

func TestIPv6(t *testing.T) {
	db, _ := OpenIPv6DB("./ipv6wry.db")
	res, _ := db.GetIPLocation("2400::1")

	println(res.Country)
	println(res.Area)
}
