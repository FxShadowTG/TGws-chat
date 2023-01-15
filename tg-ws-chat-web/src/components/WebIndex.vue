<template>
    <div class="web-index-container">
        <el-container>
            <el-aside name="aside">
                <div style="margin-top:10px">
                    <el-input style="background: transparent;width: 85%;" v-model="inputWsUrl"
                        :autosize="{ minRows: 2, maxRows: 2 }" placeholder="请输入要连接的ws地址" maxlength="100" type="textarea"
                        resize="none" />
                    <el-button type="warning" style="margin-top:10px" @click="connWsUrl"
                        :disabled="connButtonStatus">连接</el-button>
                    <el-button type="danger" style="margin-top:10px" @click="closeConn" :disabled="closeConnButtonStatus">断开</el-button>
                </div>
                <div style="margin-top:30px;">
                    当前有 {{ userLists.length }} 位主机正连接该地址
                    <p></p>
                    <p></p>
                    <ul>
                        <li v-for="user in userLists" :key="user"
                            style="list-style:none;position: relative;margin-right: 40px;background-color: whitesmoke; border-radius: 15px;">
                            {{ showMyIdTag(user) }}UID: {{ user.UID }}
                        </li>
                    </ul>
                </div>
            </el-aside>
            <el-container>
                <el-header name="header">
                    <div style="font-size: large;position: relative;margin: 0 auto;margin-top: 20px">
                        {{ wsUrl }}
                        <span style="float:right;position: relative;margin: 0 auto;">
                            <el-button type="info" plain size="small" style="border: 1px solid rgb(225, 224, 224);" @click="optionDialogVisible = true">
                                <el-icon size="17px">
                                    <Setting />
                                </el-icon>
                            </el-button>
                            </span>
                            </div>
                            
                            <el-dialog v-model="optionDialogVisible" title="设置" width="30%" center >
                                <el-input v-model="inputDefaultWsUrl" placeholder="设置你的默认连接地址" />
                                <template #footer>
                                    <span class="dialog-footer">
                                        <el-button type="primary" @click="optionDialogVisible = false,setWsUrlCookie()">
                                            保存
                                        </el-button>
                                        <el-button @click="optionDialogVisible = false">关闭</el-button>
                                    </span>
                                </template>
                            </el-dialog>
                            
                            </el-header>
                            <el-main name="main">
                    <el-scrollbar height="325px" class="main-scrollbar" style="overflow:visible;">
                        <p v-for="messageItem in messageItems" :key="messageItem" style="position:relative">
                            <el-tag class="messageContent" type="info" size="large">
                                {{ messageItem.senderUID }} 用户在 {{ messageItem.sendTime }} 发送：<div
                                    style="font-size: larger;background-color: rgb(62, 234, 93);    box-shadow: 0 2px 7px rgba(190, 190, 190, 0.8);border-bottom: 1px solid rgb(225, 224, 224);">
                                    {{ messageItem.content }}
                                </div>
                            </el-tag>
                        </p>
                    </el-scrollbar>
                </el-main>
                <el-footer name="footer">
                    <el-input style="background: transparent;" v-model="sendMessageContent" id="enterMessageArea"
                        :autosize="{ minRows: 3, maxRows: 3 }" maxlength="100" placeholder="请输入内容" type="textarea"
                        resize="none" @keyup.enter.="sendMessageEvent()" />
                    <el-button type="success" class="sendMessageButton" @click="sendMessageEvent()">发送</el-button>
                </el-footer>
            </el-container>
        </el-container>
    </div>
    <span style="color:grey;font-size:small">
        <br><br><br><br>
        Author: FxShadowTG<br>
        ⭐https://github.com/FxShadowTG⭐
    </span>
</template>

<script>
import { nanoid } from "nanoid";
import { ElMessage } from 'element-plus'

