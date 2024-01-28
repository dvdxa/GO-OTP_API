package usecase

import (
	"fmt"
	"otp_api/internal/app/pkg/logger"
	"otp_api/internal/app/pkg/models"
	"otp_api/pkg/utils"
	"strconv"
	"time"
)

const (
	PROCESSED    = 200
	CREATED      = 201
	ERROR        = 400
	NOT_FOUND    = 404
	NOT_MODIFIED = 304
	PARTIAL      = 206
)

var Invalid bool

type Service struct {
	storage IStorage
	adapter IAdapter
	log     *logger.Logger
}

func NewService(storage IStorage, adapter IAdapter, log *logger.Logger) *Service {
	return &Service{
		storage: storage,
		adapter: adapter,
		log:     log,
	}
}

func (s *Service) CreateOTP(request *models.Request, secretKey string, validateLimit int) (*models.Response, *models.ErrorResponse) {
	var (
		response models.Response
		errResp  models.ErrorResponse
		hash     models.Hash
	)

	lifetime := strconv.Itoa(int(*request.Lifetime))

	valLim := strconv.Itoa(validateLimit)

	hash.Account = request.Account
	hash.State = "201" //created
	hash.CreatedAt = time.Now().Format("02.01.2006 15:04:05 -0700")
	hash.ExpiredAt = time.Now().Add(time.Duration(*request.Lifetime) * time.Second).Format("02.01.2006 15:04:05 -0700")
	hash.LifeTime = lifetime
	hash.ValidateLimit = valLim

	mp, err := utils.StructToMap(hash)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	err = s.storage.CreateOTP(request.ID, mp)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	concat := fmt.Sprintf("%s%s%s", request.ID, hash.Account, lifetime)
	intOtp, err := utils.GenerateOTP(concat, secretKey)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = fmt.Sprintf("failed to generate otp: %v", err)
		return nil, &errResp
	}

	strOtp := strconv.Itoa(intOtp)

	s.log.Infof("[OTP] %s\n", strOtp)

	var flag bool
	var state string

	err = s.adapter.SendOTPToClient(response.Account, strOtp)
	if err != nil {
		errResp.Code = "internal error"
		errResp.Message = fmt.Sprintf("failed to send otp to client: %v", err)

		//imitate stub
		//set here flag to true
	}

	if flag {
		//failed to send otp to client
		state = "400"

		err := s.storage.UpdateOTPField(request.ID, "state", state)
		if err != nil {
			return nil, &errResp
		}

		return nil, &errResp
	}

	//sent to client, need to validate otp
	state = "206"
	err = s.storage.UpdateOTPField(request.ID, "state", state)
	if err != nil {
		errResp.Code = "internal error"
		errResp.Message = fmt.Sprintf("failed to update state after sending otp: %v", err)
		return nil, &errResp
	}

	response.ID = request.ID
	response.Account = hash.Account
	response.State = CREATED
	response.CreatedAt = hash.CreatedAt
	response.ExpiredAt = hash.ExpiredAt
	response.LifeTime = *request.Lifetime

	return &response, nil
}

