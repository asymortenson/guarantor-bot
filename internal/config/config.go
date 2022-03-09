package config

import (
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