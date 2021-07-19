package utils

import "time"

const (
	EnvProduction                = "prod"
	StatusActive                 = "active"
	StatusBlocked                = "blocked"
	StatusInactive               = "inactive"
	StatusPending                = "pending"
	StatusAccepted               = "accepted"
	StatusDeclined               = "declined"
	StatusIncomplete             = "incomplete"
	StatusCompleted              = "completed"
	StatusBooked                 = "booked"
	DateLayout                   = "2006-01-02T15:04:05.999999Z"
	ISOLayout                    = "2006-01-02T15:04:05.000Z"
	DateRangeLayout              = "02-01-2006"
	SimpleDateLayout             = "Monday, 2 Jan 2006, 03:04:05 PM"
	AuthorizationKey             = "authorization"
	DefaultExpirationPeriod      = time.Hour * 24 * 7
	AccessTokenExpirationPeriod  = time.Hour * 24
	RefreshTokenExpirationPeriod = time.Hour * 24 * 7
	KeyForSecretKey              = "Secret-Key"
	RealUserIpKey                = "X-Original-Forwarded-For"
	UsernameKey                  = "username"
	DefaultHashingIteration      = 1500
	TryAgainMessage              = "Please try again later"
	MaxOTPAttempt                = 50
	MaxLoginAttempt              = 5
	PasswordPattern              = "^([a-zA-z0-9!@#%*_=+/-]*)$"
	accessTokenKey               = "jdskhhiuewfhosfkaskfhajksfeiwhfuiowehfiwejdkfewudhuiewhjfdiu"
	refreshTokenKey              = "fdshfjdshfjhdsjlfhuoashfuherifherhfuqheruifhiquwhfukwjnfjiwhl"
	DefaultCaptchaValue          = "11111"
	LastResetEventAtKey          = "last_reset_at"
)
