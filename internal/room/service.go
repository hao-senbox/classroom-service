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
	GetAllRooms(ctx context.Context) ([]*RoomInfor, error)
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

func (s *roomService) GetAllRooms(ctx context.Context) ([]*RoomInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := s.client.getAllRooms(token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no room data found")
	}

	rooms := make([]*RoomInfor, 0, len(data))
	for _, item := range data {
		id, _ := item["id"].(string)
		name, _ := item["name"].(string)

		rooms = append(rooms, &RoomInfor{
			ID:   id,
			Name: name,
		})
	}

	return rooms, nil

}

func (c *callAPI) getAllRooms(token string) ([]map[string]interface{}, error) {

	endpoint := "/api/v1/locations?type=class"

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

	dataListRaw, ok := parse["data"].([]interface{})
	if !ok {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	rooms := make([]map[string]interface{}, 0)

	for _, item := range dataListRaw {
		rooms = append(rooms, item.(map[string]interface{}))
	}

	return rooms, nil

}
