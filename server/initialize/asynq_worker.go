package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	tasks "github.com/endless-cfcdn/shared-tasks"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscan_utils "github.com/flipped-aurora/gin-vue-admin/server/utils/cfscan"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const redisAddr = "cloud2.131433.xyz:5379"
const redisPass = "fireinrain@redis"
const redisDb = 0
const dynamicCronFile = "dynamic_cron.yml"

var astkCtx context.Context
var taskCancel context.CancelFunc

func init() {
	astkCtx, taskCancel = context.WithCancel(context.Background())
}

func AsynQClient() *asynq.Client {
	redisClientOpt := asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass, DB: redisDb}
	client := asynq.NewClient(redisClientOpt)
	return client
	//defer client.Close()
}

func AsynQInspector() *asynq.Inspector {
	redisClientOpt := asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass, DB: redisDb}
	inspector := asynq.NewInspector(redisClientOpt)
	return inspector
}

func SelfAsynQTaskClientRun() {

	go func() {
		//scheduler
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}
		scheduler := asynq.NewScheduler(
			asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass, DB: redisDb},

			&asynq.SchedulerOpts{
				Location: loc,
			},
		)
		now := time.Now()
		payload := UpdateASNInfoCIDR{
			UUID:      uuid.New().String(),
			TimeStamp: &now,
		}
		//task := asynq.NewTask(TypeBashTask, []byte("ls"), asynq.Queue("self-admin"))
		runTask, _ := NewUpdateASNInfoCIDRTask(payload, "self-admin")

		entryID, err := scheduler.Register("@every 360h", runTask)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("本地注册定时任务: scheduler entry: %q\n", entryID)
		if err := scheduler.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {

		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass, DB: redisDb},
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: 3,
				// Optionally specify multiple queues with different priority.
				Queues: map[string]int{
					"cron-scan":  6,
					"self-admin": 3,
				},
				// See the godoc for other configuration options
			},
		)
		mux := asynq.NewServeMux()
		mux.Handle(TypeUpdateASNInfoCIDR, NewUpdateCIDRRunProcessor())
		//定时分发扫描ASN的任务
		mux.Handle(TypeDistributeCronASNScan, NewDistributeCronASNScanProcessor())
		log.Println(">>> 注册本地AsynQ-Wokrer成功...")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	//动态定时任务分发器
	go func() {
		time.Sleep(5 * time.Second)
		provider := &FileBasedConfigProvider{filename: dynamicCronFile}

		mgr, err := asynq.NewPeriodicTaskManager(
			asynq.PeriodicTaskManagerOpts{
				RedisConnOpt:               asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass, DB: redisDb},
				PeriodicTaskConfigProvider: provider,         // this provider object is the interface to your config source
				SyncInterval:               10 * time.Second, // this field specifies how often sync should happen
			})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(">>> 动态定时任务管理器启动成功...")

		if err := mgr.Run(); err != nil {
			log.Fatal(err)
		}
	}()

}

type CronTaskPayload struct {
	SchedulerTaskID string `json:"scheduler_task_id"`
	CronTaskStr     string `json:"cron_task_str"`
	ScanDesc        string `json:"scanDesc"`
	ScanType        string `json:"scanType"`
	AsnNumber       string `json:"asnNumber"`
	IpbatchSize     int    `json:"ipbatchSize"`
	EnableTls       string `json:"enableTls"`
	ScanPorts       string `json:"scanPorts"`
	ScanRate        int    `json:"scanRate"`
	IpcheckThread   int    `json:"ipcheckThread"`
	EnableSpeedtest string `json:"enableSpeedtest"`
}

