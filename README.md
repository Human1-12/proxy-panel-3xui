<div align="center">

# 3X-UI 中文增强版

**免费开源 · 无授权码 · 对标 X-Panel Pro**

基于 3X-UI 的中文增强、多节点管理与一键配置面板

<p>
  <img src="https://img.shields.io/badge/license-GPLv3-blue.svg" alt="License: GPLv3">
  <img src="https://img.shields.io/badge/based%20on-3X--UI%20v3.4.2-brightgreen.svg" alt="Based on 3X-UI">
  <img src="https://img.shields.io/badge/Xray--core-v26.6.27-orange.svg" alt="Xray-core">
  <img src="https://img.shields.io/badge/X--Panel%20Pro-开源替代目标-8a2be2.svg" alt="X-Panel Pro alternative">
</p>

</div>

---

## 项目定位

**3X-UI 中文增强版**不是只替换界面文字的“汉化包”，而是在 [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui) 基础上继续开发的完整面板：

- 保留 3X-UI 的主流协议、订阅、流量管理、Telegram、API 和运维能力；
- 增加批量一键配置、自动放行端口、多节点、远程探测和中转能力；
- 所有自研功能继续采用 GPLv3 开源，不设置授权码，不做联网验权；
- 以**免费开源替代 X-Panel Pro 的常用收费能力**为产品目标。

> [!IMPORTANT]
> **当前结论（2026-08-13）：本项目已经覆盖 X-Panel Pro 的一批核心实用能力，但还不能严谨地宣称“已拥有 X-Panel Pro 收费版全部功能”。** 独立客户端限速、按设备数量限制、Telegram 多面板管理、WebSSH、线路/IP/DNS 检测、远程救援快照、积分与购买机器人等仍需开发或逐项验收。完成下方“全功能替代清单”并通过真实多机测试后，才会正式标注为“全功能替代版”。

本项目与 X-Panel 官方无隶属或授权关系。“X-Panel”仅用于功能对比说明。

---

## 已实现功能

### Xray 与面板基础能力

- 协议：VLESS、VMess、Trojan、Shadowsocks、WireGuard、Hysteria2、SOCKS、HTTP、Dokodemo-door、TUN、MTProto；
- 传输与安全：TCP、mKCP、WebSocket、gRPC、HTTPUpgrade、XHTTP，以及 TLS、XTLS、REALITY；
- 客户端流量配额、到期时间、IP 限制、在线状态、二维码、分享链接和订阅；
- 独立订阅服务器、REST API、Swagger、API Token、WebSocket；
- Telegram Bot、邮件与事件通知；
- 多语言界面、明暗主题、Fail2ban；
- SQLite 与 PostgreSQL；
- 继承 3X-UI 的备份、证书、日志、系统状态和常规运维能力。

### 中文增强版新增能力

- ⚡ **一键配置**：一次批量创建 1—100 个入站；
- 四种免证书预设：
  - VLESS + TCP + REALITY + Vision；
  - Shadowsocks-2022；
  - VMess + TCP；
  - VLESS + TCP；
- 每个节点使用独立 UUID、密钥或 subId；
- 从指定起始端口确定性分配，并跳过面板已占用端口；
- 创建完成后尽力自动放行主机防火墙端口；
- 提供明确的批量成功、失败和防火墙处理结果；
- Linux 一键安装、更新和管理脚本；
- Release 资产 SHA-256 完整性校验；
- 多节点管理、证书指纹、mTLS、探测历史和节点上下线通知；
- 中转/路由引擎基础能力。

---

## 与 X-Panel Pro 的覆盖情况

| 能力 | 当前状态 | 说明 |
|---|:---:|---|
| 主流 Xray 协议与入站管理 | ✅ | 基于现代 3X-UI，协议覆盖完整 |
| 批量 VLESS REALITY Vision | ✅ | 默认可一次生成 10 个，最多 100 个 |
| SS2022 / VMess / VLESS TCP 一键预设 | ✅ | 已在后端与前端接入 |
| 自动分配端口与放行防火墙 | ✅ | 防火墙操作为尽力执行，失败会返回说明 |
| Linux 一键安装与升级 | ✅ | 安装脚本已经可用，并从本仓库 Release 下载 |
| 无授权码、无联网验权 | ✅ | 开源版不需要购买授权 |
| 多节点、mTLS、远程探测 | ✅ / 待压测 | 代码已具备；建议继续做真实多 VPS 长时间验收 |
| 一键中转完整向导 | 🟡 | 已有节点、出站与中转引擎；仍需补齐 X-Panel 式端到端向导和生产验收 |
| 流量、到期、IP 限制 | ✅ | 继承 3X-UI 能力 |
| 按设备数量限制 | 🟡 | 当前 IP 限制不能完全等同于 X-Panel 的设备限制 |
| 独立客户端速度限制 | ❌ | 尚未发现等价实现 |
| 任意每月重置日期 | 🟡 | 已有周期重置框架，仍需补齐并验收 1—31 日配置 |
| Telegram 基础管理与通知 | ✅ | 已有机器人及节点/系统事件通知 |
| Telegram 多面板切换、远程取链接、一键配置 | 🟡 | 需要按 X-Panel Pro 的具体交互继续补齐 |
| WebSSH 安装与管理 | ❌ | 待开发 |
| 线路、IP 质量、地区 DNS 检测 | ❌ | 待开发 |
| 多套 Pro 主题 | 🟡 | 已有明暗主题与多语言，未逐项复刻其五套主题 |
| 远程快照与救援恢复 | ❌ | 常规备份不等同于远程救援快照 |
| 积分、签到、兑换码、抽奖、购买机器人 | ❌ | 可作为可选运营模块开发，不影响代理面板核心功能 |

