package models

type HomeMessage struct {
	Message string `json:"message"`
}

func GetHomeMessage() []HomeMessage {
	return []HomeMessage{
		{Message: "Backend is running!"},
	}
}
