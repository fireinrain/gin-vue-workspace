package cfscan

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscanReq "github.com/flipped-aurora/gin-vue-admin/server/model/cfscan/request"
	cfscanUtils "github.com/flipped-aurora/gin-vue-admin/server/utils/cfscan"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ScheduleTaskService struct{}

var cronManager = NewConfigFileManager(dynamicCronFile)

const typeDistributeCronASNScan = "cron-scan:cron-scan-asn"

// CreateScheduleTask 创建scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) CreateScheduleTask(scheduleTask *cfscan.ScheduleTask) (err error) {
	//初始化任务 必须是非开启 任务状态是禁用中的
	if scheduleTask.Enable != "0" || scheduleTask.TaskStatus != "0" {
		return errors.New("新建定时任务的状态必须为禁用中和非开启状态")
	}
	//设置为1970 时间戳
	unixEpoch := time.Unix(0, 0)
	scheduleTask.LastRunAt = unixEpoch
	scheduleTask.NextRunAt = unixEpoch
	err = global.GVA_DB.Create(scheduleTask).Error
	if err != nil {
		log.Printf(">>> Create schedule task failed: %s", err)
		return err
	}
	//不做直接保存到动态cron文件
	//cronConfig := Config{
	//	CronId:      scheduleTask.ID,
	//	Cronspec:    scheduleTask.CrontabStr,
	//	TaskType:    typeDistributeCronASNScan,
	//	TaskPayload: scheduleTask.TaskConfig,
	//}
	//
	//err = cronManager.InsertTask(cronConfig)
	//if err != nil {
	//	log.Printf(">>> Write cron task to dynamic config file failed: %s", err)
	//	return err
	//}
	return nil
}

// DeleteScheduleTask 删除scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) DeleteScheduleTask(ID string) (err error) {
	task := cfscan.ScheduleTask{Enable: "0", TaskStatus: "0"}
	//将enable 设置为0
	err = global.GVA_DB.Model(&cfscan.ScheduleTask{}).Where("id = ?", ID).Updates(&task).Error
	if err != nil {
		log.Printf(">>> Update schedule task failed: %s", err)
		return err
	}
	err = global.GVA_DB.Delete(&cfscan.ScheduleTask{}, "id = ?", ID).Error
	if err != nil {
		log.Printf(">>> Delete schedule task failed: %s", err)
		return err
	}
	// 将字符串解析为uint64
	num, _ := strconv.ParseUint(ID, 10, 64)
	// 将uint64转换为uint
	result := uint(num)
	err = cronManager.DeleteTaskByID(result)
	return err
}

// DeleteScheduleTaskByIds 批量删除scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) DeleteScheduleTaskByIds(IDs []string) (err error) {
	for _, id := range IDs {
		task := cfscan.ScheduleTask{Enable: "0", TaskStatus: "0"}
		//将enable 设置为0
		err = global.GVA_DB.Model(&cfscan.ScheduleTask{}).Where("id = ?", id).Updates(&task).Error
		if err != nil {
			log.Printf(">>> Update schedule task failed: %s", err)
			return err
		}
	}

	err = global.GVA_DB.Delete(&[]cfscan.ScheduleTask{}, "id in ?", IDs).Error
	if err != nil {
		log.Printf(">>> Delete schedule tasks failed: %s", err)
		return err
	}
	for _, ID := range IDs {
		// 将字符串解析为uint64
		num, _ := strconv.ParseUint(ID, 10, 64)
		// 将uint64转换为uint
		result := uint(num)
		err = cronManager.DeleteTaskByID(result)
	}
	return err
}

// UpdateScheduleTask 更新scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) UpdateScheduleTask(scheduleTask cfscan.ScheduleTask) (err error) {
	//任务在调度中或者是在运行中 无法修改
	if scheduleTask.TaskStatus == "2" || scheduleTask.TaskStatus == "3" {
		return errors.New("当前任务正在调度或运行中,无法修改,请稍后尝试")
	}
	if scheduleTask.Enable == "1" {
		scheduleTask.TaskStatus = "1"
	}
	if scheduleTask.Enable == "0" {
		scheduleTask.TaskStatus = "0"
	}
	//获取定时任务表达式 使用当前时间 + 表达式时间
	nextRunTime, err := cfscanUtils.NextRunTime(scheduleTask.CrontabStr)
	if err != nil {
		return errors.New("无法计算定时任务下一次运行时间")
	}
	scheduleTask.NextRunAt = nextRunTime
	err = global.GVA_DB.Model(&cfscan.ScheduleTask{}).Where("id = ?", scheduleTask.ID).Updates(&scheduleTask).Error
	if err != nil {
		log.Printf(">>> Update schedule task failed: %s", err)
		return err
	}

	if scheduleTask.Enable == "1" && scheduleTask.TaskStatus == "1" {
		newCronConf := Config{
			CronId:      scheduleTask.ID,
			Cronspec:    scheduleTask.CrontabStr,
			TaskType:    typeDistributeCronASNScan,
			TaskPayload: scheduleTask.TaskConfig,
		}
		err = cronManager.UpdateTaskByID(scheduleTask.ID, newCronConf)
		if err != nil && strings.Contains(err.Error(), "not found") {
			//insert new
			_ = cronManager.InsertTask(newCronConf)
			err = nil
		}
	}
	//移除定时任务
	if scheduleTask.Enable == "0" && scheduleTask.TaskStatus == "0" {
		err = cronManager.DeleteTaskByID(scheduleTask.ID)
		if err != nil && strings.Contains(err.Error(), "not found") {
			err = nil
		}
	}

	return err
}

