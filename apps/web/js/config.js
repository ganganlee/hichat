//用户服务api
const API = "http://127.0.0.1:8080";

//websocketapi
let WEBSOCKETAPI = 'ws://127.0.0.1:8081';

//发送消息api
let MSGAPI = 'http://192.168.3.2:8080';

//文件上传api
let UPLOADAPI = 'http://192.168.3.2:8089';

//emoji表情列表
const EMOJIS = [9,11,12,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,66,67,68,69,70,71,72,73,74,75,76,77,78,80,81,82,83,84,85,86,87,88,89,90,91];

//文件类型
const TYPES = {
    "img":2,
    "mp3":3,
    "mp4":4,
    "file":5
};

//聊天类型
const ChatType = {
    "chat":'privateMsg',//私聊
    "group":'groupMsg',//群聊
};

//登录用户token
let TOKEN = '';