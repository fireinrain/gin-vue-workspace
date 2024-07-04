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
	"log"
	"strings"
	"time"
)

func SelfAsynQTaskClientRun() {
	const redisAddr = "cloud2.131433.xyz:5379"

	go func() {
		//scheduler
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}
		scheduler := asynq.NewScheduler(
			asynq.RedisClientOpt{Addr: redisAddr, Password: "fireinrain@redis", DB: 0},

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
		log.Printf("registered an scheduler entry: %q\n", entryID)
		if err := scheduler.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {

		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: redisAddr, Password: "fireinrain@redis", DB: 0},
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: 3,
				// Optionally specify multiple queues with different priority.
				Queues: map[string]int{
					"self-admin": 6,
				},
				// See the godoc for other configuration options
			},
		)
		mux := asynq.NewServeMux()
		mux.Handle(TypeUpdateASNInfoCIDR, NewBashTaskRunProcessor())
		log.Println(">>> 注册AsynQ-Wokrer成功...")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

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

func NewBashTaskRunProcessor() *UpdateCIDRRunProcessor {
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
