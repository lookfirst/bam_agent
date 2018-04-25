[![Build Status](https://travis-ci.org/blockassets/bam_agent.svg?branch=master)](https://travis-ci.org/blockassets/bam_agent)

# Block Assets Manager (BAM) Agent

This is an agent that is intended to be installed on miners to help facilitate management of them via HTTP as well 
as monitoring for issues internal to the miner. The documentation below assumes you have the latest version of the
agent installed.

Currently, the agent targets the BW-L21 miner, but in the future we will target all miners. Pull requests welcome.

Thanks to [HyperBitShop.io](https://hyperbitshop.io) for sponsoring this project.

## Running (defaults):

``
./bam_agent-linux-arm -port=1111 -no-update=false
``

A config file is writen to `/etc/bam_agent.json`.

## Automatic updates

By default, the BAM Agent will automatically attempt to self update from the Github 
[latest release](https://github.com/blockassets/bam_agent/releases) tab. It chooses a random hour of the day to update. 
This way, if you have a number of machines, they will not all DDOS Github and your network. You can override the update 
behavior to not perform any updates by setting `-no-update=true` when starting the agent.

Automatic updates brings security into question. You don't want someone being able to install a binary on your
machine that has not been developed in a secure fashion. We employ a strict process around this that is completely
transparent on github.

The authors have 2FA setup on their Github accounts and the project requires 2FA. All commits are required to be 
signed. All merges to master require a PR and passing unit tests. Builds are tested to ensure the binary starts up 
and serves the /status page. Travis CI automates the updates to the release page. Updates have their downloaded 
binaries hash checked. Updates are zero downtime.

## Install onto miner

The [releases tab](https://github.com/blockassets/bam_agent/releases) has `master` binaries cross compiled for ARM 
suitable for running on the miner. These are built automatically on 
[Travis](https://travis-ci.org/blockassets/bam_agent).

Download the [latest release](https://github.com/blockassets/bam_agent/releases) and copy the gunzipped 
`bam_agent-linux-arm` binary to `/usr/bin`

```
gunzip bam_agent-linux-arm.gz
chmod ugo+x bam_agent-linux-arm
scp bam_agent-linux-arm root@MINER_IP:/usr/bin
```

Create `/etc/systemd/system/bam_agent.service`

```
ssh root@MINER_IP "echo '
[Unit]
Description=bam_agent
After=init.service

[Service]
Type=simple
ExecStart=/usr/bin/bam_agent-linux-arm
Restart=always
RestartSec=4s
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
' > /etc/systemd/system/bam_agent.service"
```

Enable the service:

```
ssh root@MINER_IP "systemctl enable bam_agent; systemctl start bam_agent"
```

## Building from source

We recommend that you download prebuilt binaries from the releases tab. However, if you would like to build your own...

* Install golang (OSX, use brew)
* Install [dep](https://github.com/golang/dep) (OSX, use brew)
* `git clone https://github.com/blockassets/bam_agent.git`
* `make dep`
* `make arm` (Builds binary for arm)

## API

### `GET /cgminer/start`

```
Starts cgminer via systemd
```

### `GET /cgminer/quit`

```
Stops cgminer via cgminer API call (systemd will restart it). Does not work with BW-L21.
```

### `GET /cgminer/restart`

```
Stops cgminer via cgminer API call (systemd will restart it)
```

### `PUT /config/frequency`

Send PUT request with json body:

```
{"frequency": "684"}
```

Restarts cgminer.

### `GET /config/pools`

```
{"pool1": "", "pool2": "", "pool3": ""}
```

### `PUT /config/pools`

Send PUT request with json body:

```
{"pool1": "", "pool2": "", "pool3": ""}
```

Restarts cgminer.

### `PUT /config/dhcp`

Updates `/usr/app/conf.default` and `/etc/network/interfaces`

```
NO BODY NECESSARY
```

Call `/reboot` to make the changes take effect

### `PUT /config/ip`

Updates `/usr/app/conf.default` and `/etc/network/interfaces`

```
{"ip": "10.10.0.11", "mask": "255.255.252.0", "gateway": "10.10.0.1", "dns": "8.8.8.8"}
```

Call `/reboot` to make the changes take effect

### `PUT /config/location`

Store the physical location of the miner. This is saved in `/etc/bam_agent.json` and exposed in the `/status` call.

```
{"facility": "", "rack": "", "row": "", "shelf": 1, "position": 1}
```

### `GET /status`

```
{
  "agent": "39892e1 2018-03-06 02:06:09",
  "miner": "value in /usr/app/version.txt",
  "uptime": "0s",
  "mac": "ab:bc:32:b2:81:79",
  "location": {
    "facility": "",
    "rack": "",
    "row": "",
    "shelf": 1,
    "position": 1
  }
}
```

### `GET /reboot`

Reboots the miner. Obviously be careful with this one. :-) Do note that we send HTTP expiry headers so that
if you do run this from your browser (not really advised), it won't get cached.

### `POST /update`

Upload a compressed archive of files and execute an enclosed shell script. This allows one to easily distribute
updates to the miners.

1. Create a folder and put a script in it called `update.sh`
1. .tar.gz the folder
1. HTTP file upload the compressed archive with the form parameter name of `file`
1. Optionally add a `script` form parameter to override the default `update.sh` name

Returns the stdout/stderr of the script.

Example: `curl -F "script=update.sh" -F "file=@/tmp/update.tar.gz" http://IP:1111/update`

## Monitors

Monitors allow us to execute code periodically.
Monitors are configured by editing the `/etc/bam_agent.json` file. This file is created when the agent first starts.

### Accepted shares

**Enabled by default.** If the miner has not accepted any shares after 5m, reboot. This works around a bug where the 
miner software stops submitting shares to the pool, yet continues doing work.

### High load

**Enabled by default.** If the 5m average load is above 5, `reboot -f` the miner. This works around a bug where the load 
spikes and the miner stops submitting shares to the pool.

### High temp

**Enabled by default.** Runs every 5m and checks to see if the temperature is over 100c. If so, it uses systemd to 
shut cgminer down. A reboot will enable things again.

### Low memory

**Enabled by default.** Runs every 1m and checks to see if the available system memory is under 140mb. If so, it reboots.
This works around a bug where cgminer will randomly eat up available memory and freeze the machine.

### Quit cgminer

Disabled by default. Periodically quit the miner app to free up memory and start fresh.

### Reboot

Disabled by default. Periodically reboot the entire miner to free up memory and start fresh.

## Prometheus exporters

[Prometheus](https://prometheus.io) is an amazing metrics collection system. Combine it with 
[Grafana](https://grafana.com) and you have an extremely powerful visualization and performance tracking tool for your miners.

Included in the agent are two [prometheus exporters](https://prometheus.io/docs/instrumenting/exporters/):

1. [node_exporter](https://github.com/prometheus/node_exporter) (operating system metrics)
2. [cgminer_exporter](https://github.com/blockassets/cgminer_exporter) (cgminer metrics)

By including them in the agent, you don't need to install those binaries separately.

In order to configure prometheus, use configuration like this:

`prometheus.yml`

```yml
  - job_name: 'cgminer_exporter'
    metrics_path: /metrics/cgminer_exporter
    file_sd_configs:
      - files:
        - 'workers.json'
  - job_name: 'node_exporter'
    metrics_path: /metrics/node_exporter
    file_sd_configs:
      - files:
        - 'workers.json'
```

`workers.json`

```
[{
        "targets": [
          "10.10.0.11:1111",
          "10.10.0.12:1111",
          "10.10.0.13:1111"
        ]
}]
```
