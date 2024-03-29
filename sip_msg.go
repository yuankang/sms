package main

//海康 大华 同网 跨网 sip消息的不同
//https://zhuanlan.zhihu.com/p/98533891
//各种sip消息 比较全
//https://blog.csdn.net/weixin_43360707/article/details/120975297

//首先IPC会发送不带账户信息的注册请求
var Sip1Rqst = `REGISTER sip:34020000002000000001@3402000000 SIP/2.0
Via: SIP/2.0/UDP 192.168.10.8:60719;rport=60719;branch=z9hG4bK1377739937
Max-Forwards: 70
Contact: <sip:34020000002000000719@192.168.10.8:60719>
To: <sip:34020000002000000719@3402000000>
From: <sip:34020000002000000719@3402000000>;tag=1534571172
Call-ID: 1566767517
CSeq: 1 REGISTER
Expires: 3600
User-Agent: IP Camera
Content-Length: 0
`

var Sip1Rsps = `SIP/2.0 401 Unauthorized
Via: SIP/2.0/UDP 192.168.10.8:60719;rport=60719;branch=z9hG4bK1377739937
To: <sip:34020000002000000719@3402000000>;tag=a9ced173
From: <sip:34020000002000000719@3402000000>;tag=1534571172
Call-ID: 1566767517
CSeq: 1 REGISTER
User-Agent: SYSZUX28181
WWW-Authenticate: Digest nonce="1577261948:e7c287426affd7479e0b017b20aa6690",algorithm=MD5,realm="3402000000",qop="auth,auth-int"
Content-Length: 0
`

var Sip2Rqst = `REGISTER sip:34020000002000000001@3402000000 SIP/2.0
Via: SIP/2.0/UDP 192.168.10.8:60719;rport=60719;branch=z9hG4bK1307786573
Max-Forwards: 70
Contact: <sip:34020000002000000719@192.168.10.8:60719>
To: <sip:34020000002000000719@3402000000>
From: <sip:34020000002000000719@3402000000>;tag=1534571172
Call-ID: 1566767517
CSeq: 2 REGISTER
Expires: 3600
User-Agent: IP Camera
Authorization: Digest username="34020000002000000719", realm="3402000000", nonce="1577261948:e7c287426affd7479e0b017b20aa6690", uri="sip:34020000002000000001@3402000000", response="d588f21ffa1efabdd5baed3950008a41", algorithm=MD5, cnonce="0a4f113b", qop=auth, nc=00000001
Content-Length: 0
`

var Sip2Rsps = `SIP/2.0 200 OK
Via: SIP/2.0/UDP 44.198.62.2:5061;branch=z9hG4bK-353632-bcee7db5e65a86c81bd997d93ef56e80
From: <sip:44190012002000000001@4419001200>;tag=fcfaada97a764f96886612835fe5341a-1635241188369
To: <sip:44190012002000000001@4419001200>;tag=810172799
CSeq: 2 REGISTER
Call-ID: 51c333de7d9c95813c08ca888a968027@44.198.62.2
User-Agent: LiveGBS v211022
Contact: <sip:44190012002000000001@44.198.62.2:5061>
Content-Length: 0
Date: 2021-10-26T17:49:05.279
Expires: 3600
`

//心跳消息, ipc发给svr
var SipKeepRqst = `MESSAGE sip:11000000122000000034@1100000012 SIP/2.0
Via: SIP/2.0/TCP 10.3.220.151:42341;rport;branch=z9hG4bK581882449
From: <sip:11010000121310000034@1100000012>;tag=57873780
To: <sip:11000000122000000034@1100000012>
Call-ID: 934168825
CSeq: 20 MESSAGE
Content-Type: Application/MANSCDP+xml
Max-Forwards: 70
User-Agent: IP Camera
Content-Length:   177

<?xml version="1.0" encoding="GB2312"?>
<Notify>
<CmdType>Keepalive</CmdType>
<SN>44</SN>
<DeviceID>11010000121310000034</DeviceID>
<Status>OK</Status>
<Info>
</Info>
</Notify>`

