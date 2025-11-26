# Grafana

```bash
# Forward grafana to localport
# You need to set up ssh key to do that
ssh -L 5000:localhost:5000 -N -T ubuntu@217.16.18.125

# Open the dashboard
firefox localhost:5000
# Note that login and passwords are set it .env
```

## Dashboards

- [Node Exporter Full](https://grafana.com/grafana/dashboards/1860-node-exporter-full/)
