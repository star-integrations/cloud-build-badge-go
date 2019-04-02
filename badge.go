package badge

import (
	"context"
	"log"
	"os"
)

// PubSubMessage ...
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// CloudBuildBadgeStatus ...
func CloudBuildBadgeStatus(ctx context.Context, m PubSubMessage) error {
	repo := os.Getenv("REPO")
	branch := os.Getenv("BRANCH")
	log.Printf("repo:[%s]", repo)
	log.Printf("branch:[%s]", branch)
	log.Printf("pub/sub msg:[%s]", string(m.Data))
	return nil
}
