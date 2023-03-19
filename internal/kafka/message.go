package kafka

type Headers map[string]string

type Message struct {
	Partition int
	Offset    int64
	Headers   Headers

	Key  string
	Body []byte
}
