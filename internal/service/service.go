package service

import (
	"fmt"
	"sort"
	"strings"
	"taskbot/internal/model"
	"taskbot/internal/repository"
)

func HandleTasks(userID int64) string {
	tasks := repository.GetAllTask()
	var builder strings.Builder
	if len(tasks) == 0 {
		return "Нет задач"
	}
	sort.Sort(model.TaskSlice(tasks))
	for _, task := range tasks {
		if task.Assignee != nil && task.Assignee.ID == userID {
			builder.WriteString(fmt.Sprintf("%d. %s by @%s\nassignee: я\n/unassign_%d /resolve_%d\n\n",
				task.ID, task.Title, task.Owner.UserName, task.ID, task.ID))
		} else if task.Assignee != nil {
			builder.WriteString(fmt.Sprintf("%d. %s by @%s\nassignee: @%s\n\n",
				task.ID, task.Title, task.Owner.UserName, task.Assignee.UserName))
		} else {
			//fmt.Println(task.Owner, task.Assignee)
			builder.WriteString(fmt.Sprintf("%d. %s by @%s\n/assign_%d\n\n",
				task.ID, task.Title, task.Owner.UserName, task.ID))
		}
	}
	return strings.TrimSpace(builder.String())
}

func HandleNew(title string, owner model.User) string {
	t := &model.Task{
		Title: title,
		Owner: &owner,
	}
	repository.AddTask(t)
	return fmt.Sprintf("Задача %q создана, id=%d", title, t.ID)
}

func HandleAssignee(newAssignee model.User, taskID int64) map[int64]string {
	task, ok := repository.GetTask(taskID)
	if !ok {
		return map[int64]string{newAssignee.ID: "Задача не найдена"}
	}
	prefAssignee := task.Assignee
	task.Assignee = &newAssignee
	ok = repository.ChangeAssignee(&newAssignee, taskID)
	if !ok {
		return map[int64]string{newAssignee.ID: "Задача не найдена"}
	}
	res := map[int64]string{newAssignee.ID: "Задача \"" + task.Title + "\" назначена на вас"}
	if prefAssignee != nil {
		res[prefAssignee.ID] = fmt.Sprintf("Задача %q назначена на @%s", task.Title, newAssignee.UserName)
	} else if task.Owner.ID != newAssignee.ID {
		res[task.Owner.ID] = fmt.Sprintf("Задача %q назначена на @%s", task.Title, newAssignee.UserName)
	}
	return res
}

func HandleUnassign(taskID int64, userID int64) map[int64]string {
	task, ok := repository.GetTask(taskID)
	if !ok {
		return map[int64]string{userID: "Задача не найдена"}
	}
	if task.Assignee == nil {
		return map[int64]string{userID: "Задача не назначена"}
	}
	if task.Assignee.ID != userID {
		return map[int64]string{userID: "Задача не на вас"}
	}
	ok = repository.ChangeAssignee(nil, task.ID)
	if !ok {
		return map[int64]string{userID: "Задача не найдена"}
	}
	return map[int64]string{
		userID:        "Принято",
		task.Owner.ID: fmt.Sprintf("Задача %q осталась без исполнителя", task.Title),
	}
}

func HandleResolve(taskID int64, userID int64) map[int64]string {
	task, ok := repository.GetTask(taskID)
	if !ok {
		return map[int64]string{userID: "Задача не найдена"}
	}
	if task.Assignee == nil {
		return map[int64]string{userID: "Задача не назначена"}
	}
	if task.Assignee.ID != userID {
		return map[int64]string{userID: "Задача не на вас"}
	}
	repository.DeleteTask(task.ID)
	return map[int64]string{
		userID:        fmt.Sprintf("Задача %q выполнена", task.Title),
		task.Owner.ID: fmt.Sprintf("Задача %q выполнена @%s", task.Title, task.Assignee.UserName),
	}
}

func HandleMy(userID int64) string {
	tasks := repository.GetMyTasks(userID)
	if len(tasks) == 0 {
		return "Нет задач"
	}
	sort.Sort(model.TaskSlice(tasks))
	var builder strings.Builder
	for _, task := range tasks {
		builder.WriteString(fmt.Sprintf("%d. %s by @%s\n/unassign_%d /resolve_%d\n\n",
			task.ID, task.Title, task.Owner.UserName, task.ID, task.ID))
	}
	return strings.TrimSpace(builder.String())
}

func HandleOwner(userID int64) string {
	tasks := repository.CreateMyTasks(userID)
	if len(tasks) == 0 {
		return "Нет задач"
	}
	sort.Sort(model.TaskSlice(tasks))
	var builder strings.Builder
	for _, task := range tasks {
		builder.WriteString(fmt.Sprintf("%d. %s by @%s\n/assign_%d\n\n",
			task.ID, task.Title, task.Owner.UserName, task.ID))
	}
	return strings.TrimSpace(builder.String())

}
