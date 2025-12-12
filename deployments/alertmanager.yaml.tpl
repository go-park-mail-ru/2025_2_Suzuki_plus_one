global:
  resolve_timeout: 5m

route:
  receiver: telegram-notifications
  group_by: ['alertname', 'service']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1000000h

  routes:
    # Only send the "too many requests vs average" alert to Telegram
    # However top receiver is set as default so it is unnecessary
    - receiver: telegram-notifications
      matchers:
        - alertname = "HighRequestSpikePerService"

receivers:
  - name: 'telegram-notifications'
    telegram_configs:
      - bot_token: ${TELEGRAM_BOT_TOKEN} # Have to replaced
        chat_id: ${TELEGRAM_CHAT_ID}     # Have to replaced
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
