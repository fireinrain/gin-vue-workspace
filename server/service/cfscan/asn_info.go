package cfscan

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AsnInfoService struct{}

// CreateAsnInfo 创建asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) CreateAsnInfo(asnInfo *cfscan.AsnInfo) (err error) {
	enable := asnInfo.Enable
	if asnInfo.AsnName == "" || *enable == 0 {
		return errors.New("ASN name cant be empty or enable cant be 0")
	}
	err = global.GVA_DB.Create(asnInfo).Error
	return err
}

func (asnInfoService *AsnInfoService) GetASNDetailByASN(asnInfo *cfscan.AsnInfo) *cfscan.AsnInfo {
	result := asnInfo
	asnName := asnInfo.AsnName
	asnName = strings.TrimSpace(asnName)
	asnNum := strings.ReplaceAll(asnName, "AS", "")
	// URL to fetch
	queryUrl := fmt.Sprintf("https://bgpview.io/asn/%s", asnNum)

	// Create a new request
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("Error creating request:", zap.Error(err))
		return nil
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	// Create a client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("Error sending request:", zap.Error(err))
		return nil
	}

	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		global.GVA_LOG.Error("Response 404 not found:", zap.Error(err))
		return nil
	}

	// Read the response body
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	global.GVA_LOG.Error("Error reading response body:", zap.Error(err))
	//	return res
	//}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTML:", err)
	}

	// CSS selector for the element you want to extract
	cssSelector := "body > div.container.main > div > div:nth-child(1) > div > div > div"

	// Find the first element that matches the CSS selector
	fullName := ""
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		fullName = strings.Join(strings.Fields(textData), " ")
		// You can also extract other attributes if needed, for example:
		// attr, exists := s.Attr("href")
		// fmt.Println(attr, exists)
		result.FullName = fullName
	})
	//判断 如果fullName中存在 Unallocated， 则说明ASNxxx不存在 直接返回
	//ipv4Counts
	headerCSS := "body > div.container.main > div > div:nth-child(2) > div > div"
	doc.Find(headerCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		cTextData := strings.Split(textData, "\n")
		var dataList []string
		for _, datum := range cTextData {
			datum = strings.TrimSpace(datum)
			if datum == "" {
				continue
			}
			dataList = append(dataList, datum)
		}
		for _, data := range dataList {
			pairData := strings.Split(data, ": ")
			key := pairData[0]
			value := pairData[1]
			if key == "IPv4 Addresses" {
				// 移除逗号
				cleanedString := strings.ReplaceAll(value, ",", "")

				// 将清理后的字符串转换为整数
				intValue, _ := strconv.Atoi(cleanedString)
				result.Ipv4Counts = &intValue
				continue
			}
			if key == "Number of Peers" {
				// 移除逗号
				cleanedString := strings.ReplaceAll(value, ",", "")

				// 将清理后的字符串转换为整数
				intValue, _ := strconv.Atoi(cleanedString)
				result.PeersCounts = &intValue
				continue
			}
			if key == "Number of Prefixes" {
				// 移除逗号
				cleanedString := strings.ReplaceAll(value, ",", "")

				// 将清理后的字符串转换为整数
				intValue, _ := strconv.Atoi(cleanedString)
				result.PrefixesCounts = &intValue
				continue
			}
			if key == "Traffic Estimation" {
				result.TrafficBandwidth = value
				continue
			}
		}

		println(dataList)
	})

	//container data
	containerCSS := "#content-info > div:nth-child(2) > div:nth-child(1) > div:nth-child(1)"
	doc.Find(containerCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		regionalRegistry := strings.Join(strings.Fields(textData), " ")
		regionalRegistry = strings.Split(regionalRegistry, ":")[1]
		result.RegionalRegistry = regionalRegistry
	})
	//status
	statusCSS := "#content-info > div:nth-child(2) > div:nth-child(1) > div:nth-child(2)"
	doc.Find(statusCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		allocationStatus := strings.Join(strings.Fields(textData), " ")
		allocationStatus = strings.Split(allocationStatus, ":")[1]
		result.AllocationStatus = allocationStatus
	})
	//ration with less table data
	//ratio
	ratioCSS := "#content-info > div:nth-child(2) > div:nth-child(2) > div:nth-child(2)"
	doc.Find(ratioCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		if strings.Contains(textData, "Website") {
			website := strings.Join(strings.Fields(textData), " ")
			website = strings.SplitN(website, ":", 2)[1]
			result.Website = website
			return
		}
		if strings.Contains(textData, "Traffic Ratio") {
			trafficRatio := strings.Join(strings.Fields(textData), " ")
			trafficRatio = strings.Split(trafficRatio, ":")[1]
			result.TrafficRatio = trafficRatio
			return
		}
	})
	ratioCSS2 := "#content-info > div:nth-child(2) > div:nth-child(2) > div:nth-child(1)"
	doc.Find(ratioCSS2).Each(func(i int, s *goquery.Selection) {
		textData := s.Text()
		if strings.Contains(textData, "Traffic Ratio") {
			trafficRatio := strings.Join(strings.Fields(textData), " ")
			trafficRatio = strings.Split(trafficRatio, ":")[1]
			result.TrafficRatio = trafficRatio
			return
		}
	})
	//date
	dateCSS := "#content-info > div:nth-child(2) > div:nth-child(1) > div:nth-child(3)"
	doc.Find(dateCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		allocationDate := strings.Join(strings.Fields(textData), " ")
		allocationDate = strings.Split(allocationDate, ":")[1]
		result.AllocationDate = allocationDate
	})
	//exchange
	//exchangeCSS := "#content-info > div:nth-child(2) > div:nth-child(2) > div:nth-child(3)"
	//doc.Find(dateCSS).Each(func(i int, s *goquery.Selection) {
	//	// Extract text content of the selected element
	//	textData := s.Text()
	//	allocationDate := strings.Join(strings.Fields(textData), " ")
	//	allocationDate = strings.Split(allocationDate, ":")[1]
	//	result.AllocationDate = allocationDate
	//})

	//website
	websiteCSS := "#content-info > div:nth-child(2) > div:nth-child(2) > div:nth-child(4)"
	doc.Find(websiteCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		website := strings.Join(strings.Fields(textData), " ")
		website = strings.SplitN(website, ":", 2)[1]
		result.Website = website
	})
	//country
	countryCSS := "#content-info > div:nth-child(2) > div:nth-child(1) > div:nth-child(4) > span > img"
	doc.Find(countryCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		src, _ := s.Attr("src")
		// 解析 URL
		parsedURL, err := url.Parse(src)
		if err != nil {
			fmt.Println("解析 URL 时出错:", err)
			src = ""
		}
		// 提取路径部分并拆分成组件
		pathParts := strings.Split(parsedURL.Path, "/")
		// 获取倒数第二个部分（国家代码）
		if len(pathParts) > 1 {
			countryCode := strings.ToUpper(strings.TrimSuffix(pathParts[len(pathParts)-1], ".png"))
			//fmt.Println("国家代码:", countryCode)
			src = countryCode
		}

		title, _ := s.Attr("title")
		result.AllocationCountry = fmt.Sprintf("%s(%s)", src, title)
	})

	//v4prefix
	v4prefix := "#content-info > div:nth-child(5) > div:nth-child(1) > div:nth-child(1)"
	doc.Find(v4prefix).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		ipv4Prefixies := strings.Join(strings.Fields(textData), " ")
		ipv4Prefixies = strings.Split(ipv4Prefixies, ":")[1]
		ipv4Prefixies = strings.ReplaceAll(ipv4Prefixies, " ", "")
		if strings.Contains(ipv4Prefixies, ",") {
			ipv4Prefixies = strings.ReplaceAll(ipv4Prefixies, ",", "")
		}
		intValue, _ := strconv.Atoi(ipv4Prefixies)
		result.Ipv4Prefixies = &intValue
	})
	//v6prefix
	v6prefix := "#content-info > div:nth-child(5) > div:nth-child(2) > div:nth-child(1)"
	doc.Find(v6prefix).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		ipv6Prefixies := strings.Join(strings.Fields(textData), " ")
		ipv6Prefixies = strings.Split(ipv6Prefixies, ":")[1]
		ipv6Prefixies = strings.ReplaceAll(ipv6Prefixies, " ", "")
		if strings.Contains(ipv6Prefixies, ",") {
			ipv6Prefixies = strings.ReplaceAll(ipv6Prefixies, ",", "")
		}
		intValue, _ := strconv.Atoi(ipv6Prefixies)
		result.Ipv6Prefixies = &intValue
	})
	//v4Peers
	v4PeersCSS := "#content-info > div:nth-child(5) > div:nth-child(1) > div:nth-child(2)"
	doc.Find(v4PeersCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		v4Peers := strings.Join(strings.Fields(textData), " ")
		v4Peers = strings.Split(v4Peers, ":")[1]
		v4Peers = strings.ReplaceAll(v4Peers, " ", "")
		if strings.Contains(v4Peers, ",") {
			v4Peers = strings.ReplaceAll(v4Peers, ",", "")
		}
		intValue, _ := strconv.Atoi(v4Peers)
		result.Ipv4Peers = &intValue
	})
	//v6Peers
	v6PeersCSS := "#content-info > div:nth-child(5) > div:nth-child(2) > div:nth-child(2)"
	doc.Find(v6PeersCSS).Each(func(i int, s *goquery.Selection) {
		// Extract text content of the selected element
		textData := s.Text()
		v6Peers := strings.Join(strings.Fields(textData), " ")
		v6Peers = strings.Split(v6Peers, ":")[1]
		v6Peers = strings.ReplaceAll(v6Peers, " ", "")
		if strings.Contains(v6Peers, ",") {
			v6Peers = strings.ReplaceAll(v6Peers, ",", "")
		}
		intValue, _ := strconv.Atoi(v6Peers)
		result.Ipv6Peers = &intValue
	})
	//println(result)
	return result
}

