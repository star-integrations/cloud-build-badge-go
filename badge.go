package badge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	RepoSource Repo `json:"repoSource"`
}

// Repo gitトリガ
type Repo struct {
	ProjectID  string `json:"projectId"`
	RepoName   string `json:"repoName"`
	BranchName string `json:"branchName"`
}

// CloudBuildBadgeStatus バッジを更新する。
func CloudBuildBadgeStatus(ctx context.Context, m PubSubMessage) error {
	log.Printf("pub/sub msg:[%#+v]", string(m.Data))
	pubsubMsg := &CloudBuildPubSubMsg{}
	if err := json.Unmarshal(m.Data, pubsubMsg); err != nil {
		log.Printf("pub/sub message unmarshal error:[%s]", err.Error())
		return nil
	}
	log.Printf("pub/sub msg:[%#+v]", pubsubMsg)

	repo := pubsubMsg.Source.RepoSource
	if repo.ProjectID == "" || repo.RepoName == "" || repo.BranchName == "" {
		log.Printf("pub/sub message repo infomation empty")
		return nil
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("cloud storage client error:[%s]", err.Error())
		return nil
	}
	badgeBucket := client.Bucket("cloud-build-my-badge")
	buildResultBucket := client.Bucket("cloud-build-my-result")

	buildResultObject := buildResultBucket.Object(fmt.Sprintf("%s/%s/%s", repo.ProjectID, repo.RepoName, repo.BranchName))

	var badge *storage.ObjectHandle
	if pubsubMsg.Status == "SUCCESS" {
		log.Printf("build success!")
		badge = badgeBucket.Object("success.svg")
	} else {
		log.Printf("build failure...")
		badge = badgeBucket.Object("failure.svg")
	}
	if _, err := buildResultObject.CopierFrom(badge).Run(ctx); err != nil {
		log.Printf("badge object copy error:[%s]", err.Error())
		return nil
	}

	log.Printf("done")
	return nil
}
