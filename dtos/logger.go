package dtos

type GraylogMessage struct {
	Version     string
	Host        string
	ShortMessage string
	FullMessage string
	Environment string
	ServiceName string
	Level       string
}