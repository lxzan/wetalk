package wetalk

type MessageHeader struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Ack   bool   `json:"ack"`
	From  string `json:"from"`
	To    string `json:"to"`
}

type Message struct {
	Header MessageHeader `json:"header"`
	Data   JSON          `json:"data"`
}

func (this *Message) Reply(data JSON) *Message {
	this.Data = data
	return this
}
