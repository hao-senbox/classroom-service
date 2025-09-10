package room

import (
	"classroom-service/pkg/constants"
	"classroom-service/pkg/consul"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type RoomService interface {
	GetRoomByID(ctx context.Context, id string) (*RoomInfor, error)
}

type roomService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "inventory-service"
)

func NewRoomService(client *api.Client) RoomService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &roomService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		fmt.Printf("Error creating service discovery: %v\n", err)
		return nil
	}

	service, err := sd.DiscoverService()
	if err != nil {
		fmt.Printf("Error discovering service: %v\n", err)
		return nil
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		fmt.Println("Running in LOCAL_TEST mode — overriding service address to localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (s *roomService) GetRoomByID(ctx context.Context, id string) (*RoomInfor, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := s.client.getRoomByID(token, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("room not found")
	}

	room := &RoomInfor{
		ID:   data["id"].(string),
		Name: data["name"].(string),
	}

	return room, nil
}


func (c *callAPI) getRoomByID(token, id string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/api/v1/storage/%s", id)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	response, err := c.client.CallAPI(c.clientServer, endpoint, "GET", nil, headers)
	if err != nil {
		return nil, err
	}

	var parse map[string]interface{}
	if err := json.Unmarshal([]byte(response), &parse); err != nil {
		return nil, err
	}

	dataRaw, ok := parse["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return dataRaw, nil
}

