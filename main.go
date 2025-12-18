package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	AccessKeyId     string
	SecretAccessKey string
	Region          string
	Endpoint        string
	Iam             bool
	Bucket          string
}

var (
	testfilename = "foo.bar"
)

func main() {
	s3cfg := ParseIni()

	var minioClient *minio.Client
	var err error

	if s3cfg.Iam {
		// log.Printf("Metadata Endpoint Env: %s", os.Getenv("AWS_EC2_METADATA_SERVICE_ENDPOINT"))
		// log.Printf("Metadata Disabled Env: %s", os.Getenv("AWS_EC2_METADATA_DISABLED"))
		err = os.Unsetenv("AWS_EC2_METADATA_SERVICE_ENDPOINT")
		fmt.Printf("Unset environment VAR: AWS_EC2_METADATA_SERVICE_ENDPOINT ==> %v\n", err)
		err = os.Unsetenv("AWS_EC2_METADATA_DISABLED")
		fmt.Printf("Unset environment VAR: AWS_EC2_METADATA_DISABLED ==> %v\n", err)

		creds := credentials.NewIAM("")
		minioClient, err = minio.New(s3cfg.Endpoint, &minio.Options{
			Creds:  creds,
			Secure: true,
		})
		log.Printf("使用IAM初始化 S3 客户端: %#v", minioClient)
	} else {
		minioClient, err = minio.New(s3cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s3cfg.AccessKeyId, s3cfg.SecretAccessKey, ""),
			Region: s3cfg.Region,
			Secure: true,
		})
		log.Printf("使用AccessKeyId及SecretAccessKey初始化 S3 客户端: %#v", minioClient)
	}

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("成功初始化 S3 客户端: %#v", minioClient)

	WriteBytes(s3cfg, minioClient)
	ReadBytes(s3cfg, minioClient)
}

func ParseIni() S3Config {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalf("Failed to read ini file: %v", err)
	}

	s3cfg := S3Config{
		AccessKeyId:     "",
		SecretAccessKey: "",
		Region:          "cn-northwest-1",
		Endpoint:        "s3.cn-northwest-1.amazonaws.com.cn",
		Iam:             false,
		Bucket:          "",
	}
	if cfg.Section("aws").HasKey("aws_access_key_id") {
		s3cfg.AccessKeyId = cfg.Section("aws").Key("aws_access_key_id").String()
	}

	if cfg.Section("aws").HasKey("aws_secret_access_key") {
		s3cfg.SecretAccessKey = cfg.Section("aws").Key("aws_secret_access_key").String()
	}

	if cfg.Section("aws").HasKey("s3_region") {
		s3cfg.Region = cfg.Section("aws").Key("s3_region").String()
	}

	if cfg.Section("aws").HasKey("s3_endpoint") {
		s3cfg.Endpoint = cfg.Section("aws").Key("s3_endpoint").String()
	}

	if cfg.Section("aws").HasKey("aws_iam") {
		s3cfg.Iam, err = cfg.Section("aws").Key("aws_iam").Bool()
		if err != nil {
			log.Panicf("parse aws_iam failed: %v\n", err)
		}
	}

	if cfg.Section("aws").HasKey("s3_bucket") {
		s3cfg.Bucket = cfg.Section("aws").Key("s3_bucket").String()
	}

	fmt.Printf("\nParsed s3 config: %+v\n", s3cfg)
	return s3cfg
}

func WriteBytes(s3cfg S3Config, minioClient *minio.Client) {
	bucketName := s3cfg.Bucket
	objectName := testfilename
	contentType := "text/plain"
	fileContent := "hello, world!"

	reader := strings.NewReader(fileContent)
	objectSize := int64(reader.Len())

	ctx := context.Background()
	info, err := minioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Successfully write object %s of size %d bytes\n", objectName, info.Size)
}

func ReadBytes(s3cfg S3Config, minioClient *minio.Client) {
	bucketName := s3cfg.Bucket
	objectName := testfilename

	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("Error retrieving object:", err)
	}
	defer object.Close() // Important: close the object when done

	fileBytes, err := io.ReadAll(object)
	if err != nil {
		log.Fatalln("Error reading object content:", err)
	}

	fileContent := string(fileBytes)

	fmt.Printf("Successfully read file: %s length=%v\n", objectName, len(fileBytes))
	fmt.Printf("\nFile content:\n%s\n", fileContent)
}

func DeleteFile(s3cfg S3Config, minioClient *minio.Client) {
	bucketName := s3cfg.Bucket
	objectName := testfilename
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln("Error reading object content:", err)
	}
	fmt.Printf("\nFile deleted: %v\n", objectName)
}
