# Channel Monitor

[Chinese 中文](README-zh.md) | English 英文

## Introduction

Channel Monitor is a tool designed for monitoring OneAPI/NewAPI channels. It directly reads the channels data table and tests the availability of each channel's models at regular intervals. Based on whether the requests are successful, it updates the available models and writes them to the data table, thus achieving automated monitoring of the channels. This ensures high availability of the overall OneAPI/NewAPI and minimizes the number of error returns.

## Features

- [x] Directly read and write to the OneAPI/NewAPI database
- [x] Test the availability of each model in the channels
- [x] Automatically fetch available models from upstream
- [x] Automatically update the available models in the database for each channel
- [x] Support exclusion of channels and models from monitoring
- [x] Support configurable intervals
- [x] Support multiple database types, including MySQL, SQLite, PostgreSQL, and SQL Server
- [x] Concurrent testing
- [x] Request rate limiting (per-second RPS and per-minute RPM)
- [x] Support Uptime Kuma, push URL during testing to visualize model availability
- [x] Support update notifications via SMTP email and Telegram Bot
- [x] Support both JSON and YAML configuration formats

## Installation

### Binary

Download the latest version of the binary file from the [Releases](https://github.com/DullJZ/ChannelMonitor/releases) page. After configuring `config.json` or `config.yaml` in the same directory, you can run it. It is recommended to use tools like `screen` or `nohup` to run it in the background.

```bash
mkdir ChannelMonitor && cd ChannelMonitor
wget https://github.com/DullJZ/ChannelMonitor/releases/download/v0.1.0/ChannelMonitor_linux_amd64
chmod +x ChannelMonitor_linux_amd64
# Download and modify the configuration file (choose JSON or YAML format)
wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.json -O config.json
# or use YAML format
# wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.yaml -O config.yaml
nano config.json  # or nano config.yaml
screen -S ChannelMonitor
./ChannelMonitor_linux_amd64
```

### Docker

```bash
docker pull dulljz/channel-monitor
# Download and modify the configuration file (choose JSON or YAML format)
wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.json -O config.json
# or use YAML format
# wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.yaml -O config.yaml
nano config.json  # or nano config.yaml
# If using the host's database, you can simply use the host mode,
# and use localhost:3306 as the database address
docker run -d --name ChannelMonitor -v ./config.json:/app/config.json --net host dulljz/channel-monitor
# If using YAML format
# docker run -d --name ChannelMonitor -v ./config.yaml:/app/config.yaml --net host dulljz/channel-monitor
# If using an SQLite database, mount the database file
# docker run -d --name ChannelMonitor -v ./config.json:/app/config.json -v /path/to/database.db:/app/database.db dulljz/channel-monitor
```

### Docker Compose

```yaml
version: '3'
services:
  channel-monitor:
    image: dulljz/channel-monitor
    volumes:
      - ./config.json:/app/config.json
      # or use YAML format
      # - ./config.yaml:/app/config.yaml
      # If using an SQLite database, mount the database file
      # - /path/to/database.db:/app/database.db
    # If using the host's database, you can simply use the host mode,
    # and use localhost:3306 as the database address
    network_mode: host
```

```bash
# Download and modify the configuration file (choose JSON or YAML format)
wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.json -O config.json
# or use YAML format
# wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.yaml -O config.yaml
nano config.json  # or nano config.yaml
docker-compose up -d
```

## Configuration

The configuration file can be either `config.json`, `config.yaml`, or `config.yml` in the same directory. The program will automatically detect and use the available configuration file in the order of `config.yaml` -> `config.yml` -> `config.json`.

### JSON Format

<details>
<summary>Click to expand/collapse JSON configuration example</summary>

```json
{
  "oneapi_type": "oneapi",
  "exclude_channel": [5],
  "exclude_model": ["advanced-voice", "minimax_s2v-01", "minimax_video-01", "minimax_video-01-live2d"],
  "models": ["gpt-3.5-turbo", "gpt-4o"],
  "force_models": false,
  "force_inside_models": false,
  "time_period": "1h",
  "max_concurrent": 5,
  "rps": 5,
  "rpm": 0,
  "timeout": 10,
  "db_type": "YOUR_DB_TYPE",
  "db_dsn": "YOUR_DB_DSN",
  "do_not_modify_db": false,
  "base_url": "http://localhost:3000",
  "system_token": "YOUR_SYSTEM_TOKEN",
  "uptime-kuma": {
    "status": "disabled",
    "model_url": {
      "gpt-3.5-turbo": "https://demo.kuma.pet/api/push/A12n43563?status=up&msg=OK&ping=",
      "gpt-4o": "https://demo.kuma.pet/api/push/ArJd2BOUJN?status=up&msg=OK&ping="
    },
    "channel_url": {
      "5": "https://demo.kuma.pet/api/push/ArJd2BOUJN?status=up&msg=OK&ping="
    }
  },
  "notification": {
    "smtp": {
      "enabled": false,
      "host": "smtp.example.com",
      "port": 25,
      "username": "your-email@example.com",
      "password": "your-password",
      "from": "sender@example.com",
      "to": "recipient@example.com"
    },
    "webhook": {
      "enabled": false,
      "type": "telegram",
      "telegram": {
        "chat_id": "YOUR_CHAT_ID",
        "retry": 3
      },
      "secret": "YOUR_WEBHOOK_SECRET"
    }
  }
}
```

</details>

### YAML Format

<details>
<summary>Click to expand/collapse YAML configuration example</summary>

```yaml
oneapi_type: oneapi
exclude_channel: [5]
exclude_model: 
  - advanced-voice
  - minimax_s2v-01
  - minimax_video-01
  - minimax_video-01-live2d
models: 
  - gpt-3.5-turbo
  - gpt-4o
force_models: false
force_inside_models: false
time_period: 1h
max_concurrent: 5
rps: 5
rpm: 0
timeout: 10
db_type: YOUR_DB_TYPE
db_dsn: YOUR_DB_DSN
do_not_modify_db: false
base_url: http://localhost:3000
system_token: YOUR_SYSTEM_TOKEN
uptime-kuma:
  status: disabled
  model_url:
    gpt-3.5-turbo: https://demo.kuma.pet/api/push/A12n43563?status=up&msg=OK&ping=
    gpt-4o: https://demo.kuma.pet/api/push/ArJd2BOUJN?status=up&msg=OK&ping=
  channel_url:
    "5": https://demo.kuma.pet/api/push/ArJd2BOUJN?status=up&msg=OK&ping=
notification:
  smtp:
    enabled: false
    host: smtp.example.com
    port: 25
    username: your-email@example.com
    password: your-password
    from: sender@example.com
    to: recipient@example.com
  webhook:
    enabled: false
    type: telegram
    telegram:
      chat_id: YOUR_CHAT_ID
      retry: 3
    secret: YOUR_WEBHOOK_SECRET
```

</details>

Configuration explanation:
- oneapi_type: Type of OneAPI, including oneapi, newapi, onehub (reserved field, currently has no effect)
- exclude_channel: IDs of channels to exclude from monitoring
- exclude_model: IDs of models to exclude from monitoring
- models: List of models, used only when unable to retrieve models from the channel (/v1/models)
- force_models: If true, only the above models will be tested, and channel models will not be fetched. Default is false
- force_inside_models: If true, only the models set in OneAPI will be tested, and the model list will not be fetched. Default is false. If force_models is true, this option is invalid.
- time_period: Interval for testing model availability, recommended not less than 30 minutes, accepts time formats s, m, h
- max_concurrency: Maximum number of concurrent tests within a channel, default is 5
- rps: Requests per second within a channel, default is 5 (when neither rps nor rpm is configured). Set to 0 to disable
- rpm: Requests per minute within a channel, default is 0 (disabled). When both rps and rpm are set, the stricter limit takes effect
- timeout: Request timeout (seconds), default is 10
- db_type: Database type, including mysql, sqlite, postgres, sqlserver
- db_dsn: Database DSN string, the format varies by database type. Examples below
- do_not_modify_db: If true, the available models in the database will not be modified. Default is false
- base_url: The base URL for OneAPI/NewAPI/OneHub. If using host mode, you can directly use http://localhost:3000. Currently, only OneHub requires this field.
- system_token: System token, currently only required for OneHub.
- uptime-kuma: Configuration for Uptime Kuma. The status can be `enabled` or `disabled`. The model_url and channel_url are the availability Push URLs for models and channels.
- notification: Configuration for update notifications, including SMTP email and Telegram Bot
- notification.smtp: SMTP email configuration, where enabled is `true` or `false`, host is the SMTP server address, port is the server port, username and password are login credentials, from is the sender's email, and to is the recipient's email
- notification.webhook: Webhook configuration, where enabled is `true` or `false`, type currently only supports `telegram`, telegram contains Telegram Bot settings, chat_id is your telegram ID, retry is the number of retry attempts, and secret is the API key

### MySQL

```json
{
  "db_type": "mysql",
  "db_dsn": "username:password@tcp(host:port)/dbname"
}
```

### SQLite

```json
{
  "db_type": "sqlite",
  "db_dsn": "/path/to/database.db"
}
```

### PostgreSQL

```json
{
  "db_type": "postgres",
  "db_dsn": "host=host port=port user=username password=password dbname=dbname sslmode=disable"
}
```

### SQL Server

```json
{
  "db_type": "sqlserver",
  "db_dsn": "sqlserver://username:password@host:port?database=dbname"
}
```
