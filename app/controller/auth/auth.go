package auth_controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	service_auth "comprix/app/services/auth"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthController struct {
	config     dto.Config
	service    service_auth.Service
	repository repositories.UserRepository
}

func InitAuthController(config dto.Config) AuthController {
	return AuthController{
		config:     config,
		service:    service_auth.InitService(config),
		repository: repositories.InitUserRepository(config.GormDB),
	}
}

func (ctr AuthController) Update(ginContext *gin.Context) {
	// Obtener formulario
	form := server.GinFormBinding(ginContext, dto.UpdateProfile{}, "Datos erroneos")
	claims := ginContext.MustGet("claims").(jwt.MapClaims)
	user, err := ctr.repository.FindByID(claims["sub"].(map[string]interface{})["id"])
	// Verificar si se encontro
	if err != nil {
		server.InternalErrorException("No se pudo actualizar.", nil)
	}
	// Inicializar datos a actualizar
	updatedUser := dao.User{
		Picture:     user.Picture,
		LastName:    form.LastName,
		FirstName:   form.FirstName,
		PhoneNumber: form.PhoneNumber,
	}
	// Verificar si existe archivo
	if form.File != nil {
		ext := strings.ToLower(filepath.Ext(form.File.Filename))
		uploadDir := fmt.Sprintf("users/%v/", user.ID)
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		fileSrc := filepath.Join(ctr.config.Settings.Storage.Path+"/"+uploadDir, fileName)
		// Si el usuario ya tiene una imagen, eliminarla antes de guardar la nueva
		if user.Picture != nil {
			oldFilePath := filepath.Join(ctr.config.Settings.Storage.Path+"/"+uploadDir, *user.Picture)
			if _, err := os.Stat(oldFilePath); err == nil {
				os.Remove(oldFilePath)
			}
		}
		// Actualizar valor
		updatedUser.Picture = &fileName
		// Guardar archivo
		if err := ginContext.SaveUploadedFile(form.File, fileSrc); err != nil {
			server.InternalErrorException("No se pudo guardar el archivo.", err.Error())
		}
	}
	// Guardar en la BD
	if ctr.repository.Update(user.ID, &updatedUser) != nil {
		server.InternalErrorException("No se pudo actualizar.", nil)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}

func (ctr AuthController) GetDetail(ginContext *gin.Context) {
	// Obtener formulario
	claims := ginContext.MustGet("claims").(jwt.MapClaims)
	user, err := ctr.repository.FindByID(claims["sub"].(map[string]interface{})["id"])
	// Verificar si se encontro
	if err != nil {
		server.InternalErrorException("No se pudo actualizar.", nil)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, user)
}