// GetScheduleTask 根据ID获取scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) GetScheduleTask(ID string) (scheduleTask cfscan.ScheduleTask, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&scheduleTask).Error
	return
}

// GetScheduleTaskInfoList 分页获取scheduleTask表记录
// Author [piexlmax](https://github.com/piexlmax)
func (scheduleTaskService *ScheduleTaskService) GetScheduleTaskInfoList(info cfscanReq.ScheduleTaskSearch) (list []cfscan.ScheduleTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cfscan.ScheduleTask{})
	var scheduleTasks []cfscan.ScheduleTask
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.TaskDesc != "" {
		db = db.Where("task_desc LIKE ?", "%"+info.TaskDesc+"%")
	}
	if info.AsnNumber != "" {
		db = db.Where("asn_number LIKE ?", "%"+info.AsnNumber+"%")
	}
	if info.AsnDesc != "" {
		db = db.Where("asn_desc LIKE ?", "%"+info.AsnDesc+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["asn_number"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	} else {
		OrderStr = "id desc"
		db = db.Order(OrderStr)
	}

	err = db.Find(&scheduleTasks).Error
	return scheduleTasks, total, err
}

//动态读写dynamic cron tasks

const dynamicCronFile = "dynamic_cron.yml"

type Config struct {
	CronId      uint   `json:"cron_id"`
	Cronspec    string `yaml:"cronspec"`
	TaskType    string `yaml:"task_type"`
	TaskPayload string `yaml:"task_payload"`
}

// ConfigFileManager 实现 增删改查 yml节点的方法
type ConfigFileManager struct {
	ConfigFilePath string `json:"config_file_path"`
}

type ConfigFile struct {
	Configs []*Config `yaml:"configs"`
}

func NewConfigFileManager(configFilePath string) *ConfigFileManager {
	return &ConfigFileManager{ConfigFilePath: configFilePath}
}

// ReadYAMLFile 读取YAML文件
func (cfm *ConfigFileManager) ReadYAMLFile() (*ConfigFile, error) {
	data, err := os.ReadFile(cfm.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	var configFile ConfigFile
	err = yaml.Unmarshal(data, &configFile)
	if err != nil {
		return nil, err
	}

	return &configFile, nil
}

// WriteYAMLFile 写入YAML文件
func (cfm *ConfigFileManager) WriteYAMLFile(configFile *ConfigFile) error {
	data, err := yaml.Marshal(configFile)
	if err != nil {
		return err
	}

	return os.WriteFile(cfm.ConfigFilePath, data, 0644)
}

// UpdateTaskByID 根据ID更新YAML文件中的任务
func (cfm *ConfigFileManager) UpdateTaskByID(id uint, newTask Config) error {
	configFile, err := cfm.ReadYAMLFile()
	if err != nil {
		return err
	}

	found := false
	for i, task := range configFile.Configs {
		if task.CronId == id {
			configFile.Configs[i] = &newTask
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return cfm.WriteYAMLFile(configFile)
}

// CreateYAMLFile 创建新的YAML文件
func (cfm *ConfigFileManager) CreateYAMLFile(tasks []Config) error {
	var r []*Config
	for _, task := range tasks {
		r = append(r, &task)
	}
	configFile := ConfigFile{Configs: r}
	return cfm.WriteYAMLFile(&configFile)
}

// DeleteTaskByID 根据ID从YAML文件中删除任务
func (cfm *ConfigFileManager) DeleteTaskByID(id uint) error {
	configFile, err := cfm.ReadYAMLFile()
	if err != nil {
		return err
	}

	newConfigs := make([]Config, 0)
	found := false
	for _, task := range configFile.Configs {
		if task.CronId != id {
			newConfigs = append(newConfigs, *task)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}
	var cs []*Config
	for _, config := range newConfigs {
		cs = append(cs, &config)
	}

	configFile.Configs = cs
	return cfm.WriteYAMLFile(configFile)
}

// InsertTask 插入新任务到YAML文件
func (cfm *ConfigFileManager) InsertTask(newTask Config) error {
	configFile, err := cfm.ReadYAMLFile()
	if err != nil {
		// 如果文件不存在，创建一个新的配置文件
		if os.IsNotExist(err) {
			return cfm.CreateYAMLFile([]Config{newTask})
		}
		return err
	}

	// 检查是否已存在相同ID的任务
	for _, task := range configFile.Configs {
		if task.CronId == newTask.CronId {
			return fmt.Errorf("task with ID %d already exists", newTask.CronId)
		}
	}

	// 添加新任务
	configFile.Configs = append(configFile.Configs, &newTask)

	return cfm.WriteYAMLFile(configFile)
}

// QueryTaskByID 根据ID查询任务
func (cfm *ConfigFileManager) QueryTaskByID(id uint) (*Config, error) {
	configFile, err := cfm.ReadYAMLFile()
	if err != nil {
		return nil, err
	}

	for _, task := range configFile.Configs {
		if task.CronId == id {
			return task, nil
		}
	}

	return nil, fmt.Errorf("task with ID %d not found", id)
}
