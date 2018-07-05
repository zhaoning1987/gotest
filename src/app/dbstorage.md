<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [1. 人脸入库任务](#1-%E4%BA%BA%E8%84%B8%E5%85%A5%E5%BA%93%E4%BB%BB%E5%8A%A1)
  - [1.1 创建任务](#11-%E5%88%9B%E5%BB%BA%E4%BB%BB%E5%8A%A1)
  - [1.2 启动任务](#12-%E5%90%AF%E5%8A%A8%E4%BB%BB%E5%8A%A1)
  - [1.3 停止任务](#13-%E5%81%9C%E6%AD%A2%E4%BB%BB%E5%8A%A1)
  - [1.4 删除任务](#14-%E5%88%A0%E9%99%A4%E4%BB%BB%E5%8A%A1)
  - [1.5 获得任务](#15-%E8%8E%B7%E5%BE%97%E4%BB%BB%E5%8A%A1)
  - [1.6 查询任务信息](#16-%E6%9F%A5%E8%AF%A2%E4%BB%BB%E5%8A%A1%E4%BF%A1%E6%81%AF)
  - [1.7 查询任务日志](#17-%E6%9F%A5%E8%AF%A2%E4%BB%BB%E5%8A%A1%E6%97%A5%E5%BF%97)
  - [1.7 下载任务日志](#17-%E4%B8%8B%E8%BD%BD%E4%BB%BB%E5%8A%A1%E6%97%A5%E5%BF%97)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

基本语义说明
* 资源表示方式（URI）。通过统一方式定位、获取资源（图片、二进制数据等）
	* HTTP， 网络资源，形如：`http://host/path`、`https://host/path`
	* Data，Data URI Scheme形态的二进制文件，形如：`data:application/octet-stream;base64,xxx`。

# 1. 人脸入库任务

## 1.1 创建任务

> 新建任务，单次请求创建单个任务。新建任务时需指定所属的group（即对哪个group进行人脸入库）
>
> 还可指定该次任务的相关配置（可选）
>
> 新建任务时还需传入一csv文件，文件内包含所有欲入库人脸的uri，tag，description

**Request**

```
POST /v1/tasks/create
Content-Type: multipart/form-data

form参数：
"group_name"
"pos_pitch"
"pos_yaw"
"pos_rol"
"width"
"height"
"blur_degree"
"thread_num"
"mode"
"file"

```

**Response**

```
200 OK
Content-Type: application/json
{
   "id" : ""
}
```

接口字段说明：

| 字段               | 取值   | 说明                     |
| :----------------  | :----- | :----------------------- |
| group_name         | string | 任务所属的人脸库group名，必选|
| pos_pitch          | int    | 允许的人脸最大俯仰角，超过则该人脸不会入库，可选，默认为0（即不限制） |
| pos_yaw            | int    | 允许的人脸最大偏航角，超过该值则该人脸不会入库，可选，默认为0（即不限制） |
| pos_rol            | int    | 允许的人脸最大翻滚角，超过该值则该人脸不会入库，可选，默认为0（即不限制） |
| width              | int    | 允许的人脸最小宽度，低于该值则该人脸不会入库，可选，默认为0（即不限制） |
| height             | int    | 允许的人脸最小高度，低于该值则该人脸不会入库，可选，默认为0（即不限制） |
| thread_num         | int    | 该任务起用的线程数，可选，默认为10 |
| mode               | int    | 入库时遇见错误图片时的处理方式，0为跳过继续执行，1为将任务停止，可选，默认为0 |
| file               | int    | 表单提交的csv文件，必选 |
| id                 | string | 创建成功返回的任务id |



## 1.2 启动任务

> 启动入库任务，若任务处于已完成状态则返回错误

**Request**

```
GET /v1/tasks/<task_id>/start
```

**Response**

```
200 OK
```

接口字段说明：

| 字段 | 取值 | 说明 |
| :--- | :--- | :--- |



## 1.3 停止任务

> 停止入库任务，若任务不处于启动状态则返回错误
> 停止后若再次调用启动任务接口，则任务会从上一次中断点继续执行
> 若服务意外停止，则重启服务时会将所有之前处于启动状态的任务置为停止状态

**Request**

```
GET /v1/tasks/<task_id>/stop
```

**Response**

```
200 OK
```

接口字段说明：

| 字段 | 取值 | 说明 |
| :--- | :--- | :--- |



## 1.4 删除任务

> 无论任务处于何种状态，将停止并删除任务，且不能恢复（之后也无法再继续启动）

**Request**

```
GET /v1/tasks/<task_id>/delete
```

**Response**

```
200 OK
```

接口字段说明：

| 字段 | 取值 | 说明 |
| :--- | :--- | :--- |



## 1.5 获得任务

> 基于人脸库group名，获得属于该group下所有task的id

**Request**

```
POST /v1/tasks/list
Content-Type: application/json

{
    "group_name": ""
    "status": ""
}
```

**Response**

```
200 OK
Content-Type: application/json

{
    "ids": [
        "AAAO054X4wzh",
        "AAAHdaa3wwzh"
    ]
}

```

接口字段说明：

| 字段       | 取值   | 说明                |
| :-----     | :----- | :-----------------|
| group_name | string | 人脸库group名，必选      |
| status     | int    | 只查询指定状态的任务，可选，默认为0（即查找所有）  |
| ids.[]     | string | 入库任务的id        |



## 1.6 查询任务信息

> 获得任务的相关信息，包括所属的人脸库名、配置信息、当前状态、该任务所需处理的总数、该任务当前已处理的数目

**Request**

```
GET /v1/tasks/<task_id>/detail
```

**Response**

```
200 OK
Content-Type: application/json
{
    "group_name" : "",
    "config":{
        "pos_pitch": "",
        "pos_yaw": "",
        "pos_rol": "",
        "width": "",
        "height": "",
        "mode": "",
        "thread_num": ""
    },
    "total_num" : "",
    "handled_num": ""
    "state": ""
    "last_error": ""
}
```

接口字段说明：

| 字段               | 取值   | 说明                     |
| :----------------  | :----- | :----------------------- |
| group_name         | string | 任务所属的人脸库group名|
| config.pos_pitch   | int    | 配置的允许的人脸最大俯仰角 |
| config.pos_yaw     | int    | 配置的允许的人脸最大偏航角 |
| config.pos_rol     | int    | 配置的允许的人脸最大翻滚角 |
| config.width       | int    | 配置的允许的人脸最小宽度 |
| config.height      | int    | 配置的允许的人脸最小高度 |
| config.mode        | int    | 配置的入库模式 |
| config.thread_num  | int    | 配置的线程数 |
| total_num          | int    | 该任务一共需要处理的人脸数 |
| handled_num        | int    | 该任务当前已处理完的人脸数 |
| state              | int    | 该任务当前状态，1为已创建未执行，2为正在执行中，3为已停止，4为已完成 |
| last_error         | string | 当入库模式为出错即停止时，该字段返回造成停止的人脸图片信息即错误原因 |



## 1.7 查询任务日志

> 每当入库人脸出现错误时，系统会将错误记录到日志
> 可随时查询任务日志，无论任务是否已结束

**Request**

```
GET /v1/tasks/<task_id>/log
```

**Response**

```
200 OK
Content-Type: application/json
{
    "logs": [
        {
            "uri":"",
            "code":"",
            "message":""
        }
    ]
}
```

接口字段说明：

| 字段       | 取值   | 说明                |
| :-----     | :----- | :-----------------|
| uri        | string | 出错的图片uri      |
| code       | int    | 错误码  |
| message    | string | 错误信息        |


错误码说明：

| 错误码 | 说明                |
| :-----| :-----------------|
| 101   | 图片uri不存在，无法下载  |
| 102   | 图片无法打开  |
| 201   | 图片不包含人脸  |
| 202   | 人脸俯仰角超过阈值  |
| 203   | 人脸偏航角超过阈值  |
| 204   | 人脸翻滚角超过阈值  |
| 205   | 人脸宽度小于阈值  |
| 206   | 人脸高度小于阈值  |
| 301   | 人脸信息服务无响应  |
| 302   | 人脸入库服务无响应  |
| 303   | 系统错误  |


## 1.7 下载任务日志

**待定，看了源码，貌似restrpc框架不能很好的支持下载，直接操作response.Writer的话，框架依然会在末尾带上一个空json"{}"**

> 除了可以查询日志外，系统还提供下载日志，下载文件为csv文件。格式为：人脸图片uri，错误码，错误信息

**Request**

```
GET /v1/tasks/<task_id>/download_log
```

**Response**

```
200 OK
Content-Type: application/octet-stream
***文件内容***
```