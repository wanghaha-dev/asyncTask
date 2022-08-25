package asyncTask

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Config struct {
	Addr     string
	DB       int
	Password string
}

func NewTask(ctx context.Context, config Config) (*Task, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})

	task := &Task{
		RedisClient: rdb,
		ctx:         ctx,
	}

	return task, nil
}

type Map map[string]interface{}

type Task struct {
	RedisClient *redis.Client
	ctx         context.Context
}

func (t *Task) putTask(taskId string, list string, data Map) error {
	d := struct {
		TaskId string `json:"task_id"`
		Data   Map    `json:"data"`
	}{
		taskId,
		data,
	}

	taskContent, err := json.Marshal(d)
	if err != nil {
		return err
	}

	err = t.RedisClient.LPush(t.ctx, list, taskContent).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) PutNormalTask(taskId string, data Map) error {
	return t.putTask(taskId, "normal", data)
}
func (t *Task) PutSuccessTask(taskId string, data Map) error {
	return t.putTask(taskId, "success", data)
}
func (t *Task) PutFailTask(taskId string, data Map) error {
	return t.putTask(taskId, "fail", data)
}

func (t *Task) takeNormalTask(list string) (Map, error) {
	sliceCmd := t.RedisClient.BRPop(t.ctx, 0, list)

	err := sliceCmd.Err()
	if err != nil {
		return nil, err
	}

	result, err := sliceCmd.Result()
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(result[1]), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *Task) TakeNormalTask() (Map, error) {
	data, err := t.takeNormalTask("normal")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *Task) TakeSuccessTask() (Map, error) {
	data, err := t.takeNormalTask("success")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *Task) TakeFailTask() (Map, error) {
	data, err := t.takeNormalTask("fail")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *Task) getTaskLength(list string) (int64, error) {
	lLen := t.RedisClient.LLen(t.ctx, list)
	err := lLen.Err()

	if err != nil {
		return -1, err
	}

	return lLen.Result()
}

func (t *Task) GetFailTaskLength() (int64, error) {
	return t.getTaskLength("fail")
}

func (t *Task) GetSuccessTaskLength() (int64, error) {
	return t.getTaskLength("success")
}

func (t *Task) GetNormalTaskLength() (int64, error) {
	return t.getTaskLength("normal")
}

func (t *Task) getFailTaskList(list string) ([]string, error) {
	sliceCmd := t.RedisClient.LRange(t.ctx, list, 0, -1)

	err := sliceCmd.Err()
	if err != nil {
		return []string{}, err
	}
	return sliceCmd.Result()
}

func (t *Task) GetFailTaskList() ([]string, error) {
	return t.getFailTaskList("fail")
}

func (t *Task) GetSuccessTaskList() ([]string, error) {
	return t.getFailTaskList("success")
}

func (t *Task) GetNormalTaskList() ([]string, error) {
	return t.getFailTaskList("normal")
}

func EachStrings(data []string) {
	for _, item := range data {
		fmt.Println(item)
	}
}

func Each(count int, f func()) {
	for i := 0; i < count; i++ {
		f()
	}
}

func Datetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
