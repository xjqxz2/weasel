const PACK_PING_HEALTH = "PING/HEALTH"
const PACK_PONG_HEALTH = "PONG/HEALTH"

const WSClient = function (serialNo, serverDomain, onReceiver) {
    this.serialNo = serialNo
    this.serverDomain = serverDomain
    this.onReceiver = onReceiver || null
    this.websocket = null
    this.connectRetry = null
    this.isRetry = true
    this.health = null

    this.connect()
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

WSClient.prototype.connect = function (prototype) {
    let that = this
    prototype = prototype || "ws://"

    //  Open the retry status
    this.setRetryStatus(true)
    this.websocket = new WebSocket(
        prototype + this.serverDomain + "/dev/conn?serial_no=" + this.serialNo + "&serial_name=UnKnowDevice"
    )

    //  Open the websocket
    this.websocket.onopen = function () {
        console.log("设备 " + that.serialNo + "已成功连接到服务器")

        //  Create Heart Pack
        if (that.health == null) {
            that.health = setInterval(function () {
                that.websocket.send(PACK_PING_HEALTH)
            }, 15000)

            console.log("已启用心跳检测")
        }

        if (that.connectRetry != null) {
            clearInterval(that.connectRetry)
            that.connectRetry = null
            console.log("已关闭重连定时器")
        }
    }

    this.websocket.onclose = function () {
        console.log("设备 " + that.serialNo + "已断开连接")

        that.initResource()

        if (that.connectRetry == null && that.isRetry) {
            //  开启重连模式
            that.connectRetry = setInterval(function () {
                that.connect()
            }, 1500)

            console.log("已开启重连模式...")
        }
    }

    this.websocket.onmessage = function (e) {
        if (e.data === PACK_PONG_HEALTH) {
            console.log("收到心跳回复 PONG")
            return
        }

        that.onReceiver(e)
    }

    this.websocket.onerror = function (e) {
        console.log("Websocket 发生错误 -> " + e)
    }
}

WSClient.prototype.close = function () {
    //  close retry connection
    this.setRetryStatus(false)
    this.initResource()

    this.websocket.close()
    this.websocket = null
}

WSClient.prototype.initResource = function () {
    if (this.health != null) {
        clearInterval(this.health)
        this.health = null
        console.log("已停用心跳检测")
    }

    if (this.connectRetry != null) {
        clearInterval(this.connectRetry)
        this.connectRetry = null
    }
}