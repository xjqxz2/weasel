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

    //  设置健康检查相关函数
    this.health = setInterval(function () {
        that.websocket.send(PACK_PING_HEALTH)
    }, 15000)

    this.connectionWS()
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

WSClient.prototype.connectionWS = function (prototype) {
    prototype = prototype || "ws://"

    this.websocket = new WebSocket(
        prototype + that.serverDomain + "/dev/conn?serial_no=" + that.serialNo + "&serial_name=UnKnowDevice"
    )

    //  Open the websocket
    this.websocket.onopen = function () {
        console.log("设备 " + this.serialNo + "已成功连接到服务器")

        if (this.connectRetry != null) {
            clearInterval(this.connectRetry)
            this.connectRetry = null
            console.log("已关闭重连定时器")
        }
    }

    this.websocket.onclose = function () {
        console.log("设备 " + this.serialNo + "已断开连接")

        //  断开连接时清除心跳监测
        clearInterval(that.health)


        if (this.connectRetry == null && this.isRetry) {
            //  开启重连模式
            this.connectRetry = setInterval(function () {
                this.connectionWS()
            }, 1500)

            console.log("已开启重连模式...")
        }
    }

    this.websocket.onmessage = function (e) {
        if (e.data === PACK_PONG_HEALTH) {
            console.log("收到心跳回复 PONG")
            return
        }

        this.onReceiver(e)
    }

    this.websocket.onerror = function (e) {
        console.log("Websocket 发生错误 -> " + e)
    }
}

WSClient.prototype.close = function () {
    //  close retry connection
    this.setRetryStatus(false)

    clearInterval(this.health)

    if (this.connectRetry != null)
        clearInterval(this.connectRetry)

    this.websocket.close()
    this.websocket = null
}