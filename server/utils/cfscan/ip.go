package cfscan

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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

// GetCIDRByASN
//
//	@Description: 通过ASN name获取cidr字符串列表
//	@param asnName
//	@return []string
//	@return error
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

// GetCIDRByASN2File
//
//	@Description: 用asn编号获取cidr并写入文件
//	@param asnName
//	@param filePath
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
//	@Description: 从字符串中抽取ipv4 列表
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

// ExtractIPv4CIDRAddresses
//
//	@Description: 从CIDR字符串中抽取ipv4 CIDR列表
//	@param input
//	@return []string
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

type CFIPTestResult struct {
	IP         string `json:"ip"`          // IP地址
	Port       int    `json:"port"`        // 端口
	DataCenter string `json:"data_center"` // 数据中心
	Region     string `json:"region"`      // 地区
	City       string `json:"city"`        // 城市
	//CF下载延迟
	Latency string `json:"latency"` // 延迟
	//CF代理判断延迟
	TcpDuration   time.Duration `json:"tcp_duration"`   // TCP请求延迟
	EnableTLS     bool          `json:"enable_tls"`     //是否开启TLS
	DownloadSpeed string        `json:"download_speed"` // 下载速度
}

// CleanResultJson
//
//	@Description: 清理测试结果json字符串 包含去重
//	@param taskResults
//	@return string
//	@return error
func CleanResultJson(taskResults []string) (string, error) {
	//去除重复的IP
	seenIPs := make(map[string]bool)
	var uniqueRecords []CFIPTestResult
	for _, result := range taskResults {
		var cfResult []CFIPTestResult
		err := json.Unmarshal([]byte(result), &cfResult)
		if err != nil {
			log.Printf("Error on unmarshall value: %s\n", result)
			continue
		}
		//去重标准是 ip+端口
		for _, r := range cfResult {
			port := strconv.Itoa(r.Port)
			if !seenIPs[r.IP+port] {
				seenIPs[r.IP+port] = true
				uniqueRecords = append(uniqueRecords, r)
			}
		}
	}
	//转化为json列表
	marshal, err := json.Marshal(uniqueRecords)
	if err != nil {
		log.Printf("Error on marshall value: %v\n", uniqueRecords)
		return "", err
	}
	return string(marshal), nil

}

const (
	earthRadiusKm = 6371 // 地球平均半径（公里）
	//深圳福田 北纬 东经
	localGeo = "22.5455,114.0683" //使用者本地GEO
)

// 将角度转换为弧度
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// CaculateGeoIPDist
//
//	@Description: 使用经纬度计算两地之间的物理距离
//	@param sourceGeoStr
//	@param targetGeoStr
//	@return float64
func CaculateGeoIPDist(sourceGeoStr string, targetGeoStr string) float64 {
	geo1 := strings.Split(sourceGeoStr, ",")
	lat1, _ := strconv.ParseFloat(strings.TrimSpace(geo1[0]), 64)
	lon1, _ := strconv.ParseFloat(strings.TrimSpace(geo1[1]), 64)

	geo2 := strings.Split(targetGeoStr, ",")
	lat2, _ := strconv.ParseFloat(strings.TrimSpace(geo2[0]), 64)
	lon2, _ := strconv.ParseFloat(strings.TrimSpace(geo2[1]), 64)

	lat1 = toRadians(lat1)
	lon1 = toRadians(lon1)
	lat2 = toRadians(lat2)
	lon2 = toRadians(lon2)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKm * c
	return distance
}

//
// CaculateGeoIPDistFromLocal
//  @Description: 计算ip到本地的物理距离
//  @param targetGeoStr
//  @return float64
//

func CaculateGeoIPDistFromLocal(targetGeoStr string) float64 {
	return CaculateGeoIPDist(localGeo, targetGeoStr)
}

// ipinfo.io/[IP address]?token=958278ec91a07c
const IPInfoIOToken = "958278ec91a07c"

// node info
// 启动时 自动获取hostname
// 并且请求ipinfo.io 获取当前ip geo的一些信息
type IPGeoInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func FetchIPGeoInfoByIP(ipStr string) (*IPGeoInfo, error) {
	// 发出 HTTP 请求
	resp, err := http.Get("https://ipinfo.io/" + ipStr + "?token=" + IPInfoIOToken)
	if err != nil {
		fmt.Println("Error fetch ip geo info:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应数据
	var ipInfo IPGeoInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		fmt.Println("Error on decode json to struct:", err)
		return nil, err
	}
	return &ipInfo, nil
}
