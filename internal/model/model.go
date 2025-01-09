package model

type User struct {
	ID       int64
	UserName string
}

type Task struct {
	ID       int64
	Title    string
	Owner    *User
	Assignee *User
}

// TaskSlice - тип для сортировки задач
type TaskSlice []Task

func (ts TaskSlice) Len() int           { return len(ts) }
func (ts TaskSlice) Less(i, j int) bool { return ts[i].ID < ts[j].ID }
func (ts TaskSlice) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
