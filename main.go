package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"math/rand"
	"time"
)

const maxRetry = 3

func main() {
	ctx := context.Background()
	// クライアントオブジェクトを作って
	client, err := makeClient(ctx)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(client)

	// 送信内容を作って
	message := makeMessage("", "Hello")

	fmt.Println(message)
	// 指数バックオフで送信。とりあえず3回くらいはリトライしてみる
	var response string
	for i := 0; i < maxRetry; i++ {
		response, err = client.Send(ctx, message)
		if err == nil {
			break
		}
		// 待ち時間はどんどん倍にしていく
		waitTime := (1<<i)*1000 + rand.Intn(1000)
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}
	if err != nil {
		fmt.Printf("Failed to send message: %s", err)
		return
	}
	fmt.Printf("response=%s\n", response)
}

func makeClient(ctx context.Context) (*messaging.Client, error) {
	opt := option.WithCredentialsFile(".json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	return app.Messaging(ctx)
}

func makeMessage(token, body string) *messaging.Message {
	return &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Body: body,
		},
		Android: &messaging.AndroidConfig{},
		APNS:    &messaging.APNSConfig{},
		Webpush: &messaging.WebpushConfig{},
	}
}
