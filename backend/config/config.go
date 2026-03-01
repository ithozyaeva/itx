package config

import (
	"log"
	"time"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	Database           DatabaseConfig
	JwtSecret          []byte
	Port               string
	TelegramToken      string
	TelegramMainChatID int64
	PublicDomain       string
	BackendDomain      string
	AllowedOrigins     string
	S3                 S3Config

	AlertReminderIntervalMinutes       int64
	AlertReminderFirstIntervalMinutes  int64
	AlertReminderSecondIntervalMinutes int64
	AlertReminderThirdIntervalMinutes  int64
	AlertScheduledTime                 string
	AlertScheduledHour                 int
	AlertScheduledMinute               int
}

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

var CFG *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	alertReminderInterval := viper.GetInt64("ALERT_REMINDER_INTERVAL_MINUTES")
	if alertReminderInterval == 0 {
		alertReminderInterval = 1440
	}
	
	alertReminderFirst := viper.GetInt64("ALERT_REMINDER_FIRST_INTERVAL_MINUTES")
	if alertReminderFirst == 0 {
		alertReminderFirst = 10080
	}
	
	alertReminderSecond := viper.GetInt64("ALERT_REMINDER_SECOND_INTERVAL_MINUTES")
	if alertReminderSecond == 0 {
		alertReminderSecond = 1440
	}
	
	alertReminderThird := viper.GetInt64("ALERT_REMINDER_THIRD_INTERVAL_MINUTES")
	if alertReminderThird == 0 {
		alertReminderThird = 60
	}

	var alertScheduledTime string
	var alertScheduledHour, alertScheduledMinute int
	
	if viper.IsSet("ALERT_SCHEDULED_TIME") {
		alertScheduledTime = viper.GetString("ALERT_SCHEDULED_TIME")
		parsedTime, err := time.Parse("15:04", alertScheduledTime)
		if err != nil {
			log.Printf("Warning: ALERT_SCHEDULED_TIME=%s is invalid (expected HH:MM format), using default 12:00", alertScheduledTime)
			alertScheduledTime = "12:00"
			alertScheduledHour = 12
			alertScheduledMinute = 0
		} else {
			alertScheduledHour = parsedTime.Hour()
			alertScheduledMinute = parsedTime.Minute()
		}
	} else {
		alertScheduledTime = "12:00"
		alertScheduledHour = 12
		alertScheduledMinute = 0
	}

	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "jwt_secret"
	}

	allowedOrigins := viper.GetString("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	CFG = &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JwtSecret:          []byte(jwtSecret),
		Port:               viper.GetString("PORT"),
		TelegramToken:      viper.GetString("TELEGRAM_BOT_TOKEN"),
		TelegramMainChatID: viper.GetInt64("TELEGRAM_MAIN_CHAT_ID"),
		PublicDomain:       viper.GetString("PUBLIC_DOMAIN"),
		BackendDomain:      viper.GetString("BACKEND_DOMAIN"),
		AllowedOrigins:     allowedOrigins,
		AlertReminderIntervalMinutes:       alertReminderInterval,
		AlertReminderFirstIntervalMinutes:  alertReminderFirst,
		AlertReminderSecondIntervalMinutes: alertReminderSecond,
		AlertReminderThirdIntervalMinutes:  alertReminderThird,
		AlertScheduledTime:                 alertScheduledTime,
		AlertScheduledHour:                 alertScheduledHour,
		AlertScheduledMinute:               alertScheduledMinute,
		S3: S3Config{
			Endpoint:  viper.GetString("S3_ENDPOINT"),
			Region:    viper.GetString("S3_REGION"),
			AccessKey: viper.GetString("S3_ACCESS_KEY"),
			SecretKey: viper.GetString("S3_SECRET_KEY"),
			Bucket:    viper.GetString("S3_BUCKET"),
		},
	}
}
