package cfscan

import (
	"fmt"
	"testing"
)

func TestGetCIDRByASN(t *testing.T) {
	asn, err := GetCIDRByASN("AS906")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(asn)
}

func TestGetCIDRByASN2File(t *testing.T) {
	GetCIDRByASN2File("AS906", "asn906.txt")
}
