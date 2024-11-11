package gov

type Entry struct {
	Message string `json:"message"`
}

func GetForeignAidData() Entry {
	return Entry{Message: "hello"}
}
