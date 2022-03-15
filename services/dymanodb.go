package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

const (
	AWS_S3_REGION2            = "ap-northeast-1"
	AWS_S3_ACCESS_KEY_ID2     = "AKIAZRF7Z3PJVVHW3DNY"
	AWS_S3_SECRET_ACCESS_KEY2 = "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"
)

type D struct {
	ID             string `json:"id"`
	TimezoneOffset string `json:"timezone_offset"`
	Lon            string `json:"lon"`
	Lat            string `json:"lat"`
	Timezone       string `json:"timezone"`
	CreatedAt      string `json:"createdAt"`
}

type Ho struct {
	Hourly string `json:"hourly"`
	Info   *D
}

type Da struct {
	Daily string `json:"daily"`
	Info  *D
}

type Na struct {
	Alerts string `json:"alerts"`
	Info   *D
}

type Hi struct {
	Hourly  string `json:"hourly"`
	Current string `json:"current"`
	Info    *D
}

type AP struct {
	ID        string `json:"id"`
	List      string `json:"list"`
	Coord     string `json:"coord"`
	CreatedAt string `json:"createdAt"`
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
		d = append(d, &Ho{
			Hourly: v["hourly"].String(),
			Info: &D{
				ID:             v["id"].String(),
				TimezoneOffset: v["timezone_offset"].String(),
				Lon:            v["lon"].String(),
				Lat:            v["lat"].String(),
				Timezone:       v["timezone"].String(),
				CreatedAt:      v["createdAt"].String(),
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
		d = append(d, &Da{
			Daily: v["daily"].String(),
			Info: &D{
				ID:             v["id"].String(),
				TimezoneOffset: v["timezone_offset"].String(),
				Lon:            v["lon"].String(),
				Lat:            v["lat"].String(),
				Timezone:       v["timezone"].String(),
				CreatedAt:      v["createdAt"].String(),
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
		d = append(d, &Na{
			Alerts: v["alerts"].String(),
			Info: &D{
				ID:             v["id"].String(),
				TimezoneOffset: v["timezone_offset"].String(),
				Lon:            v["lon"].String(),
				Lat:            v["lat"].String(),
				Timezone:       v["timezone"].String(),
				CreatedAt:      v["createdAt"].String(),
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
		d = append(d, &Hi{
			Hourly:  v["hourly"].String(),
			Current: v["current"].String(),
			Info: &D{
				ID:             v["id"].String(),
				TimezoneOffset: v["timezone_offset"].String(),
				Lon:            v["lon"].String(),
				Lat:            v["lat"].String(),
				Timezone:       v["timezone"].String(),
				CreatedAt:      v["createdAt"].String(),
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
		d = append(d, &AP{
			ID:        v["id"].String(),
			List:      v["list"].String(),
			Coord:     v["coord"].String(),
			CreatedAt: v["createdAt"].String(),
		})
	}

	return d, err
}
