package repository

import (
	"sync"
	"taskbot/internal/model"
)

type Memory struct {
	sync.RWMutex
	Users map[int64]*model.User
	Tasks map[int64]*model.Task
	ID    int64
}

func AddUser(user *model.User) {
	memory.Lock()
	defer memory.Unlock()
	memory.Users[user.ID] = user
}

func AddTask(task *model.Task) {
	memory.Lock()
	defer memory.Unlock()
	task.ID = memory.ID + 1
	memory.Tasks[memory.ID+1] = task
	memory.ID++
}

func GetUser(id int64) (model.User, bool) {
	memory.RLock()
	defer memory.RUnlock()
	u, ok := memory.Users[id]
	return *u, ok
}

func GetTask(id int64) (model.Task, bool) {
	memory.RLock()
	defer memory.RUnlock()
	t, ok := memory.Tasks[id]
	if !ok {
		return model.Task{}, false
	}
	return *t, ok
}

func GetAllTask() []model.Task {
	memory.RLock()
	defer memory.RUnlock()
	var tasks []model.Task
	for _, task := range memory.Tasks {
		tasks = append(tasks, *task)
	}
	return tasks
}

func DeleteTask(id int64) {
	memory.Lock()
	defer memory.Unlock()
	delete(memory.Tasks, id)
}

func ChangeAssignee(user *model.User, ID int64) bool {
	memory.Lock()
	defer memory.Unlock()
	t, ok := memory.Tasks[ID]
	if !ok {
		return false
	}
	t.Assignee = user
	memory.Tasks[ID] = t
	return true
}

// GetMyTasks Задачи, которые были назначены на меня
func GetMyTasks(userID int64) []model.Task {
	memory.RLock()
	defer memory.RUnlock()
	var tasks []model.Task
	for _, task := range memory.Tasks {
		if task.Assignee == nil {
			continue
		}
		if task.Assignee.ID == userID {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}

// CreateMyTasks задачи созданные мной
func CreateMyTasks(userID int64) []model.Task {
	memory.RLock()
	defer memory.RUnlock()
	var tasks []model.Task
	for _, task := range memory.Tasks {
		if task.Owner.ID == userID {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}

var memory *Memory = &Memory{
	Tasks: make(map[int64]*model.Task, 0),
	Users: make(map[int64]*model.User, 0),
}
