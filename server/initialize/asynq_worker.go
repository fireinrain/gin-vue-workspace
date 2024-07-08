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
		provider := &FileBasedConfigProvider{filename: "dynamic_cron.yml"}

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
	CronId      string `json:"cron_id"`
	Cronspec    string `yaml:"cronspec"`
	TaskType    string `yaml:"task_type"`
	TaskPayload string `yaml:"task_payload"`
}

//custom self run task here

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
