package http

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

// Error implements go error interface.
func (msg *ErrorMessage) Error() string {
	return fmt.Sprintf("API Error: %s", msg.Message)
}

type UserClient struct {
	client *req.Client
}

func NewUserHTTPClient(log *zap.Logger) *UserClient {
	return &UserClient{
		client: req.C().
			SetBaseURL("https://api.github.com").
			SetCommonErrorResult(&ErrorMessage{}).
			EnableDumpEachRequest().
			OnAfterResponse(func(client *req.Client, resp *req.Response) error {
				if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
					return nil
				}
				if errMsg, ok := resp.ErrorResult().(*ErrorMessage); ok {
					resp.Err = errMsg // Convert api error into go error
					return nil
				}
				if !resp.IsSuccessState() {
					// Neither a success response nor a error response, record details to help troubleshooting
					resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
				}
				return nil
			}),
	}
}

type UserProfile struct {
	Name string `json:"name"`
	Blog string `json:"blog"`
}

// GetUserProfile_Style1 returns the user profile for the specified user.
// Github API doc: https://docs.github.com/en/rest/users/users#get-a-user
func (u *UserClient) GetUserInfoByPwd(ctx context.Context, username, password string) (user *UserProfile, err error) {

	//组装参数
	formData := map[string]string{"username": username, "pwd": password}

	// 请求
	_, err = u.client.R().
		SetContext(ctx).
		SetFormData(formData).
		SetSuccessResult(&user).
		Post("/users/pwd")

	return
}