// DeleteAsnInfo 删除asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) DeleteAsnInfo(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.AsnInfo{}, "id = ?", ID).Error
	return err
}

// DeleteAsnInfoByIds 批量删除asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) DeleteAsnInfoByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.AsnInfo{}, "id in ?", IDs).Error
	return err
}

// UpdateAsnInfo 更新asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) UpdateAsnInfo(asnInfo cfscan.AsnInfo) (err error) {
	err = global.GVA_DB.Model(&cfscan.AsnInfo{}).Where("id = ?", asnInfo.ID).Updates(&asnInfo).Error
	return err
}

// GetAsnInfo 根据ID获取asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) GetAsnInfo(ID string) (asnInfo cfscan.AsnInfo, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&asnInfo).Error
	return
}

// GetAsnInfoInfoList 分页获取asnInfo表记录
// Author [piexlmax](https://github.com/piexlmax)
func (asnInfoService *AsnInfoService) GetAsnInfoInfoList(info cfscanReq.AsnInfoSearch) (list []cfscan.AsnInfo, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.AsnInfo{})
	var asnInfos []cfscan.AsnInfo
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	//模糊查询
	if info.AsnName != "" {
		db = db.Where("asn_name LIKE ?", "%"+info.AsnName+"%")
	}
	if info.FullName != "" {
		db = db.Where("full_name LIKE ?", "%"+info.FullName+"%")
	}
	if info.AllocationCountry != "" {
		db = db.Where("allocation_country LIKE ?", "%"+info.AllocationCountry+"%")

	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	OrderStr := "id desc"
	err = db.Order(OrderStr).Find(&asnInfos).Error
	return asnInfos, total, err
}
