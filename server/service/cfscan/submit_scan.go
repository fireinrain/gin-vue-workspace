package cfscan

import (
	"context"
	tasks "github.com/endless-cfcdn/shared-tasks"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
	cfscan_utils "github.com/flipped-aurora/gin-vue-admin/server/utils/cfscan"
	"log"
	"strings"
	"sync"
	"time"
)

var astkCtx context.Context
var taskCancel context.CancelFunc

func init() {
	astkCtx, taskCancel = context.WithCancel(context.Background())
}

var asnInfoService = new(AsnInfoService)

type SubmitScanService struct{}

// CreateSubmitScan 创建submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) CreateSubmitScan(submitScan *cfscan.SubmitScan) (err error) {
	submitScan.ScanStatus = "0"
	r := global.GVA_DB.Create(submitScan)
	if r.Error != nil {
		return err
	}
	//scan single ASN
	if submitScan.ScanType == "1" {
		err = DoASNScanBackground(submitScan)
		if err != nil {
			return err
		}
	}
	//scan ASNS
	if submitScan.ScanType == "2" {
		err := DoASNSScanBackground(submitScan)
		if err != nil {
			return err
		}
		return nil
	}
	//scan single ip or ips
	if submitScan.ScanType == "3" || submitScan.ScanType == "4" {
		err := DoIPScanBackground(submitScan)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func DoIPScanBackground(submitScan *cfscan.SubmitScan) error {
	ctx, cancel := context.WithTimeout(astkCtx, 12*time.Hour)
	newRecordId := submitScan.ID

	var batchedCIDRS [][]string
	//scan IP
	//IP类型 单ip类型
	if submitScan.IpinfoType == "1" {
		addresses := cfscan_utils.ExtractIPv4Addresses(submitScan.IpinfoList)
		batchedCIDRS = [][]string{addresses}
	}

	//CIDR类型
	if submitScan.IpinfoType == "2" {
		addresses := cfscan_utils.ExtractIPv4CIDRAddresses(submitScan.IpinfoList)
		//这里需要限制 最大批量数为25
		batchedCIDRS = cfscan_utils.SplitCIDRs(addresses, *submitScan.IpbatchSize)
	}

	//update scan status--> 未运行到运行中
	subscan := cfscan.SubmitScan{
		ScanStatus: "1",
	}
	//update to db
	err := global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&subscan).Error
	if err != nil {
		log.Printf("Error on update submit scan for scan status: %s\n", err.Error())
		return err
	}
	taskBatchSize := len(batchedCIDRS)
	log.Printf("本次扫描任务将分成: %d 个小扫描任务\n", taskBatchSize)
	var wg sync.WaitGroup
	results := make(chan string, taskBatchSize)
	for index, cidr := range batchedCIDRS {
		wg.Add(1)
		go func(index int, cidr []string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// 超时或被取消
				log.Printf("当前任务: %d 已超时,分发任务已取消运行...\n", index)
				return
			default:
				var enableSpeedtest int = 1
				if submitScan.EnableSpeedtest == "0" {
					enableSpeedtest = 0
				}
				payload := tasks.IPScanCFPayload{
					AsnNumber:       submitScan.AsnNumber,
					EnableTls:       submitScan.EnableTls,
					ScanPorts:       submitScan.ScanPorts,
					ScanRate:        *submitScan.ScanRate,
					IpcheckThread:   *submitScan.IpcheckThread,
					EnableSpeedtest: enableSpeedtest,
					CIDRList:        cidr,
					IPBatchSize:     *submitScan.IpbatchSize,
				}
				runTask, _ := tasks.NewIPScanCFTask(payload)

				info, err := global.AsynQClient.EnqueueContext(ctx, runTask)
				if err != nil {
					log.Printf("Enqueue error: %v", err)
					return
				}
				log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

				// 等待任务完成
				for {
					select {
					case <-ctx.Done():
						// 超时或被取消
						log.Printf("当前任务: %s已超时,任务等待结果已取消运行...\n", info.ID)
						return
					default:
						taskInfo, err := global.AsynQInspector.GetTaskInfo("default", info.ID)
						result := taskInfo.Result

						if result == nil {
							time.Sleep(5 * time.Second)
							continue
						}
						if err != nil {
							log.Printf("Error getting task info: %v", err)
						}
						log.Printf("Batch Task(%d/%d): %s,运行完成：\n", index, taskBatchSize, info.ID)
						results <- string(result)
						return
					}
				}
			}
		}(index, cidr)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	go func() {
		defer cancel()
		// 收集结果
		var finalResult []string
		for {
			select {
			case <-ctx.Done():
				// 超时或被取消，更新扫描状态为超时
				sub := cfscan.SubmitScan{
					ScanStatus: "3", // 假设3表示超时状态
					ScanResult: strings.Join(finalResult, ","),
				}
				err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&sub).Error
				if err != nil {
					log.Printf("Error on update submit scan for timeout: %s\n", err.Error())
				}
				log.Printf("Batch Task cancled due to the timeout on waiting result: RecordID: %s", newRecordId)
				return
			case result, ok := <-results:
				if !ok {
					// 通道已关闭，所有结果已收集完毕
					// ... (保持原有的结果处理逻辑不变)
					//convert sub json list to big json
					// 创建一个切片来存储去掉方括号的 JSON 对象
					var jsonObjects []string

					// 遍历每个 JSON 列表并去掉方括号
					for _, jsonList := range finalResult {
						jsonObjects = append(jsonObjects, strings.TrimSuffix(strings.TrimPrefix(jsonList, "["), "]"))
					}

					// 将所有 JSON 对象用逗号连接起来，并包裹在方括号中
					mergedJSON := "[" + strings.Join(jsonObjects, ",") + "]"
					sub := cfscan.SubmitScan{
						ScanResult: mergedJSON,
						ScanStatus: "2",
					}
					//save to db
					err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&sub).Error
					if err != nil {
						log.Printf("Error on update submit scan for result data: %s\n", err.Error())
					}
					log.Println("Submit Task finished: Submit RecordID", newRecordId)
					return
				}
				if result != "EmptyResult" {
					finalResult = append(finalResult, result)
				}
			}
		}

	}()

	return nil
}

