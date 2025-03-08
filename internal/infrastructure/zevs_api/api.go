package zevs_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/url"
	"os"
	"zevsbot/internal/application/services"
	"zevsbot/internal/domain/entities"
	"zevsbot/internal/infrastructure/messages"
	"zevsbot/internal/infrastructure/zevs_api/dto"
)

type zevsApi struct {
	authUrl   string
	searchUrl string
	ctx       context.Context
}

func Init(ctx context.Context) services.ZevsApi {
	return &zevsApi{
		authUrl:   "",
		searchUrl: "",
		ctx:       ctx,
	}
}

func (za zevsApi) Auth(email string, pass []byte) (*entities.User, error) {
	var webAuthResp dto.WebAuthResp
	requestURL := os.Getenv("ZEVS_WEBSITE_URL_AUTH")
	req, err := http.NewRequestWithContext(za.ctx, http.MethodPost, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = url.Values{
		"email":    {email},
		"password": {string(pass)},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(messages.GetMessage("authorization_failed"))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &webAuthResp)
	if err != nil {
		return nil, err
	}

	user := webAuthResp.ToEntityUser()
	return user, nil
}

func (za zevsApi) Logout(email string) error {
	//TODO
	// delete token from db
	return nil
}

func (za zevsApi) SearchRemote(token, query string) ([]entities.Item, error) {

	request := &dto.WebSearchReq{
		Token: token,
		Data: struct {
			Query string `json:"query"`
		}{Query: query},
	}
	requestJSONData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var webSearchResp dto.WebSearchResp

	requestURL := os.Getenv("ZEVS_WEBSITE_URL_SEARCH")
	req, err := http.NewRequestWithContext(za.ctx, http.MethodPost, requestURL, bytes.NewReader(requestJSONData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			log.Error().Msg(err.Error())
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var failedAnswer *json.UnmarshalTypeError
	err = json.Unmarshal(resBody, &webSearchResp)
	if err != nil {
		if errors.As(err, &failedAnswer) {
			return nil, nil
		}
		return nil, err
	}
	var searchRes []entities.Item
	for key, val := range webSearchResp {
		searchRes = append(searchRes, val.ToEntityItem(key))
	}

	return searchRes, nil
}
