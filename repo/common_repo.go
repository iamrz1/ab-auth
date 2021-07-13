package repo

import (
	"fmt"
	"github.com/iamrz1/auth/infra"
	infraCache "github.com/iamrz1/auth/infra/cache"
	"github.com/iamrz1/auth/logger"
	"github.com/iamrz1/auth/utils"
	"log"
	"time"
)

type CommonRepo struct {
	DB    infra.DB
	Cache *infraCache.Redis
	Log   logger.StructLogger
}

func NewCommonRepo(db infra.DB, table string, cache *infraCache.Redis, log logger.StructLogger) *CommonRepo {
	return &CommonRepo{
		DB:    db,
		Cache: cache,
		Log:   log,
	}
}


func (pr *CommonRepo) SaveOTP(username, otp string) error {
	key := fmt.Sprintf("otp.%s",username)
	scmd := pr.Cache.Client.Set(key,otp,time.Minute*5)

	return scmd.Err()
}

func (pr *CommonRepo) VerifyOTP(username, process, otp string) error {
	otpAttemptKey := fmt.Sprintf("otp.attempt.%s.%s", process,username)


	pr.Cache.Client.Incr(otpAttemptKey)
	gcmd := pr.Cache.Client.Get(otpAttemptKey)
	if gcmd.Err() != nil{
		log.Println(gcmd.Err())
	}

	count, err := gcmd.Int()
	if err != nil{
		log.Println(err)
		return fmt.Errorf("%s",utils.TryAgainMessage)
	}

	if count > utils.MaxOTPAttempt{
		time.Sleep(time.Duration(count)*time.Second) // to jam connection
		return fmt.Errorf("%s",utils.TryAgainMessage)
	}


	bcmd := pr.Cache.Client.SetNX(otpAttemptKey,1,time.Second*5)
	if bcmd.Err()!= nil {
		log.Println(bcmd.Err())
		return fmt.Errorf("%s",utils.TryAgainMessage)
	}
	boolVal, err := bcmd.Result()
	if err != nil || !boolVal{
		return fmt.Errorf("%s",utils.TryAgainMessage)
	}

	invalidOTP :="Invalid OTP"
	key := fmt.Sprintf("otp.%s",username)
	scmd := pr.Cache.Client.Get(key)
	if scmd.Err()!= nil {
		log.Println(scmd.Err())
		return fmt.Errorf("%s",invalidOTP)
	}

	cachedOtp, err := scmd.Result()
	if err != nil || !boolVal{
		log.Println(err)
		return fmt.Errorf("%s",invalidOTP)
	}

	if cachedOtp!= otp{
		return fmt.Errorf("%s",invalidOTP)
	}

	return nil
}