func DoASNScanBackground(submitScan *cfscan.SubmitScan) error {
	ctx, cancel := context.WithTimeout(astkCtx, 12*time.Hour)
	newRecordId := submitScan.ID

	//scan ASN

	var asnInfo cfscan.AsnInfo
	tx := global.GVA_DB.Where("asn_name = ?", submitScan.AsnNumber).Find(&asnInfo)
	if tx.RowsAffected == 0 {
		//创建一个
		var enable = 1
		info := cfscan.AsnInfo{AsnName: submitScan.AsnNumber, Enable: &enable}
		asnInfoService.CreateAsnInfo(&info)
		_ = global.GVA_DB.Where("asn_name = ?", submitScan.AsnNumber).Find(&asnInfo)
	}
	waitForProcess := strings.Split(asnInfo.Ipv4CIDR, "\n")
	//这里需要限制 最大批量数为25
	batchedCIDRS := cfscan_utils.SplitCIDRs(waitForProcess, *submitScan.IpbatchSize)

	//update scan status--> 未运行到运行中
	subscan := cfscan.SubmitScan{
		ScanStatus: "1",
	}
	//update to db
	err := global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&subscan).Error
	if err != nil {
		log.Printf("Error on update submit scan for scan status: %s\n", err.Error())
		return err
	}
	taskBatchSize := len(batchedCIDRS)
	log.Printf("本次扫描任务将分成: %d 个小扫描任务\n", taskBatchSize)

	var wg sync.WaitGroup
	results := make(chan string, taskBatchSize)
	for index, cidr := range batchedCIDRS {
		wg.Add(1)
		go func(index int, cidr []string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// 超时或被取消
				log.Printf("当前任务: %d 已超时,分发任务已取消运行...\n", index)
				return
			default:
				var enableSpeedtest int = 1
				if submitScan.EnableSpeedtest == "0" {
					enableSpeedtest = 0
				}
				payload := tasks.ASNScanCFPayload{
					AsnNumber:       submitScan.AsnNumber,
					EnableTls:       submitScan.EnableTls,
					ScanPorts:       submitScan.ScanPorts,
					ScanRate:        *submitScan.ScanRate,
					IpcheckThread:   *submitScan.IpcheckThread,
					EnableSpeedtest: enableSpeedtest,
					CIDRList:        cidr,
					IPBatchSize:     *submitScan.IpbatchSize,
				}
				runTask, _ := tasks.NewASNScanCFTask(payload)

				info, err := global.AsynQClient.EnqueueContext(ctx, runTask)
				if err != nil {
					log.Printf("Enqueue error: %v", err)
					return
				}
				log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

				// 等待任务完成
				for {
					select {
					case <-ctx.Done():
						// 超时或被取消
						log.Printf("当前任务: %s已超时,任务等待结果已取消运行...\n", info.ID)
						return
					default:
						taskInfo, err := global.AsynQInspector.GetTaskInfo("default", info.ID)
						result := taskInfo.Result

						if result == nil {
							time.Sleep(5 * time.Second)
							continue
						}
						if err != nil {
							log.Printf("Error getting task info: %v", err)
						}
						log.Printf("Batch Task(%d/%d): %s,运行完成：\n", index, taskBatchSize, info.ID)
						results <- string(result)
						return
					}
				}
			}
		}(index, cidr)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	go func() {
		defer cancel()
		// 收集结果
		var finalResult []string
		for {
			select {
			case <-ctx.Done():
				// 超时或被取消，更新扫描状态为超时
				sub := cfscan.SubmitScan{
					ScanStatus: "3", // 假设3表示超时状态
					ScanResult: strings.Join(finalResult, ","),
				}
				err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&sub).Error
				if err != nil {
					log.Printf("Error on update submit scan for timeout: %s\n", err.Error())
				}
				log.Printf("Batch Task cancled due to the timeout on waiting result: SubmitRecordID: %s", newRecordId)
				return
			case result, ok := <-results:
				if !ok {
					// 通道已关闭，所有结果已收集完毕
					// ... (保持原有的结果处理逻辑不变)
					//convert sub json list to big json
					// 创建一个切片来存储去掉方括号的 JSON 对象
					var jsonObjects []string

					// 遍历每个 JSON 列表并去掉方括号
					for _, jsonList := range finalResult {
						jsonObjects = append(jsonObjects, strings.TrimSuffix(strings.TrimPrefix(jsonList, "["), "]"))
					}

					// 将所有 JSON 对象用逗号连接起来，并包裹在方括号中
					mergedJSON := "[" + strings.Join(jsonObjects, ",") + "]"
					sub := cfscan.SubmitScan{
						ScanResult: mergedJSON,
						ScanStatus: "2",
					}
					//save to db
					err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&sub).Error
					if err != nil {
						log.Printf("Error on update submit scan for result data: %s\n", err.Error())
					}

					log.Println("Submit Task finished: Submit RecordID", newRecordId)

					return
				}
				if result != "EmptyResult" {
					finalResult = append(finalResult, result)
				}
			}
		}

	}()

	return nil
}

