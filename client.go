package wkcfirebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Client struct {
	ConfigPath string
	Logger     *log.Logger
	Context    context.Context
	fsclient   *firestore.Client
}

func (c *Client) Setup() {
	opt := option.WithCredentialsFile(c.ConfigPath)
	app, err := firebase.NewApp(c.Context, nil, opt)
	if err != nil {
		c.Logger.Fatal(fmt.Sprintf("error initializing firebase: %v", err))
	}

	c.fsclient, err = app.Firestore(c.Context)

	if err != nil {
		c.Logger.Fatal(err)
	}
}

func (c *Client) TearDown() {
	err := c.fsclient.Close()
	if err != nil {
		c.Logger.Fatal(err)
		return
	}
}

func (c *Client) SetupStateChangeListener(listener StateChangedListener) {
	snapshotIteration := c.fsclient.Collection("State").Doc("State").Snapshots(context.Background())
	c.Logger.Println("snapshot listener setup")
	for {
		snap, err := snapshotIteration.Next()

		if status.Code(err) == codes.DeadlineExceeded {
			c.Logger.Println("Snapshot deadlineExceeded")
			return
		}

		if err != nil {
			c.Logger.Printf("Snapshot.Next error: %s\n", err.Error())
			return
		}

		if !snap.Exists() {
			c.Logger.Fatalln("document nog longer exists")
			return
		}

		s := &StateDocument{}
		e := snap.DataTo(s)

		if e != nil {
			c.Logger.Printf("error parsing snapshot data: %s\n", e.Error())
		}

		// set state
		listener(s)
	}
}

func (c *Client) SetState(state *StateDocument) error {
	_, err := c.fsclient.Collection("State").Doc("State").Set(context.Background(), state)

	if err != nil {
		return err
	}

	return nil
}
