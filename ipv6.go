package ipGeo

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
)

func OpenIPv6DB(path string) (*IPv6DB, error) {
	db := &IPv6DB{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	db.data = data
	db.IndexStart = binary.LittleEndian.Uint64(data[16:24])
	db.ipSize = uint64(data[7])
	db.refSize = uint64(data[6])
	db.indexSize = db.ipSize + db.refSize
	db.RecordsCount = binary.LittleEndian.Uint64(data[8:16])
	return db, nil
}

func (db *IPv6DB) findIPRecord(ip uint64) uint64 {
	var start uint64 = 0
	var end uint64 = db.RecordsCount
	for end-start > 1 {
		mid := (start + end) / 2
		midIndex := db.IndexStart + mid*db.indexSize
		midIP := binary.LittleEndian.Uint64(db.data[midIndex : midIndex+db.ipSize])
		if ip < midIP {
			end = mid
		} else {
			start = mid
		}
	}
	resIndex := db.IndexStart + start*db.indexSize
	return bytesToUint64(db.data[resIndex+db.ipSize : resIndex+db.indexSize])
}

func (db *IPv6DB) readRecord(recordIndex uint64) *Result {
	var nextArea uint64
	head := db.data[recordIndex]
	result := &Result{}
	if head == 1 {
		next := bytesToUint64(db.data[recordIndex+1 : recordIndex+1+db.refSize])
		return db.readRecord(next)
	} else {
		result.Country = db.readIndirectString(recordIndex)
		if head == 2 {
			nextArea = recordIndex + 1 + db.refSize
		} else {
			nextArea = recordIndex + uint64(len(result.Country)) + 1
		}
		result.Area = db.readIndirectString(nextArea)
	}
	return result
}

func (db *IPv6DB) readIndirectString(current uint64) string {
	flag := db.data[current]
	if flag == 1 || flag == 2 {
		next := bytesToUint64(db.data[current+1 : current+1+db.refSize])
		return db.readString(next)
	} else {
		return db.readString(current)
	}
}

func (db *IPv6DB) readString(stringIndex uint64) string {
	endIndex := stringIndex
	for db.data[endIndex] != 0 {
		endIndex += 1
	}
	return string(db.data[stringIndex:endIndex])
}

func (db *IPv6DB) GetIPLocation(ip string) (*Result, error) {
	var ipBytes []byte
	if ipBytes = net.ParseIP(ip).To16(); ipBytes == nil {
		return nil, errors.New("malformed ip address")
	}
	return db.GetIPNumLocation(ipBytes)
}

func (db *IPv6DB) GetIPNumLocation(ipBytes net.IP) (*Result, error) {
	ipNum := binary.BigEndian.Uint64(ipBytes[:8])
	recordIndex := db.findIPRecord(ipNum)
	res := db.readRecord(recordIndex)
	res.IP = ipBytes
	return res, nil
}

func bytesToUint64(input []byte) uint64 {
	var res = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(input); i++ {
		res[i] = input[i]
	}
	return binary.LittleEndian.Uint64(res)
}