func StartAllCronTaskFromDB() {
	// 开始事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println(">>> transaction rolled back due to panic:", r)
		}
	}()

	//重置提交任务 中的状态
	upResult0 := tx.Model(&cfscan.SubmitScan{}).Where("scan_status = ?", "1").Update("scan_status", "3")
	if upResult0.Error != nil {
		log.Println(">>> Failed to update SubmitScan scan_status:", upResult0.Error)
		tx.Rollback()
		return
	}
	//设置scheduleTaskHist 全部设置为超时(3) 因为服务器进行了重启

	// 更新所有 HistStatus 为 "1" 的记录，将其修改为 "3"
	// 获取当前时间
	now := time.Now()

	// 更新所有 hist_status 为 "1" 的记录
	var taskHists []cfscan.ScheduleTaskHist
	err := tx.Model(&cfscan.ScheduleTaskHist{}).Where("hist_status = ?", "1").Find(&taskHists).Error
	if err != nil {
		log.Fatalf("查询错误: %v", err)
	}

	for _, taskHist := range taskHists {
		taskHist.HistStatus = "3"
		taskHist.EndTime = now
		taskHist.CostTime = int(now.Sub(taskHist.StartTime).Seconds())

		err = tx.Model(&cfscan.ScheduleTaskHist{}).Where("id = ?", taskHist.ID).Updates(&taskHist).Error
		if err != nil {
			log.Printf("更新ScheduleTaskHist错误: %v", err)
			// 处理错误，例如记录错误日志或回滚事务
			break
		}
	}

	if err != nil {
		log.Println(">>> Failed to update ScheduleTaskHist hist_status:", err)
		tx.Rollback()
		return
	}
	//重置定时任务 中的状态
	upResult2 := tx.Model(&cfscan.ScheduleTask{}).Where("task_status = ? OR task_status = ?", "2", "3").Update("task_status", "1")
	if upResult2.Error != nil {
		log.Println(">>> Failed to update ScheduleTask task_status:", upResult2.Error)
		tx.Rollback()
		return
	}
	//更新task_status=1 的下次运行时间
	var tasks []cfscan.ScheduleTask
	if err := tx.Where("enable = ? AND task_status = ?", "1", "1").Find(&tasks).Error; err != nil {
		tx.Rollback()
		log.Println(">>> Failed to find ScheduleTask records:", err)
		return
	}

	for _, task := range tasks {
		nextRunAt, err := NextRunTime(task.CrontabStr)
		if err != nil {
			tx.Rollback()
			log.Println(">>> Failed to calculate next run at:", err)
			return
		}

		if err := tx.Model(&cfscan.ScheduleTask{}).Where("id = ?", task.ID).Update("next_run_at", nextRunAt).Error; err != nil {
			tx.Rollback()
			log.Println(">>> Failed to update record:", err)
			return
		}
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Println(">>> Failed to commit update scheduleTaskHist and submit scan status:", err)
		tx.Rollback()
		return
	}
	log.Printf(">>> 服务重启成功,已将进行中的的定时任务记录设置为超时,数量: %d\n", len(taskHists))

	var scheduleTasks []cfscan.ScheduleTask
	db := global.GVA_DB.Model(&cfscan.ScheduleTask{})
	db = db.Where("enable = ?", 1)
	result := db.Find(&scheduleTasks)
	if result.RowsAffected == 0 {
		global.GVA_LOG.Info("数据库中没有开启的Schedule Task...")
		return
	}
	configFileManager := NewConfigFileManager(dynamicCronFile)
	var configs []Config
	for _, scheduleTask := range scheduleTasks {
		//marshall the payload json
		var cronTaskPayload CronTaskPayload
		err := json.Unmarshal([]byte(scheduleTask.TaskConfig), &cronTaskPayload)
		if err != nil {
			log.Printf(">>> Unmarshal json to struct failed: %v\n", err)
			continue
		}
		cronTaskPayload.SchedulerTaskID = strconv.Itoa(int(scheduleTask.ID))
		cronTaskPayload.CronTaskStr = scheduleTask.CrontabStr
		cronTaskPayloadStr, _ := json.Marshal(cronTaskPayload)
		//查询出所有开启的schedule task 然后写入dynamic_cron.yml
		c := Config{
			CronId:      scheduleTask.ID,
			Cronspec:    scheduleTask.CrontabStr,
			TaskType:    TypeDistributeCronASNScan,
			TaskPayload: string(cronTaskPayloadStr),
		}
		configs = append(configs, c)
	}
	configFileManager.CreateYAMLFile(configs)
	log.Println(">>> 初始化扫描ASN定时任务完成...")

}

//动态任务管理器

// FileBasedConfigProvider 实现（与之前相同）
// FileBasedConfigProvider implements asynq.PeriodicTaskConfigProvider interface.
type FileBasedConfigProvider struct {
	filename string
}

// GetConfigs Parses the yaml file and return a list of PeriodicTaskConfigs.
func (p *FileBasedConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	data, err := os.ReadFile(p.filename)
	if err != nil {
		return nil, err
	}
	var c PeriodicTaskConfigContainer
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	var configs []*asynq.PeriodicTaskConfig
	for _, cfg := range c.Configs {
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: cfg.Cronspec, Task: asynq.NewTask(cfg.TaskType, []byte(cfg.TaskPayload), asynq.Queue("cron-scan"))})
	}
	return configs, nil
}

type PeriodicTaskConfigContainer struct {
	Configs []*Config `yaml:"configs"`
}

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

