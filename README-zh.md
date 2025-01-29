# Channel Monitor 渠道监控

Chinese 中文 | [English 英文](README.md)

## 介绍

Channel Monitor 是一个用于监控OneAPI/NewAPI渠道的工具，它直接读取channels数据表，每间隔一段时间测试每个渠道的模型可用性，根据请求是否成功更新可用模型，写入到数据表中，以此来实现渠道的自动监控，保证整体OneAPI/NewAPI的高可用，尽可能减少错误返回次数。

## 特性

- [x] 直接读写OneAPI/NewAPI的数据库
- [x] 测试渠道的每个模型可用性
- [x] 自动向上游获取可用模型
- [x] 自动更新数据库中的每个渠道可用模型
- [x] 支持排除不予监控的渠道和模型
- [x] 支持间隔时间配置
- [x] 支持多种数据库类型（MySQL、SQLite、PostgreSQL、SQL Server）
- [x] 并发测试
- [x] 支持Uptime Kuma， 在测试时Push URL来可视化模型可用性
- [x] 支持更新推送，包括SMTP邮件和Telegram Bot


## 安装

### 二进制

从[Releases](https://github.com/DullJZ/ChannelMonitor/releases)页面下载最新版本的二进制文件，在同一目录下配置`config.json`后运行即可。建议使用`screen`或`nohup`等工具后台运行。
注意⚠️：如果你需要使用SQLite数据库，请使用docker方案或自行启用CGO编译。

```bash
mkdir ChannelMonitor && cd ChannelMonitor
wget https://github.com/DullJZ/ChannelMonitor/releases/download/v0.1.0/ChannelMonitor_linux_amd64
chmod +x ChannelMonitor_linux_amd64
# 下载并修改配置文件
wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.json -O config.json
nano config.json
screen -S ChannelMonitor
./ChannelMonitor_linux_amd64
```

### Docker

```bash
docker pull dulljz/channel-monitor
# 下载并修改配置文件
wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.json -O config.json
nano config.json
# 如果使用宿主机的数据库，可以简单使用host模式，
# 并使用localhost:3306作为数据库地址
docker run -d --name ChannelMonitor -v ./config.json:/app/config.json -net host dulljz/channel-monitor
# 如果使用SQLite数据库，挂载数据库文件
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
      # 如果使用SQLite数据库，挂载数据库文件
      # - /path/to/database.db:/app/database.db
    # 如果使用宿主机的数据库，可以简单使用host模式，
    # 并使用localhost:3306作为数据库地址
    network_mode: host
```

```bash
# 下载并修改配置文件
wget https://raw.githubusercontent.com/DullJZ/ChannelMonitor/refs/heads/main/config_example.json -O config.json
nano config.json
docker-compose up -d
```

## 配置

配置文件使用同级目录下的`config.json`文件，格式如下：

```json
{
  "oneapi_type": "oneapi",
  "exclude_channel": [5],
  "exclude_model": ["advanced-voice", "minimax_s2v-01", "minimax_video-01", "minimax_video-01-live2d"],
  "models": ["gpt-3.5-turbo", "gpt-4o"],
  "force_models": false,
  "time_period": "1h",
  "db_type": "YOUR_DB_TYPE",
  "db_dsn": "YOUR_DB_DSN",
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

配置说明：
- oneapi_type: OneAPI的类型，包括oneapi、newapi、onehub（保留字段，暂时无影响）
- exclude_channel: 排除不予监控的渠道ID
- exclude_model: 排除不予监控的模型ID  
- models: 模型列表，仅当获取不到渠道的模型(/v1/models)时使用
- force_models: 如果为true，将强制只测试上述模型，不再获取渠道的模型，默认为false
- time_period: 模型可用性测试的时间间隔，建议不小于30分钟，接收的时间格式为s、m、h
- db_type: 数据库类型，包括mysql、sqlite、postgres、sqlserver
- db_dsn: 数据库DSN字符串，不同数据库类型的DSN格式不同，示例如下
- base_url: OneAPI/NewAPI/OneHub的基础URL，如果使用host模式，可以直接使用http://localhost:3000，目前只有OneHub需要填写
- system_token: 系统Token，目前只有OneHub需要填写
- uptime-kuma: Uptime Kuma的配置，status为`enabled`或`disabled`，model_url和channel_url为模型和渠道的可用性Push URL
- notification: 更新推送的配置，包括SMTP邮件和Telegram Bot
- notification.smtp: SMTP邮件配置，enabled为`true`或`false`，host为SMTP服务器地址，port为端口，username和password为登录凭证，from为发件人，to为收件人
- notification.webhook: Webhook配置，enabled为`true`或`false`，type目前仅支持`telegram`，telegram为Telegram Bot的配置，chat_id为聊天ID（填写你的telegram id），retry为重试次数，secret为API密钥

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

