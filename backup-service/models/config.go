package models

type Config struct {
	AWS_ENDPOINT_URL_S3   string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	BUCKET_NAME           string
	DB_PREFIX             string
	DISCORD_WEBHOOK_URL   string
	DISCORD_SECRET_KEY    string
}
