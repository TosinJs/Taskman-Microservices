package firebaseinit

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func Init(ctx context.Context) (*messaging.Client, error) {
	opt := option.WithCredentialsFile("./taskman-firebase-adminsdk.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Println("Unable to Connect To Firebase", err)
		return nil, err
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		fmt.Println("Unable to Connect To Firebase", err)
		return nil, err
	}

	return fcmClient, nil
}
