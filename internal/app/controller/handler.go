package controller

import (
	"encoding/json"
	"io/ioutil"
	"otp_api/conf"
	"otp_api/internal/app/pkg/logger"
	"otp_api/internal/app/pkg/models"
	"otp_api/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     *conf.Config
	log     *logger.Logger
	usecase usecase.IService
}

func NewHandler(cfg *conf.Config, log *logger.Logger, usecase *usecase.Service) *Handler {
	return &Handler{
		cfg:     cfg,
		log:     log,
		usecase: usecase,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(usecase.PROCESSED, map[string]interface{}{
		"message": "app is ok",
	})
}

func (h *Handler) CreateOTP(c *gin.Context) {
	var (
		request models.Request
		errResp models.ErrorResponse
	)

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Errorf("[ERROR] %v\n", err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		c.JSON(usecase.ERROR, errResp)
		return
	}

	h.log.Printf("[REQUEST] %v\n", string(b))

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.log.Error("[ERROR] ", err)
		errResp.Code = "invalid input"
		errResp.Message = "resource not found"

		h.log.Infof("[RESPONSE] %v\n", errResp)
		c.JSON(usecase.ERROR, errResp)
		return
	}

	if !validate(request.ID, request.Account) || request.Lifetime == nil {
		errResp.Code = "invalid input"
		errResp.Message = "required fields are missing"
		c.JSON(usecase.ERROR, errResp)
		return
	}

	response, errR := h.usecase.CreateOTP(&request, h.cfg.Sever.SecretKey, int(h.cfg.Sever.ValidateLimit))
	if errR != nil {
		h.log.Infof("[RESPONSE] %v\n", errR)
		c.JSON(usecase.ERROR, errR)
		return
	}

	h.log.Infof("[RESPONSE] %v\n", response)
	c.JSON(usecase.CREATED, response)
}

func (h *Handler) GetOTP(c *gin.Context) {
	var (
		errResp models.ErrorResponse
	)

	id := c.Param("id")

	h.log.Infof("[ID] %s\n", id)

	if id == "" {
		errResp.Code = "invalid input"
		errResp.Message = "failed to get id"
		c.JSON(usecase.NOT_FOUND, errResp)
		return
	}

	response, errR := h.usecase.GetOTP(id, h.cfg.Sever.SecretKey)
	if errR != nil {
		h.log.Error("[ERROR] ", errR)
		c.JSON(usecase.NOT_FOUND, errR)
		return
	}

	h.log.Infof("[RESPONSE] %v\n", response)
	c.JSON(usecase.PROCESSED, response)
}

func (h *Handler) ValidateOTP(c *gin.Context) {
	var (
		request models.Request
		errResp models.ErrorResponse
	)

	id := c.Param("id")

	h.log.Infof("[ID] %s\n", id)

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Errorf("[ERROR] %v\n", err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		c.JSON(usecase.NOT_FOUND, errResp)
		return
	}

	h.log.Printf("[REQUEST] %v\n", string(b))

	err = json.Unmarshal(b, &request)
	if err != nil {
		h.log.Error("[ERROR] ", err)
		errResp.Code = "invalid input"
		errResp.Message = "resource not found"

		h.log.Infof("[RESPONSE] %v\n", errResp)
		c.JSON(usecase.NOT_FOUND, errResp)
		return
	}

	if id == "" || request.Value == "" {
		errResp.Code = "invalid input"
		errResp.Message = "required fields are missing"
		c.JSON(usecase.NOT_FOUND, errResp)
		return
	}

	response, errR := h.usecase.ValidateOTP(id, request.Value, h.cfg.Sever.SecretKey, int(h.cfg.Sever.RetryAfter), int(h.cfg.Sever.BanDuration))
	if errR != nil {
		h.log.Errorf("[RESPONSE] %v\n", errR)
		status := usecase.NOT_FOUND

		if usecase.Invalid {
			usecase.Invalid = false
			status = usecase.NOT_MODIFIED
		}

		c.JSON(status, errR)
		return
	}

	h.log.Infof("[RESPONSE] %v\n", response)
	c.JSON(usecase.PROCESSED, response)
}

func validate(fields ...string) bool {
	for i := range fields {
		if fields[i] == "" {
			return false
		}
	}
	return true
}
