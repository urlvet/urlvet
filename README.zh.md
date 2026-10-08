<div align="center">

# url.vet

**链接看起来可疑？就用 url.vet 检测一下。**

**网站：[url.vet](https://url.vet)**

开源 URL 安全检测引擎 — 粘贴任意链接，实时获取信任评分、完整可解释的检测报告，以及带有页面预览的可分享安全报告。

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-AGPL--3-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/urlvet/urlvet?style=social)](https://github.com/urlvet/urlvet)
[![Last Commit](https://img.shields.io/github/last-commit/urlvet/urlvet)](https://github.com/urlvet/urlvet/commits/main)

[English](./README.md) | 中文

[⚡ 快速开始](#快速开始) · [⚙️ 检测引擎](#检测引擎) · [🏛 架构](#架构) · [📚 文档](#文档) · [🤝 贡献](#贡献)

<sub>英文名 URLvet · 原名 SafeSurf</sub>

</div>

---

## 钓鱼检测演示

> 粘贴 URL → 实时获取**信任评分、检测结论与详细报告**。

![钓鱼检测演示](assets/demo.gif)

在 [url.vet](https://url.vet) 上免费使用，无需注册。

## 快速开始

```bash
git clone https://github.com/urlvet/urlvet.git
cd urlvet
make start
```

打开 Web UI：**[localhost:3000](http://localhost:3000)**

详细配置指南：[docs/setup.md](docs/setup.md)

## 功能概览

- 实时扫描，即时出结果
- 21 个分析器，39 个信号，完全可解释
- HTTP API + Web UI + Chrome 扩展
- 可解释评分（无黑盒 ML 模型）
- 简单的 Docker 一键部署，完全可自托管

## 横向对比

| 功能 | url.vet | VirusTotal | Google Safe Browsing | URLScan.io | CheckPhish |
|------|---------|------------|----------------------|------------|------------|
| 实时抓取，即时结果 | ✅ | 部分 | ❌ | 部分 | 部分 |
| 解释每条检测结论 | ✅ | 部分 | ❌ | 部分 | 部分 |
| 对新手友好的界面 | ✅ | 部分 | 部分 | 部分 | 部分 |
| 凭据表单检测 | ✅ | ❌ | ❌ | 部分 | ✅ |
| 跟踪重定向链 | ✅ | ✅ | ❌ | ✅ | ✅ |
| 详细技术洞察 | ✅ | ❌ | ❌ | ✅ | 部分 |
| 页面实时预览 | ✅ | ❌ | ❌ | ✅ | ✅ |
| AI/ML 检测 | ❌ | ✅ | ✅ | 部分 | ✅ |
| 已知钓鱼数据库覆盖 | ✅¹ | ✅ | ✅ | 部分 | 部分 |
| 批量扫描多个 URL | ❌ | ✅ | ✅ | ✅ | ✅ |
| 浏览器保护 | ✅ | ✅ | ✅ | ✅ | ❌ |
| 开源 | ✅ | ❌ | ❌ | ❌ | ❌ |

¹ 检查 PhishTank、URLhaus 和 Google Safe Browsing。自行部署时，URLhaus 和 Safe Browsing 需要免费的 API 密钥。

快速扫描器（如 Google Safe Browsing）仅做数据库查询，无法解释结论也不支持实时扫描。深度爬虫（如 URLScan.io）速度太慢。url.vet 填补了这一空白——实时分析、逐信号解释，而且完全开源。

## 适用人群

- 普通用户：检测可疑链接
- 注重隐私的用户：无需账号，没有 Cookie 和追踪器，日志中不记录 IP 和被扫描的链接。你检测的网站只会看到我们的服务器，永远看不到你
- 开发者：集成 URL 分析能力
- 安全团队：构建检测流水线
- 安全研究人员

## API 示例

通过 HTTP 分析 URL：

```bash
curl "http://localhost:8080/api/v1/analyze?url=https://example.com"
```

**示例响应：**

```json
{
  "url": "https://example.com",
  "trust_score": 100,
  "verdict": "Safe",
  "reasons": {
    "good_reasons": [...]
  }
}
```

完整响应结构 → [docs/api.md#example](docs/api.md#example)

## 检测引擎

**21 个并发 goroutine** 运行于 **7 个信号类别**，产生 **39 个独立信号**。每项检查都会输出一条原因字符串（正面、负面或中性），因此最终评分始终完全可解释，没有黑盒判定。

评分公式：`finalScore = clamp(50 + (trustScore − riskScore) × 0.5)` → **高风险** < 30 · **可疑** 30–64 · **安全** ≥ 65

短链接会先被追踪（包括连续的多个短链接），实际扫描的是它最终指向的页面。

> 50 为中性基线——没有任何信号的 URL 得分恰好为 50（可疑），这是对未知 URL 的正确默认值。信任信号拉高评分，风险信号拉低评分，各以 0.5× 加权，两者互不压制。公式运算前，两个分数分别被限制在 0–100 之间，防止单一灾难性信号淹没其他所有信息。

**URL 信号**（8 项检测）

1. 以原始 IP 地址作为主机名（常见规避手法）
2. Punycode / IDN 编码（同形域名欺骗）
3. URL 缩短服务（追踪至目标页面；无法追踪的最高只判为可疑）
4. URL 长度过长（超过 150 个字符；知名网站不计）
5. URL 路径深度过深（超过 6 级路径；知名网站不计）
6. URL 路径中包含钓鱼关键词（login、verify、secure、update……）
7. 子域名数量过多
8. 主机名中包含非 ASCII Unicode 字符（IDN 同形攻击，例如使用西里尔字母 а 的 аpple.com）

**HTTP / 网络**（4 项检测，单次 HTTP 请求）

9. 重定向链跳转次数（仅在跳出本站时计入）
10. 跨域重定向（最终目标与源域名不同，同一品牌的域名除外）
11. HSTS 支持（读取被扫描网站本身的响应，而非最后一跳）
12. HTTP 状态码

**DNS**（3 项检测）

13. NS 记录有效性
14. MX 记录有效性
15. IP 解析

**TLS / SSL**（2 项检测，单次 TLS 握手）

16. TLS 存在性，以及证书是否签发给其他主机名
17. 证书链——有效性、过期时间、签发机构、CT 日志状态、已知恶意指纹

**域名情报**（6 项检测）

18. 域名排名（在全球 Top 100 万列表中的位置）
19. TLD 信任度 / 风险度 / ICANN 状态（托管平台与建站服务上的子域名单独评估，不继承平台的声誉）
20. 通过 RDAP/WHOIS 获取域名注册时间（90 天内 = 风险；超过 1、3、5 年 = 信任逐步增加）
21. DNSSEC（DNS 响应的加密完整性）
22. Shannon 熵值评分（标记算法生成的域名）
23. 仿冒域名检测（访问量前 5,000 的网站）与组合仿冒检测（514 个精选品牌）

**内容分析**（12 项检测）

24. 未排名或新注册域名上的登录表单
25. 支付表单（信用卡、CVV 字段）
26. 个人信息表单
27. 隐藏的 `<iframe>`（凭据窃取 / 点击劫持载体）
28. 跟踪像素（1×1 隐藏图片）
29. 页面标题中声称的品牌与托管域名对比（514 个品牌；同形字母会先还原）
30. 表单提交至外部域名
31. 未加密 HTTP 上的密码字段
32. 脚本或 meta refresh 跳转（跳至原始 IP 或其他网站）
33. 托管商或 CDN 的警告页面（Cloudflare 钓鱼拦截页、451 下架）
34. 位于 WordPress 系统目录中的页面（被入侵网站上的钓鱼工具包）
35. 程序下载（在任何人都可上传的平台上最高只判为可疑）

**威胁情报**（4 项检测）

36. 本地威胁列表（默认 PhishTank；URLhaus、OpenPhish 需手动启用）
37. 同一单一所有者主机上的其他页面已被列入
38. Google Safe Browsing / Web Risk（需手动启用；Safe Browsing 只接收哈希前缀）
39. PhishTank API（仅在其列表未加载时使用）

![url.vet 分析器流水线](assets/pipeline.png)

## 局限性

- 基于启发式的检测可能产生误报
- 无 ML 模型（有意为之，优先保证可解释性与可审计性）

本工具不提供安全保证，请与其他防御手段配合使用。

## 架构

四个容器化服务运行于共享的 Docker bridge 网络中。Go 后端是唯一对外部 API 发起出站请求的服务——前端、Chrome 和缓存均为严格内部通信。

![url.vet 架构图](assets/architecture.png)

| 服务 | 职责 |
|---|---|
| `urlvet-web` | SvelteKit UI — :3000（生产）· :5173（开发） |
| `urlvet-backend` | Go REST API 与分析引擎 — :8080 |
| `urlvet-chrome` | 无头 Chrome — WebSocket :9222 |
| `urlvet-valkey` | Valkey（兼容 Redis）— :6379，LRU 缓存，持久化卷 |

### 请求生命周期

1. 通过 UI 或 REST API 提交 URL
2. 后端验证并规范化 URL（缺少协议头时自动推断）
3. 检查 Valkey 缓存——命中则立即返回完整结果，无需重新分析
4. 短链接会被追踪至目标页面，改为扫描该页面
5. 未命中时：21 个 goroutine 通过 `sync.WaitGroup` 并发启动；每个任务的 panic 单独恢复，不影响整体请求
6. 收集结果 → 汇总评分 → 分配结论
7. 完整结果缓存至 Valkey（TTL 24 小时）并记录到扫描历史
8. 返回响应——信任评分、检测结论、逐信号原因、重定向链、页面截图、各任务耗时

```text
server/
  cmd/urlvet/           入口点
  internal/analyzer/    goroutine 运行器、任务定义、评分汇总
  internal/service/
    checks/             各个独立分析器实现
    screenshot/         无头 Chrome 集成
    cache/              Valkey 客户端
    threatfeeds/        威胁列表、PhishTank、Safe Browsing、Web Risk
    typosquat/          品牌相似度引擎
web/website/            SvelteKit UI
web/chrome-extension/   浏览器扩展
docker/                 开发与生产 Compose 配置
docs/                   API、部署、架构、安全文档
```

## 文档

| | |
|---|---|
| [部署配置](docs/setup.md) | 本地与 Docker 部署、Makefile 命令 |
| [环境变量](docs/configuration.md) | 所有环境变量说明 |
| [生产部署](docs/deployment.md) | VPS、反向代理、防火墙 |
| [API 参考](docs/api.md) | 接口、频率限制、响应示例 |
| [架构说明](docs/architecture.md) | 服务、请求生命周期、检测引擎 |
| [安全](docs/security.md) | 管理员认证、密码哈希 |
| [性能](docs/performance.md) | 延迟、资源占用、调优 |
| [设计决策](docs/design-decisions.md) | 为什么这样构建 |
| [运维](docs/maintenance.md) | 缓存、日志、备份 |
| [词汇表](docs/glossary.md) | 术语与缩写 |

交互式 API 文档（Swagger UI）：[api.url.vet/swagger/index.html](https://api.url.vet/swagger/index.html)

## 引用

如果您在学术或研究工作中使用了本项目，请引用 — 参见 [CITATION.cff](CITATION.cff)。

## 许可证

Copyright (C) 2021–2026 Abhishek K P

url.vet 采用双重许可：

- **社区版** — [GNU Affero General Public License v3.0](LICENSE)。可免费使用、修改和自托管。任何通过网络运行的修改版本必须向用户开放其源代码。
- **商业版** — 为无法遵守 AGPL-3.0 的组织（如闭源 SaaS）提供单独的[商业许可证](COMMERCIAL.md)。

## 贡献

- 发现 Bug？→ [提交 Issue](https://github.com/urlvet/urlvet/issues)
- 有问题或想法？→ [发起讨论](https://github.com/urlvet/urlvet/discussions)
- 想贡献代码？→ [CONTRIBUTING.md](.github/CONTRIBUTING.md)

如果这个项目对您有帮助，欢迎点个 Star ⭐