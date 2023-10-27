package models

import (
	"github.com/ip2location/ip2location-go/v9"
)

const (
	LiteDB1Path = "./assets/ip2location/IP2LOCATION-LITE-DB1.BIN"
	LiteDB3Path = "./assets/ip2location/IP2LOCATION-LITE-DB3.BIN"
)

// GetIp2Location 获取IP4地址
func GetIp2Location(ip4 string) (*ip2location.IP2Locationrecord, error) {
	db, err := ip2location.OpenDB(LiteDB3Path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	results, err := db.Get_all(ip4)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
