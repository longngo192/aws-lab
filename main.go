package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	DEFAULT_TIME     = "2006012150405"
	CONFIG_FILE_NAME = "s3_infor.txt"
	REGION           = "us-west-2"
)

var (
	S3Client   *s3.Client
	bucketName string
)

func checkErr(err error) {
	if err != nil {
		log.Panicf("%s", err)
	}
}
func getConfigFile(fileName string) string {
	content, isEmpty := readFile(fileName)
	if isEmpty {
		createNewConfigFile()
		content = getConfigFile(fileName)
	}
	return content
}
func readFile(fileName string) (string, bool) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDONLY, 0644)
	checkErr(err)
	content, err := io.ReadAll(f)
	checkErr(err)
	f.Sync()
	if len(content) == 0 {
		fmt.Println("empty")
		return "", true
	}
	defer f.Close()
	return string(content), false
}
func createNewConfigFile() {
	bucketName1 := fmt.Sprintf("%s-%s", "bucket", time.Now().Format(DEFAULT_TIME))
	f, err := os.OpenFile(CONFIG_FILE_NAME, os.O_APPEND|os.O_WRONLY, 0644)
	checkErr(err)
	_, err1 := f.WriteString(bucketName1)
	checkErr(err1)
	f.Sync()
	f.Close()
}
func init() {
	bucketName = getConfigFile(CONFIG_FILE_NAME)
	fmt.Printf("\nbkname:%s\n", bucketName)
	fmt.Printf("%s", bucketName)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(REGION))
	checkErr(err)
	S3Client = s3.NewFromConfig(cfg)
}
func isExitedBucket(bucketName2 string) bool {
	listBucket, err := S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	checkErr(err)
	// fmt.Printf("%v", listBucket)
	fmt.Printf("%v, %v", &listBucket.Buckets[0].Name, *listBucket.Buckets[0].Name)
	for _, b := range listBucket.Buckets {
		if *b.Name == bucketName2 {
			return true
		}
	}
	return false
}
func getOrCreateBucket(bucketName3 string) {
	exited := isExitedBucket(bucketName3)
	if !exited {
		_, err := S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{Bucket: &bucketName3, CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraintUsWest2}})
		checkErr(err)
	}
	output, err := S3Client.GetBucketAcl(context.TODO(), &s3.GetBucketAclInput{Bucket: &bucketName3})
	checkErr(err)
	fmt.Printf("%s", output)
}
func main() {
	getOrCreateBucket(bucketName)

	// AlcOutput, err := S3Client.GetBucketAcl(context.TODO(), &s3.GetBucketAclInput{Bucket: aws.String(bucketName)})
	// checkErr(err)
	// fmt.Printf("%v", AlcOutput)

}
