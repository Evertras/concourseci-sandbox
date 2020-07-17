package crudlib

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestAddsTwoPlusTwo(t *testing.T) {
	result := TwoPlusTwo()
	expected := 4

	if result != expected {
		t.Errorf("Got %d but expected %d", result, expected)
	}
}

func TestAddsDogAndRetrievesAge(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("TEST_INTEGRATION environment variable not set to true")
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatal("testcontainers.GenericContainer():", err)
	}

	defer redisContainer.Terminate(ctx)

	host, err := redisContainer.Host(ctx)

	if err != nil {
		t.Fatal("redisContainer.Host():", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379")

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port.Port()), redis.DialConnectTimeout(time.Second))
		},
	}

	client := NewClient(pool)

	// Approximately 9 years
	expectedAge := time.Hour * 24 * 365 * 9
	dogName := "Genji"

	err = client.CreateDog(dogName, "tan", expectedAge)

	if err != nil {
		t.Fatal("client.CreateDog():", err)
	}

	age, err := client.GetDogAge(dogName)

	if age != expectedAge {
		t.Errorf("Got %v but expected %v", age, expectedAge)
	}
}
