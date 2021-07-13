package utils

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	rnd "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	hashers "github.com/meehow/go-django-hashers"
	"io"
	"log"
	"math"
	"math/rand"
	"reflect"
	"sync/atomic"
	"time"
)

type NameValuePairs struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

var AlphaNumerics = [...]byte{'r', 'e', 'z', 'o', 'a', 'n', 't', 'm', 'l', '0','1','2','3','4','5','6','7','8','9'}
var Digits = [...]byte{'0','1','2','3','4','5','6','7','8','9'}

func BoolP(in bool) *bool {
	return &in
}

func ExistsInSlice(arr []string, a string) bool {
	for _, v := range arr {
		if v == a {
			return true
		}
	}

	return false
}

func FieldRequiredErr(field string) error {
	errString := fmt.Sprintf("%s field is required", field)
	return fmt.Errorf(errString)
}

func CharLimitErr(text string, limit int) error {
	errString := fmt.Sprintf("%s should not exceed %d characters", text, limit)
	return fmt.Errorf(errString)
}

func StringToBoolP(value string) *bool {
	var res bool
	if value == "true" {
		res = true
	} else {
		res = false
	}
	return &res
}

func BoolPToString(value *bool) string {
	var res string
	if *value == true {
		res = "true"
	} else {
		res = "false"
	}
	return res
}

func DeepCopyStruct(in interface{}) interface{} {
	b, err := json.Marshal(in)
	if err != nil {
		log.Println(err)
		return nil
	}

	var out interface{}
	err = json.Unmarshal(b, &out)
	if err != nil {
		log.Println(err)
		return nil
	}

	return out
}

var (
	reqid uint64
)

func CustomJsonMarshal(data interface{}, tag string) ([]byte, error) {
	var json = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 tag,
	}.Froze()

	return json.Marshal(data)
}

func GetTracingID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

func SetTracingID(ctx context.Context) context.Context {
	uid := uuid.New().String()
	myid := atomic.AddUint64(&reqid, 1)
	requestID := fmt.Sprintf("%s-%06d", uid, myid)
	ctx = context.WithValue(ctx, middleware.RequestIDKey, requestID)
	return ctx
}

const myCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandStr(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = myCharset[seededRand.Intn(len(myCharset))]
	}
	return string(b)
}

func DecodeInterface(input, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

const float64EqualityThreshold = 1e-9

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
func TrueP() *bool {
	t := true
	return &t
}

func FalseP() *bool {
	t := false
	return &t
}

func GetCmFromFootInch(feet, inch int) float32 {
	inches := (feet * 12) + inch
	//log.Println(inch)
	//inchValue := unit.Inch * unit.Length((feet*12)+inch)
	//log.Println("inchValue:",inchValue)
	return 2.54 * float32(inches)
}

func GetFootInchFromCm(cm float32) (int, int) {
	inchValue := cm * 0.393701
	feet := int(inchValue / 12)
	inch := int(inchValue) % 12
	return feet, inch
}

func ValidateStruct(s interface{}) error {
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field validationTag value
		validationTag := field.Tag.Get("validate")
		jsonTag := field.Tag.Get("json")
		if validationTag == "nonzero" {
			if !isNonzero(field) {
				return FieldRequiredErr(field.Tag.Get(jsonTag))
			}
		}

		fmt.Printf("%d. %v (%v), validationTag: '%v'\n", i+1, field.Name, field.Type.Name(), validationTag)
	}

	return nil
}

func isNonzero(field reflect.StructField) bool {
	return false

}

func GetTimeFromISOString(in string) time.Time {
	bd, err := time.Parse(ISOLayout, in)
	if err != nil {
		return time.Time{}
	}
	return bd
}

func GetMd5SumString(data []byte) string {
	b := md5.Sum(data)
	return hex.EncodeToString(b[:])
}

func GetStringToByteArray(in string) []byte {
	return []byte(in)
}

func GetPasswordHash(in string) string {
	hashers.Iter = DefaultHashingIteration
	hashPass, err := hashers.MakePassword(in)
	if err != nil {
		return ""
	}

	return hashPass
}

func GenerateTimedRandomAlphaNumerics(key string, count, ttlMinutes int) (string, string, error) {
	//get the fixed 5 digit random number
	if len(key) < 11 {
		return "", "", fmt.Errorf("invalid phone number, timed digit generation failed")
	}

	ttl := time.Now().UTC().Add(time.Duration(ttlMinutes) * time.Minute).Unix()

	randomText := getRandomText(count, AlphaNumerics[:]) //get 5 digit OTP token
	fullHash := fmt.Sprintf("%v.%v", generateHash(fmt.Sprintf("%v.%v.%v", key, randomText, ttl)), ttl)

	return randomText, fullHash, nil
}

func GenerateTimedRandomDigits(key string, count, ttlMinutes int) (string, string, error) {
	//get the fixed 5 digit random number
	ttl := time.Now().UTC().Add(time.Duration(ttlMinutes) * time.Minute).Unix()

	randomDigits := getRandomText(count, Digits[:]) //get 5 digit OTP token
	fullHash := fmt.Sprintf("%v.%v", generateHash(fmt.Sprintf("%v.%v.%v", key, randomDigits, ttl)), ttl)

	return  randomDigits, fullHash, nil
}

func GetRandomDigits(len int) string {
	return getRandomText(len, Digits[:])
}


func getRandomText(count int, seed []byte) string {
	b := make([]byte, count)
	io.ReadAtLeast(rnd.Reader, b, count)

	for i := 0; i < len(b); i++ {
		b[i] = seed[int(b[i])%len(seed)]
	}
	return string(b)
}

func generateHash(data string) string {
	secretKey := "uJ9eZCNwg9NxJ2rmY6nN"
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
