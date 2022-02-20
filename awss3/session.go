package awss3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

var (
	sessionS3  *s3.S3
	bucketName string
)

func init() {
	viper.AutomaticEnv()
	bucketName = viper.GetString("BUCKET_NAME_AWS")
	// credenciales tomadas desde env.sh o variables del sistema
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("REGION_AWS")),
		Credentials: credentials.NewStaticCredentials(
			viper.GetString("FELIX_AWS_ACCESS_KEY_ID"),
			viper.GetString("FELIX_AWS_SECRET_ACCESS_KEY"),
			""),
	})

	fmt.Println(viper.GetString("FELIX_AWS_ACCESS_KEY_ID"))
	fmt.Println(viper.GetString("FELIX_AWS_SECRET_ACCESS_KEY"))

	if err != nil {
		log.Println(err)
	}

	// Iniciamos
	sessionS3 = s3.New(sess)
}

// basepath debe terminar en /
func GuardarImagen(basepath, filename string, data *multipart.FileHeader) {

	file, err := data.Open()
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(basepath + filename),
		Body:          bytes.NewReader(fileData),
		ContentLength: aws.Int64(data.Size),
		ContentType:   aws.String(data.Header.Get("Content-Type")),
	}

	fmt.Println(params)

	resp, err := sessionS3.PutObject(params)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(awsutil.Prettify(resp))

}

// basepath debe terminar en /
func GetImage(basepath, filename string) (*s3.GetObjectOutput, error) {
	resp, err := sessionS3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(basepath + filename),
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}
