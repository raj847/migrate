package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	DBDriver = GetEnv("DB_DRIVER", "postgres")
	DBName   = GetEnv("DB_NAME", "local")
	DBHost   = GetEnv("DB_HOST", "localhost")
	DBPort   = GetEnv("DB_PORT", "5432")
	DBUser   = GetEnv("DB_USER", "root")
	DBPass   = GetEnv("DB_PASS", "")
	SSLMode  = GetEnv("SSL_MODE", "disable")

	DBNameCloud = GetEnv("DB_NAME_CLOUD", "local")
	DBHostCloud = GetEnv("DB_HOST_CLOUD", "localhost")
	DBPortCloud = GetEnv("DB_PORT_CLOUD", "5432")
	DBUserCloud = GetEnv("DB_USER_CLOUD", "root")
	DBPassCloud = GetEnv("DB_PASS_CLOUD", "")

	REDISHost = GetEnv("REDIS_HOST")
	REDISPass = GetEnv("REDIS_PASS")
	REDISPort = GetEnv("REDIS_PORT")

	REDISHostLocal = GetEnv("REDIS_HOST_LOCAL")
	REDISPassLocal = GetEnv("REDIS_PASS_LOCAL")
	REDISPortLocal = GetEnv("REDIS_PORT_LOCAL")

	MONGOHost = GetEnv("MONGO_HOST")
	MONGOPort = GetEnv("MONGO_PORT")
	MONGODB   = GetEnv("MONGO_DB")
	MONGOUser = GetEnv("MONGO_USER")
	MONGOPass = GetEnv("MONGO_PASS")

	DEFAULTChannelRedisParking        = GetEnv("DEFAULT_CHANNEL_REDIS_PARKING")
	DEFAULTChannelRedisParkingDeposit = GetEnv("DEFAULT_CHANNEL_REDIS_PARKING_DEPOSIT")

	AMQPServerUrl = GetEnv("AMQP_SERVER_URL")

	SALT_KEY              = GetEnv("SALT_KEY")
	MERCHANT_KEY          = GetEnv("MERCHANT_KEY")
	MERCHANT_KEY_APPS2PAY = GetEnv("MERCHANT_KEY_APPS2PAY")
	URL_MASTER_PPKGW      = GetEnv("URL_MASTER_PPKGW")
	UrlPaymentA2P         = GetEnv("URL_PAYMENT_A2P")
	ADDRESS               = GetEnv("ADDRESS")
	QRISPayment           = GetEnv("QRIS_PAYMENT", "N")

	BASEParkingURL              = GetEnv("BASE_PARKING_URL")
	UrlServer                   = GetEnv("URL_SERVER")
	UrlWebPaymentOnline         = GetEnv("URL_WEB_PAYMENT_ONLINE")
	BaseParkingURLPaymentOnline = GetEnv("BASE_PARKING_URL_PAYMENT_ONLINE")
	EXCSyncProductCode          = GetEnv("EXC_SYNC_PRODUCT_CODE")
	MEMBER_BY                   = GetEnv("MEMBER_BY")
	NORMAL_SF                   = GetEnv("NORMAL_SF", "Y")
	SPECIAL_PRODUCT_CODE        = GetEnv("WHITELIST_PRODUCT_CODE")
	BasicAuth                   = GetEnv("BASIC_AUTH")
)

func GetEnv(key string, value ...string) string {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error Load file .env not found")
	}

	if os.Getenv(key) != "" {
		return os.Getenv(key)
	} else {
		if len(value) > 0 {
			return value[0]
		}
		return ""
	}
}
