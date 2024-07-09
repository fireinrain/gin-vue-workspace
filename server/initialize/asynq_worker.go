package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cfscan"
	cfscan_utils "github.com/flipped-aurora/gin-vue-admin/server/utils/cfscan"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
	"time"
)

const redisAddr = "cloud2.131433.xyz:5379"
const redisPass = "fireinrain@redis"
const redisDb = 0
const dynamicCronFile = "dynamic_cron.yml"

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

func StartAllCronTaskFromDB() {
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
		//查询出所有开启的schedule task 然后写入dynamic_cron.yml
		c := Config{
			CronId:      scheduleTask.ID,
			Cronspec:    scheduleTask.CrontabStr,
			TaskType:    TypeDistributeCronASNScan,
			TaskPayload: scheduleTask.TaskConfig,
		}
		configs = append(configs, c)
	}
	configFileManager.CreateYAMLFile(configs)
	global.GVA_LOG.Info("初始化扫描ASN定时任务完成...")
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
	//var payload UpdateASNInfoCIDR
	//if err := json.Unmarshal(task.Payload(), &payload); err != nil {
	//	return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	//}
	//
	//log.Printf("Doing update CIDR task=%s--%s", payload.UUID, payload.TimeStamp)
	//// Do scan stuff
	//_, err := DoBashRun(ctx, task, payload)
	//if err != nil {
	//	return fmt.Errorf("failed to process task: %v", err)
	//}
	log.Println("接收到定时任务参数: ", string(task.Payload()))
	return nil
}

func NewDistributeCronASNScanProcessor() *DistributeCronASNScanProcessor {
	return &DistributeCronASNScanProcessor{}
}
