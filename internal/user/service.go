package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"classroom-service/pkg/constants"
	"classroom-service/pkg/consul"

	"github.com/hashicorp/consul/api"
)

type UserService interface {
	GetUserInfor(ctx context.Context, userID string) (*UserInfor, error)
	GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetTeacherInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetStaffInfor(ctx context.Context, studentID string) (*UserInfor, error)
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

	data, err := u.client.getUserInfor(userID, token)
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

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil
}

func (u *userService) GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStudentInfor(studentID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", studentID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil

}

func (u *userService) GetTeacherInfor(ctx context.Context, studentID string) (*UserInfor, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getTeacherInfor(studentID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", studentID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil

}

func (u *userService) GetStaffInfor(ctx context.Context, studentID string) (*UserInfor, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStaffInfor(studentID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", studentID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil

}

func (c *callAPI) getUserInfor(userID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/users/%s", userID)

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

func (c *callAPI) getStudentInfor(studentID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/students/%s", studentID)

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

func (c *callAPI) getTeacherInfor(studentID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/teachers/%s", studentID)

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

func (c *callAPI) getStaffInfor(studentID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/staffs/%s", studentID)

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

func castToInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	default:
		return 0
	}
}

func castToBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1"
	default:
		return false
	}
}
