package userService

import (
	"fmt"
	"net/http"
	"togrpc/constans"
	"togrpc/helpers"
	"togrpc/models"
	"togrpc/services"
	"togrpc/utils"

	"github.com/labstack/echo"
)

type loginService struct {
	Service services.UsecaseService
}

func NewLoginService(service services.UsecaseService) loginService {
	return loginService{
		Service: service,
	}
}

func (svc loginService) AuthLogin(ctx echo.Context) error {
	var result models.Response
	var response models.ResponseAuth
	var rolesName string = constans.EMPTY_VALUE
	request := new(models.RequestAuthLogin)
	if err := helpers.BindValidateStruct(ctx, request); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	dataLogin, err := svc.Service.UserMongoRepo.FindUserByIndex(request.Username)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		fmt.Println("Result Request FindUserByIndex ", utils.ToString(result))
		return ctx.JSON(http.StatusBadRequest, result)
	}

	token := models.RequestAuthDashboard{
		Username: request.Username,
		Password: request.Password,
	}

	// check pass
	match := utils.CheckPasswordHash(token.Password, dataLogin.User.Password)
	if !match {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Wrong Password", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	if dataLogin.User.Active != constans.YES {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "User not actived", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	if dataLogin.User.RolesName != constans.EMPTY_VALUE {
		rolesName = dataLogin.User.RolesName
	}

	resultDevice, err := svc.Service.UserMongoRepo.FindDeviceIdByIndex(request.Username, request.DeviceId)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	accessToken, refreshToken, err := utils.GenerateTokenDashboard(dataLogin.User.Id, dataLogin.User.Username, rolesName, dataLogin.User.IsAdmin, dataLogin.User.IsInternal, dataLogin.MerchantKeyParking, dataLogin.User.PolicyDefaultId, dataLogin.User.OuDefaultId)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	response.User = dataLogin.User
	response.User.Password = constans.EMPTY_VALUE
	response.Token = accessToken
	response.RefreshToken = refreshToken
	response.FlagProgressive = resultDevice.FlgProgressive
	response.MerchantKey = dataLogin.MerchantKeyParking
	response.AdditionalInfo = nil
	response.MKey = resultDevice.MerchantKey

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, response)
	return ctx.JSON(http.StatusOK, result)
}
