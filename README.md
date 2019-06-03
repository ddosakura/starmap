# starmap

星图（Star Map）

## What's Star-Map

+ 领域（Domain）/星座（Constellation）/星团（Star Cluster）
    + 领域关系：包含/相交/相离
+ 知识点（Point）/星（Star）
    + 知识点关系：依赖/被依赖
+ 知识（Knowledge）/摘要（Summary）/链接（URL）/博文（blog）
    + 知识属性：难易度

## services

### api (open api) - 在代码中注释 API 参数

+ [ ] auth 认证接口
    + [x] user 用户
    + [x] role 角色 (only 权限查询)
    + [ ] oauth2 授权
+ [ ] sys 后台接口
    + [x] user 用户
    + [x] role 角色
    + [x] permission 权限
+ [ ] star 知识网接口
    + [ ] domain 领域
    + [ ] point 知识点
    + [ ] detail 知识
    + [ ] spider 爬虫
+ [ ] task 计划任务接口
+ [ ] fs 云储存/文件系统接口
+ [ ] mp 小程序接口

### gate - 准备尝试用非 go-micro-api 重做 API 网关

+ [ ] auth 认证接口
    + [ ] user 用户
    + [ ] role 角色 (only 权限查询)
    + [ ] oauth2 授权
+ [ ] sys 后台接口
    + [ ] user 用户
    + [ ] role 角色
    + [ ] permission 权限
+ [ ] star 知识网接口
    + [ ] domain 领域
    + [ ] point 知识点
    + [ ] detail 知识
    + [ ] spider 爬虫
+ [ ] task 计划任务接口
+ [ ] fs 云储存/文件系统接口
+ [ ] mp 小程序接口

### srv - 在 protobuf 中注释参数说明

+ [ ] auth 认证服务 (Authentication Service)
    + [x] user 用户
    + [x] role 角色
    + [x] permission 权限
    + [ ] authorization 授权
+ [ ] star 知识网服务
    + [ ] domain 领域
    + [ ] point 知识点
    + [ ] detail 知识
    + [ ] review 审核（包括自动审核）
+ [ ] spider 爬虫服务
    + [ ] target 目标
    + [ ] rule 策略
    + [ ] crawler 爬虫

### web (报废方案 - 转为前端路由的计划表)

+ [ ] manage 后台前端 - /manage
    + [ ] sys 后台管理页&后台系统接口 - /manage/sys & /manage/api/sys
        + [ ] user 用户
        + [ ] role 角色
        + [ ] permission 权限
        + [ ] authorization 授权
    + [ ] stat 系统状态 - /manage/stat & /manage/api/stat
        + [ ] log 日志
    + [ ] starmap 星图　- /manage/starmap & /manage/api/starmap
        + [ ] domain、point、detail 知识网管理页
        + [ ] spider 爬虫管理页
        + [ ] review 审核页
+ [ ] starmap 星图前端 - /
    + [ ] /user 个人主页
    + [ ] /user/login & /user/register 登录注册页
    + [ ] /user/:uid 用户主页
    + [ ] /map 星图页
    + [ ] /map/star 知识点页
    + [ ] /map/star/:id 知识页
    + [ ] /map/scan 扫描（爬虫）配置页

## Build Warning

1. `GOPATH` in Makefile
    > 开发时由于 $GOPTH 配置了多个路径，所以单独使用了环境变量 $ORIGIN_GOPATH 来获取了第一个 GOPATH
    ```
    #GOPATH:=$(shell go env GOPATH)
    GOPATH:=${shell echo $$ORIGIN_GOPATH}
    ```
