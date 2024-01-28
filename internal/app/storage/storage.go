package storage

import (
	"errors"
	"otp_api/internal/app/pkg/logger"

	"github.com/go-redis/redis"
)

type Storage struct {
	log *logger.Logger
	db  *redis.Client
}

func NewStorage(log *logger.Logger, db *redis.Client) *Storage {
	return &Storage{
		log: log,
		db:  db,
	}
}

func (s *Storage) CreateOTP(id string, hash map[string]interface{}) error {
	err := s.db.HMSet(id, hash).Err()
	if err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

func (s *Storage) GetOTP(key string) (map[string]string, error) {
	hash, err := s.db.HGetAll(key).Result()
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	if len(hash) == 0 {
		s.log.Error("invalid key: hash not found")
		return nil, errors.New("notFound")
	}

	return hash, nil
}

func (s *Storage) UpdateOTPField(hashKey string, key string, value interface{}) error {
	_, err := s.db.HSet(hashKey, key, value).Result()
	if err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

func (s *Storage) UpdateMultipleFields(hashKey string, fields map[string]interface{}) error {
	err := s.db.HMSet(hashKey, fields).Err()
	if err != nil {
		s.log.Error(err)
		return err
	}
	return nil
}

func (s *Storage) BanAccountForTime(account string, field string, duration string) error {
	err := s.db.HSet(account, field, duration).Err()
	if err != nil {
		s.log.Error(err)
		return err
	}
	return nil
}

func (s *Storage) GetAccount(account string) (map[string]string, error) {
	mp, err := s.db.HGetAll(account).Result()
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return mp, nil
}

func (s *Storage) RemoveAccountFromBan(account string) error {
	err := s.db.Del(account).Err()
	if err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}
