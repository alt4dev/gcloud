package service

import (
	"cloud.google.com/go/logging"
	vkit "cloud.google.com/go/logging/apiv2"
	"context"
	"google.golang.org/api/option"
	"strings"
)

var client *logging.Client
var pbClient *vkit.Client

func setupClient() {
	// Don't create a client during testing or silent Modes
	if options.Mode == ModeTesting || options.Mode == ModeSilent {
		return
	}

	var err error
	client, err = logging.NewClient(context.Background(), project)
	if err != nil {
		panic(err)
	}

	client.OnError = func(err error) {
		emitError.Println(err)
	}

	pbClient, err = newProtoClient(context.Background())

	if err != nil {
		panic(err)
	}
}

func makeParent(parent string) string {
	if !strings.ContainsRune(parent, '/') {
		return "projects/" + parent
	}
	return parent
}

// newProtoClient
func newProtoClient(ctx context.Context, opts ...option.ClientOption) (*vkit.Client, error) {
	opts = append([]option.ClientOption{
		option.WithScopes(logging.WriteScope),
	}, opts...)
	c, err := vkit.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	c.SetGoogleClientInfo("alt4.dev", "go-gcloud")

	return c, err
}
