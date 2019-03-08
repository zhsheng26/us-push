package socket

type HandleReceivedMsg func(body string)

type PushMessage struct {
	Content string
	Code    int
}
