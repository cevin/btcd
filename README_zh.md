[English](https://github.com/cevin/btcd/blob/main/README.md) | 简体中文

# 比特币离线交易处理

简单的自托管比特币离线交易、支付工具

包含生成地址（含MultiSig、Bech32）、解析WIF、生成离线交易、签名离线交易（含Bech32地址的交易和MultiSig钱包交易）

# 使用说明

## 启动程序

`./btcd -addr localhost:8000`

## 地址

### 生成地址

#### 普通地址（旧地址）

`GET /address/new` `POST /address/new`

正确响应

```json
{
  "code": 200,
  "address": "18JS3qcfWmmGQiwuvWQhc19jU3ML71zAmj",
  "bech32_address": "bc1q2q2de9qwewgfzwve8scxdlgw5c57e456fdtf5z",
  "wif": "L2ihBrcpk2bwJnrbBd9kv6yCFymewj1oNpuDNKWYwaBFEHR9qCQZ",
  "private_key_hex": "a418caa27fe90edf145de85389bc0521a9af750ff9f64a2e5815820d10755009",
  "public_key_hex": "030f9849b7ba1a8e935663ab68e4fee6ac53f5ff853cf76fd8c4e921958b7bd2f9"
}
```

#### 多重签名地址

```text
GET /address/new-multi-sig?public_key_hexes=031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a,02623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d6,023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f&required=2
```

```text
POST /address/new-multi-sig

# json
{
    "required": 2,
    "public_key_hexes": "031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a,02623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d6,023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f"
}
```

正确响应

```json
{
  "code": 200,
  "address": "34P2BVSiDrHNHmmediJzPnp3rcztKPAytS",
  "script": "5221023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f2102623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d621031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a53ae",
  "asm": "2 023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f 02623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d6 031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a 3 OP_CHECKMULTISIG",
  "type": "multisig",
  "reqSigs": 2,
  "addresses": [
    "1NDcMyNeyZmLRTkAiAgDUvhy9GLU7S2npd",
    "1LNXg6moCRSDPvAwTEaRfSA2TxrjqJQk9J",
    "1PAfPcvnwMJ8Lok14u1UYGyPawxzNZmLrv"
  ]
}
```


#### 解析地址

##### 压缩后的公钥（hex格式编码）

`GET /address/parse?public_key_hex=023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f`

```text
POST /address/parse

# json
{"public_key_hex":"023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f"}

# form
public_key_hex=023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f
```
正确响应

```json
{
  "code": 200,
  "address": "1NDcMyNeyZmLRTkAiAgDUvhy9GLU7S2npd",
  "public_key_hex": "023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f"
}
```

##### 解析WIF私钥

`GET /address/parse?wif=L4G8o9MCkkQnL48yw1GnkQ6SzDZeDVn3GennEHHVqB4aqWjtBymJ`

```text
POST /address/parse

# json
{"wif":"L4G8o9MCkkQnL48yw1GnkQ6SzDZeDVn3GennEHHVqB4aqWjtBymJ"}

# form
wif=L4G8o9MCkkQnL48yw1GnkQ6SzDZeDVn3GennEHHVqB4aqWjtBymJ
```

正确响应

```json
{
  "code": 200,
  "address": "1NDcMyNeyZmLRTkAiAgDUvhy9GLU7S2npd",
  "bech32_address": "bc1qaz7jd0lg23ppfjkme8lhdesacykln2j06aww2u",
  "wif": "L4G8o9MCkkQnL48yw1GnkQ6SzDZeDVn3GennEHHVqB4aqWjtBymJ",
  "private_key_hex": "d21be41dc8e073e02d7ad36ca6f13c952752a28bc81b1c1140bd323e9fb7d933",
  "public_key_hex": "023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f"
}
```

##### 多重签名地址的兑付脚本（hex编码）

```text
GET /address/parse?script=5221023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f2102623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d621031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a53ae
```

```text
POST /address/parse

# json
{
    "script": "5221023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f2102623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d621031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a53ae"
}

# form
script=5221023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f2102623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d621031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a53ae
```

正确响应

```json
{
  "code": 200,
  "address": "34P2BVSiDrHNHmmediJzPnp3rcztKPAytS",
  "script": "5221023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f2102623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d621031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a53ae",
  "asm": "2 023b95f449bbe5b01e9f791a19266cc9f94d98e09488bde3afa4a75f823efd751f 02623db1f69a896b1f6bf54ab786773ce5f9385b58eb8eed58696fc7a4111d22d6 031af2a6177e1179cd119657112aab53725155ee305a4199a722c8b3097782d21a 3 OP_CHECKMULTISIG",
  "type": "multisig",
  "reqSigs": 2,
  "addresses": [
    "1NDcMyNeyZmLRTkAiAgDUvhy9GLU7S2npd",
    "1LNXg6moCRSDPvAwTEaRfSA2TxrjqJQk9J",
    "1PAfPcvnwMJ8Lok14u1UYGyPawxzNZmLrv"
  ]
}
```

## 交易

### 解析一个16进制编码的交易

`GET /transaction/decode?tx=0100000001fe75a438b72fdc302b80cc216d66d5e3bbb0359bce3bb4cecf743f5fda1f4eb101000000fdfd000048304502210096b617a5b2bd676ee8d3f8d8d91bf60c599e16382d1e12a61a1f9562c35b2cb102204379706a55c07bb45d20336159f80ebe9786938e34b9309e49ed422e6d2a44470147304402201550a8bb0c28107098289fe6fe64488bdee46800d28bfbb0b0a1e1b2d64b9fb4022004684015095b999185b3da1a23d239452ad73b199a032f71978760f8ae42313f014c6952210265e6f7fb614a369c9230912a3bb09c33c5c5be2e1bcfc2293ecaed46708e0b5c2103f546edf7b434b50aa0115c1c82a0f9a96505d9eff55d2fe3b848c4b51c06b6432102908375f301c7ea583f7e113939eab1164abda4ac27898b1cf78abf1c82f02da953aeffffffff01f8a70000000000001976a914bd63bf79e39f4cd52361c092c3fba9264662285688ac00000000`

```text
POST /transaction/decode

# json
{
    "tx":"0100000001fe75a438b72fdc302b80cc216d66d5e3bbb0359bce3bb4cecf743f5fda1f4eb101000000fdfd000048304502210096b617a5b2bd676ee8d3f8d8d91bf60c599e16382d1e12a61a1f9562c35b2cb102204379706a55c07bb45d20336159f80ebe9786938e34b9309e49ed422e6d2a44470147304402201550a8bb0c28107098289fe6fe64488bdee46800d28bfbb0b0a1e1b2d64b9fb4022004684015095b999185b3da1a23d239452ad73b199a032f71978760f8ae42313f014c6952210265e6f7fb614a369c9230912a3bb09c33c5c5be2e1bcfc2293ecaed46708e0b5c2103f546edf7b434b50aa0115c1c82a0f9a96505d9eff55d2fe3b848c4b51c06b6432102908375f301c7ea583f7e113939eab1164abda4ac27898b1cf78abf1c82f02da953aeffffffff01f8a70000000000001976a914bd63bf79e39f4cd52361c092c3fba9264662285688ac00000000"
}

# form
tx=0100000001fe75a438b72fdc302b80cc216d66d5e3bbb0359bce3bb4cecf743f5fda1f4eb101000000fdfd000048304502210096b617a5b2bd676ee8d3f8d8d91bf60c599e16382d1e12a61a1f9562c35b2cb102204379706a55c07bb45d20336159f80ebe9786938e34b9309e49ed422e6d2a44470147304402201550a8bb0c28107098289fe6fe64488bdee46800d28bfbb0b0a1e1b2d64b9fb4022004684015095b999185b3da1a23d239452ad73b199a032f71978760f8ae42313f014c6952210265e6f7fb614a369c9230912a3bb09c33c5c5be2e1bcfc2293ecaed46708e0b5c2103f546edf7b434b50aa0115c1c82a0f9a96505d9eff55d2fe3b848c4b51c06b6432102908375f301c7ea583f7e113939eab1164abda4ac27898b1cf78abf1c82f02da953aeffffffff01f8a70000000000001976a914bd63bf79e39f4cd52361c092c3fba9264662285688ac00000000
```

正确响应

```json
{
  "code": 200,
  "transaction": {
    "txid": "c7d1582d4cf85fbd10732002c5bb06068d4b86cfd5cca151ef88104c6702435a",
    "size": 340,
    "vsize": 340,
    "weight": 0,
    "version": 1,
    "locktime": 0,
    "vin": [
      {
        "txid": "b14e1fda5f3f74cfceb43bce9b35b0bbe3d5666d21cc802b30dc2fb738a475fe",
        "vout": 1,
        "scriptSig": {
          "asm": "0 304502210096b617a5b2bd676ee8d3f8d8d91bf60c599e16382d1e12a61a1f9562c35b2cb102204379706a55c07bb45d20336159f80ebe9786938e34b9309e49ed422e6d2a444701 304402201550a8bb0c28107098289fe6fe64488bdee46800d28bfbb0b0a1e1b2d64b9fb4022004684015095b999185b3da1a23d239452ad73b199a032f71978760f8ae42313f01 52210265e6f7fb614a369c9230912a3bb09c33c5c5be2e1bcfc2293ecaed46708e0b5c2103f546edf7b434b50aa0115c1c82a0f9a96505d9eff55d2fe3b848c4b51c06b6432102908375f301c7ea583f7e113939eab1164abda4ac27898b1cf78abf1c82f02da953ae",
          "hex": "0048304502210096b617a5b2bd676ee8d3f8d8d91bf60c599e16382d1e12a61a1f9562c35b2cb102204379706a55c07bb45d20336159f80ebe9786938e34b9309e49ed422e6d2a44470147304402201550a8bb0c28107098289fe6fe64488bdee46800d28bfbb0b0a1e1b2d64b9fb4022004684015095b999185b3da1a23d239452ad73b199a032f71978760f8ae42313f014c6952210265e6f7fb614a369c9230912a3bb09c33c5c5be2e1bcfc2293ecaed46708e0b5c2103f546edf7b434b50aa0115c1c82a0f9a96505d9eff55d2fe3b848c4b51c06b6432102908375f301c7ea583f7e113939eab1164abda4ac27898b1cf78abf1c82f02da953ae"
        },
        "sequence": 4294967295
      }
    ],
    "vout": [
      {
        "value": 0.00043,
        "n": 0,
        "scriptPubKey": {
          "asm": "OP_DUP OP_HASH160 bd63bf79e39f4cd52361c092c3fba92646622856 OP_EQUALVERIFY OP_CHECKSIG",
          "hex": "76a914bd63bf79e39f4cd52361c092c3fba9264662285688ac",
          "reqSigs": 1,
          "type": "pubkeyhash",
          "addresses": [
            "1JGQCsBmRqksmJArpEqVKSJyao3SZuZru3"
          ]
        }
      }
    ]
  }
}
```

### 生成一个离线交易

> 可以使用在线公开服务来获取指定地址的未使用Input, 如 https://blockchain.info/unspent?active=钱包地址

> 所有Input累计未使用金额 减去 将要发送的累计金额 等于 全网交易（确认）手续费
>
> 假设要发送部分金额，需要把自己的地址也作为发送目标，接收找零金额
> 
> 如：地址A有1BTC，想发送0.1给地址B，想要支付0.0000001BTC作为交易手续费，则pay_to_addresses可能为:
> 
> [{"address":地址B, "amount":0.1}, {"address":地址A, "amount":0.8999999]

```text
POST /transaction/create

# json
{
    "txin": [
        {
            "txid": "未使用的交易ID",
            "vout": 整数, # input的vout值
        }
    ],
    "pay_to_addresses": [
        {
            "address": "任意正确的比特币地址",
            "amount": float，最大支持精确到小数点后8位
        }
    ]
}
```

正确响应

```json
{
  "code": 200,
  "raw": "未签名的hex编码的交易"
}
```

### 对一个交易进行签名

来自 MultiSig 多重签名地址的交易，可能需要使用兑付脚本（redeem script）和对应的私钥多次签名

所有的Input中，有任意一个Input是SegWit交易的，所有Input都必须包含Input对应的收入金额

```text
POST /transaction/sign

# json
{
    "raw": "未签名的hex编码的交易",
    "txin": [
        {
            "txid": "未使用的交易ID",
            "wif": "WIF格式的私钥",
            "redeem-script": "多重签名地址的兑付脚本（redeem script）16进制编码",
            "segwit": bool 是否来自bech32格式的地址的交易
            "amount": float64 (可选 segwit 时必须)
        }
    ]
}
```

正确响应

```json
{
    "code": 200,
    "raw": "签名后的交易16进制编码"
}
```

### 创建并签名一个交易

`POST /transaction/create-and-sign`

```text
{
    "txin": [
        {
            "txid": "未使用的交易ID",
            "vout": int, # input vout
            "wif": "WIF格式的私钥",
            "redeem-script": "多重签名地址的兑付脚本（redeem script）16进制编码",
            "segwit": bool 是否来自bech32格式的地址的交易
        }
    ],
    "pay_to_addresses": [
        {
            "address": "任意正确的比特币地址",
            "amount": float，最大支持精确到小数点后8位
        }
    ]
}
```

正确响应

```json
{
  "code": 200,
  "raw": "签名后的交易16进制编码" 
}
```


## 广播交易

签名后的交易需要广播到比特币网络中，被其他节点确认后才能生效

可以使用在线服务或您自己的钱包工具进行广播，在线服务如下：

https://coinb.in/#broadcast

https://explorer.btc.com/tools/tx/publish


---

捐助

比特币地址: 1DMHiyzcjNzuYhWCbB4tx3wKfcvid1qgC4

![wxp___f2f06nUUOuf340I1qnfmHgzyZJjKBnfrMyncNWZXxKBtkD8.png](https://s2.loli.net/2023/10/01/9uEwOAx3XGa2vjk.png)