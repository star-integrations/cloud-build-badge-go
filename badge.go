package badge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

// PubSubMessage pub/subメッセージを受信する
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// CloudBuildPubSubMsg Pub/Subメッセージを表す
type CloudBuildPubSubMsg struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Status    string `json:"status"`
	Source    Source `json:"source"`
}

// Source トリガ元
type Source struct {
	RepoSource    Repo    `json:"repoSource"`
	StorageSource Storage `json:"storageSource"`
}

// Repo gitトリガ
type Repo struct {
	ProjectID  string `json:"projectId"`
	RepoName   string `json:"repoName"`
	BranchName string `json:"branchName"`
}

// Storage CloudStorageトリガ
type Storage struct {
	Bucket     string `json:"bucket"`
	Object     string `json:"object"`
	Generation string `json:"generation"`
}

// CloudBuildBadgeStatus バッジを更新する。
func CloudBuildBadgeStatus(ctx context.Context, m PubSubMessage) error {
	// pub/submessage Unmarshal
	pubsubMsg := &CloudBuildPubSubMsg{}
	if err := json.Unmarshal(m.Data, pubsubMsg); err != nil {
		log.Printf("pub/sub message unmarshal error:[%s]", err.Error())
		return nil
	}

	// CloudStorage client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("cloud storage client error:[%s]", err.Error())
		return nil
	}

	// build badge bucket
	badgeBucketName := os.Getenv("BUILD_RESULT_BADGE_BUCKET")
	if badgeBucketName == "" {
		log.Printf("badge bucket empty")
		return nil
	}
	badgeBucket := client.Bucket(badgeBucketName)

	// build result badge object
	var buildResultObject *storage.ObjectHandle
	source := pubsubMsg.Source
	if source.RepoSource.ProjectID != "" && source.RepoSource.RepoName != "" && source.RepoSource.BranchName != "" {
		repo := source.RepoSource
		badgeObject := fmt.Sprintf("%s/%s/%s/badge.svg", repo.ProjectID, repo.RepoName, repo.BranchName)
		buildResultObject = badgeBucket.Object(badgeObject)
	} else if source.StorageSource.Bucket != "" && source.StorageSource.Generation != "" {
		log.Printf("pub/sub message repo infomation empty")
		repo := source.StorageSource
		badgeObject := fmt.Sprintf("%s/%s/badge.svg", repo.Bucket, repo.Generation)
		buildResultObject = badgeBucket.Object(badgeObject)
	}

	// build badge object
	var badgeObject *storage.ObjectHandle
	if pubsubMsg.Status == "SUCCESS" {
		log.Printf("build success!")
		badgeObject = badgeBucket.Object("success.svg")
	} else {
		log.Printf("build failure...")
		badgeObject = badgeBucket.Object("failure.svg")
	}

	// result badge copy from bucket
	if _, err := buildResultObject.CopierFrom(badgeObject).Run(ctx); err != nil {
		log.Printf("badge object copy error:[%s]", err.Error())
		return nil
	}

	if err := buildResultObject.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Printf("fail set public url to result object:[%s]", err.Error())
		return nil
	}

	log.Printf("done")
	return nil
}