func DoASNSScanBackground(submitScan *cfscan.SubmitScan) error {
	ctx, cancel := context.WithTimeout(astkCtx, 12*3*time.Hour)
	newRecordId := submitScan.ID

	//scan ASN
	asnNumbers := submitScan.AsnNumber
	asnList := strings.Split(asnNumbers, ",")
	var allAsnCIDRList []string
	for _, asn := range asnList {
		log.Println("Process ASN: ", asn)
		var asnInfo cfscan.AsnInfo
		tx := global.GVA_DB.Where("asn_name = ?", asn).Find(&asnInfo)
		if tx.RowsAffected == 0 {
			//创建一个
			var enable = 1
			info := cfscan.AsnInfo{AsnName: asn, Enable: &enable}
			asnInfoService.CreateAsnInfo(&info)
			_ = global.GVA_DB.Where("asn_name = ?", asn).Find(&asnInfo)
		}
		waitForProcess := strings.Split(asnInfo.Ipv4CIDR, "\n")
		allAsnCIDRList = append(allAsnCIDRList, waitForProcess...)
	}

	//这里需要限制 最大批量数为25
	batchedCIDRS := cfscan_utils.SplitCIDRs(allAsnCIDRList, *submitScan.IpbatchSize)

	//update scan status--> 未运行到运行中
	subscan := cfscan.SubmitScan{
		ScanStatus: "1",
	}
	//update to db
	err := global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&subscan).Error
	if err != nil {
		log.Printf("Error on update submit scan for scan status: %s\n", err.Error())
		return err
	}
	taskBatchSize := len(batchedCIDRS)
	log.Printf("本次扫描任务将分成: %d 个小扫描任务\n", taskBatchSize)

	var wg sync.WaitGroup
	results := make(chan string, len(batchedCIDRS))
	for index, cidr := range batchedCIDRS {
		wg.Add(1)
		go func(index int, cidr []string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// 超时或被取消
				log.Printf("当前任务: %d 已超时,分发任务已取消运行...\n", index)
				return
			default:
				var enableSpeedtest int = 1
				if submitScan.EnableSpeedtest == "0" {
					enableSpeedtest = 0
				}
				payload := tasks.ASNScanCFPayload{
					AsnNumber:       submitScan.AsnNumber,
					EnableTls:       submitScan.EnableTls,
					ScanPorts:       submitScan.ScanPorts,
					ScanRate:        *submitScan.ScanRate,
					IpcheckThread:   *submitScan.IpcheckThread,
					EnableSpeedtest: enableSpeedtest,
					CIDRList:        cidr,
					IPBatchSize:     *submitScan.IpbatchSize,
				}
				runTask, _ := tasks.NewASNScanCFTask(payload)

				info, err := global.AsynQClient.EnqueueContext(ctx, runTask)
				if err != nil {
					log.Printf("Enqueue error: %v", err)
					return
				}
				log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

				// 等待任务完成
				for {
					select {
					case <-ctx.Done():
						// 超时或被取消
						log.Printf("当前任务: %s已超时,任务等待结果已取消运行...\n", info.ID)
						return
					default:
						taskInfo, err := global.AsynQInspector.GetTaskInfo("default", info.ID)
						result := taskInfo.Result

						if result == nil {
							time.Sleep(5 * time.Second)
							continue
						}
						if err != nil {
							log.Printf("Error getting task info: %v", err)
						}
						log.Printf("Batch Task(%d/%d): %s,运行完成：\n", index, taskBatchSize, info.ID)
						results <- string(result)
						return
					}
				}
			}
		}(index, cidr)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	go func() {
		defer cancel()
		// 收集结果
		var finalResult []string
		for {
			select {
			case <-ctx.Done():
				// 超时或被取消，更新扫描状态为超时
				sub := cfscan.SubmitScan{
					ScanStatus: "3", // 假设3表示超时状态
					ScanResult: strings.Join(finalResult, ","),
				}
				err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&sub).Error
				if err != nil {
					log.Printf("Error on update submit scan for timeout: %s\n", err.Error())
				}
				log.Printf("Batch Task cancled due to the timeout on waiting result: SubmitRecordID: %s", newRecordId)
				return
			case result, ok := <-results:
				if !ok {
					// 通道已关闭，所有结果已收集完毕
					// ... (保持原有的结果处理逻辑不变)
					//convert sub json list to big json
					// 创建一个切片来存储去掉方括号的 JSON 对象
					var jsonObjects []string

					// 遍历每个 JSON 列表并去掉方括号
					for _, jsonList := range finalResult {
						jsonObjects = append(jsonObjects, strings.TrimSuffix(strings.TrimPrefix(jsonList, "["), "]"))
					}

					// 将所有 JSON 对象用逗号连接起来，并包裹在方括号中
					mergedJSON := "[" + strings.Join(jsonObjects, ",") + "]"
					sub := cfscan.SubmitScan{
						ScanResult: mergedJSON,
						ScanStatus: "2",
					}
					//save to db
					err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", newRecordId).Updates(&sub).Error
					if err != nil {
						log.Printf("Error on update submit scan for result data: %s\n", err.Error())
					}

					log.Println("Submit Task finished: Submit RecordID", newRecordId)

					return
				}
				if result != "EmptyResult" {
					finalResult = append(finalResult, result)
				}
			}
		}

	}()

	return nil

}

