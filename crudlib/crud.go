package crudlib

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SayHello() string {
	return "Hello"
}
