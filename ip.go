package gohack

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

// FakeIP 按照指定的location（省或市）生成随机IP
func FakeIP(location string) string {
	if location != "" {
		for k, v := range ChinaIP {
			if strings.Contains(k, location) {
				cidr := v[NewRand().Intn(len(v))]
				return RandInCidr(cidr)
			}
		}
	}
	for _, v := range ChinaIP {
		cidr := v[NewRand().Intn(len(v))]
		return RandInCidr(cidr)
	}
	return ""
}

// Int2IP 数字转换IP地址
func Int2IP(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// IP2Int IP地址转数字
func IP2Int(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

// RandInCidr 在CIDR地址范围内生成一个随机IP
func RandInCidr(cidr string) string {
	minIP, maxIP := Cidr2Range(cidr)
	return RandInRange(minIP, maxIP)
}

// RandInRange 在给定的IP地址范围内随机生成一个IP，minIP和maxIP可以是.分隔的IP表示，也可以是十进制数字IP表示
func RandInRange(minIP, maxIP string) string {
	if minIP == maxIP {
		return minIP
	}
	var ip1, ip2 int64
	if strings.Contains(minIP, ".") {
		ip1 = IP2Int(minIP)
	} else {
		ip1 = int64(IntVal(minIP))
	}
	if strings.Contains(maxIP, ".") {
		ip2 = IP2Int(maxIP)
	} else {
		ip2 = int64(IntVal(maxIP))
	}
	return Int2IP(ip1 + NewRand().Int63n(ip2-ip1))
}

// Cidr2Range 将cidr类型的地址转换成IP地址段表示
func Cidr2Range(cidr string) (string, string) {
	ipMask := strings.Split(cidr, "/")
	if ipMask[1] == "32" {
		return ipMask[0], ipMask[0]
	}
	ipSegs := strings.Split(ipMask[0], ".")
	maskLen, _ := strconv.Atoi(ipMask[1])
	seg3MinIP, seg3MaxIP := getIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, seg4MaxIP := getIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	return ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP),
		ipPrefix + strconv.Itoa(seg3MaxIP) + "." + strconv.Itoa(seg4MaxIP)
}

// GetCidrHostNum 计算得到CIDR地址范围内可拥有的主机数量
func GetCidrHostNum(maskLen int) uint {
	cidrIPNum := uint(0)
	i := uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIPNum += 1 << i
	}
	return cidrIPNum
}

//得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIPSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIP, _ := strconv.Atoi(ipSegs[2])
		return segIP, segIP
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return getIPSegRange(uint8(ipSeg), uint8(24-maskLen))
}

//得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIPSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIP, segMaxIP := getIPSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIP + 1, segMaxIP
}

//根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func getIPSegRange(userSegIP, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIP := ipSegMax << offset
	segMinIP := netSegIP & userSegIP
	segMaxIP := userSegIP&(255<<offset) | ^(255 << offset)
	return int(segMinIP), int(segMaxIP)
}
