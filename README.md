# 叮咚买菜助手

## 抢菜

### 通过抓包获取叮咚小程序的网络请求, 模仿接口调用, 配置文件中 header 和 form 里的必要参数需要根据自己抓包到的数据进行修改

### 配置文件 config.json 设置

```json
  "snap_up": 1, // 0: 关闭, 1: 6点抢, 2: 8点半抢, 3: 6点和8点半都抢
  "pick_up_needed": false, // 闲时捡漏开关 false关闭 true打开
  "monitor_needed": false, // 监视器开关 监视是否有可配送时段
  "notify_needed": true, // 通知开关 发现有可配送时段时通知大家有可购商品
  "notify_interval": 5, // 通知间隔 单位: 分钟
```

## 可派送时段监听

<img src="/assets/effect.jpeg" width="300" alt="effect" />

### 当有可派送时段时, 发送通知到手机(只支持 ios)

### 注意事项

#### 1.安装 bark 得到自己的 barkID, 并将其写入配置文件中

<img src="/assets/user.jpeg" width="300" alt="user" />

#### 2.配置文件中 users 为一组需要通知的 bark userID

#### 3.bark 需打开允许通知

<img src="/assets/notify.jpeg" width="300" alt="notify" />

### 版权说明

**本项目为 GPL3.0 协议，请所有进行二次开发的开发者遵守 GPL3.0 协议，并且不得将代码用于商用。**
