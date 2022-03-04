package backuper

import (
	"fmt"

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

func (a *AlibabaOSS) UploadFromPath(srcPath, fileName string) error {
	bucket, err := a.cfg.Client.Bucket(a.cfg.Bucket)
	if err != nil {
		return fmt.Errorf("get bucket (%s): %w", a.cfg.Bucket, err)
	}

	err = bucket.PutObjectFromFile(srcPath, fileName)
	if err != nil {
		return fmt.Errorf("put object (srcPath:%s fileName:%s) : %w", srcPath, fileName, err)
	}

	return nil
}
