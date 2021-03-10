package alfred

func sendChatResponse(msg chan *Message, id int) {
	msg <- formatChatMessage("ChatChannel", id)
}