// DeleteSubmitScan 删除submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) DeleteSubmitScan(ID string) (err error) {
	err = global.GVA_DB.Delete(&cfscan.SubmitScan{}, "id = ?", ID).Error
	return err
}

// DeleteSubmitScanByIds 批量删除submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) DeleteSubmitScanByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cfscan.SubmitScan{}, "id in ?", IDs).Error
	return err
}

// UpdateSubmitScan 更新submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) UpdateSubmitScan(submitScan cfscan.SubmitScan) (err error) {
	err = global.GVA_DB.Model(&cfscan.SubmitScan{}).Where("id = ?", submitScan.ID).Updates(&submitScan).Error
	return err
}

// GetSubmitScan 根据ID获取submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) GetSubmitScan(ID string) (submitScan cfscan.SubmitScan, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&submitScan).Error
	return
}

// GetSubmitScanInfoList 分页获取submitScan表记录
// Author [piexlmax](https://github.com/piexlmax)
func (submitScanService *SubmitScanService) GetSubmitScanInfoList(info cfscanReq.SubmitScanSearch) (list []cfscan.SubmitScan, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.SubmitScan{})
	var submitScans []cfscan.SubmitScan
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.ScanDesc != "" {
		db = db.Where("scan_desc LIKE ?", "%"+info.ScanDesc+"%")
	}
	if info.ScanType != "" {
		db = db.Where("scan_type = ?", info.ScanType)
	}
	if info.AsnNumber != "" {
		db = db.Where("asn_number LIKE ?", "%"+info.AsnNumber+"%")
	}
	if info.ScanStatus != "" {
		db = db.Where("scan_status = ?", info.ScanStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&submitScans).Error
	return submitScans, total, err
}
