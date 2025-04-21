<div align=center>
	<img src="docs/images/favicon.ico"/>
</div>

中文 | [English](./README.md)

## 介绍
Scope Sentry是一款具有分布式资产测绘、子域名枚举、信息泄露检测、漏洞扫描、目录扫描、子域名接管、爬虫、页面监控功能的工具，通过构建多个节点，自由选择节点运行扫描任务。当出现新漏洞时可以快速排查关注资产是否存在相关组件。

分布式搭建参考文章：[https://mp.weixin.qq.com/s/xfgRxUjljoQ8KzacblktxA](https://mp.weixin.qq.com/s/xfgRxUjljoQ8KzacblktxA)

服务器推荐：[lightnode](https://www.lightnode.com/?inviteCode=CQ11JU&promoteWay=LINK)

## 编程语言
服务端：python - FastApi

扫描端：go

前端：vue - vue-element-plus-admin

## 网址

- 官网&安装教程：[https://www.scope-sentry.top](https://www.scope-sentry.top)
- Github: [https://github.com/Autumn-27/ScopeSentry](https://github.com/Autumn-27/ScopeSentry)
- 扫描端源码：[https://github.com/Autumn-27/ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan)
- 前端源码：[https://github.com/Autumn-27/ScopeSentry-UI](https://github.com/Autumn-27/ScopeSentry-UI)
- 插件市场: [插件市场](https://plugin.scope-sentry.top/)
- 插件模板：[https://github.com/Autumn-27/ScopeSentry-Plugin-Template](https://github.com/Autumn-27/ScopeSentry-Plugin-Template)



## 插件流程图

<img src="流程图.svg"/>


## 目前功能
- 插件系统（通过拓展的方式加入任何工具）
- 子域名枚举
- 子域名接管检测
- 端口扫描
- ICP自动化收集
- APP自动化收集
- 小程序自动化收集
- 资产识别
- 目录扫描
- 漏洞扫描
- 敏感信息泄露检测(支持扫描pdf)
- URL提取
- 爬虫
- 页面监控
- 自定义WEB指纹
- POC导入
- 资产分组
- 多节点扫描
- webhook
- 数据导出

## To DO
- 弱口令爆破
- 关系图

## 安装

安装教程见[官网](https://www.scope-sentry.top)

## 交流
联系方式见下方

## 赞助 && 合作
如果项目对您有帮助，赞助作者一杯咖啡吧~

<img src="docs/images/zfb.png" alt="WX" width="200"/>
<img src="docs/images/wx.jpg" alt="WX" width="200"/>


## 截图

### 登录

![alt text](docs/images/login.png)

### 首页面板
![alt text](docs/images/index-cn.png)

## 插件系统
![alt text](docs/images/plugin-cn.png)
![alt text](docs/images/plugin-m-cn.png)

## 资产数据
### 资产
![alt text](docs/images/asset-cn.png)
![alt text](docs/images/screenshot1.png)
![alt text](docs/images/screenshot2.png)
![alt text](docs/images/asset-change.png)
### 快捷语法搜索：
![alt text](docs/images/search.gif)

## 根域名
![alt text](docs/images/rootdomain-cn.png)
### 子域名
![alt text](docs/images/subdomain-cn.png)

### 子域名接管
![alt text](docs/images/subt-cn.png)

### APP
![alt text](docs/images/app-cn.png)

### 小程序
![alt text](docs/images/mp-cn.png)

### URL
![alt text](docs/images/url-cn.png)

### 爬虫
![alt text](docs/images/craw-cn.png)

### 敏感信息
![alt text](docs/images/sns-cn.png)

### 目录扫描
![alt text](docs/images/dir-cn.png)

### 漏洞
![alt text](docs/images/vul-cn.png)

### 页面监控
![alt text](docs/images/page-cn.png)
![alt text](docs/images/page-change.png)
## 项目

![](docs/images/project-cn.png)


## 项目资产聚合
### 面板-概况
![](docs/images/project-dsh.png)
### 子域名
![](docs/images/project-subdomain.png)
### 端口
![](docs/images/project-port.png)
### 服务
![](docs/images/project-server.png)

## 任务

![](docs/images/create-task-cn.png)

## 任务进度

![](docs/images/task-pg-cn.png)

## 节点

![](docs/images/node-cn.png)


Discord:

[https://discord.gg/GWVwSBBm48](https://discord.gg/GWVwSBBm48)

QQ:

<img src="docs/images/qq.jpg" alt="QQ" width="200"/>

WX:

<img src="docs/images/wx-2.jpg" alt="WX" width="200"/>

群满可以关注公众号后台私信拉群
<img src="docs/images/wx.png" alt="WX" width="200"/>

# 许可证
该项目所有分支遵循AGPL-3.0，另外需要遵循附加条款：
1. 本软件的商业用途需要单独的商业许可。
2. 公司、组织和营利性实体在使用、分发或修改本软件之前必须获得商业许可。
3. 个人和非营利组织可以根据 AGPL-3.0 的条款自由使用本软件。
4. 如有商业许可查询，请联系 rainy-autumn@outlook.com。
