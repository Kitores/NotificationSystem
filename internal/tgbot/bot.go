package tgbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendNotificationToBot(text string, sub int32) error {
	//for _, sub := range subs {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TOKEN"))
	message := map[string]interface{}{
		"chat_id": sub,
		"text":    text,
	}
	fmt.Println(message, url)
	jsonDate, err := json.Marshal(message)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonDate))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}
	//	}
	return nil
}
