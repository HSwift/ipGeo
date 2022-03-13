package ipGeo

import (
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"os"
)

func OpenIPv4DB(path string) (*IPv4DB, error) {
	db := &IPv4DB{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	db.data = data
	db.IndexStart = binary.LittleEndian.Uint32(data[0:4])
	end := binary.LittleEndian.Uint32(data[4:8])
	db.ipSize = 4
	db.refSize = 3
	db.indexSize = db.ipSize + db.refSize
	db.RecordsCount = (end - db.IndexStart) / (db.ipSize + db.refSize)
	return db, nil
}

func (db *IPv4DB) findIPRecord(ip uint32) uint32 {
	var start uint32 = 0
	var end uint32 = db.RecordsCount
	for end-start > 1 {
		mid := (start + end) / 2
		midIndex := db.IndexStart + mid*db.indexSize
		midIP := binary.LittleEndian.Uint32(db.data[midIndex : midIndex+db.ipSize])
		if ip < midIP {
			end = mid
		} else {
			start = mid
		}
	}
	resIndex := db.IndexStart + start*db.indexSize
	return bytesToUint32(db.data[resIndex+db.ipSize : resIndex+db.indexSize])
}

func (db *IPv4DB) readRecord(recordIndex uint32) *Result {
	var nextArea uint32
	head := db.data[recordIndex]
	result := &Result{}
	enc := simplifiedchinese.GBK.NewDecoder()
	if head == 1 {
		next := bytesToUint32(db.data[recordIndex+1 : recordIndex+1+db.refSize])
		return db.readRecord(next)
	} else {
		result.Country = db.readIndirectString(recordIndex)
		if head == 2 {
			nextArea = recordIndex + 1 + db.refSize
		} else {
			nextArea = recordIndex + uint32(len(result.Country)) + 1
		}
		result.Area = db.readIndirectString(nextArea)
	}
	result.Country, _ = enc.String(result.Country)
	result.Area, _ = enc.String(result.Area)
	return result
}

func (db *IPv4DB) readIndirectString(current uint32) string {
	flag := db.data[current]
	if flag == 1 || flag == 2 {
		next := bytesToUint32(db.data[current+1 : current+1+db.refSize])
		return db.readString(next)
	} else {
		return db.readString(current)
	}
}

func (db *IPv4DB) readString(stringIndex uint32) string {
	endIndex := stringIndex
	for db.data[endIndex] != 0 {
		endIndex += 1
	}
	return string(db.data[stringIndex:endIndex])
}

func (db *IPv4DB) GetIPLocation(ip string) (*Result, error) {
	var ipBytes []byte
	if ipBytes = net.ParseIP(ip).To4(); ipBytes == nil {
		return nil, errors.New("malformed ip address")
	}
	return db.GetIPNumLocation(ipBytes)
}

func (db *IPv4DB) GetIPNumLocation(ipBytes net.IP) (*Result, error) {
	ipNum := binary.BigEndian.Uint32(ipBytes)
	recordIndex := db.findIPRecord(ipNum)
	res := db.readRecord(recordIndex + 4)
	res.IP = ipBytes
	return res, nil
}

func bytesToUint32(input []byte) uint32 {
	var res = []byte{0, 0, 0, 0}
	for i := 0; i < len(input); i++ {
		res[i] = input[i]
	}
	return binary.LittleEndian.Uint32(res)
}
