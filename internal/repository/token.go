package repository

type Bucket string

const (
	AccessTokens  Bucket = "access_tokens"
	RequestTokens Bucket = "request_tokens"
)

type TokenRepository interface {
	Put(chatId int64, token string, bucket Bucket) error
	Get(chatId int64, bucket Bucket) (string, error)
}

func GetBuckets() []Bucket {
	return []Bucket{
		AccessTokens,
		RequestTokens,
	}
}