var SipKeepRsps = `SIP/2.0 200 OK
Via: SIP/2.0/UDP 44.198.62.2:5061;branch=z9hG4bK-353632-82aa13ec871fe910f87da0aad9404808;X-HwDim=4
From: <sip:44190012002000000001@4419001200>;tag=5ec9a354a2294398b6ebc451b4862cbf-1635241248370
To: <sip:34020000002000000001@3402000000>;tag=328213176
CSeq: 1 MESSAGE
Call-ID: 02638cc1b4fa0e3863f2b9cfb3204eda@44.198.62.2
User-Agent: LiveGBS v211022
Content-Length: 0
`

//播放请求消息 这里有SDP信息, svr发给ipc
var SipPlayRqst = `INVITE sip:34020059001320100485@4419001200 SIP/2.0
Via: SIP/2.0/UDP 44.198.63.22:15060;rport;branch=z9hG4bK434173967
From: <sip:34020000002000000001@3402000000>;tag=699173967
To: <sip:34020059001320100485@4419001200>
Call-ID: 367173963
CSeq: 3 INVITE
Content-Type: APPLICATION/SDP
Contact: <sip:34020000002000000001@44.198.63.22:15060>
Max-Forwards: 70
User-Agent: LiveGBS v211022
Subject: 34020059001320100485:0200590485,34020000002000000001:0
Content-Length: 258

v=0
o=34020059001320100485 0 0 IN IP4 44.198.63.22
s=Play
c=IN IP4 44.198.63.22
t=0 0
m=video 30002 TCP/RTP/AVP 96 97 98
a=recvonly
a=rtpmap:96 PS/90000
a=rtpmap:97 MPEG4/90000
a=rtpmap:98 H264/90000
a=setup:active
a=connection:new
y=0200590485`

//有时候ipc会先发一个tring，然后再发送200ok
var SipTryingRsps = `SIP/2.0 100 Trying
Via: SIP/2.0/UDP 44.198.63.22:15060;branch=z9hG4bK434173967;rport=15060
Call-ID: 367173963
From: <sip:34020000002000000001@a3402000000>;tag=699173967
To: <sip:34020059001320100485@a4419001200>
CSeq: 3 INVITE
Content-Length: 0
`

//这里有SDP信息
var SipPlayRsps = `SIP/2.0 200 OK
Via: SIP/2.0/UDP 44.198.63.22:15060;branch=z9hG4bK434173967
Call-ID: 367173963
From: <sip:34020000002000000001@3402000000>;tag=699173967
To: <sip:34020059001320100485@4419001200>;tag=c3db5071e0794fb9b11c6af5c31d1665-1635241209281
CSeq: 3 INVITE
Contact: <sip:44.198.62.2:5061>
Content-Length: 269
Content-Type: application/sdp

v=0
o=34020059001320100485 0 0 IN IP4 44.198.62.25
s=Play
c=IN IP4 44.198.62.25
t=0 0
m=video 5489 TCP/RTP/AVP 96 108
a=sendonly
a=setup:passive
a=connection:new
a=rtpmap:96 PS/90000
a=rtpmap:108 H265/90000
y=0200590485
f=v/5a///
b=RS:0
b=RR:0`

//ACK消息, ???
var SipAckRqst = `ACK sip:34020059001320100485@4419001200 SIP/2.0
Via: SIP/2.0/UDP 44.198.63.22:15060;rport;branch=z9hG4bK455174084
From: <sip:34020000002000000001@3402000000>;tag=699173967
To: <sip:34020059001320100485@4419001200>;tag=c3db5071e0794fb9b11c6af5c31d1665-1635241209281
Call-ID: 367173963
CSeq: 3 ACK
Contact: <sip:34020000002000000001@44.198.63.22:15060>
Max-Forwards: 70
User-Agent: LiveGBS v211022
Content-Length: 0
`

