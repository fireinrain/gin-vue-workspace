package initialize

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

// NextRunTime 计算并返回当前时间加上 cron 表达式表示的时间
func NextRunTime(cronExpr string) (time.Time, error) {
	// 创建一个新的 cron 解析器
	// 六位crontab表达式
	//parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

	//五位cron 表达式
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	// 解析 cron 表达式
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("无效的 cron 表达式: %v", err)
	}

	// 获取当前时间
	now := time.Now()

	// 计算下次运行时间
	next := schedule.Next(now)

	// 计算时间差
	duration := next.Sub(now)

	// 返回当前时间加上时间差
	return now.Add(duration), nil
}

// PreviousRunTime 计算并返回当前时间减去 cron 表达式表示的时间间隔
func PreviousRunTime(cronExpr string) (time.Time, error) {
	// 创建一个新的 cron 解析器
	//parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// 解析 cron 表达式
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("无效的 cron 表达式: %v", err)
	}

	// 获取当前时间
	now := time.Now()

	// 找到上一次运行时间
	var prev time.Time
	for t := now.Add(-time.Hour * 24); t.Before(now); t = t.Add(time.Second) {
		next := schedule.Next(t)
		if next.After(now) {
			prev = t
			break
		}
	}

	if prev.IsZero() {
		return time.Time{}, fmt.Errorf("无法确定上一次运行时间")
	}

	return prev, nil
}
