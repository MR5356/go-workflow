package main

import (
	"fmt"
	"github.com/MR5356/go-workflow"
	"math/rand"
	"strconv"
	"time"
)

var logger = workflow.UseLogger("sleep")

type SleepTask struct {
	sleep int

	workflow.UnimplementedITask
}

func (t *SleepTask) SetParams(params *workflow.TaskParams) error {
	sleepStr := params.Get("sleep")
	sleep, err := strconv.Atoi(sleepStr)
	if err != nil {
		sleep = rand.Intn(5) + 1
		logger.Info("invalid sleep time, use default value: " + strconv.Itoa(sleep))
	}

	t.sleep = sleep
	return nil
}

func (t *SleepTask) Run() error {
	logger.Info(fmt.Sprintf("sleep %d seconds", t.sleep))
	time.Sleep(time.Second * time.Duration(t.sleep))
	return nil
}

func main() {
	workflow.Serve(&SleepTask{})
}