//BYE消息
var SipByeRqst = `BYE sip:34020059001320100485@4419001200 SIP/2.0
Via: SIP/2.0/UDP 44.198.63.22:15060;rport;branch=z9hG4bK115255475
From: <sip:34020000002000000001@3402000000>;tag=896204437
To: <sip:34020059001320100485@4419001200>;tag=aa5a0c6d9964497586f4657b6bdd1157-1635241239699
Call-ID: 752204436
CSeq: 29 BYE
Contact: <sip:34020000002000000001@44.198.63.22:15060>
Max-Forwards: 70
User-Agent: LiveGBS v211022
Content-Length: 0
`

var SipByeRsps = `SIP/2.0 200 OK
Via: SIP/2.0/UDP 44.198.63.22:15060;branch=z9hG4bK115255475
Call-ID: 752204436
From: <sip:34020000002000000001@3402000000>;tag=896204437
To: <sip:34020059001320100485@4419001200>;tag=aa5a0c6d9964497586f4657b6bdd1157-1635241239699
CSeq: 29 BYE
Content-Length: 0
`

//设备目录查询消息
var SipClRqst = `MESSAGE sip:44190012002000000001@4419001200 SIP/2.0
Via: SIP/2.0/UDP 44.198.63.22:15060;rport;branch=z9hG4bK817173798
From: <sip:34020000002000000001@3402000000>;tag=112173799
To: <sip:44190012002000000001@4419001200>
Call-ID: 153173799
CSeq: 1 MESSAGE
Content-Type: Application/MANSCDP+xml
Max-Forwards: 70
User-Agent: LiveGBS v211022
Content-Length: 157

<?xml version="1.0" encoding="GB2312"?>
<Query>
  <CmdType>Catalog</CmdType>
  <SN>553173798</SN>
  <DeviceID>44190012002000000001</DeviceID>
</Query>
`

//摄像头先回复200ok, 再回复目录信息
var SipClRsps = `SIP/2.0 200 OK
Via: SIP/2.0/UDP 44.198.63.22:15060;branch=z9hG4bK817173798
Call-ID: 153173799
From: <sip:34020000002000000001@3402000000>;tag=112173799
To: <sip:44190012002000000001@4419001200>;tag=25fda1ea1c5843768a53d84c546a9e38-1635241209019
CSeq: 1 MESSAGE
Content-Length: 0
`

//摄像头回复的目录信息
var SipClInfoRqst = `MESSAGE sip:34020000002000000001@44.198.63.22:15060;transport=udp SIP/2.0
Via: SIP/2.0/UDP 44.198.62.2:5061;branch=z9hG4bK-353632-9942f723e9768c20eef038a9a21b67a6;X-HwDim=4
Call-ID: 9529032ee389aa82ddd1cd630c274407@44.198.62.2
From: <sip:44190012002000000001@4419001200>;tag=2e77af9f7456496ba96392d054e338a1-1635241209026
To: <sip:34020000002000000001@3402000000>
CSeq: 1 MESSAGE
Max-Forwards: 70
Content-Length: 301
Content-Type: application/MANSCDP+xml

<?xml version="1.0" encoding="UTF-8"?>
<Response>
<CmdType>Catalog</CmdType>
<SN>553173798</SN>
<DeviceID>44190012002000000001</DeviceID>
<Result>OK</Result>
<SumNum>7</SumNum>   //这里是7个通道
<DeviceList Num="1">
<Item>
<DeviceID>340200</DeviceID>
<Name>......</Name>
</Item>
</DeviceList>
</Response>
`

//服务端对目录信息的回复
var SipClInfoRsps = `SIP/2.0 200 OK
Via: SIP/2.0/UDP 44.198.62.2:5061;branch=z9hG4bK-353632-9942f723e9768c20eef038a9a21b67a6;X-HwDim=4
From: <sip:44190012002000000001@4419001200>;tag=2e77af9f7456496ba96392d054e338a1-1635241209026
To: <sip:34020000002000000001@3402000000>;tag=138173826
CSeq: 1 MESSAGE
Call-ID: 9529032ee389aa82ddd1cd630c274407@44.198.62.2
User-Agent: LiveGBS v211022
Content-Length: 0
`
