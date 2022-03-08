package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
	ApiKey        string
	Fee float64
	Wallet string
	AdminChannel int64
	ExchangeChannel int64
	RequestLink string `mapstructure:"request_link"`
	Db struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Messages Messages
	Buttons Buttons
	PublicPages PublicPages
}

type Messages struct {
	Responses
}


type Buttons struct {
	Contacts 			string `mapstructure:"contacts"`
	BuyAd 				string `mapstructure:"buy_ad"`
	BackToPrevious		 string `mapstructure:"back_to_previous"`
	Paid                   string `mapstructure:"paid_button"`
	DeclinePaid            string `mapstructure:"decline_paid_button"`
	ChoosePublicPage       string `mapstructure:"choose_public_page_button"`
	CreateRequest          string `mapstructure:"create_request_button"`
	ApprovePublicPage      string `mapstructure:"approve_public_page_button"`
	DeclinePublicPage      string `mapstructure:"decline_public_page_button"`
	Approved 			  string `mapstructure:"approved_button"`
	Rejected 			  string `mapstructure:"rejected_button"`
}

type PublicPages struct {
	Programmer 			  string `mapstructure:"programmer_button"`
	AboutTon 			  string `mapstructure:"aboutton_button"`
}


type Responses struct {
	Contacts 				 	string `mapstructure:"contacts"`
	PutPublicPage             string `mapstructure:"put_public_page"`
	AfterSubmittingPublicPage string `mapstructure:"after_submitting_public_page"`
	RejectPublicPage          string `mapstructure:"reject_public_page"`
	AfterPaymentResponse      string `mapstructure:"after_payment_response"`
	FailedPaymentResponse     string `mapstructure:"failed_payment_response"`
	SuccessPaymentResponse    string `mapstructure:"success_payment_response"`
	PaymentMessage            string `mapstructure:"payment_message"`
	AlreadyPaid 			  string `mapstructure:"already_paid"`
	Signature 			  	  string `mapstructure:"signature"`

}

func Init() (*Config, error) {
	var cfg Config
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("buttons", &cfg.Buttons); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("public_pages", &cfg.PublicPages); err != nil {
		return nil, err
	}


	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {

	os.Setenv("TOKEN", "5124531442:AAGIiZ73NJZXUSA_epez5BZyLdugjf0pwcw")
	os.Setenv("API_KEY", "498f0c3860ee799bb60650d5a7c0adfa2cf5adc877dc430cdd1adef1bf402888")
	os.Setenv("WALLET", "EQBtZgocpA8H1S3U1uCj-Uoz1BvuDESm2laaNVEexQZ7nQPp")
	os.Setenv("FEE", "0.2")
    os.Setenv("REQUEST_LINK", "@abouttonbot")
	os.Setenv("ADMIN_CHANNEL", "-1001688233561")
	os.Setenv("EXCHANGE_CHANNEL", "-1001523233085")

	os.Setenv("DB_DSN", "postgresql://doadmin:gkfa8123qylxVwAv@app-cc86af8b-3b5b-41f8-af11-e750fc477993-do-user-10807910-0.b.db.ondigitalocean.com:25060/telegrambot?sslmode=require")

	os.Setenv("DB_MAX_OPEN_CONNS", "25")
	os.Setenv("DB_MAX_IDLE_CONNS", "25")
	os.Setenv("DB_MAX_IDLE_TIME", "15m")




	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("api_key"); err != nil {
		return err
	}

	if err := viper.BindEnv("wallet"); err != nil {
		return err
	}

	if err := viper.BindEnv("fee"); err != nil {
		return err
	}
	if err := viper.BindEnv("admin_channel"); err != nil {
		return err
	}

	if err := viper.BindEnv("request_link"); err != nil {
		return err
	}

	if err := viper.BindEnv("exchange_channel"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_dsn"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_max_open_conns"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_max_idle_conns"); err != nil {
		return err
	}

	if err := viper.BindEnv("db_max_idle_time"); err != nil {
		return err
	}

	
	cfg.TelegramToken = viper.GetString("token")
	cfg.ApiKey = viper.GetString("api_key")
	cfg.Wallet = viper.GetString("wallet")
	cfg.ExchangeChannel = viper.GetInt64("exchange_channel")
	cfg.AdminChannel = viper.GetInt64("admin_channel")
	cfg.Fee = viper.GetFloat64("fee")
	cfg.Db.Dsn = viper.GetString("db_dsn")
	cfg.RequestLink = viper.GetString("request_link")
	cfg.Db.MaxOpenConns = viper.GetInt("db_max_open_conns")
	cfg.Db.MaxIdleConns = viper.GetInt("db_max_idle_conns")
	cfg.Db.MaxIdleTime = viper.GetString("db_max_idle_time")

	return nil
}