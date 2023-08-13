package consts

const (
	SQSSignalsQueueUrl    = "SQS_SIGNALS_QUEUE_URL"
	SNSSignalsTopic       = "SNS_SIGNALS_TOPIC"
	SNSPriceRequestTopic  = "SNS_PRICE_REQUEST_TOPIC"
	SignalsTableName      = "SIGNALS_TABLE_NAME"
	TransactionsTableName = "TRANSACTIONS_TABLE_NAME"
	ConfigTableName       = "CONFIG_TABLE_NAME"

	BinanceAPIKey    = "BINANCE_API_KEY"
	BinanceSecretKey = "BINANCE_SECRET_KEY"
	BinanceUrl       = "BINANCE_URL"

	AwsRegion          = "AWS_REGION"
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	PulumiAccessToken  = "PULUMI_ACCESS_TOKEN"

	ProjectDir        = "PROJECT_DIR"
	ProjectOutputDir  = "PROJECT_OUTPUT_DIR"
	CloudfareApiToken = "CLOUDFLARE_API_TOKEN"
	DnsZone           = "DNS_ZONE"
	DnsRecord         = "DNS_RECORD"

	CloudfareApiKey            = "CLOUDFLARE_API_KEY"
	CloudfareApiEmail          = "CLOUDFLARE_API_EMAIL"
	GptApiKey                  = "GPT_API_KEY"
	TelegramBotTextToken       = "TELEGRAM_BOT_TEXT_TOKEN"
	TelegramBotImageToken      = "TELEGRAM_BOT_IMAGE_TOKEN"
	CacheTable                 = "CACHE_TABLE"
	SendImageQueue             = "SEND_IMAGE_QUEUE"
	TelegramWebhookTokenHeader = "x-telegram-bot-api-secret-token"
)
