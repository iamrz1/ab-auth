package repo

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/iamrz1/ab-auth/infra"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/utils"
	"log"
	"time"
)

type CommonRepo struct {
	DB    infra.DB
	Cache *infraCache.Redis
	Log   logger.StructLogger
}

func NewCommonRepo(db infra.DB, cache *infraCache.Redis, log logger.StructLogger) *CommonRepo {
	return &CommonRepo{
		DB:    db,
		Cache: cache,
		Log:   log,
	}
}

func (cmr *CommonRepo) GetOTP(username, service string, limit, limitDuration, lockDuration int) (string, error) {
	otp := utils.GetRandomDigits(5)
	ok, err := cmr.LockKey(fmt.Sprintf("%s_%s_otp_gen", username, service), lockDuration)
	if err != nil || !ok {
		return "", fmt.Errorf("%s", "Can not request multiple OTPs at once")
	}

	if !cmr.EnsureUsageLimit(fmt.Sprintf("%s_%s_otp_gen_limit", username, service), limit, limitDuration) {
		return "", fmt.Errorf("%s", "Please try again in 24 hours")
	}

	return otp, nil
}

func (cmr *CommonRepo) SetOTP(username, service, otp string, durationSec int) error {
	scmd := cmr.Cache.Client.Set(fmt.Sprintf("%s_%s_otp", username, service), otp, time.Duration(durationSec))
	if scmd.Err() != nil {
		return fmt.Errorf("%s", "OTP request failed")
	}

	return nil
}

func (cmr *CommonRepo) MatchOTP(username, service, otp string) error {
	scmd := cmr.Cache.Client.Get(fmt.Sprintf("%s_%s_otp", username, service))
	if scmd.Err() != nil {
		return fmt.Errorf("%s", "OTP match failed")
	}

	if scmd.Val() != otp {
		return fmt.Errorf("%s", "Incorrect OTP")
	}

	return nil
}

func (cmr *CommonRepo) LockKey(key string, durationSec int) (bool, error) {
	res := cmr.Cache.Client.SetNX(key, 1, time.Second*time.Duration(durationSec))
	if res.Err() != nil {
		log.Println(res.Err())
		return false, res.Err()
	}

	return res.Result()
}

func (cmr *CommonRepo) EnsureUsageLimit(key string, limit, durationSec int) bool {
	//pipe := cmr.Cache.Client.TxPipeline()
	usedLimit := 0
	scmd := cmr.Cache.Client.Get(key)
	if scmd.Err() != nil {
		if scmd.Err() != redis.Nil {
			cmr.Log.Errorf("EnsureUsageLimit", "", scmd.Err().Error())
			return false
		} else {
			cmr.Cache.Client.Set(key, 1, time.Minute*time.Duration(durationSec))
			return true
		}
	} else {
		n, err := scmd.Int()
		if err != nil {
			cmr.Log.Errorf("EnsureUsageLimit", "", scmd.Err().Error())
			return false
		}
		usedLimit = n
	}

	if usedLimit > limit {
		return false
	}

	cmr.Cache.Client.Incr(key)

	return true
}
