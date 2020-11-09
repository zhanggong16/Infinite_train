package config

import "time"

// PollingConfig ...
type PollingConfig struct {
	PollingTimeoutMin			 time.Duration `toml:"polling_timeout_min"`
	PollingTimeoutLow            time.Duration `toml:"polling_timeout_low"`
	PollingTimeoutMedium         time.Duration `toml:"polling_timeout_medium"`
	PollingTimeoutHigh           time.Duration `toml:"polling_timeout_high"`
	PollingTimeoutLarge          time.Duration `toml:"polling_timeout_large"`
	PollingTimeoutHuge           time.Duration `toml:"polling_timeout_huge"`
	PollingCreateBackupTimeout   time.Duration `toml:"polling_create_backup_timeout"`
	PollingDownloadBackupTimeout time.Duration `toml:"polling_download_backup_timeout"`
	PollingIntervalLow           time.Duration `toml:"polling_interval_low"`
	PollingIntervalMedium        time.Duration `toml:"polling_interval_medium"`
	PollingIntervalHigh          time.Duration `toml:"polling_interval_high"`
	PollingCreateBackupInterval  time.Duration `toml:"polling_create_backup_interval"`
}

// RetryConfig ...
type RetryConfig struct {
	RetryCountLow       int           `toml:"retry_count_low"`
	RetryCountMedium    int           `toml:"retry_count_medium"`
	RetryCountHigh      int           `toml:"retry_count_high"`
	RetryCountLarge     int           `toml:"retry_count_large"`
	RetryCountHuge      int           `toml:"retry_count_huge"`
	RetryCountMax       int           `toml:"retry_count_max"`
	RetryIntervalLow    time.Duration `toml:"retry_interval_low"`
	RetryIntervalMedium time.Duration `toml:"retry_interval_medium"`
	RetryIntervalHigh   time.Duration `toml:"retry_interval_high"`
}