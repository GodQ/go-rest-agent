package agent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godq/go-rest-agent/utils"
)

const (
	AgentFinishFlag = "Agent Command Execution Finished"
)

type Task struct {
	TaskId         int           `json:"task_id"`
	Command        string        `json:"command"`
	TimeoutSeconds int           `json:"timeout_seconds"`
	RetryCount     int           `json:"retry_count"`
	Extra          string        `json:"extra"`
	Key            string        `json:"key"`
	CallbackUrl    string        `json:"callback_url"`
	Status         string        `json:"status"`
	ReturnCode     int           `json:"return_code"`
	ErrorMsg       string        `json:"error_msg"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Duration       time.Duration `json:"duration"`
	LogRedisKey    string        `json:"log_redis_key"`
	LogRedisUrl    string        `json:"log_redis_url"`
	ThreadName     string        `json:"thread_name"`
	Stdout         string        `json:"stdout"`
}

type controller struct {
	TaskList      map[int]*Task
	CurrentTaskId int
}

func NewController() *controller {
	ctl := controller{
		TaskList:      map[int]*Task{},
		CurrentTaskId: 1,
	}
	return &ctl
}

func (ctl *controller) useNewTaskId() int {
	t := ctl.CurrentTaskId
	ctl.CurrentTaskId += 1
	return t
}

func (ctl *controller) addToTaskList(t *Task) (*Task, error) {
	t.TaskId = ctl.useNewTaskId()
	ctl.TaskList[t.TaskId] = t
	t.Status = "pending"
	go ctl.runTask(t)
	return t, nil
}

func (ctl *controller) runTask(t *Task) {
	result := utils.NewCmdResult()
	t.Status = "doing"
	utils.LogInfof("start to run task %d", t.TaskId)
	fmt.Println("start", t.TaskId)
	t.StartTime = time.Now()
	result, err := utils.RunShellWithRetry(t.Command, result, t.TimeoutSeconds, t.RetryCount)
	utils.LogInfof("end to run task %d", t.TaskId)
	fmt.Println("end", t.TaskId)
	t.EndTime = time.Now()
	t.Duration = time.Duration(t.EndTime.Sub(t.StartTime).Seconds())
	if err != nil {
		t.Status = "error"
		t.ErrorMsg = err.Error()
	}
	if result.Code != 0 {
		t.Status = "failed"
	} else {
		t.Status = "done"
	}
	t.ReturnCode = result.Code
	t.Stdout = result.Output()
}

func (ctl *controller) CreateTask(c *gin.Context) {
	var task Task

	if err := c.Bind(&task); err == nil {
		t, err := ctl.addToTaskList(&task)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, t)
		return
	} else {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func (ctl *controller) ListTasks(c *gin.Context) {
	task_id_str := c.Query("task_id")
	if task_id_str == "" {
		c.JSON(200, ctl.TaskList)
		return
	} else {
		task_id, err := strconv.Atoi(task_id_str)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		task, found := ctl.getTask(task_id)
		if !found {
			c.JSON(404, gin.H{"error": "task id not found " + task_id_str})
			return
		}
		resp := map[int]*Task{
			task.TaskId: task,
		}
		c.JSON(200, resp)
	}
}

func (ctl *controller) getTask(task_id int) (*Task, bool) {
	task, found := ctl.TaskList[task_id]
	return task, found
}

func (ctl *controller) GetTask(c *gin.Context) {
	task_id_str := c.Param("task_id")
	task_id, err := strconv.Atoi(task_id_str)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	task, found := ctl.getTask(task_id)
	if !found {
		c.JSON(404, gin.H{"error": "task id not found " + task_id_str})
		return
	}
	c.JSON(200, task)
}

func (ctl *controller) SendFile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Create Task",
	})
}

func (ctl *controller) GetFile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get File",
	})
}