---

## 使用流程

完整链路：**安装面板 → 打开管理页 → 一键或手动创建入站 → 导出链接/订阅 → 导入客户端**。

1. 安装完成后，以脚本输出的访问地址、用户名和密码登录。
2. 进入“入站”页面：
   - 点击“⚡ 一键配置”，选择协议、数量和起始端口；
   - 或点击“添加入站”，手动配置任意受支持协议。
3. 从入站或客户端操作中复制标准分享链接、二维码或订阅地址。
4. 导入 v2rayN、v2rayNG、NekoBox 等兼容客户端。

> REALITY 预设可在没有自有域名和证书的情况下部署，但目标站点、网络环境及相关使用行为仍需符合当地法律和服务条款。

---

## Linux VPS 一键安装

需要一台支持的 Linux VPS 和 root 权限。安装脚本会读取本仓库最新 Release，并在支持的环境中配置服务：

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Human1-12/proxy-panel-3xui/main/install.sh)
```

安装完成后：

- 终端会显示访问地址、用户名和密码；
- 安装结果同时写入 `/etc/x-ui/install-result.env`，权限为 `600`；
- 使用 `x-ui` 打开管理菜单；
- 可使用 `x-ui status`、`x-ui log`、`x-ui restart`、`x-ui update` 等命令。

> [!WARNING]
> 不要继续使用旧文档中的固定 `admin/admin` 和固定 `2053` 端口假设。以安装脚本本次实际输出为准，并妥善保存凭据。

---

## Windows 源码构建

开发环境需要 Go 1.26.5+、Node.js 22+ 和支持 CGO 的 C 编译器。

```powershell
git clone https://github.com/Human1-12/proxy-panel-3xui.git
cd proxy-panel-3xui
copy .env.example .env
.\build.ps1
.\run-dev.ps1
```

开发模式下还需要按项目配置准备对应平台的 Xray Core、`geoip.dat` 与 `geosite.dat`。

---

## 一键配置 API

```http
POST /panel/api/inbounds/oneclick/reality
Authorization: Bearer <API Token>
Content-Type: application/json

{
  "count": 10,
  "portStart": 20000,
  "protocol": "reality",
  "remarkPrefix": "",
  "dest": "www.microsoft.com:443"
}
```

`protocol` 支持：

| 值 | 生成配置 |
|---|---|
| `reality` | VLESS + TCP + REALITY + Vision |
| `ss2022` | Shadowsocks 2022-blake3-aes-256-gcm |
| `vmess` | VMess + TCP，无 TLS |
| `vlessTcp` | VLESS + TCP，security=none |

空请求体会使用默认值：10 个节点、起始端口 20000、REALITY 预设。每个节点拥有独立凭据，端口会自动避开数据库中已经使用的入站端口和 Xray API 保留端口。

---

## 全功能替代路线图

### P0：达到“X-Panel Pro 核心功能完整替代”

- [x] 批量 REALITY / SS2022 / VMess / VLESS TCP 一键配置；
- [x] 自动端口分配与主机防火墙放行；
- [x] Linux 一键安装、升级和 Release 完整性校验；
- [x] 多节点、mTLS、远程探测和中转引擎；
- [ ] 完成多节点与中转的真实多 VPS 端到端测试；
- [ ] 独立客户端上传/下载速度限制；
- [ ] 区分 IP 限制与设备数量限制；
- [ ] 每月 1—31 日任意日期流量重置；
- [ ] Telegram 多面板、远程节点链接和一键配置；
- [ ] WebSSH 与远程安装；
- [ ] 线路、IP 质量和地区 DNS 检测；
- [ ] 远程快照、回滚与救援恢复。

### P1：体验与运维增强

- [ ] 完整的一键中转图形向导；
- [ ] 多套中文主题；
- [ ] BBR、FQ、TCP Fast Open 等系统调优向导；
- [ ] 更完整的通知模板和定时报告；
- [ ] 全功能对照表自动化回归测试。

### P2：可选运营模块

- [ ] 积分、签到与兑换码；
- [ ] 抽奖模块；
- [ ] 购买机器人。

> 只有 P0 清单完成、并经过真实服务器验收后，项目才会把宣传语升级为“拥有 X-Panel Pro 全部核心功能”。

---

## 技术栈

- 后端：Go 1.26.5、Gin、GORM；
- 前端：React 19、TypeScript 6、Ant Design 6、Vite 8；
- 核心：Xray Core v26.6.27；
- 存储：SQLite / PostgreSQL。

主要增强代码包括：

- `internal/web/service/oneclick.go`：批量一键配置；
- `internal/web/controller/oneclick.go`：一键配置 REST 接口；
- `internal/web/service/node.go`：节点、探测、mTLS 与中转能力；
- `frontend/src/pages/inbounds/InboundsPage.tsx`：一键配置界面。

---

## 跟随上游 3X-UI

建议将原项目配置为 `upstream`，定期合并并人工解决增强代码冲突：

```bash
git remote add upstream https://github.com/MHSanaei/3x-ui.git
git fetch upstream
git merge upstream/main
```

---

## License 与致谢

本项目采用 **GPLv3**。基于 [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui)，代理核心来自 [XTLS/Xray-core](https://github.com/XTLS/Xray-core)。感谢所有上游作者与贡献者。

本项目仅供合法的自建、自用、运维和学习研究。使用者应自行遵守所在地法律、云服务商规则和网络服务条款。
