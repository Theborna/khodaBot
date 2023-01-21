package services

import "github.com/sausheong/goreplicate"

type SDClient struct{ *goreplicate.Client }

const (
	version = "f178fa7a1ae43a9a9af01b833b9d2ecf97b1bcb0acfd2dc5dd04895e042863f1"
	owner   = "stability-ai"
	name    = "stable-diffusion"
)

func NewSDClient(apiKey string) (*SDClient, error) {
	model := goreplicate.NewModel(owner, name, version)
	model.Input["prompt"] = "An astronaut riding a horse in photorealistic style"
	model.Input["num_outputs"] = 4
	client := SDClient{goreplicate.NewClient(apiKey, model)}
	// err := client.Create()
	// if err != nil {
	// 	return &client, err
	// }
	return &client, nil
}
