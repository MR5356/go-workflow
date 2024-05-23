package main

import (
	"github.com/MR5356/go-workflow"
	"math/rand"
	"strconv"
	"time"
)

type SleepTask struct {
	sleep int

	workflow.UnimplementedITask
}

func (t *SleepTask) SetParams(params *workflow.TaskParams) error {
	sleepStr := params.Get("sleep")
	sleep, err := strconv.Atoi(sleepStr)
	if err != nil {
		sleep = rand.Intn(5) + 1
		workflow.Logger.Info("invalid sleep time, use default value: " + strconv.Itoa(sleep))
	}

	t.sleep = sleep
	return nil
}

func (t *SleepTask) Start() error {
	workflow.Logger.Info("sleep " + string(rune(t.sleep)))
	time.Sleep(time.Second * time.Duration(t.sleep))
	return nil
}

func main() {
	workflow.Serve(&SleepTask{})
}
