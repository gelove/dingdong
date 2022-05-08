# 使用教程

> 由于叮咚APP不允许走代理, 我们只能通过手机端全局代理APP来抓包

## 安装手机端代理工具—***Potatso Lite***

申请一个美区Apple ID并登录, 在AppStore中免费下载

[美国Apple ID申请注册教程](https://www.bilibili.com/read/cv5569420)

[(小红书)如何注册美区APPLE ID](https://www.xiaohongshu.com/discovery/item/5ddbf92600000000010017e6)

## ***安装Charles***

在电脑上[***安装Charles***](https://www.charlesproxy.com/download/)

### 注册Charles

![](https://upload-images.jianshu.io/upload_images/28036656-9babfc8ad6c10534.png?imageMogr2/auto-orient/strip|imageView2/2/w/173/format/webp)

> Registered Name：**ddmc**
>
> License Key：**8DEA943587A9B04870**

### 在电脑上安装证书

![](https://upload-images.jianshu.io/upload_images/28036656-626be007ae690bec.png?imageMogr2/auto-orient/strip|imageView2/2/w/967/format/webp)

![](https://upload-images.jianshu.io/upload_images/28036656-55385810f76d3088.png?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp)

### 配置Charles

> 查看电脑IP

![](https://i0.hdslb.com/bfs/album/ed1606f433ee1c8a2adc156c5c829d66a13b91ce.png)

> 配置代理
![](https://upload-images.jianshu.io/upload_images/28036656-b3c30094b80cd213.png?imageMogr2/auto-orient/strip|imageView2/2/w/483/format/webp)
![](https://upload-images.jianshu.io/upload_images/28036656-fea2f7095a491a1b.png?imageMogr2/auto-orient/strip|imageView2/2/w/774/format/webp)

> 添加 ssl 代理
>
> 点击菜单栏 <kbd>Proxy</kbd> -> <kbd>SSL Proxying Settings</kbd> -> <kbd>Add</kbd>
> 
> Host 与 Port 都填写 * , 点击 <kbd>OK</kbd>
![](https://i0.hdslb.com/bfs/album/c109f8543a71dd180021edc7a99fc9d33e7163a4.png)

> 安装手机端证书
![](https://upload-images.jianshu.io/upload_images/28036656-81ca3578c2d5cb5f.png?imageMogr2/auto-orient/strip|imageView2/2/w/969/format/webp)

### 手机添加证书

> 手机打开浏览器, 访问 chls.pro/ssl

点击 <kbd>设置</kbd>–><kbd>通用</kbd>–><kbd>VPN与设备管理</kbd>–><kbd>配置描述文件</kbd>–><kbd>安装</kbd>

![](https://upload-images.jianshu.io/upload_images/28036656-15a27351739215ee.png?imageMogr2/auto-orient/strip|imageView2/2/w/679/format/webp)

点击 <kbd>设置</kbd>–><kbd>通用</kbd>–><kbd>关于本机</kbd>–><kbd>证书信任设置</kbd>–><kbd>选中下载好的证书</kbd>

### 添加手机代理

1. 安装完打开 Potatso Lite
2. 点击右上角 <kbd>+</kbd>, 点击 <kbd>添加</kbd>
3. 类型选中 Socks5, 并输入服务器地址(你电脑的IP地址)和端口(8999)
4. 点击 <kbd>完成</kbd>
5. 回到主页, 点击右下角 <kbd>▶️</kbd> (注意: 此时手机只能通过Charles代理访问互联网, 不用时及时关闭 Potatso)


## **抓包**

打开叮咚APP点一下个人刷新一下, 再点一下购物车

### 抓包内容

### 获取 Session

**如果无法找到所列出的请求，请参见后文 iOS 设备 Charles 抓包帮助**

1. 在iOS设备上启动叮咚买菜APP
2. 完成登录
3. 启动Charles并完成抓包配置（需要配置SSL抓包）
4. 点击“购物车”并刷新
5. 在请求中找到 https://maicai.api.ddxq.mobi/cart/index
6. 右击该请求，选择Export Session，保存到项目 session 文件夹下，文件类型请选择JSON Session File (.chlsj)

![](https://upload-images.jianshu.io/upload_images/28036656-3c7984d1c105bd3d.png?imageMogr2/auto-orient/strip|imageView2/2/w/342/format/webp)

右击上图红圈点击导出session 为 json session file 格式

### 获取 im_secret

1. 点击“我的”并刷新
2. 在请求中找到https://sunquan.api.ddxq.mobi/api/v1/user/detail
3. 左击该请求，选择Contents选项卡，在下半部分选项卡中选择JSON Text视图
4. 找到 user_info 下的 im_secret 字段，复制其值到配置文件中

![](https://upload-images.jianshu.io/upload_images/28036656-ea59a51813e1cb41.png?imageMogr2/auto-orient/strip|imageView2/2/w/977/format/webp)

在上图红圈处单击 JSON Text, 找到 user_info 下的 im_secret

## **运行**

只需下载你的系统对应的 dingdong 可执行程序, 无需下载任何环境

### MacOS

1. 进入到你下载的程序目录
2. [在文件夹中打开终端](https://zhuanlan.zhihu.com/p/162748665)
3. 在打开的终端里输入 ***./dingdong***
4. 信任此程序 <kbd>系统设置</kbd>–><kbd>安全与隐私</kbd>–><kbd>通用</kbd>–><kbd>总是打开此程序</kbd>

### Windows

1. 进入到你下载的程序目录
2. 在文件管理器地址栏输入 <kbd>cmd</kbd> 然后回车
3. 在打开的命令行程序里输入 ***dingdong.exe***

![](https://upload-images.jianshu.io/upload_images/28036656-8b79370ea2f6b5d6.png?imageMogr2/auto-orient/strip|imageView2/2/w/404/format/webp)
