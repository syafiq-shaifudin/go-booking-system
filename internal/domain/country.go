package domain

import (
	"time"
)

type Country struct {
	ID                         uint       `gorm:"primaryKey;column:id" json:"id"`
	Name                       *string    `gorm:"column:name;type:varchar(255)" json:"name,omitempty"`
	Shortname                  *string    `gorm:"column:shortname;type:varchar(255)" json:"shortname,omitempty"`
	CountryCode                *int       `gorm:"column:country_code" json:"country_code,omitempty"`
	CurrencyName               *string    `gorm:"column:currency_name;type:varchar(255)" json:"currency_name,omitempty"`
	CurrencyCode               *string    `gorm:"column:currency_code;type:varchar(255)" json:"currency_code,omitempty"`
	CurrencySymbol             *string    `gorm:"column:currency_symbol;type:varchar(255)" json:"currency_symbol,omitempty"`
	CurrencyRate               *float64   `gorm:"column:currency_rate" json:"currency_rate,omitempty"`
	Latitude                   *float64   `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude                  *float64   `gorm:"column:longitude" json:"longitude,omitempty"`
	TimezoneName               *string    `gorm:"column:timezone_name;type:varchar(255)" json:"timezone_name,omitempty"`
	NameVariant                *string    `gorm:"column:name_variant;type:varchar(255)" json:"name_variant,omitempty"`
	CreatedTime                *time.Time `gorm:"column:created_time" json:"created_time,omitempty"`
	LastModified               *string    `gorm:"column:last_modified;type:varchar(255)" json:"last_modified,omitempty"`
	PaypalTransactionSurcharge *float64   `gorm:"column:paypal_transaction_surcharge" json:"paypal_transaction_surcharge,omitempty"`
	PaypalPayoutFee            *float64   `gorm:"column:paypal_payout_fee" json:"paypal_payout_fee,omitempty"`
	NameZhCn                   *string    `gorm:"column:name_zh_cn;type:varchar(255)" json:"name_zh_cn,omitempty"`
	NameZhTw                   *string    `gorm:"column:name_zh_tw;type:varchar(255)" json:"name_zh_tw,omitempty"`
	ReferralReward             *int       `gorm:"column:referral_reward" json:"referral_reward,omitempty"`
	SearchRadius               *int       `gorm:"column:search_radius" json:"search_radius,omitempty"`
	GMT                        *string    `gorm:"column:gmt;type:varchar(255)" json:"gmt,omitempty"`
	NameCsCz                   *string    `gorm:"column:name_cs_cz;type:varchar(255)" json:"name_cs_cz,omitempty"`
	NameTh                     *string    `gorm:"column:name_th;type:varchar(255)" json:"name_th,omitempty"`
	NameJaJp                   *string    `gorm:"column:name_ja_jp;type:varchar(255)" json:"name_ja_jp,omitempty"`
	NameKoKr                   *string    `gorm:"column:name_ko_kr;type:varchar(255)" json:"name_ko_kr,omitempty"`
	NoDecimalCurrency          *int       `gorm:"column:no_decimal_currency" json:"no_decimal_currency,omitempty"`
	HasVatGst                  *int       `gorm:"column:has_vat_gst" json:"has_vat_gst,omitempty"`
	DepositOnly                *int       `gorm:"column:deposit_only" json:"deposit_only,omitempty"`
	PointsMultiplier           *int       `gorm:"column:points_multiplier" json:"points_multiplier,omitempty"`
	AllowShare                 *int       `gorm:"column:allow_share" json:"allow_share,omitempty"`
	ReferralHostReward         *int       `gorm:"column:referral_host_reward" json:"referral_host_reward,omitempty"`
	HostExpectedEarnings       *int       `gorm:"column:host_expected_earnings" json:"host_expected_earnings,omitempty"`
	PlatformFeePercent         *int       `gorm:"column:platform_fee_percent" json:"platform_fee_percent,omitempty"`
	Hits                       *int       `gorm:"column:hits" json:"hits,omitempty"`
}

func (Country) TableName() string {
	return "country"
}