// TypeUpdateASNInfoCIDR custom self run task here
const TypeUpdateASNInfoCIDR = "self-admin:updateCIDR"

type UpdateASNInfoCIDR struct {
	UUID      string     `json:"uuid"`
	TimeStamp *time.Time `json:"time_stamp"`
}

func NewUpdateASNInfoCIDRTask(bashPayload UpdateASNInfoCIDR, queueName string) (*asynq.Task, error) {
	payload, err := json.Marshal(bashPayload)
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	//最大重试三次，超时时间60分钟
	return asynq.NewTask(TypeUpdateASNInfoCIDR, payload, asynq.Queue(queueName)), nil
}

// UpdateCIDRRunProcessor implements asynq.Handler interface.
type UpdateCIDRRunProcessor struct {
}

func (processor *UpdateCIDRRunProcessor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload UpdateASNInfoCIDR
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Doing update CIDR task=%s--%s", payload.UUID, payload.TimeStamp)
	// Do scan stuff
	_, err := DoBashRun(ctx, task, payload)
	if err != nil {
		return fmt.Errorf("failed to process task: %v", err)
	}
	return nil
}

func NewUpdateCIDRRunProcessor() *UpdateCIDRRunProcessor {
	return &UpdateCIDRRunProcessor{}
}

func DoBashRun(ctx context.Context, task *asynq.Task, payload UpdateASNInfoCIDR) ([]byte, error) {
	var asnInfos []cfscan.AsnInfo
	db := global.GVA_DB.Model(&cfscan.AsnInfo{})
	db = db.Where("enable = ?", 1)
	result := db.Find(&asnInfos)
	if result.RowsAffected == 0 {
		global.GVA_LOG.Info("数据库中没有开启的ASN号码或不存在ASN信息记录")
		return nil, nil
	}
	for _, asnInfo := range asnInfos {
		asn, err := cfscan_utils.GetCIDRByASN(asnInfo.AsnName)
		if err != nil {
			global.GVA_LOG.Error("Error on getting CIDR data of: ", zap.Error(err))
			continue
		}
		cidrData := strings.Join(asn, "\n")
		asnInfo.Ipv4CIDR = cidrData
		now := time.Now()
		asnInfo.LastCIDRUpdate = &now
		if err = global.GVA_DB.Model(&cfscan.AsnInfo{}).Where("id = ?", asnInfo.ID).Updates(&asnInfo).Error; err != nil {
			global.GVA_LOG.Error("update CIDR data failed: ", zap.Error(err))
			continue
		}

	}
	return nil, nil
}

// TypeDistributeCronASNScan 定时任务分发 任务
const TypeDistributeCronASNScan = "cron-scan:cron-scan-asn"

// DistributeCronASNScanProcessor implements asynq.Handler interface.
type DistributeCronASNScanProcessor struct {
}

