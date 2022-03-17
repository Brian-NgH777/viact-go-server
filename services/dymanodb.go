package services

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"strings"
)

const (
	AWS_S3_REGION2            = "ap-northeast-1"
	AWS_S3_ACCESS_KEY_ID2     = "AKIAZRF7Z3PJVVHW3DNY"
	AWS_S3_SECRET_ACCESS_KEY2 = "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"
)

type Result []interface{}
type Object interface{}

type D struct {
	ID             *string `json:"id"`
	TimezoneOffset *string `json:"timezone_offset"`
	Lon            *string `json:"lon"`
	Lat            *string `json:"lat"`
	Timezone       *string `json:"timezone"`
	CreatedAt      *string `json:"createdAt"`
}

type Ho struct {
	Hourly interface{}
	Info   *D
}

type Da struct {
	Daily interface{}
	Info  *D
}

type Na struct {
	Alerts interface{}
	Info   *D
}

type Hi struct {
	Hourly  interface{}
	Current interface{}
	Info    *D
}

type AP struct {
	ID        *string `json:"id"`
	List      interface{}
	Coord     interface{}
	CreatedAt *string `json:"createdAt"`
}

var se *session.Session

func init() {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_S3_REGION2),
		Credentials: credentials.NewStaticCredentials(
			AWS_S3_ACCESS_KEY_ID2,
			AWS_S3_SECRET_ACCESS_KEY2,
			""),
	})
	if err != nil {
		log.Fatal(err)
	}
	se = s
}

func FindHourlyForecast2Days() (d []*Ho, err error) {
	svc := dynamodb.New(se)
	input := &dynamodb.ScanInput{
		TableName: aws.String("hourly-forecast-2days"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	for _, v := range result.Items {
		var arr Result
		str := *v["hourly"].S
		str = strings.ReplaceAll(str, "'", "\"")
		if err = json.Unmarshal([]byte(str), &arr); err != nil {
			return nil, err
		}
		d = append(d, &Ho{
			Hourly: arr,
			Info: &D{
				ID:             v["id"].S,
				TimezoneOffset: v["timezone_offset"].N,
				Lon:            v["lon"].S,
				Lat:            v["lat"].S,
				Timezone:       v["timezone"].S,
				CreatedAt:      v["createdAt"].S,
			},
		})
	}

	return d, err
}

func FindDailyForecast7days() (d []*Da, err error) {
	svc := dynamodb.New(se)
	input := &dynamodb.ScanInput{
		TableName: aws.String("daily-forecast-7days"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	for _, v := range result.Items {
		var arr Result
		str := *v["daily"].S
		str = strings.ReplaceAll(str, "'", "\"")
		if err = json.Unmarshal([]byte(str), &arr); err != nil {
			return nil, err
		}
		d = append(d, &Da{
			Daily: arr,
			Info: &D{
				ID:             v["id"].S,
				TimezoneOffset: v["timezone_offset"].N,
				Lon:            v["lon"].S,
				Lat:            v["lat"].S,
				Timezone:       v["timezone"].S,
				CreatedAt:      v["createdAt"].S,
			},
		})
	}

	return d, err
}

func FindNationalWeatherAlerts() (d []*Na, err error) {
	svc := dynamodb.New(se)
	input := &dynamodb.ScanInput{
		TableName: aws.String("national-weather-alerts"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	for _, v := range result.Items {
		var arr Result
		str := *v["alerts"].S
		str = strings.ReplaceAll(str, "'", "\"")
		if err = json.Unmarshal([]byte(str), &arr); err != nil {
			return nil, err
		}
		d = append(d, &Na{
			Alerts: arr,
			Info: &D{
				ID:             v["id"].S,
				TimezoneOffset: v["timezone_offset"].N,
				Lon:            v["lon"].S,
				Lat:            v["lat"].S,
				Timezone:       v["timezone"].S,
				CreatedAt:      v["createdAt"].S,
			},
		})
	}

	return d, err
}

func FindHistoricalWeather5Days() (d []*Hi, err error) {
	svc := dynamodb.New(se)
	input := &dynamodb.ScanInput{
		TableName: aws.String("previous-5days"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	for _, v := range result.Items {
		var arrH Result
		str := *v["hourly"].S
		str = strings.ReplaceAll(str, "'", "\"")
		if err = json.Unmarshal([]byte(str), &arrH); err != nil {
			return nil, err
		}

		var objC Object
		str2 := *v["current"].S
		str2 = strings.ReplaceAll(str2, "'", "\"")
		if err = json.Unmarshal([]byte(str2), &objC); err != nil {
			return nil, err
		}
		d = append(d, &Hi{
			Hourly:  arrH,
			Current: objC,
			Info: &D{
				ID:             v["id"].S,
				TimezoneOffset: v["timezone_offset"].N,
				Lon:            v["lon"].S,
				Lat:            v["lat"].S,
				Timezone:       v["timezone"].S,
				CreatedAt:      v["createdAt"].S,
			},
		})
	}

	return d, err
}

func FindAirPollution() (d []*AP, err error) {
	svc := dynamodb.New(se)
	input := &dynamodb.ScanInput{
		TableName: aws.String("air-pollution-forecast"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}
	}

	for _, v := range result.Items {
		var arrL Result
		str := *v["list"].S
		str = strings.ReplaceAll(str, "'", "\"")
		if err = json.Unmarshal([]byte(str), &arrL); err != nil {
			return nil, err
		}

		var objC Object
		str2 := *v["coord"].S
		str2 = strings.ReplaceAll(str2, "'", "\"")
		if err = json.Unmarshal([]byte(str2), &objC); err != nil {
			return nil, err
		}
		d = append(d, &AP{
			ID:        v["id"].S,
			List:      arrL,
			Coord:     objC,
			CreatedAt: v["createdAt"].S,
		})
	}

	return d, err
}