func (s *Service) GetOTP(id string, secretKey string) (*models.Response, *models.ErrorResponse) {
	var (
		hash     models.Hash
		response models.Response
		errResp  models.ErrorResponse
	)

	mp, err := s.storage.GetOTP(id)
	fmt.Printf("%+v\n", mp)
	if err != nil {
		if err.Error() == "notFound" {
			errResp.Code = "invalid input"
			errResp.Message = "invalid resource ID"
			return nil, &errResp
		}
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	err = utils.MapToStruct(mp, &hash)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	state, err := strconv.Atoi(hash.State)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	lifetime, err := strconv.Atoi(hash.LifeTime)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	response.ID = id
	response.State = int64(state)
	response.Account = hash.Account
	response.CreatedAt = hash.CreatedAt
	response.ExpiredAt = hash.ExpiredAt
	response.LifeTime = int64(lifetime)

	return &response, nil
}

func (s *Service) ValidateOTP(id string, otp string, secretKey string, retryAfter int, banDur int) (*models.Response, *models.ErrorResponse) {
	var (
		response models.Response
		errResp  models.ErrorResponse
		hash     models.Hash
	)

	mp, err := s.storage.GetOTP(id)

	if err != nil || len(mp) == 0 {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = "failed to get otp by ID"
		return nil, &errResp
	}

	err = utils.MapToStruct(mp, &hash)

	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = fmt.Sprintf("failed to convert map to struct: %v\n", err)
		return nil, &errResp
	}

	ok, err := s.isAccountBanned(hash, retryAfter)
	if err != nil {
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	if err == nil && ok {
		errResp.Code = "try later"
		errResp.Message = "account is banned"
		return nil, &errResp
	}

	if hash.ValidateLimit == "0" {
		Invalid = true
		mp := map[string]interface{}{
			"state":    "304",
			"ban_time": time.Now().Format("02.01.2006 15:04:05 -0700"),
		}
		//ban account
		err := s.storage.UpdateMultipleFields(id, mp)
		if err != nil {
			s.log.Error(err)
			return nil, &models.ErrorResponse{}
		}

		dur := strconv.Itoa(banDur)

		err = s.storage.BanAccountForTime(hash.Account, "ban_time", dur)
		if err != nil {
			s.log.Error(err)
			return nil, &models.ErrorResponse{}
		}
	}

	if hash.State == "200" {
		Invalid = true
		err := s.storage.UpdateOTPField(id, "state", "304")
		if err != nil {
			s.log.Error(err)
			return nil, &models.ErrorResponse{}
		}
		return nil, &models.ErrorResponse{}
	}

	ok, err = s.checkExpiration(hash)
	if err != nil {
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	if err == nil && ok {
		Invalid = true
		s.log.Infof("otp is expired")
		return nil, &models.ErrorResponse{}
	}

	concat := fmt.Sprintf("%s%s%s", id, hash.Account, hash.LifeTime)

	intOtp, err := utils.GenerateOTP(concat, secretKey)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = fmt.Sprintf("failed to validate otp: %v", err)
		return nil, &errResp
	}

	strOtp := strconv.FormatInt(int64(intOtp), 10)

	if strOtp != otp {
		s.log.Errorf("[ERROR] OTP's are not equal: sent by client: %s and actual otp on server: %s\n", otp, strOtp)

		Invalid = true

		validateLim, err := strconv.Atoi(hash.ValidateLimit)
		if err != nil {
			s.log.Error(err)
			return nil, &models.ErrorResponse{}
		}

		validateLim--

		err = s.storage.UpdateOTPField(id, "validate_limit", validateLim)
		if err != nil {
			s.log.Error(err)
			return nil, &models.ErrorResponse{}
		}

		return nil, &models.ErrorResponse{}
	}

	lifetime, err := strconv.Atoi(hash.LifeTime)
	if err != nil {
		s.log.Error(err)
		errResp.Code = "internal error"
		errResp.Message = err.Error()
		return nil, &errResp
	}

	response.ID = id
	response.Account = hash.Account
	response.State = PROCESSED
	response.CreatedAt = hash.CreatedAt
	response.ExpiredAt = hash.ExpiredAt
	response.LifeTime = int64(lifetime)

	return &response, nil
}

func (s *Service) isAccountBanned(hash models.Hash, retryAfter int) (bool, error) {

	//check if account banned
	mp, err := s.storage.GetAccount(hash.Account)
	if err != nil {
		return false, err
	}

	//account in ban
	if len(mp) != 0 {

		layout := "02.01.2006 15:04:05 -0700"
		tm, err := time.Parse(layout, hash.BanTime)
		if err != nil {
			s.log.Error(err)
			return false, fmt.Errorf("failed to parse ban time: %v", err)
		}

		elapsed := time.Since(tm)
		if elapsed <= time.Duration(retryAfter)*time.Second {
			s.log.Infof("account %v is banned\n", hash.Account)
			return true, nil
		}

		err = s.storage.RemoveAccountFromBan(hash.Account)
		if err != nil {
			s.log.Error(err)
			return false, err
		}
	}

	return false, nil
}

func (s *Service) checkExpiration(hash models.Hash) (bool, error) {
	layout := "02.01.2006 15:04:05 -0700"
	tm, err := time.Parse(layout, hash.ExpiredAt)
	if err != nil {
		s.log.Error(err)
		return false, err
	}

	lifetime, err := strconv.Atoi(hash.LifeTime)
	if err != nil {
		s.log.Error(err)
		return false, err
	}

	//otp is expired
	elapsed := time.Since(tm)
	if elapsed >= time.Duration(lifetime)*time.Second {
		return true, nil
	}
	return false, nil
}
