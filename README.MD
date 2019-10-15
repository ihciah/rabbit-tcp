# Rabbit TCP

![Build badge](https://github.com/ihciah/rabbit-tcp/workflows/Build/badge.svg) ![Docker pull badge](https://img.shields.io/docker/pulls/ihciah/rabbit)

A multi-connection TCP forwarder created for accelerating TCP connections

English | [中文](README_ZH.MD)

## Introduction

This project is created to support and forward TCP connections. It will split bytes into blocks and reconstruct them during forwarding.
All upper connections visible to all users are carried by `N` underlying physical connections, and even a single upper connection is forwarded through all underlying connections.

Due to the use of multiple connections, similar to the multi-threaded downloader, this project can accelerate the connection without any more traffic consumption (bilateral deployment is required).
This project the extension of the accelerating method of [ShadowSocks-magic](https://github.com/ihciah/go-shadowsocks-magic).

![Rabbit TCP](.github/resources/rabbit-tcp.svg)

## Get Started
This project can work in standalone proxy mode or inline in other Golang code.

Here are two common usage examples. For detailed documentation, please go to the [Wiki](https://github.com/ihciah/rabbit-tcp/wiki).

### Accelerate any TCP service in a standalone proxy mode
#### Server side
1. Install [Docker](https://docs.docker.com/install/linux/docker-ce/debian/#install-using-the-convenience-script) and [docker-compose](https://docs.docker.com/compose/install/)
2. Download and edit `docker-compose-server.yml` ([Link](https://github.com/ihciah/rabbit-tcp/raw/master/docker-compose-server.yml))
    ```yaml
    rabbit-server:
      image: ihciah/rabbit
      ports:
        - "$RABBIT_PORT:9891/tcp"
      environment:
        - MODE=s
        - PASSWORD=$RABBIT_PASSWORD
        - RABBITADDR=:9891
        - VERBOSE=2
      restart: always
    ```
   - `$RABBIT_PORT` replaced with RabbitTCP service port, which should be the same as the client side
   - `$RABBIT_PASSWORD`  replaced with RabbitTCP service password, which should be the same as the client side
3. Run in the same directory `docker-compose -f docker-compose-server.yml up -d`

#### Client side
1. Install [Docker](https://docs.docker.com/install/linux/docker-ce/debian/#install-using-the-convenience-script) and [docker-compose](https://docs.docker.com/compose/install/)
2. Download and edit `docker-compose-client.yml` ([Link](https://github.com/ihciah/rabbit-tcp/raw/master/docker-compose-client.yml))
    ```yaml
    rabbit-client:
      image: ihciah/rabbit
      ports:
        - "$CLIENT_LISTEN_PORT:9892/tcp"
      environment:
        - MODE=c
        - PASSWORD=$RABBIT_PASSWORD
        - RABBITADDR=$RABBIT_ADDR
        - LISTEN=:9892
        - DEST=$SERVICE_ADDR
        - TUNNELN=6
        - VERBOSE=2
      restart: always
    ```
   - `$RABBIT_ADDR` replaced with `ip_or_domain:port` , which should be the same as the server side
   - `$RABBIT_PASSWORD` replaced with RabbitTCP service password, which should be the same as the server side
   - `$SERVICE_ADDR` replaced with the `ip_or_domain:port` of the target service
   - `$CLIENT_LISTEN_PORT` replaced with local listen port, which can be dialed when connecting to the target service with acceleration.
   - Modify `TUNNELN` if you want to change count of physical connections
3. Run in the same directory `docker-compose -f docker-compose-server.yml up -d`

### Accelerate ShadowSocks service in a standalone proxy mode with plugin
The server-side configuration is the same as above. Please note that except for Rabbit TCP server, you have to [run ShadowSocks service](https://github.com/shadowsocks/shadowsocks-libev/blob/master/docker/alpine/docker-compose.yml) too.

Configure client-side as above is useful too if you want to run a forwarder on a server inside the blocked area.

However, run docker or daemon service on client devices is not so user-friendly. What recommended is the plugin built for ShadowSocks: [Rabbit Plugin](https://github.com/ihciah/rabbit-plugin).

1. [Download the latest Rabbit Plugin release](https://github.com/ihciah/rabbit-plugin/releases), unzip it and put it into the ShadowSocks folder(same as ShadowSocks.exe).
2. Fill in ShadowSocks client configuration(Take [ShadowSocks Windows](https://github.com/shadowsocks/shadowsocks-windows) as an example)
    - Server Addr: Rabbit TCP service ip/domain
    - Server Port: Rabbit TCP service port
    - Password: ShadowSocks password
    - Encryption: ShadowSocks Encryption
    - Plugin Program: The binary filename(without `.exe`)
    - Plugin Options: `serviceAddr=SHADOWSOCKS_ADDR;password=RABBIT_PASSWORD;tunnelN=4`
        - `SHADOWSOCKS_ADDR` is ShadowSocks server address(ip:port or domain:port), for example: `10.10.10.10:443`
        - `RABBIT_PASSWORD` is Rabbit TCP password
        - Modify `tunnelN` to change count of physical connections
3. Save

## Speedup Result

Environment:

- `Chrome <--> ShadowsocksWindows <--> RabbitTCP <==[ISP]==> RabbitTCP <--> ShadowsocksLibev`
- Local ISP: China Unicom - Shanghai
- Remote ISP: Amazon LightSail - Tokyo
- Physical Connections Count: 4


Speedup with Rabbit TCP([Link](https://www.speedtest.net/result/8667412671)):

![Speed with rabbit-tcp](.github/resources/SpeedWithRabbit.jpg)

Original ShadowSocks libev([Link](https://www.speedtest.net/result/8667415664)):

![Speed without rabbit-tcp](.github/resources/SpeedWithoutRabbit.jpg)
