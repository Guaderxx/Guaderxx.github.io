---
title: miband6-heart-rate
date: 2024-10-06 13:04:08
tags:
- Rust
categories:
- Daily
keywords:
- Mi Band 6
- Rust
- BLE
copyright: Guader
copyright_author_href:
copyright_info:
---

> 前一阵子看 vtuber 感觉旁边显式的心跳插件挺好玩的。  
> 不过当时没找到手环  
> 最近找到了，就看了下   
>
> 设备：小米手环6 (MI Smart Band 6)


关于这个心跳数据，看到了两种方案获取，一种是抓应用的包，一种是抓蓝牙的包。

这里我用的是用蓝牙获取数据。


## 获取小米手环 AuthKey

> 对于支持广播的设备，可以直接跳到下面  

1. 向小米手环请求随机数
2. 接收到随机数后，使用该手环的 Auth Key 对随机数进行 AES 对称加密。
3. 将加密后的信息发回给手环。
4. 验证通过。

### Android

首先通过 [freemyband][freemyband] 页面上的指引，下载魔改后的小米运动APP，打开与手环配对，然后在 `/sdcard/freemyband` 路径下获取手环的 AuthKey.


### Rooted IPhone

配对好后用 ssh 连接上，在 `/var/mobile/Containers/Data/Application/<MiFit_App_UUID>/Documents` 下面找到 `HMDBDeviceInfoDataBaseV2.sqlite` 这样的数据库，然后从 `device_info` 表中 `deviceOauthKey` 获取手环的 AuthKey.


### IPhone

根据 [huamitoken][huamitoken] 这个项目获取 AuthKey.


## Note

**这个方法需要手环连接电脑而不是手机**

感兴趣的可以对照 [e99p1ant][e99p1ant] 这篇文章继续。  

我使用的是下面的方法。


# 通过广播获取

看到一半的时候我发现这个手环设置里有 **蓝牙广播** 和 **运动心率广播** 。  
我感觉通过让它广播我再读数据会更简单点。  
我这是 windows, 所以下载的是微软的 [BTP][btp]。  

*Mac可以下载 bluetility*  
`brew install --cask bluetility`  

BTP 会自动解压到 `C:\BTP\v1.14.0\x86\btvs.exe` ，直接双击就行，然后它会打开 Wireshark 进行转包。  
然后弹窗的 `Full Packet Logging` 按钮也要记得点击。

抓包时需要一直打开电脑蓝牙的扫描，也就是 *添加设备*。  
且需要打开运动模式，比如 *自由训练* 。

会有一些其他的数据也被获取到，可以通过 `btcommon.eir_ad.entry.company_id == 0x0157` 该过滤器过滤掉不相干的设备。

通过 Wireshark 的 `Bluetooth HCI Event - LE Meta -> Advertising Data -> Manufacturer Specific -> Company ID` 获取的过滤器。

Data 部分是 24 字节，比较多个数据包可以发现，只有第四字节是变动的，因此可以推测，第四字节就是心跳。

接下来就是写程序验证了。

```rust
use std::error::Error;

use bluest::{Adapter,AdvertisingDevice};
use futures_lite::StreamExt;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let adapter = Adapter::default()
        .await
        .ok_or("Bluetooth adapter not found")?;
    adapter.wait_available().await?;
    println!("Starting scan");

    let mut scan = adapter.scan(&[]).await?;
    println!("Scan started");

    while let Some(discovered_device) = scan.next().await {
        if let Some(manufacturer_data) = discovered_device.adv_data.manufacturer_data {
            if manufacturer_data.company_id != 0x0157 {
                continue;
            }
            let name = discovered_device
                .device
                .name()
                .unwrap_or("unknown".to_string());
            let rssi = discovered_device.rssi.unwrap_or_default();
            let heart_rate = manufacturer_data.data[3];
            println!("{name} ({rssi}dBm) Heart Rate: {heart_rate:?}");
        }
    }
    Ok(())
}
```


# H5 Bluetooth

顺带看了眼，H5 这还属于实验性吧。能看到设备，但是没连上。

```js
// https://webbluetoothcg.github.io/web-bluetooth/#introduction-examples
// EXAMPLE 1
navigator.bluetooth.requestDevice({
  filters: [{
    services: ['heart_rate'],
  }]
}).then(device => device.gatt.connect())
.then(server => server.getPrimaryService('heart_rate'))
.then(service => {
  chosenHeartRateService = service;
  return Promise.all([
    service.getCharacteristic('body_sensor_location')
      .then(handleBodySensorLocationCharacteristic),
    service.getCharacteristic('heart_rate_measurement')
      .then(handleHeartRateMeasurementCharacteristic),
  ]);
});
```

这段代码里 `services: ['heart_rate']` 换成 `manufacturerData: [{ companyIdentifier: 0x0157 }]` 倒是能找到手环。

然后在 `device.gatt.connect` 这部分，这里和 `https://webbluetoothcg.github.io/web-bluetooth/#dom-bluetoothremotegattserver-connect` 这里的内容也有出入。

总而言之，过阵子再看吧。





## Related Project

[wuhanbeat][wuhanbeat]： [e99p1ant][e99p1ant] 这篇文章对应的项目, 在 MacOS 上监听心跳，支持 2,3,4,5,6   Golang
[miband4][miband4]： 一个 4 的连接器  Nodejs
[heartrate][heartrate]：在 windows 上监听心跳，支持 2,3,4  C#



[freemyband]: http://www.freemyband.com/
[huamitoken]: https://github.com/argrento/huami-token
[wuhanbeat]: https://github.com/wuhan005/mebeats
[miband4]: https://github.com/satcar77/miband4
[heartrate]: https://github.com/Eryux/miband-heartrate
[btp]: https://learn.microsoft.com/zh-cn/windows-hardware/drivers/bluetooth/testing-btp-setup-package
[e99p1ant]: https://github.red/miband-heart-rate/
