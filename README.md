[![Go Reference](https://pkg.go.dev/badge/github.com/adminlove520/vulnDb-Notifier.svg)](https://pkg.go.dev/github.com/adminlove520/vulnDb-Notifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/adminlove520/vulnDb-Notifier)](https://goreportcard.com/report/github.com/adminlove520/vulnDb-Notifier)

# vulnDb Notifier

- 该工具从 [vuldb.com](https://vuldb.com/?) 抓取 CVE  feed，根据关键词进行过滤，并通过 Slack 或 Discord 通知你所关注的技术或产品的最新 CVE 漏洞。

## 功能特性

- 使用 [gofeed](https://github.com/mmcdole/gofeed) 解析来自 [vuldb.com](https://vuldb.com/?rss.recent) 的 RSS feed。
- 根据定义的关键词过滤 feed。
- 将过滤后的 CVE 存储在数据库中。
- 每当有新的 CVE 插入数据库时，发送 Slack 或 Discord 通知。

## 安装

确保 Go 环境已正确配置
```
go install github.com/adminlove520/vulnDb-Notifier/cmd/CVENotifier@latest
```

## 使用方法

1.  设置环境变量，配置你的 webhook URL。例如：

    ```bash
    # Slack webhook（可选）
    export SLACK_WEBHOOK=https://hooks.slack.com/services/<id>/<id>
    # Discord webhook（可选）
    export DISCORD_WEBHOOK=https://discord.com/api/webhooks/<id>/<token>
    ```
    
    必须设置至少一个 webhook（Slack 或 Discord）。

2.  在 `config.yaml` 中设置配置：

    ```yaml
    # 推送模式：daily（推送所有RSS内容）或 keyword（根据关键词过滤）
    # 默认为 daily
    push_mode: daily

    keywords:
      - Floodlight
      - wordpress
    ```

3.  定期运行该工具（例如每几小时），以获取最新的 feed 并接收新 CVE 的通知。建议为此设置一个 cron 作业。

    ```bash
    vulnDb-Notifier -config config.yaml
    ```

cron 作业示例
```
0 * * * * user vulnDb-Notifier -config config.yaml 2>&1 | tee -a vulnDb-Notifier.log
```

## Slack 通知
![Slack notification](slack.png)

## 开发计划

- [x] 从 https://vuldb.com/?rss.recent 获取 RSS feed
- [x] 如果标题中包含任何关键词，则过滤 feed
- [x] 如果标题中找到关键词，则将数据存储在数据库中
- [x] 如果插入操作成功，则发送 Slack 消息
- [x] 如果插入操作成功，则发送带有 embed 格式的 Discord 消息

## 包结构

项目现在组织为以下包：

*   `cmd/CVENotifier`：包含主应用程序逻辑。
*   `internal/config`：包含配置加载逻辑。
*   `internal/rss`：包含 RSS feed 解析逻辑。
*   `internal/slack`：包含 Slack 通知逻辑。
*   `internal/discord`：包含带有 embed 支持的 Discord 通知逻辑。
*   `internal/util`：包含实用函数，如 HTML 标签移除。
*   `internal/db`：包含数据库操作逻辑。
*   `internal/errors`：包含自定义错误类型。

## 错误处理

项目现在使用自定义错误类型来提供更具描述性的错误消息。

## 配置文件说明

`config.yaml` 文件用于配置工具的行为，包括推送模式和关注的关键词。

### 推送模式配置

```yaml
# 推送模式：daily（推送所有RSS内容）或 keyword（根据关键词过滤）
# 默认为 daily
push_mode: daily
```

- **daily 模式**：推送所有从 RSS 获取的 CVE 内容，不进行关键词过滤
- **keyword 模式**：只推送标题中包含指定关键词的 CVE 内容

### 关键词配置

```yaml
keywords:
  - wordpress
  - drupal
  - apache
  - nginx
```

当 `push_mode` 设置为 `keyword` 时，工具会检查每个 CVE 标题是否包含这些关键词，如果包含，则会存储并通知你。

### 完整配置示例

```yaml
# 推送模式配置
push_mode: daily

# 关键词配置（当 push_mode 为 keyword 时生效）
keywords:
  - wordpress
  - Floodlight
  - apache
  - nginx
  - drupal
  - mysql
  - postgresql
  - redis
  - mongodb
  - jenkins
  - kubernetes
  - docker
  - aws
  - azure
  - google cloud
  - linux
  - windows
  - macos
```