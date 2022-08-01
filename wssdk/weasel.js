const PACK_PING_HEALTH = "PING/HEALTH"
const PACK_PONG_HEALTH = "PONG/HEALTH"

const WSClient = function (serialNo, serverDomain, onReceiver) {
    let that = this;

    this.serialNo = serialNo
    this.serverDomain = serverDomain
    this.onReceiver = onReceiver || null
    this.websocket = null
    this.connectRetry = null
    this.isRetry = true

    let connectionWS = function (prototype) {
        prototype = prototype || "ws://"

        that.websocket = new WebSocket(
            prototype + that.serverDomain + "/dev/conn?serial_no=" + that.serialNo + "&serial_name=UnKnowDevice"
        )

        //  Open the websocket
        that.websocket.onopen = function () {
            console.log("设备 " + that.serialNo + "已成功连接到服务器")

            if (that.connectRetry != null) {
                clearInterval(that.connectRetry)
                that.connectRetry = null
                console.log("已关闭重连定时器")
            }
        }

        that.websocket.onclose = function () {
            console.log("设备 " + that.serialNo + "已断开连接")

            //  断开连接时清除心跳监测
            clearInterval(that.health)


            if (that.connectRetry == null && that.isRetry) {
                //  开启重连模式
                that.connectRetry = setInterval(function () {
                    connectionWS()
                }, 1500)

                console.log("已开启重连模式...")
            }
        }

        that.websocket.onmessage = function (e) {
            if (e.data === PACK_PONG_HEALTH) {
                console.log("收到心跳回复 PONG")
                return
            }

            that.onReceiver(e)
        }

        that.websocket.onerror = function (e) {
            console.log("Websocket 发生错误 -> " + e)
        }
    }

    //  设置健康检查相关函数
    this.health = setInterval(function () {
        that.websocket.send(PACK_PING_HEALTH)
    }, 30000)


    connectionWS()
}

WSClient.prototype.send = function (serialNo, message) {
    const request = {
        "serial_no": serialNo,
        "message": message
    }

    fetch("http://" + this.serverDomain + "/msg/broadcast", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(request)
    }).then(function (res) {
        console.log(res)
        console.log("消息发送成功")
    }).catch(function (e) {
        console.log("发送消息失败" + e)
    })
}

WSClient.prototype.setRetryStatus = function (b) {
    this.isRetry = b
}