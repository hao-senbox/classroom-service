package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net/http"
	"os"
	"classroom-service/pkg/constants"
	"classroom-service/pkg/consul"
)

type UserService interface {
	GetUserInfor(ctx context.Context, userID string) (*UserInfor, error)
}

type userService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "go-main-service"
)

func NewUserService(client *api.Client) UserService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &userService{
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

func (u *userService) GetUserInfor(ctx context.Context, userID string) (*UserInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.GetUserInfor(userID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", userID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var roles []Role
	if rolesRaw, ok := innerData["roles"].([]interface{}); ok && len(rolesRaw) > 0 {
		for _, r := range rolesRaw {
			if roleMap, ok := r.(map[string]interface{}); ok {
				roles = append(roles, Role{
					RoleID:   fmt.Sprintf("%v", roleMap["id"]),   
					RoleName: fmt.Sprintf("%v", roleMap["role"]), 
				})
			}
		}
	} else {
		roles = nil 
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["username"]),
		FullName: fmt.Sprintf("%v", innerData["fullname"]),
		Avartar:  fmt.Sprintf("%v", innerData["avatar"]),
		Roles:    roles,
	}, nil
	
}

func safeString(val interface{}) string {
	if val == nil {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func (c *callAPI) GetUserInfor(userID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/user/%s", userID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var userData interface{}

	err = json.Unmarshal([]byte(res), &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := userData.(map[string]interface{})

	return myMap, nil
}