func (processor *DistributeCronASNScanProcessor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("定时任务: %s触发,接收到定时任务参数: %s \n", task.Type(), string(task.Payload()))
	ctx, cancel := context.WithTimeout(astkCtx, 12*time.Hour)

	var cronTaskPayload CronTaskPayload
	//payload to struct
	_ = json.Unmarshal(task.Payload(), &cronTaskPayload)
	cronTaskId := cronTaskPayload.SchedulerTaskID
	cronTaskIdInt, _ := strconv.Atoi(cronTaskId)
	//如果此时有cronTaskID 相同 但是未完成的扫描 直接返回
	var taskHist cfscan.ScheduleTaskHist

	result := global.GVA_DB.Model(&cfscan.ScheduleTaskHist{}).
		Where("schedule_task_id = ? and hist_status = ?", cronTaskIdInt, "1").
		Order("id desc").First(&taskHist)
	if result.RowsAffected != 0 {
		log.Println("当前定时任务ID存在未完成的任务调度记录,请仔细检查...")
		return nil
	}
	//修改schedule task 任务状态为进行中 并且将下一次运行的时间更新上
	// 开始事务
	tx := global.GVA_DB.Begin()

	// 确保事务在执行完后关闭
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	nextRunTime, _ := NextRunTime(cronTaskPayload.CronTaskStr)
	updateTask := cfscan.ScheduleTask{
		TaskStatus: "3",
		NextRunAt:  nextRunTime,
	}
	// 更新操作
	if err := tx.Model(&cfscan.ScheduleTask{}).Where("id = ?", cronTaskId).Updates(&updateTask).Error; err != nil {
		tx.Rollback()
		log.Fatalf(">>> Update ScheduleTask failed: %v", err)
	}

	// 插入操作
	newScheduleTaskHist := cfscan.ScheduleTaskHist{
		ASNName:        cronTaskPayload.AsnNumber,
		ScheduleTaskId: cronTaskIdInt,
		StartTime:      time.Now(),
		HistStatus:     "1", //运行中
	}
	if err := tx.Model(&cfscan.ScheduleTaskHist{}).Create(&newScheduleTaskHist).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Insert failed: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Fatalf(">>> Transaction commit failed: %v", err)
	}
	log.Println(">>> Transaction(update scheduleTask and insert scheduleTaskHist) completed successfully")

	//分发ASN扫描任务 等待结果 更新数据
	var asnInfo cfscan.AsnInfo
	tx = global.GVA_DB.Where("asn_name = ?", cronTaskPayload.AsnNumber).Find(&asnInfo)
	if tx.RowsAffected == 0 {
		log.Printf(">>> 当前ASN编号不在数据库中,请确保ASN编号已经存在ASNInfo表中")
		cronTaskIdUInt := uint(cronTaskIdInt)
		err := updateScheduleTaskStatusAndTime(global.GVA_DB, cronTaskIdUInt, cronTaskPayload.CronTaskStr)
		log.Printf(">>> 更新定时任务和定时任务记录状态失败: %s\n", err)
		return err
	}
	waitForProcess := strings.Split(asnInfo.Ipv4CIDR, "\n")
	//这里需要限制 最大批量数为25
	batchedCIDRS := cfscan_utils.SplitCIDRs(waitForProcess, cronTaskPayload.IpbatchSize)

	taskBatchSize := len(batchedCIDRS)
	log.Printf("本次扫描任务将分成: %d 个小扫描任务\n", taskBatchSize)

	//time.Sleep(5 * time.Second) 模拟运行完
	//重新设置scheduleTask状态为1,并更新最后运行时间 下一次运行时间
	//查找最新一条 Hist记录 并且更新scheduleTaskHist记录为完成(2)

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
				if cronTaskPayload.EnableSpeedtest == "0" {
					enableSpeedtest = 0
				}
				payload := tasks.ASNScanCFPayload{
					AsnNumber:       cronTaskPayload.AsnNumber,
					EnableTls:       cronTaskPayload.EnableTls,
					ScanPorts:       cronTaskPayload.ScanPorts,
					ScanRate:        cronTaskPayload.ScanRate,
					IpcheckThread:   cronTaskPayload.IpcheckThread,
					EnableSpeedtest: enableSpeedtest,
					CIDRList:        cidr,
					IPBatchSize:     cronTaskPayload.IpbatchSize,
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
				cronTaskIdUInt := uint(cronTaskIdInt)

				err := updateScheduleTaskStatusAndTime(global.GVA_DB, cronTaskIdUInt, cronTaskPayload.CronTaskStr)
				if err != nil {
					log.Printf("Error on update ScheduleTask status and time")
				}

				log.Printf("Batch Task cancled due to the timeout on waiting result: ScheduleTask: %d", cronTaskIdUInt)
				return
			case result, ok := <-results:
				if !ok {
					// 通道已关闭，所有结果已收集完毕
					// ... (保持原有的结果处理逻辑不变)
					//convert sub json list to big json
					// 创建一个切片来存储去掉方括号的 JSON 对象
					cleanResultJson, err2 := cfscan_utils.CleanResultJsonE(finalResult)
					if err2 != nil {
						log.Printf("Clean Result Json Error: %s", err2)
						log.Printf("Schedule Task finished: %d, %s\n", cronTaskPayload.SchedulerTaskID, cronTaskPayload.CronTaskStr)
						return
					}
					// 开始事务
					tx := global.GVA_DB.Begin()
					// 更新scheduleTask 状态
					var task cfscan.ScheduleTask
					if err := tx.Model(&cfscan.ScheduleTask{}).First(&task, cronTaskIdInt).Error; err != nil {
						tx.Rollback()
						return
					}

					// 设置新的LastRunAt和NextRunAt时间
					nextRunAt, _ := NextRunTime(cronTaskPayload.CronTaskStr)

					// 更新字段
					task.TaskStatus = "1"
					task.LastRunAt = task.NextRunAt
					task.NextRunAt = nextRunAt

					if err := tx.Model(&cfscan.ScheduleTask{}).Where("id = ?", task.ID).Updates(&task).Error; err != nil {
						tx.Rollback()
						return
					}
					//更新ScheduleTaskHist
					//更新scheduleTaskHist
					// 查询 ScheduleTaskId 对应的最新一条记录
					var taskHist cfscan.ScheduleTaskHist
					cronTaskIdUInt := uint(cronTaskIdInt)
					err := tx.Model(&cfscan.ScheduleTaskHist{}).Where("schedule_task_id = ?", cronTaskIdUInt).
						Order("id DESC").
						First(&taskHist).Error
					if err != nil {
						tx.Rollback()
						return
					}

					// 更新 HistStatus 为 "3"，设置 EndTime 并计算 CostTime
					now := time.Now()
					taskHist.HistStatus = "2"
					taskHist.EndTime = now
					taskHist.CostTime = int(now.Sub(taskHist.StartTime).Seconds())
					taskHist.TaskResult = cleanResultJson

					// 保存更新
					err = tx.Model(&cfscan.ScheduleTaskHist{}).Where("id = ?", taskHist.ID).Updates(&taskHist).Error
					if err != nil {
						tx.Rollback()
						return
					}
					// 提交事务
					if err := tx.Commit().Error; err != nil {
						fmt.Printf(">>> Failed to commit ScheduleTask transaction: %s", err)
						return
					}

					mergedJSON := cleanResultJson

					//保存到反代IP数据库
					var speedTestResultES []cfscan.SpeedTestResultE
					err2 = json.Unmarshal([]byte(mergedJSON), &speedTestResultES)
					if err != nil {
						log.Printf("Error on unmarshal merged result json: %s\n", err.Error())
					}
					//append to proxy_ips
					var proxy_ips []cfscan.ProxyIps
					for _, r := range speedTestResultES {
						enableTLS := "0"
						if r.EnableTLS {
							enableTLS = "1"
						}
						p := cfscan.ProxyIps{
							AsnNumber:     r.AsnNumber,
							Ip:            r.IP,
							Port:          r.Port,
							EnableTls:     enableTLS,
							DataCenter:    r.DataCenter,
							Region:        r.Region,
							City:          r.City,
							Latency:       r.Latency,
							TcpDuration:   strconv.FormatInt(int64(r.TcpDuration), 10),
							DownloadSpeed: r.DownloadSpeed,
						}
						proxy_ips = append(proxy_ips, p)
					}
					// 使用 ON CONFLICT DO NOTHING 跳过冲突的记录
					if err := global.GVA_DB.Model(&cfscan.ProxyIps{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&proxy_ips).Error; err != nil {
						log.Println("Failed to insert asn scan result to proxy_ips:", err)
					}
					log.Printf("Schedule Task finished: %d, %s\n", cronTaskPayload.SchedulerTaskID, cronTaskPayload.CronTaskStr)

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

func updateScheduleTaskStatusAndTime(db *gorm.DB, cronTaskID uint, cronStr string) error {
	// 开始事务
	tx := db.Begin()

	var task cfscan.ScheduleTask
	if err := tx.Model(&cfscan.ScheduleTask{}).First(&task, cronTaskID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf(">>> Failed to find ScheduleTask record: %w", err)
	}

	// 设置新的LastRunAt和NextRunAt时间
	now := time.Now()
	nextRunAt, err := NextRunTime(cronStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf(">>> Failed to calculate next run at: %w", err)
	}

	// 更新字段
	task.TaskStatus = "1"
	task.LastRunAt = now
	task.NextRunAt = nextRunAt

	if err := tx.Model(&cfscan.ScheduleTask{}).Where("id = ?", task.ID).Updates(&task).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf(">>> Failed to update ScheduleTask record: %w", err)
	}

	//更新scheduleTaskHist
	// 查询 ScheduleTaskId 对应的最新一条记录
	var taskHist cfscan.ScheduleTaskHist
	err = db.Model(&cfscan.ScheduleTaskHist{}).Where("schedule_task_id = ?", cronTaskID).
		Order("id DESC").
		First(&taskHist).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("查询定时任务历史错误: %v", err)
	}

	// 更新 HistStatus 为 "3"，设置 EndTime 并计算 CostTime
	now = time.Now()
	taskHist.HistStatus = "3"
	taskHist.EndTime = now
	taskHist.CostTime = int(now.Sub(taskHist.StartTime).Seconds())

	// 保存更新
	err = db.Model(&cfscan.ScheduleTaskHist{}).Save(&taskHist).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新定时任务历史错误: %v", err)
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf(">>> Failed to commit ScheduleTask transaction: %w", err)
	}

	return nil
}

func NewDistributeCronASNScanProcessor() *DistributeCronASNScanProcessor {
	return &DistributeCronASNScanProcessor{}
}
