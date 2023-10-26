package helpers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/raj847/togrpc/config"
	"github.com/raj847/togrpc/constans"
)

func GetCallAPI(suffixUrl string) (*string, error) {
	url := config.BASEParkingURL + suffixUrl
	log.Println("CONFIG BASE URL :", config.BASEParkingURL)
	log.Println("GetCallAPI:", url)

	httpReq, err := http.Get(url)
	if err != nil {
		log.Println("Err http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpReq.Body.Close()

	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		log.Fatalln(err)
	}

	stringBuilder := string(body)

	return &stringBuilder, nil
}

func GetCallAPITrxPaymentOnline(suffixUrl string) (*string, error) {
	url := config.BaseParkingURLPaymentOnline + suffixUrl
	log.Println("GetCallAPI:", url)

	httpReq, err := http.Get(url)
	if err != nil {
		log.Println("Err http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpReq.Body.Close()

	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		log.Fatalln(err)
	}

	stringBuilder := string(body)

	return &stringBuilder, nil
}

func CallHttpCloudServer(method string, body []byte, suffixUrl, BearerToken, basicAuth, urlServer string) (*[]byte, error) {
	url := urlServer + suffixUrl

	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Err http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpReq.Body.Close()

	httpReq.Close = true
	httpReq.Header.Add("Content-Type", "application/json")
	if BearerToken != constans.EMPTY_VALUE {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", BearerToken))
	} else if basicAuth != constans.EMPTY_VALUE {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", "Basic", basicAuth))

	}
	httpReq.Header.Set("Connection", "close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	log.Println("URL:", url)
	log.Println("Request:", string(body))

	response, err := client.Do(httpReq)
	if err != nil {
		log.Println("Err - client.Do :", err.Error())
		return nil, err
	}

	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {

		log.Println("Err ReadAll :", err.Error())
		return nil, err
	}

	bodyString := string(body)
	log.Println("Response:", bodyString)

	return &body, nil
}

func CallHttpCloudServerParkingPayment(method string, body []byte, suffixUrl string) (*[]byte, error) {
	url := config.BASEParkingURL + suffixUrl

	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Err http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpReq.Body.Close()

	httpReq.Close = true
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic bWtwbW9iaWxlOm1rcG1vYmlsZTEyMw==")
	httpReq.Header.Set("Connection", "close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	log.Println("URL:", url)
	log.Println("Request:", string(body))

	response, err := client.Do(httpReq)
	if err != nil {
		log.Println("Err - client.Do :", err.Error())
		return nil, err
	}

	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {

		log.Println("Err ReadAll :", err.Error())
		return nil, err
	}

	return &body, nil
}

func CallHttpCloudServerParking(method string, body []byte, suffixUrl string) (*string, error) {
	url := config.BASEParkingURL + suffixUrl

	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Err http.NewRequest ", url, " :", err.Error())
		return nil, err
	}
	defer httpReq.Body.Close()

	httpReq.Close = true
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic bWtwbW9iaWxlOm1rcG1vYmlsZTEyMw==")
	httpReq.Header.Set("Connection", "close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	log.Println("URL:", url)
	log.Println("Request:", string(body))

	response, err := client.Do(httpReq)
	if err != nil {
		log.Println("Err - client.Do :", err.Error())
		return nil, err
	}

	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {

		log.Println("Err ReadAll :", err.Error())
		return nil, err
	}

	bodyString := string(body)
	log.Println("Response:", bodyString)

	return &bodyString, nil
}
