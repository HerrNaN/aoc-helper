package aoc

import "net/http"

type Helper struct {
	client *http.Client
}

func NewHelper() *Helper {
	return &Helper{
		client: http.DefaultClient,
	}
}

func (h *Helper) WithClient(client *http.Client) {
	h.client = client
}
