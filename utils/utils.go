package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"github.com/raj847/togrpc/constans"
	"github.com/raj847/togrpc/models"

	"google.golang.org/protobuf/reflect/protoreflect"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

func GetIsAdminToken(ctx echo.Context) string {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["isAdmin"].(string)
}

func GetIsInternalToken(ctx echo.Context) string {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["isInternal"].(string)
}

func Timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func DateNow() string {
	return time.Now().Format("2006-01-02")
}

func TimeFormatString(tm time.Time) string {
	return tm.Format("2006-01-02 15:04:05")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func GetUsernameToken(ctx echo.Context) string {
	// user := ctx.Get("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)

	// return claims["username"].(string)
	return "DASHBOARD MKP"
}

func GetOuDefaultIdForToken(ctx echo.Context) float64 {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["ouDefaultId"].(float64)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(err.Error()))
			return
		}

		token := fields[1]
		claims := &jwt.MapClaims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("testjwt"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(models.NewErrorResponse(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse(err.Error()))
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Token not valid"))
			return
		}

		// Set claims ke context
		ctxWithClaims := context.WithValue(r.Context(), "claims", claims)

		// Jika Anda ingin mengambil 'isAdmin' dan menyimpannya di context
		uName, ok := (*claims)["username"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse("'username' claim not found or is not a string"))
			return
		}
		isAdmin, ok := (*claims)["isAdmin"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse("'isAdmin' claim not found or is not a string"))
			return
		}
		ctxWithIsAdmin := context.WithValue(ctxWithClaims, "isAdmin", isAdmin)
		ctxWithuname := context.WithValue(ctxWithClaims, "username", uName)

		next.ServeHTTP(w, r.WithContext(ctxWithIsAdmin))
		next.ServeHTTP(w, r.WithContext(ctxWithuname))
	})
}

func DBTransaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rollback Panic
		} else if err != nil {
			tx.Rollback() // err is not nill
		} else {
			err = tx.Commit() // err is nil
		}
	}()
	err = txFunc(tx)
	return err
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytes2 = "123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandStringBytesMaskImprSrcChr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes2) {
			b[i] = letterBytes2[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

type ConvTime struct {
	Date1   time.Time
	Date2   time.Time
	Days    int
	Hours   int
	Minutes int
	Second  int
}

func (c ConvTime) ProtoReflect() protoreflect.Message {
	panic("implement me")
}

var (
	TransTime     = ""
	FullTransTime = ""
)

func ConvDiffTime(d1, d2 time.Time) ConvTime {

	days, hours, minutes, seconds := getDifference(d1, d2)
	TransTime = fmt.Sprintf("%v:%v:%v", hours, minutes, seconds)
	FullTransTime = fmt.Sprintf("%v Days, %v Hours %v Minutes %v Second", days, hours, minutes, seconds)

	return ConvTime{Date1: d1, Date2: d2, Days: days, Hours: hours, Minutes: minutes, Second: seconds}
}

func convOvernight24H(startDate time.Time, nextDate time.Time, currDatetime time.Time) int64 {
	log.Println("Start Datetime:", startDate, "Next Datetime:", nextDate, "Current Datetime:", currDatetime, "24H")
	getMinutesOvertime := (nextDate.Unix() - startDate.Unix()) / 60
	getMinutesTimeNow := (currDatetime.Unix() - startDate.Unix()) / 60

	currSecond := int64(currDatetime.Second() * 1000)
	gapMinutes := (getMinutesOvertime - getMinutesTimeNow) * constans.MILLISECOND
	return gapMinutes - currSecond
}

func leapYears(date time.Time) (leaps int) {

	y, m, _ := date.Date()

	if m <= 2 {
		y--
	}
	leaps = y/4 + y/400 - y/100
	return leaps
}

func getDifference(a, b time.Time) (days, hours, minutes, seconds int) {

	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()

	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()

	totalDays1 := y1*365 + d1

	for i := 0; i < (int)(m1)-1; i++ {
		totalDays1 += monthDays[i]
	}

	totalDays1 += leapYears(a)

	totalDays2 := y2*365 + d2

	for i := 0; i < (int)(m2)-1; i++ {
		totalDays2 += monthDays[i]
	}

	totalDays2 += leapYears(b)

	days = totalDays2 - totalDays1

	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1

	if seconds < 0 {
		seconds += 60
		minutes--
	}

	if minutes < 0 {
		minutes += 60
		hours--
	}

	if hours < 0 {
		hours += 24
		days--
	}

	return days, hours, minutes, seconds
}

func ValidateMember(val models.ValidateMember) (bool, string) {
	timeNow := time.Now().Format("2006-01-02")
	timeParse, _ := time.Parse("2006-01-02", timeNow)
	DateFrom, _ := time.Parse("2006-01-02", val.DateFrom)
	DateTo, _ := time.Parse("2006-01-02", val.DateTo)
	reqDateFrom, _ := time.Parse("2006-01-02", val.ReqDateFrom)
	reqDateTo, _ := time.Parse("2006-01-02", val.ReqDateTo)

	if val.ReqPartnerCode == val.PartnerCode {
		if timeParse.After(DateFrom) && timeParse.Before(DateTo) ||
			timeParse.Equal(DateFrom) || timeParse.Equal(DateTo) {
			return false, fmt.Sprintf("Members Are Still Active")
		}

		if reqDateFrom.Equal(DateFrom) || reqDateTo.Equal(DateTo) ||
			reqDateFrom.Equal(DateTo) || reqDateFrom.Before(DateTo) ||
			reqDateFrom.After(DateFrom) && reqDateTo.Before(DateTo) {

			if val.ReqPartnerCode == val.PartnerCode &&
				val.ReqProductId == val.ProductId {
				return false, fmt.Sprintf("Member Already Exists")
			}

			if val.ReqCardNumber == val.CardNumber &&
				val.ReqProductId == val.ProductId {
				return false, fmt.Sprintf("Card Already Exists")
			}
		}
	}

	return true, ""
}

func CurrDatetimeNow() time.Time {
	currentTime := time.Now()
	format := currentTime.Format("2006-01-02 15:04:05")
	currDatetimeNow, _ := time.Parse("2006-01-02 15:04:05", format)
	tm := time.Unix(currDatetimeNow.Unix(), 0).UTC()

	return tm
}

func GetUsernameForToken(ctx echo.Context) string {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["username"].(string)
}

func ValBlankOrNull(request interface{}, keyName ...string) error {
	var params interface{}
	_ = json.Unmarshal([]byte(ToString(request)), &params)
	paramsValue := params.(map[string]interface{})

	for idx := range keyName {
		name := keyName[idx]
		if paramsValue[name] == nil {
			return errors.New(fmt.Sprintf("%s must be filled", name))
		}

		switch reflect.TypeOf(paramsValue[name]).Kind() {
		case reflect.String:
			if len(strings.TrimSpace(paramsValue[name].(string))) == 0 {
				return errors.New(fmt.Sprintf("%s must be filled", name))
			}
		}

	}

	return nil
}

func IsConnected() bool {

	_, err := http.Get("https://www.google.com")
	if err != nil {
		return false
	}
	//_, err := http.Get("http://clients3.google.com/generate_204")
	//if err != nil {
	//	return false
	//}
	return true
}

// GenerateTokenDashboard ...
func GenerateTokenDashboard(id int64, username string, rolesName string, isAdmin string, isInternal, merchantKey string, policyDefaultId, ouDefaultId int64) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["sub"] = 1
	claims["username"] = username
	claims["rolesName"] = rolesName
	claims["isAdmin"] = isAdmin
	claims["isInternal"] = isInternal
	claims["merchantKey"] = merchantKey
	claims["ouDefaultId"] = ouDefaultId
	claims["policyDefaultId"] = policyDefaultId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString([]byte(constans.GET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = id
	rtClaims["sub"] = 1
	rtClaims["username"] = username
	rtClaims["rolesName"] = rolesName
	rtClaims["isAdmin"] = isAdmin
	rtClaims["isInternal"] = isInternal
	rtClaims["merchantKey"] = merchantKey
	rtClaims["ouDefaultId"] = ouDefaultId
	rtClaims["policyDefaultId"] = policyDefaultId
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rToken, err := refreshToken.SignedString([]byte(constans.GET_KEY))
	if err != nil {
		return "", "", err
	}

	return accessToken, rToken, nil
}

// Check hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateBackdate(val models.ValidateRegis) (bool, string) {
	reqDateFrom, _ := time.Parse("2006-01-02 15", val.RequestDateFrom)
	reqDateTo, _ := time.Parse("2006-01-02 15", val.RequestDateTo)
	//dateNow, _ := time.Parse("2006-01-02", DateNow())

	//if reqDateFrom.Before(dateNow) {
	//	return false, "Start date member not valid"
	//}

	if reqDateTo.Before(reqDateFrom) {
		return false, fmt.Sprintf("Invalid Member Registration Time")
	}

	return true, constans.EMPTY_VALUE
}

func MapStringInterfaces(i interface{}) map[string]interface{} {
	var params interface{}
	_ = json.Unmarshal([]byte(ToString(i)), &params)
	paramsValue := params.(map[string]interface{})

	return paramsValue
}
