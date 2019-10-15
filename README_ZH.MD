# Rabbit TCP

![Build badge](https://github.com/ihciah/rabbit-tcp/workflows/Build/badge.svg) ![Docker pull badge](https://img.shields.io/docker/pulls/ihciah/rabbit)

一个为加速而生的多连接 TCP 转发器

[English](README.MD) | 中文

## 简介

这个项目用来承载并转发 TCP 流量，并在转发时对流量进行分拆与重组。
所有用户可见的上层连接均由 `N` 条底层物理连接承载，即便是单条用户连接也会经所有底层连接进行转发。

由于使用了多条连接，与多线程下载器类似，本项目可以在不多消耗流量的情况下对连接进行加速（需双边部署）。
这个项目是 [ShadowSocks-magic](https://github.com/ihciah/go-shadowsocks-magic) 中加速方案的扩展。

![Rabbit TCP](.github/resources/rabbit-tcp.svg)

## 开始使用
本项目既可以以独立代理模式工作，也可以内嵌在其他 Golang 代码中提供服务。

以下列举两种常用用法，详细的文档请移步 [Wiki](https://github.com/ihciah/rabbit-tcp/wiki)。

### 以独立代理模式加速任意 TCP 服务
#### 服务端
1. 安装 [Docker](https://docs.docker.com/install/linux/docker-ce/debian/#install-using-the-convenience-script) 和 [docker-compose](https://docs.docker.com/compose/install/)
2. 下载并修改 `docker-compose-server.yml` 文件([链接](https://github.com/ihciah/rabbit-tcp/raw/master/docker-compose-server.yml))
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
   - `$RABBIT_PORT` 替换为 RabbitTCP 服务端口，与 Client 保持一致即可
   - `$RABBIT_PASSWORD` 替换为 RabbitTCP 服务密码，与 Client 保持一致即可
3. 在该文件同目录运行 `docker-compose -f docker-compose-server.yml up -d`

#### 客户端
1. 安装 [Docker](https://docs.docker.com/install/linux/docker-ce/debian/#install-using-the-convenience-script) 和 [docker-compose](https://docs.docker.com/compose/install/)
2. 下载并修改 `docker-compose-client.yml` 文件([链接](https://github.com/ihciah/rabbit-tcp/raw/master/docker-compose-client.yml))
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
   - `$RABBIT_ADDR` 替换为 RabbitTCP 服务 `地址:端口` ，端口与服务端保持一致即可
   - `$RABBIT_PASSWORD` 替换为 RabbitTCP 服务密码，与服务端保持一致即可
   - `$SERVICE_ADDR` 替换为待加速的服务 `地址:端口` 
   - `$CLIENT_LISTEN_PORT` 替换为本地监听端口，后续若要加速访问目标服务，连接该端口即可
   - 修改 `TUNNELN` 可以控制底层物理连接数
3. 在该文件同目录运行 `docker-compose -f docker-compose-server.yml up -d`

### 以独立代理 + 插件模式加速 ShadowSocks 服务
服务端配置同上。请注意 Rabbit TCP 并不包含 ShadowSocks，你仍然需要[运行 ShadowSocks](https://github.com/shadowsocks/shadowsocks-libev/blob/master/docker/alpine/docker-compose.yml)。

客户端可以按照上述配置，如运行在国内服务器上用作中转。

在用户终端上，运行 Docker 或直接运行后台服务较为麻烦。推荐使用为 ShadowSocks 定制的插件 [Rabbit Plugin](https://github.com/ihciah/rabbit-plugin)。

1. [下载最新 Rabbit Plugin 发布版](https://github.com/ihciah/rabbit-plugin/releases)并解压缩后放入 ShadowSocks 同目录
2. 在 ShadowSocks 客户端中填写信息(以 [ShadowSocks Windows](https://github.com/shadowsocks/shadowsocks-windows) 为例)
    - 服务器地址填写 Rabbit TCP 服务器
    - 服务器端口填写 Rabbit TCP 端口
    - 密码填写 ShadowSocks 密码
    - 加密方式填写 ShadowSocks 加密方式
    - 插件程序填写刚刚放入的插件程序文件名(不包含 `.exe` 后缀)
    - 插件选项填写 `serviceAddr=SHADOWSOCKS_ADDR;password=RABBIT_PASSWORD;tunnelN=4`
        - `SHADOWSOCKS_ADDR` 替换为 ShadowSocks 服务地址（包含端口），如 `10.10.10.10:443`
        - `RABBIT_PASSWORD` 替换为 Rabbit TCP 密码，与服务端保持一致
        - 修改 `tunnelN` 对应数值控制底层物理连接个数
3. 保存即可

## 加速效果

测试环境:

- `Chrome <--> ShadowsocksWindows <--> RabbitTCP <==[ISP]==> RabbitTCP <--> ShadowsocksLibev`
- 本地运营商: 中国联通 - 上海
- 远程运营商: Amazon LightSail - 东京
- 底层连接数: 4


使用 Rabbit TCP 加速([Link](https://www.speedtest.net/result/8667412671)):

![Speed with rabbit-tcp](.github/resources/SpeedWithRabbit.jpg)

使用原版 ShadowSocks-libev([Link](https://www.speedtest.net/result/8667415664)):

![Speed without rabbit-tcp](.github/resources/SpeedWithoutRabbit.jpg)

