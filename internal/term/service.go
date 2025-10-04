package term

import (
	"classroom-service/pkg/constants"
	"classroom-service/pkg/consul"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type TermService interface {
	GetTermByID(ctx context.Context, id string) (*TermInfor, error)
}

type termService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "term-service"
)

func NewTermService(client *api.Client) TermService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &termService{
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

func (s *termService) GetTermByID(ctx context.Context, id string) (*TermInfor, error) {
	
	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := s.client.getTermByID(token, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("term not found")
	}

	term := &TermInfor{
		ID:        data["id"].(string),
		StartDate: data["start_date"].(string),
		EndDate:   data["end_date"].(string),
	}

	return term, nil
}

func (c *callAPI) getTermByID(token, id string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/api/v1/gateway/terms/%s", id)

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
