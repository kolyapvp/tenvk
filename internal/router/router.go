package router

import (
	"strconv"
	"strings"
	"taskbot/internal/model"
	"taskbot/internal/repository"
	"taskbot/internal/service"
)

func Route(text string, userID int64, userName string) map[int64]string {
	repository.AddUser(&model.User{userID, userName})
	switch {
	case text == "/tasks":
		return map[int64]string{userID: service.HandleTasks(userID)}
	case text == "/my":
		return map[int64]string{userID: service.HandleMy(userID)}
	case text == "/owner":
		return map[int64]string{userID: service.HandleOwner(userID)}
	case strings.HasPrefix(text, "/new "):
		// Создаем новую задачу
		title := strings.TrimPrefix(text, "/new ")
		return map[int64]string{userID: service.HandleNew(title, model.User{ID: userID, UserName: userName})}
	case strings.HasPrefix(text, "/assign_"):
		// Назначаем задачу на пользователя
		taskID, err := strconv.ParseInt(strings.TrimPrefix(text, "/assign_"), 10, 64)
		if err != nil {
			return map[int64]string{userID: "Некорректный ID задачи"}
		}
		return service.HandleAssignee(model.User{ID: userID, UserName: userName}, taskID)

	case strings.HasPrefix(text, "/unassign_"):
		// Снимаем задачу с пользователя
		taskID, err := strconv.ParseInt(strings.TrimPrefix(text, "/unassign_"), 10, 64)
		if err != nil {
			return map[int64]string{userID: "Некорректный ID задачи"}
		}
		return service.HandleUnassign(taskID, userID)

	case strings.HasPrefix(text, "/resolve_"):
		// Выполняем и удаляем задачу
		taskID, err := strconv.ParseInt(strings.TrimPrefix(text, "/resolve_"), 10, 64)
		if err != nil {
			return map[int64]string{userID: "Некорректный ID задачи"}
		}
		return service.HandleResolve(taskID, userID)

	default:
		// Неизвестная команда
		return map[int64]string{userID: "Неизвестная команда"}
	}
}
