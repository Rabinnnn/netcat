package utils

// DisplayChats sends chat history to the specified client
func DisplayChats(client *client) {
	mChatHistory.Lock()
	defer mChatHistory.Unlock()

	if len(chatHistory) > 0 {
		for _, chat := range chatHistory {
			client.connection.Write([]byte(chat))
		}
		client.connection.Write([]byte("\n"))
	}
}
