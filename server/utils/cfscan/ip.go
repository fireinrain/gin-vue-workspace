package cfscan

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

//
// CountIPs
//  @Description: 获取CIDR代表的IP数量
//  @param cidr
//  @return int
//

// 计算CIDR表示的IP数量
func CountIPs(cidr string) int {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	ones, bits := ipnet.Mask.Size()
	return 1 << (bits - ones)
}

// splitCIDRs
//
//	@Description: 按照指定的ip批次大小，切分CIDR数据
//	@param cidrs
//	@param batchSize
//	@return [][]string
//
// 分割CIDR列表
func SplitCIDRs(cidrs []string, batchSize int) [][]string {
	type cidrWithCount struct {
		cidr  string
		count int
	}

	// 计算每个CIDR的IP数量
	cidrsWithCounts := make([]cidrWithCount, len(cidrs))
	for i, cidr := range cidrs {
		cidrsWithCounts[i] = cidrWithCount{cidr, CountIPs(cidr)}
	}

	// 按IP数量排序
	sort.Slice(cidrsWithCounts, func(i, j int) bool {
		return cidrsWithCounts[i].count > cidrsWithCounts[j].count
	})

	var result [][]string
	currentBatch := []string{}
	currentBatchSize := 0

	for _, cidr := range cidrsWithCounts {
		if currentBatchSize+cidr.count > batchSize && len(currentBatch) > 0 {
			result = append(result, currentBatch)
			currentBatch = []string{}
			currentBatchSize = 0
		}

		currentBatch = append(currentBatch, cidr.cidr)
		currentBatchSize += cidr.count
	}

	// 添加最后一个批量
	if len(currentBatch) > 0 {
		result = append(result, currentBatch)
	}

	// 如果拆分结果的批次数量超过25，忽略batchSize，按照总数不超过25的原则重新拆分
	if len(result) > 25 {
		result = SplitIntoMaxBatches(cidrs, 25)
	}

	return result
}

// splitIntoMaxBatches
//
//	@Description: 将CIDR列表尽可能拆分成总数不超过指定数量的批次
//	@param cidrs
//	@param maxBatches
//	@return [][]string
func SplitIntoMaxBatches(cidrs []string, maxBatches int) [][]string {
	var result [][]string
	batchSize := (len(cidrs) + maxBatches - 1) / maxBatches // 计算每个批次的最大CIDR数量

	for i := 0; i < len(cidrs); i += batchSize {
		end := i + batchSize
		if end > len(cidrs) {
			end = len(cidrs)
		}
		result = append(result, cidrs[i:end])
	}

	return result
}

func GetCIDRByASN(asnName string) ([]string, error) {
	asn := asnName // 替换为实际的ASN值
	asn = strings.ReplaceAll(asn, "AS", "")
	url := fmt.Sprintf("https://whois.ipip.net/AS%s", asn)
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get data: %s\n", resp.Status)
		return nil, err
	}

	scanner := bufio.NewScanner(resp.Body)
	re := regexp.MustCompile(`<a href="/AS` + asn + `/[^"]*"[^>]*>([^<]*)</a>`)
	output := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			ipRange := matches[1]
			if !strings.Contains(ipRange, ":") { // 排除IPv6地址
				output = append(output, ipRange)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return output, nil
}

func GetCIDRByASN2File(asnName string, filePath string) {
	byASNData, err := GetCIDRByASN(asnName)
	if err != nil {
		panic(err)
	}
	outFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	for _, line := range byASNData {
		_, err := outFile.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Data written to %s\n", filePath)
}

// ExtractIPv4Addresses
//
//	@Description: 从字符串中抽取ip
//	@param input
//	@return []string
func ExtractIPv4Addresses(input string) []string {
	// IPv4 地址的正则表达式模式
	// 这个模式匹配 0-255 范围内的数字，由点分隔，总共四组
	pattern := `\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 查找所有匹配项
	matches := re.FindAllString(input, -1)

	return matches
}

func ExtractIPv4CIDRAddresses(input string) []string {
	// IPv4 CIDR 地址的正则表达式模式
	// 这个模式匹配 IPv4 地址，后跟一个斜杠和 0-32 之间的数字
	pattern := `\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(?:[0-9]|[12][0-9]|3[0-2])\b`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 查找所有匹配项
	matches := re.FindAllString(input, -1)

	return matches
}
