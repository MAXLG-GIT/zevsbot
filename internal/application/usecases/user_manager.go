package usecases

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"zevsbot/internal/application/app_interfaces"
	"zevsbot/internal/application/services"
	"zevsbot/internal/domain/domain_interfaces"
	"zevsbot/internal/infrastructure/messages"
	"zevsbot/internal/zevs_errors"
)

type userManager struct {
	zevsApi services.ZevsApi
	repo    domain_interfaces.RepoService
}

func InitUserManager(repo domain_interfaces.RepoService, api services.ZevsApi) app_interfaces.UserManager {

	return &userManager{zevsApi: api, repo: repo}
}

func (um userManager) CheckTgUserAuth(id int) (bool, error) {

	token, err := um.repo.Get(strconv.Itoa(id))
	if err != nil {
		log.Error().Msg(err.Error())
		return false, err
	}
	if token == "" {
		return false, nil
	}

	return true, nil
}

func (um userManager) Auth(id int, email string, password []byte) error {
	if email == "" || len(password) < 1 {
		return &zevs_errors.PublicError{Text: messages.GetMessage("incorrect_login_pass")}

	}
	user, err := um.zevsApi.Auth(email, password)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	if user == nil {
		return &zevs_errors.PublicError{Text: messages.GetMessage("user_not_found")}
	}
	err = um.repo.Delete(strconv.Itoa(id))
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	err = um.repo.Save(strconv.Itoa(id), user.Token, user.TokenExp)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	fmt.Println(user)

	return nil
}

func (um userManager) Logout(id int) error {
	err := um.repo.Delete(strconv.Itoa(id))
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}
