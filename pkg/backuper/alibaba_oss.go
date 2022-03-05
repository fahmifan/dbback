package backuper

import (
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AlibabaOSSCfg struct {
	Client *oss.Client
	Bucket string
}

type AlibabaOSS struct {
	cfg *AlibabaOSSCfg
}

func NewAlibabaOSS(cfg *AlibabaOSSCfg) *AlibabaOSS {
	return &AlibabaOSS{cfg: cfg}
}

func (a *AlibabaOSS) Upload(key string, rd io.Reader) error {
	bucket, err := a.cfg.Client.Bucket(a.cfg.Bucket)
	if err != nil {
		return fmt.Errorf("UploadFromPath: get bucket (%s): %w", a.cfg.Bucket, err)
	}

	err = bucket.PutObject(key, rd)
	if err != nil {
		return fmt.Errorf("UploadFromPath: put object (key: %s) : %w", key, err)
	}

	return nil
}
