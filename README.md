# WiFi Pineapple Demo (aka. ssid-jungle)

This project is a demonstration using the [Wifi Pineapple Mk7](https://shop.hak5.org/products/wifi-pineapple) to show just how much data exists floating around us.

The basic demonstration shows SSID broadcasts being captured by the WiFi Pineapple reconnaissance scan.

> _Note: this demonstration is designed to use the Hak5 [MK7AC WiFi Adapter](https://shop.hak5.org/products/mk7ac-wifi-adapter) attached to the WiFi Pineapple. Demonstration will need to be reconfigured to solely run in 2.4 GHz mode without it (see [backend/pineapple.go](backend/pineapple.go#L56))._

## Prerequisites

To install you're going to need docker to run [Docker](https://docs.docker.com/engine/install/).

### WiFi Pineapple Configuration

See [Configure the WiFi Pineapple](#configure-the-wifi-pineapple)

### macOS

For macOS, the WiFi Pineapple's ethernet driver needs to be installed/patched. **This is not officially supported by the WiFi Pineapple.** If you really want to do this, see [this article](https://docs.hak5.org/wifi-pineapple/faq/macos-support).

## Usage

### Build

You will need to build the docker containers locally. To do this, simply run:

```shell
docker compose build
```

### Running

You will need to plug in the WiFi Pineapple to your local machine and wait for it to start up. The best way to check if the WiFi Pineapple has booted and API is accessible is to navigate to the web UI at [172.16.42.1:1471](http://172.16.42.1:1471/). If this is accessible you're ready to go.

Then you can start up the docker containers with:

```shell
docker compose up
# ...or to start in the background
docker compose up -d
```

## Configuration

### SSID Filtering

The SSID Filtering is hacked together for now (please submit a PR with an update!). To modify the SSID filter list, go to [backend/pineapple.go (line 14)](backend/pineapple.go#L14) and modify the filter map to pick which SSIDs should be _hidden_.

You will have to rebuild the `backend` docker container to see the changes to this.

### Configure the WiFi Pineapple

Please ensure the WiFi Pineapple credentials are set to `root:root`. Otherwise, you can override this in [main.go](main.go#L28).

## Authors

- [BradHacker](https://github.com/BradHacker)

> _Disclaimer: This project is created for educational demonstration purposes only. This project's maintainer's, RIT, or anyone affiliated with the RIT CyberCorps&copy; Scholarship for Service program is not liable for misuse of this project._
