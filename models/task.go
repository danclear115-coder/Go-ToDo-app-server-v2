package models

type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Priority string `json:"priority"`
	IsCompleted bool `json:"is_comlpeted"`
	UserId int `json:"user_id"`
}

type CreateTaskRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Title string `json:"title"`
	Content string `json:"content"`
	Priority string `json:"priority"`
}

type GetTasksRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CompleteTaskRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id 	int	`json:"id"`
}

type ChangeTaskRequest struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Priority string `json:"priority"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ModifiedTask struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type DeleteTaskRequest struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
