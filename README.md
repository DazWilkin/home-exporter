# Google Home Exporter

![Noise Level Graph](./noiselevel.png)

See [Metrics](#metrics) for the supported metrics

```
# HELP home_device_info_noise_level TBD
# TYPE home_device_info_noise_level gauge
home_device_info_noise_level -89
# HELP home_device_info_signal_level TBD
# TYPE home_device_info_signal_level gauge
home_device_info_signal_level -49
# HELP home_device_info_uptime TBD
# TYPE home_device_info_uptime gauge
home_device_info_uptime 34250.670738
# HELP home_exporter_build_info A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter
# TYPE home_exporter_build_info counter
home_exporter_build_info{git_commit="",go_version="go1.12.5",os_version=""} 1
# HELP home_exporter_start_time Exporter start time in Unix epoch seconds
# TYPE home_exporter_start_time gauge
home_exporter_start_time 1.578172839e+09
```

## Run

The Exporter exposes metrics on `:9999` and `/metrics` by default

You'll need to provide the Exporter with your device's endpoint. The Google Home app settings is one way to identify the device's IP address. The API is served on port `:8008`.

### Docker

```bash
DEVICE="[YOUR-DEVICE-IP]:8008"
PORT="9999"
docker run \
--rm --interactive --tty \
--name=${DEVICE} \
--publish=${PORT}:9999
dazwilkin/home-exporter:aadf69552f21993caca7d419ab81b4b80b31f05a
  --device=192.168.86.25:8008 \
  --endpoint=:9999 \
  --metricsPath=/metrics
```

You can then browse `http://localhost:9999/metrics`

### Docker Compose

Includes the Exporter (exposed on the host's `:9999`), Prometheus (`:9090`) and cAdvisor (`:8085`).

You'll need to configure `docker-compose.yml` with your device's configuration. Replace the `${VARIABLE}` values in the following. `9999` refers to the endpoint on which the Exporter (!) will serve `/metrics`. You may change this with `--endpoint=${ENDPOINT}`. If you change the port, you'll need to change this in the 2 other occurrences of `9999`. The port `9997` refers to the host port where the exporter is published. You may change this value too.

```YAML
${EXPORTER_NAME}:
  image: dazwilkin/home-exporter:aadf69552f21993caca7d419ab81b4b80b31f05a
  container_name: ${EXPORTER_NAME}
  command:
  - "--device=${DEVICE_IP}:8008"
  - "--endpoint=:9999"
  expose:
  - "9999"
  ports:
  - 9997:9999
```

Then:

```bash
docker-compose up
```

You may:

```bash
docker-compose ps
      Name                    Command                  State                    Ports              
---------------------------------------------------------------------------------------------------
cadvisor           /usr/bin/cadvisor -logtostderr   Up (healthy)   0.0.0.0:8085->8080/tcp          
prometheus         /bin/prometheus --config.f ...   Up             0.0.0.0:9090->9090/tcp          
exporter-1         /home-exporter --device=19 ...   Up             9402/tcp, 0.0.0.0:9998->9999/tcp
exporter-2         /home-exporter --device=19 ...   Up             9402/tcp, 0.0.0.0:9997->9999/tcp
```

And:

```bash
docker-compose logs exporter-1
Attaching to exporter-1
exporter-1    | 2020/01/05 00:00:00 [main] Exporting metrics for Google Home device (192.168.1.24:8008)
exporter-1    | 2020/01/05 00:00:00 [main] Server starting (:9999)
exporter-1    | 2020/01/05 00:00:00 [main] metrics served on: /metrics
```

Then:

+ [Exporter](http://localhost:9999)
+ [Prometheus](http://localhost:9090)
+ [cAdvisor](http://localhost:8085)

## Metrics

| Name | Type | Help
| ---- | ---- | ----
| `home_device_info_up`           | Gauge   | If 1 the Home device is accessible, 0 otherwise
| `home_device_info_noise_level`  | Gauge   | Noise Level dB
| `home_device_info_signal_level` | Gauge   | Signal Level db
| `home_device_info_uptime`       | Gauge   | Device Uptime in seconds
| `home_exporter_start_time`      | Gauge   | Exporter start time in Unix epoch seconds
| `home_exporter_build_info`      | Counter | A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter

### Labels

In addition to the default labels (`instance`, `job`), the Exporter provides:

+ `name`
+ `build_version`
+ `version`
