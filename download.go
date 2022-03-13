package ipGeo

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func readHTTP(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func DownloadIPv4DB(path string) error {
	copywrite, err := readHTTP("http://update.cz88.net/ip/copywrite.rar")
	if err != nil {
		return err
	}
	qqwry, err := readHTTP("http://update.cz88.net/ip/qqwry.rar")
	if err != nil {
		return err
	}
	key := binary.LittleEndian.Uint32(copywrite[5*4:])
	for i := 0; i < 0x200; i++ {
		key = key * 0x805
		key++
		key = key & 0xff
		qqwry[i] = byte(uint32(qqwry[i]) ^ key)
	}
	reader, err := zlib.NewReader(bytes.NewReader(qqwry))
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	if _, err = io.Copy(f, reader); err != nil {
		return err
	}
	return nil
}

func DownloadIPv6DB(path string) error {
	ipv6wry, err := readHTTP("https://ip.zxinc.org/ipv6wry.db")
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	f.Write(ipv6wry)
	return nil
}
