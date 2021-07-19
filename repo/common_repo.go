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

func (cmr *CommonRepo) VerifyOTP(username, process, otp string) error {
	bcmd := cmr.Cache.Client.SetNX(fmt.Sprintf("otp.lock.%s", username), 1, time.Second*5)
	if bcmd.Err() != nil {
		log.Println(bcmd.Err())
		return fmt.Errorf("%s", utils.TryAgainMessage)
	}
	boolVal, err := bcmd.Result()
	if err != nil || !boolVal {
		log.Println(err, boolVal)
		return fmt.Errorf("%s", utils.TryAgainMessage)
	}

	otpAttemptKey := fmt.Sprintf("otp.attempt.%s.%s", process, username)
	cmr.Cache.Client.Incr(otpAttemptKey)
	gcmd := cmr.Cache.Client.Get(otpAttemptKey)
	if gcmd.Err() != nil {
		log.Println(gcmd.Err())
		return fmt.Errorf("%s", utils.TryAgainMessage)
	}

	count, err := gcmd.Int()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("%s", utils.TryAgainMessage)
	}

	log.Println("otpAttempt count:", count)

	if count > utils.MaxOTPAttempt {
		time.Sleep(time.Duration(count) * time.Second) // to jam connection
		return fmt.Errorf("%s", utils.TryAgainMessage)
	}

	invalidOTP := "Invalid OTP"
	key := fmt.Sprintf("otp.%s", username)
	scmd := cmr.Cache.Client.Get(key)
	if scmd.Err() != nil {
		log.Println(scmd.Err())
		return fmt.Errorf("%s", invalidOTP)
	}

	cachedOtp, err := scmd.Result()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("%s", invalidOTP)
	}

	if cachedOtp != otp {
		log.Println(cachedOtp, otp)
		return fmt.Errorf("%s", invalidOTP)
	}

	cmr.Cache.Client.Del(key)

	return nil
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
