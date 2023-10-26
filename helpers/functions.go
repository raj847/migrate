package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"togrpc/constans"
	"togrpc/models"
	"togrpc/proto/trx"
	"togrpc/services"
	"togrpc/utils"
	"togrpc/validate"

	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	service services.UsecaseService
)

func InArray(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}

func BindValidateStruct(i interface{}) error {
	var r *http.Request
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		return err
	}

	if err := validate.Validate(i); err != nil {
		return err
	}

	return nil
}

func convDifferenceTime(startDate time.Time, currDatetime time.Time) int64 {
	log.Println("Start Datetime:", startDate, "Current Datetime:", currDatetime)
	getMinutesOvertime := (startDate.Unix() - currDatetime.Unix()) / 60
	currSecond := int64(currDatetime.Second() * 1000)
	return (getMinutesOvertime * constans.MILISECOND) - currSecond
}

func ConvDifferenceTimeForOvernight(checkinDatetimeStr, checkoutDatetimeStr string, currDatetimeNow time.Time) int64 {
	checkInDatetime, _ := time.Parse("2006-01-02 15:04", checkinDatetimeStr)
	checkOutDatetime, _ := time.Parse("2006-01-02 15:04", checkoutDatetimeStr)

	convTime := utils.ConvDiffTime(checkInDatetime, checkOutDatetime)

	checkInTime := convTime.Date1.Format("15:04")
	checkOutDate := convTime.Date2.Format("2006-01-02")

	lastDatetime, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", checkOutDate, checkInTime))
	lastDatetime = lastDatetime.Add(1 * time.Minute)

	return convDifferenceTime(lastDatetime, currDatetimeNow)
}

func ConvNextInvoiceForTime(checkinDatetimeStr, checkoutDatetimeStr string, currDatetimeNow time.Time) int64 {
	millisecondTime := int64(0)
	checkInDatetime, _ := time.Parse("2006-01-02 15:04", checkinDatetimeStr)
	checkOutDatetime, _ := time.Parse("2006-01-02 15:04", checkoutDatetimeStr)
	durationTime := (checkOutDatetime.Unix() - checkInDatetime.Unix()) / 60

	if durationTime%1440 >= 0 {
		checkOutDatetime = checkOutDatetime.AddDate(0, 0, 1)
		checkInTime := checkInDatetime.Format("15:04")
		checkOutDate := checkOutDatetime.Format("2006-01-02")

		lastDatetime, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", checkOutDate, checkInTime))
		lastDatetime = lastDatetime.Add(1 * time.Minute)

		millisecondTime = convDifferenceTime(lastDatetime, currDatetimeNow)
	}

	return millisecondTime
}

func ResponseJSON(success bool, code string, msg string, result *anypb.Any) *trx.Response {
	tm := timestamppb.New(time.Now())
	response := &trx.Response{
		Success:          success,
		StatusCode:       code,
		Result:           result,
		Message:          msg,
		ResponseDatetime: tm,
	}

	return response
}

func ConvOvernight24H(startDate time.Time, nextDate time.Time, currDatetime time.Time) int64 {
	log.Println("Start Datetime:", startDate, "Next Datetime:", nextDate, "Current Datetime:", currDatetime, "24H")
	getMinutesOvertime := (nextDate.Unix() - startDate.Unix()) / 60
	getMinutesTimeNow := (currDatetime.Unix() - startDate.Unix()) / 60

	currSecond := int64(currDatetime.Second() * 1000)
	gapMinutes := (getMinutesOvertime - getMinutesTimeNow) * 60000
	return gapMinutes - currSecond
}

func CalculateTrxInvoiceDetail(trxMemberList []models.TrxMember, trxInvoiceDetailItemList []models.TrxInvoiceDetailItem) float64 {
	var invoiceAmount float64

	for _, dataTrx := range trxInvoiceDetailItemList {
		for _, data := range trxMemberList {
			exists := InBetween(dataTrx.CreatedDate, data.DateFrom, data.DateTo)
			if exists {
				invoiceAmount += dataTrx.InvoiceAmount
			}
		}
	}

	return invoiceAmount
}

func InBetween(invoiceDate, dateFrom, dateTo string) bool {
	invoiceDateFmt, _ := strconv.ParseInt(strings.ReplaceAll(invoiceDate, "-", ""), 10, 64)
	dateFromFmt, _ := strconv.ParseInt(strings.ReplaceAll(dateFrom, "-", ""), 10, 64)
	dateToFmt, _ := strconv.ParseInt(strings.ReplaceAll(dateTo, "-", ""), 10, 64)

	result := false
	if (invoiceDateFmt >= dateFromFmt) && (invoiceDateFmt <= dateToFmt) {
		result = true
	}

	return result
}

func GetPrivateIPLocal() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return constans.EMPTY_VALUE, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return constans.EMPTY_VALUE, nil
}

func FilterTrxInvoiceByIndex(docNo string, productCode string, collection []models.TrxInvoiceItem) (trxInvoiceItem models.TrxInvoiceItem, err error) {

	for _, rows := range collection {
		if rows.DocNo == docNo && rows.ProductCode == productCode {
			return rows, nil
		}
	}

	return trxInvoiceItem, errors.New(fmt.Sprintf("%s %s %s %s", "Invoice not found with document no", docNo, "product code", productCode))
}
