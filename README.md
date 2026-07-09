<p align="center">
  <h1 align="center">proxy-panel-3xui</h1>
  <p align="center">一个<strong>完全免费、无授权码、持续开源</strong>的 Xray 代理管理面板</p>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/license-GPLv3-blue.svg" alt="License: GPLv3">
  <img src="https://img.shields.io/badge/based%20on-3x--ui%20v3.4.2-brightgreen.svg" alt="Based on 3x-ui">
  <img src="https://img.shields.io/badge/xray--core-v26.6.27-orange.svg" alt="xray-core">
</p>

---

## 📌 这是什么

基于 [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui)(4 万+ star 的主流开源面板)二次开发,把 **X-Panel 付费锁定的功能免费开源实现**。

| | X-Panel | 本项目 |
|---|---|---|
| 一键配置(批量生成 REALITY 节点) | 💰 付费(¥100/授权码) | ✅ **免费** |
| 授权码 / 联网验证 | 需要 | ❌ 不需要 |
| 源码 | 付费版闭源 | ✅ 全开源 GPLv3 |

> 底座 3x-ui 本身已支持全部主流协议;本项目在其上补齐 X-Panel 的**便捷向导**,并保证一切免费无授权码。

---

## ✨ 功能特性

**协议**(继承自 3x-ui,全都有):VLESS · VMess · Trojan · Shadowsocks · WireGuard · Hysteria2 · SOCKS · HTTP · Dokodemo-door · TUN · MTProto

**传输 & 安全**:TCP · mKCP · WebSocket · gRPC · HTTPUpgrade · XHTTP + TLS / XTLS / **REALITY**

**面板能力**:
- ⚡ **一键配置**(本项目新增)—— 一键批量生成 N 个节点,可选 **VLESS + REALITY + Vision** 或 **Shadowsocks-2022**(两者都免域名/证书),每个独立密钥
- 每客户端流量配额 / 到期 / IP 限制 · 在线状态 · 二维码 · 分享链接 · 订阅
- 多节点集群 · 独立订阅服务器 · Telegram bot · REST API + Swagger
- 13 种界面语言 · 明暗主题 · Fail2ban · SQLite / PostgreSQL
- 资源占用极低:面板 + xray 引擎 ≈ **116 MB 内存**

---

## 📖 使用流程(和 X-Panel 一样)

> 完整链路:**打开面板 → 加节点 → 导出链接 → 导入客户端**

1. **打开面板**
   浏览器访问 `http://<你的IP>:2053`,登录(默认 `admin` / `admin`,**首次登录请立即改密码**)。

2. **加节点** — 进「入站」页,两种方式:
   - ⚡ 点 **一键配置** → 填数量 → 立即生成一批 REALITY 节点(最省事)
   - 或点 **添加入站** → 手动选任意协议(VLESS / VMess / Trojan / …)

3. **导出链接 / 订阅**
   每个节点右侧有「复制链接 / 二维码 / 订阅」;也可一键「导出全部链接」。
   生成的是标准分享链接,例如:
   ```
   vless://<uuid>@<你的IP>:20000?flow=xtls-rprx-vision&security=reality&pbk=<公钥>&sni=www.microsoft.com&type=tcp#reality-01
   ```

4. **导入客户端(V2rayN / v2rayNG / Clash 等)**
   - 复制链接 → 在 **V2rayN** 主界面按 `Ctrl+V`(或菜单「服务器 → 从剪贴板导入批量URL」)
   - 或用**订阅**:客户端里贴一个订阅地址,节点变了自动更新

> 💡 REALITY 节点**不需要你自己的域名和 SSL 证书**——它借用大网站(如 www.microsoft.com)的真证书来伪装,裸 IP 就能跑。

---

## 🚀 安装

> 目前提供**源码构建**方式;VPS 一键安装脚本开发中(见[路线图](#-路线图))。

### A. 从源码构建(Windows,开发 / 自用)

需要:**Go 1.26+**、**Node 22+**、**C 编译器**(sqlite 驱动是 cgo,Windows 装 [MinGW-w64](https://winlibs.com/))

```powershell
git clone https://github.com/Human1-12/proxy-panel-3xui.git
cd proxy-panel-3xui
copy .env.example .env      # 本地开发配置:数据库写到 ./x-ui,端口 2053
.\build.ps1                 # 构建前端 + 后端 → 生成 xui.exe
.\run-dev.ps1               # 启动;浏览器开 http://127.0.0.1:2053  (admin/admin)
```

首次运行还需把 **Xray 引擎二进制**放进 `x-ui\` 目录:
`xray-windows-amd64.exe` + `geoip.dat` + `geosite.dat`(从 [Xray-core releases](https://github.com/XTLS/Xray-core/releases) 下载**对应版本 v26.6.27**)。

### B. 部署到 Linux VPS

你只需要一台 VPS(公网 IPv4 + root)。交叉编译出 Linux 版二进制拷到服务器运行即可;**一键安装脚本在路线图中**。REALITY 节点不需要域名,裸 IP 可用。

---

## 🧩 一键配置 API

```http
POST /panel/api/inbounds/oneclick/reality
Authorization: Bearer <API Token>
Content-Type: application/json

{ "count": 10, "portStart": 20000, "protocol": "reality", "remarkPrefix": "", "dest": "www.microsoft.com:443" }
```

`protocol` 可选 `reality`(默认,VLESS + TCP + REALITY + Vision)或 `ss2022`(Shadowsocks 2022-blake3-aes-256-gcm);选 `ss2022` 时忽略 `dest`。一次生成 `count` 个入站,各自独立的密钥 / UUID / subId,从已用端口之上**确定性分配**并自动跳过占用端口。`remarkPrefix` 留空则按协议命名(reality / ss)。面板 UI 上的「⚡ 一键配置」按钮即调用此接口。

---

## 🛠️ 技术栈

**后端** Go 1.26(Gin + GORM)· **前端** React 19 + Ant Design 6 + Vite 8 · **核心** xray-core v26.6.27 · **存储** SQLite / PostgreSQL

主要新增代码:
- `internal/web/service/oneclick.go` — 批量生成引擎
- `internal/web/controller/oneclick.go` — REST 接口
- `frontend/src/pages/inbounds/InboundsPage.tsx` — 一键配置按钮 + 表单

---

## 🔄 跟随上游 3x-ui 更新

```bash
git fetch upstream            # upstream = MHSanaei/3x-ui
git merge upstream/main       # 合并 3x-ui 的新版本
```

---

## 🗺️ 路线图

- [x] ⚡ 一键配置(批量生成节点)—— 后端 + UI,支持 **REALITY / Shadowsocks-2022** 协议选择
- [ ] 🔀 一键中转(入口机 ↔ 落地机,Xray 出站链式)
- [ ] 一键配置扩展更多协议 / 传输(VMess、gRPC、XHTTP…)
- [ ] VPS 一键安装脚本
- [ ] 界面文字多语言化

---

## 📄 License & 致谢

**GPLv3**。本项目基于 [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui);代理内核为 [XTLS/Xray-core](https://github.com/XTLS/Xray-core)。向上游作者致谢。

> 本项目仅供合法的自建自用与学习研究。