export default {
    name: "webIndex",
    data() {
        return {
            //默认连接地址
            inputDefaultWsUrl: "ws://localhost:8888/ws",
            //选项对话框是否显示
            optionDialogVisible: false,
            //连接状态
            connStatus: false,
            //ws地址错误连接状态
            webSocketErrorStatus: false,
            //关闭连接按钮状态
            closeConnButtonStatus: true,
            //连接按钮状态
            connButtonStatus: false,
            //主机数量
            hostCount: 0,
            //输入框里的URL
            inputWsUrl: "",
            //实际WsURL
            wsUrl: "",
            //发送消息输入框
            sendMessageContent: "",
            //主webSock
            webSock: null,
            //用户UID
            user: {
                UID: null,
            },
            //主机列表
            userLists: [],
            //消息列表体
            messageItems: [],
            //发送消息模型
            sendMessageModel: {
                type: "",
                senderUID: "",
                sendTime: "",
                content: "",
            }
        };
    },
    mounted(){
        //校验默认地址是否存在
        let result = this.$cookies.isKey("default_ws_Url")   
        //如果存在就设置为默认地址
        if(result == true){
            this.inputWsUrl = this.getWsUrlCookie()
        }
        return
    },
    unmounted() {
        // js提供的clearInterval方法用来清除定时器
        clearInterval(this.FlushgetHostCountTimer);
        this.webSocketClose();
    },
    methods: {
        initSys() {
            //初始化网页配置
            //使用nanoid初始化用户UID
            this.user.UID = nanoid(5);
            //循环调用刷新主机列表
            this.startFlushgetHostCount();
        },
        startFlushgetHostCount() {
            //开始循环请求刷新主机列表
            this.FlushgetHostCountTimer = setInterval(() => {
                //如果连接状态为关闭
                if(this.connStatus == false){
                    //则不执行
                    return
                }
                //在这里先清除主机列表，不然有bug
                this.userLists.splice(0, this.userLists.length);
                this.getHostCount();
            }, 5000);
        },
        endFlushgetHostCount() {
            //结束循环请求刷新主机列表
            clearInterval(this.FlushgetHostCountTimer);
            this.FlushgetHostCountTimer = null;
        },
        // 发送websockwt请求
        initWebSocket() {
            // WebSocket与普通的请求所用协议有所不同，ws等同于http，wss等同于https
            let websocketUrl = this.wsUrl;
            this.webSock = new WebSocket(websocketUrl);
            this.webSock.onopen = this.webSocketOnOpen;
            this.webSock.onerror = this.webSocketOnError;
            this.webSock.onmessage = this.webSocketOnMessage;
            this.webSock.onclose = this.webSocketClose;
        },
        webSocketOnOpen() {
            console.log("WebSocket连接成功并发送了UID到ws");
            //  连接成功发送UID过去
            this.webSocketSend(this.user.UID);
        },
        webSocketOnMessage(data) {
            //接收数据
            //将接收的数据转为json格式
            var dataJSON = JSON.parse(data.data);
            //如果类型为message则推送至聊天窗口
            if (dataJSON.type === "message") {
                this.messageItems.push(dataJSON);
            }
            //另起，如果类型为hostCount则更新主机连接数量
            else if (dataJSON.type === "hostCount") {
                //删除主机列表里的所有元素
                //把新列表里的主机放进主机列表内
                for (let i = 0; i < dataJSON.count; i++) {
                    if (dataJSON.UIDList[i] != null) {
                        this.userLists.push({ "UID": dataJSON.UIDList[i] });
                    }
                }
                //因为type是HostCount，所以不打印出来了，直接return
                return;
            }
            console.log("WebSocket接受数据：");
            console.log(dataJSON);
        },
        reconnect() {
            // 断开重连操作
            let _this = this;
            //没连接上会一直重连，设置延迟避免请求过多
            setTimeout(function () {
                _this.initWebSocket();
            }, 2000);
        },
        webSocketClose(data) {
            this.webSock.close()
            console.log("断开连接", data);
        },
        webSocketOnError(data) {
            //开启错误状态
            this.webSocketErrorStatus = true;
            console.log("连接出错");
            console.log("报错信息如下：", data);
        },
        webSocketSend(data) {
            //发送数据
            this.webSock.send(data);
        },
        sendMessageEvent() {
            //发送消息事件（处理后再正式发送给服务端）

            //检测到未连接，则无法发送
            if(this.connStatus == false){
                ElMessage.error('请连接上服务端后再发送消息')
                return
            }
            //检测消息内容是否为空，是的话就直接return
            if(this.sendMessageContent == "" || this.sendMessageContent == null) {
                return;
            }

            //设定发送类型为Message
            this.sendMessageModel.type = "message";
            //将用户UID放入发送消息的模型里
            this.sendMessageModel.senderUID = this.user.UID;
            //将发送时间放入发送消息的模型里
            let localTime = new Date().toLocaleString();
            this.sendMessageModel.sendTime = localTime;
            //将消息放入发送消息的模型里
            this.sendMessageModel.content = this.sendMessageContent;
            console.log(this.sendMessageModel);
            //转换为String传给后端
            let sendMessageModelString = JSON.stringify(this.sendMessageModel);
            //正式发送数据到websocket
            this.webSocketSend(sendMessageModelString);
            //清除输入栏内容
            this.sendMessageContent = "";
        },
        getHostCount() {
            //获取用户数连接数量（直接使用ws，不用axios）
            //设定发送类型为hostCount
            this.sendMessageModel.type = "hostCount";
            //将用户UID放入发送消息的模型里
            this.sendMessageModel.senderUID = this.user.UID;
            //将发送时间放入发送消息的模型里
            let localTime = new Date().toLocaleString();
            this.sendMessageModel.sendTime = localTime;
            //清除消息用户体的content
            this.sendMessageModel.content = "";
            //转换为String传给后端
            let sendMessageModelString = JSON.stringify(this.sendMessageModel);
            //正式发送数据到websocket
            this.webSocketSend(sendMessageModelString);
        },
        handleRemove(file, fileList) {
            console.log(file, fileList);
        },
        handlePreview(file) {
            console.log(file);
        },
        handleExceed(files, fileList) {
            this.$message.warning(`当前限制选择 3 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
        },
        showMyIdTag(user) {
            if (user.UID === this.user.UID) {
                return "(我) "
            }
        },
        connWsUrl() {
            //校验是否为ws开头的地址
            var matchResult = this.inputWsUrl.substring(0,2)
            
            //如果不是则返回
            if(matchResult != "ws"){
                ElMessage({
                    message: 'URL地址有误，请检查后重新输入',
                    type: 'warning',
                })
                return
            }

            //将输入框里的url换成实际的url
            this.wsUrl = this.inputWsUrl

            //初始化ws
            this.initSys();
            this.initWebSocket();

            //先设置按钮为不可点击状态
            this.connButtonStatus = true

            //设置定时操作错误状态
            setTimeout(() => {
                //如果ws连接错误状态为开启
                if (this.webSocketErrorStatus == true) {
                    //如果开启
                    //错误提示
                    ElMessage.error('连接失败，请检查ws地址后重试')

                    //重置实际url地址
                    this.wsUrl == "";
                    //重置错误状态
                    this.webSocketErrorStatus = false;

                    //恢复按钮
                    this.connButtonStatus = false
                    return
                }

                //如果错误状态为关闭
                //连接成功提示
                ElMessage({
                    message: '连接成功',
                    type: 'success',
                })

                //连接成功则把连接按钮设置为不可点击状态
                this.connButtonStatus = true

                //连接成功则把断开按钮设置为点击状态
                this.closeConnButtonStatus = false

                //设置连接状态为开启
                this.connStatus = true

                //连接成功立即向服务器获取在线人数
                this.getHostCount();

            }, 1000)
        },
        closeConn(){
            this.webSocketClose();
            //断开后恢复连接连接按钮
            this.connButtonStatus = false;

            //连接成功则把断开按钮设置为不可点击状态
            this.closeConnButtonStatus = true;
            //关闭连接状态
            this.connStatus = false;
            //关闭后清空在线人数列表
            this.userLists.splice(0, this.userLists.length);
        },
        setWsUrlCookie(){
            
            ElMessage({
                message: '修改成功，刷新网页后生效',
                type: 'success',
            })

            this.$cookies.set("default_ws_Url",this.inputDefaultWsUrl,-1);
        },
        getWsUrlCookie(){
            let defaultWsUrl = this.$cookies.get("default_ws_Url")
            return defaultWsUrl
        },
    },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.web-index-container {
    background-color: rgb(255, 255, 255);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    top: 50px;
    width: 1000px;
    height: 550px;
    box-shadow: 0 2px 7px rgba(190, 190, 190, 0.8);
    margin: 0 auto;
}

.web-index-container aside {
    background-color: #e6e6e6;
    width: 250px;
    height: 550px;
    border-right: 1px solid rgb(225, 224, 224);
}

.web-index-container header {
    background-color: #f2f2f2;
    border-bottom: 1px solid rgb(225, 224, 224);
}

.web-index-container main {
    background-color: #f2f2f3;
}

.web-index-container footer {
    background-color: #f2f2f2;
    height: 125px;
    border-top: 1px solid rgb(225, 224, 224);
}

.sendMessageButton {
    position: relative;
    float: right;
    margin-top: 10px
}

.messageContent {
    background-color: rgb(250, 250, 250);
    color: #4a4a4a;
    font-size: medium;
}
</style>
