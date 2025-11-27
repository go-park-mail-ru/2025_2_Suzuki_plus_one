global:
  resolve_timeout: 5m

route:
  receiver: telegram-notifications
  group_by: ['alertname', 'service']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h

  routes:
    # Only send the "too many requests vs average" alert to Telegram
    - receiver: telegram-notifications
      matchers:
        - alertname = "HighRequestSpikePerService"

receivers:
  - name: 'telegram-notifications'
    telegram_configs:
      - bot_token: ${TELEGRAM_BOT_TOKEN}
        chat_id: ${TELEGRAM_CHAT_ID}
        parse_mode: 'Markdown'
        send_resolved: true
        message: |-
          🚨 *{{ .Status | toUpper }}* 🚨

          *Alert:* {{ .CommonLabels.alertname }}
          *Severity:* {{ .CommonLabels.severity }}

          {{ range .Alerts }}
            *Summary:* {{ .Annotations.summary }}
            *Description:* {{ .Annotations.description }}
            *Starts at:* {{ .StartsAt }}
          {{ end }}
