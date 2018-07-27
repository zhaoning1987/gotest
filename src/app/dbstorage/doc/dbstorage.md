<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [1. 人脸入库服务](#1-%E4%BA%BA%E8%84%B8%E5%85%A5%E5%BA%93%E6%9C%8D%E5%8A%A1)
  - [1.1 创建任务](#11-%E5%88%9B%E5%BB%BA%E4%BB%BB%E5%8A%A1)
  - [1.2 启动任务](#12-%E5%90%AF%E5%8A%A8%E4%BB%BB%E5%8A%A1)
  - [1.3 停止任务](#13-%E5%81%9C%E6%AD%A2%E4%BB%BB%E5%8A%A1)
  - [1.4 删除任务](#14-%E5%88%A0%E9%99%A4%E4%BB%BB%E5%8A%A1)
  - [1.5 获得任务](#15-%E8%8E%B7%E5%BE%97%E4%BB%BB%E5%8A%A1)
  - [1.6 查询任务信息](#16-%E6%9F%A5%E8%AF%A2%E4%BB%BB%E5%8A%A1%E4%BF%A1%E6%81%AF)
  - [1.7 查询任务日志](#17-%E6%9F%A5%E8%AF%A2%E4%BB%BB%E5%8A%A1%E6%97%A5%E5%BF%97)
  - [1.8 下载任务日志](#18-%E4%B8%8B%E8%BD%BD%E4%BB%BB%E5%8A%A1%E6%97%A5%E5%BF%97)
- [2. 人脸入库工具](#2-%E4%BA%BA%E8%84%B8%E5%85%A5%E5%BA%93%E5%B7%A5%E5%85%B7)
  - [2.1 依赖](#21-%E4%BE%9D%E8%B5%96)
  - [2.2 功能说明](#22-%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E)
  - [2.3 实现细节](#23-%E5%AE%9E%E7%8E%B0%E7%BB%86%E8%8A%82)
  - [2.4 字段录入](#24-%E5%AD%97%E6%AE%B5%E5%BD%95%E5%85%A5)
  - [2.5 使用说明](#25-%E4%BD%BF%E7%94%A8%E8%AF%B4%E6%98%8E)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

基本语义说明
* 资源表示方式（URI）。通过统一方式定位、获取资源（图片、二进制数据等）
	* HTTP， 网络资源，形如：`http://host/path`、`https://host/path`
	* Data，Data URI Scheme形态的二进制文件，形如：`data:application/octet-stream;base64,xxx`。

# 1. 人脸入库服务

## 1.1 创建任务

> 新建任务，单次请求创建单个任务。新建任务时需指定所属的group（即对哪个group进行人脸入库）   
> 还可指定该次任务的相关配置（可选）   
> 新建任务时还需传入一csv文件，文件内包含所有欲入库人脸的uri，tag，description

**Request**

```
POST /v1/task/new
Content-Type: multipart/form-data

form参数：
"group_name"
"pos_pitch"
"pos_yaw"
"pos_rol"
"width"
"height"
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
| group_name         | string    | 任务所属的人脸库group名，必选|
| pos_pitch          | string    | 允许的人脸最大俯仰角，超过则该人脸不会入库，可选，取值为[1, 90], 默认为空（即不限制） |
| pos_yaw            | string    | 允许的人脸最大偏航角，超过该值则该人脸不会入库，可选，取值为[1, 90], 默认为空（即不限制） |
| pos_rol            | string    | 允许的人脸最大翻滚角，超过该值则该人脸不会入库，可选，取值为[1, 90], 默认为空（即不限制） |
| width              | string    | 允许的人脸最小宽度，低于该值则该人脸不会入库，可选，取值为大于0的整数，默认为空（即不限制） |
| height             | string    | 允许的人脸最小高度，低于该值则该人脸不会入库，可选，取值为大于0的整数，默认为空（即不限制） |
| mode               | string    | 入库时遇见错误图片时的处理方式，0为跳过继续执行，1为将任务停止，可选，默认为0 |
| file               | string    | 表单提交的csv文件，必选 |
| id                 | string    | 创建成功返回的任务id |



## 1.2 启动任务

> 启动入库任务，若任务处于已完成状态则返回错误

**Request**

```
POST /v1/task/<task_id>/start
```

**Response**

```
200 OK
```

接口字段说明：

| 字段    | 取值    | 说明 |
| :------ | :----- | :----- |
| task_id | string | 任务id |



## 1.3 停止任务

> 停止入库任务，若任务不处于启动状态则返回错误   
> 停止后若再次调用启动任务接口，则任务会从上一次中断点继续执行   
> 若服务意外停止，则重启服务时会将所有之前处于启动状态的任务置为停止状态

**Request**

```
POST /v1/task/<task_id>/stop
```

**Response**

```
200 OK
```

接口字段说明：

| 字段    | 取值    | 说明 |
| :------ | :----- | :----- |
| task_id | string | 任务id |



## 1.4 删除任务

> 当任务不处于运行状态时（防止误操作），删除任务，且不能恢复（之后也无法再继续启动）

**Request**

```
POST /v1/task/<task_id>/delete
```

**Response**

```
200 OK
```

接口字段说明：

| 字段    | 取值    | 说明 |
| :------ | :----- | :----- |
| task_id | string | 任务id |



## 1.5 获得任务

> 获得指定人脸库下所有任务的id，查询时还可指定任务状态

**Request**

```
GET /v1/task/<group_name>/list?status=<status>
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
| group_name | string | 人脸库名称   |
| status     | string | 只查询指定状态的任务，可选 |
| ids.[]     | string | 入库任务的id        |



## 1.6 查询任务信息

> 获得任务的相关信息，包括所属的人脸库名、配置信息、当前状态、该任务所需处理的总数、该任务当前已处理的数目

**Request**

```
GET /v1/task/<task_id>/detail
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
    },
    "total_num" : "",
    "handled_num": ""
    "status": ""
    "last_error": ""
}
```

接口字段说明：

| 字段               | 取值   | 说明                     |
| :----------------  | :----- | :----------------------- |
| task_id            | string | 任务id |
| group_name         | string | 任务所属的人脸库group名|
| config.pos_pitch   | int    | 配置的允许的人脸最大俯仰角 |
| config.pos_yaw     | int    | 配置的允许的人脸最大偏航角 |
| config.pos_rol     | int    | 配置的允许的人脸最大翻滚角 |
| config.width       | int    | 配置的允许的人脸最小宽度 |
| config.height      | int    | 配置的允许的人脸最小高度 |
| config.mode        | int    | 配置的入库模式 |
| total_count        | int    | 该任务一共需要处理的人脸数 |
| handled_count      | int    | 该任务当前已处理完的人脸数 |
| status             | int    | 该任务当前状态，1:已创建未执行，2:在队列中等待执行，3:正在执行中，4:正在停止中，5:已停止，6:已完成 |
| error              | string | 当系统停止任务时（比如任务初始化错误、入库模式为出错即停止时出错，该字段返回造成停止的错误原因 |



## 1.7 查询任务日志

> 每当入库人脸出现错误时，系统会将错误记录到日志
> 可随时查询任务日志，无论任务是否已结束

**Request**

```
GET /v1/task/<task_id>/log
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
| task_id    | string | 任务id   |
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


## 1.8 下载任务日志

> 除了可以查询日志外，系统还提供下载日志，下载文件为csv文件。格式为：人脸图片uri，错误码，错误信息

**Request**

```
GET /v1/task/<task_id>/download_log
```

**Response**

```
200 OK
Content-Type: application/octet-stream
***文件内容***
```

接口字段说明：

| 字段    | 取值    | 说明 |
| :------ | :----- | :----- |
| task_id | string | 任务id |


# 2. 人脸入库工具

## 2.1 依赖

> 该工具依赖feature_group_private的api接口，具体接口为    
1. 创建group接口: /v1/face/groups/[group_name]    
2. 人脸入库接口: /v1/face/groups/[group_name]/add    
3. 人脸姿态接口: /v1/eval/facex-pose

## 2.2 功能说明

> 读取指定文件夹下所有图片（包括子文件夹下图片）或指定文件中的url列表，调用接口入库，同时实现：

1. 去除重复文件，即重复文件不上传
2. 支持中断续传
3. 列出错误文件及错误原因（ex：非图片文件，图片损坏，图片不包含人脸，etc）

## 2.3 实现细节

> 多线程入库，且处理多线程下的中断续传（场景为海量文件，非大容量文件，重点为多线程模式下记录上次中断点）

1. 按filepath.walk()默认的方式获取文件夹内所有图片，按文件名字典排序。或逐行读取csv文件内的url信息
2. 将每个图片包装成job，扔给任务池，多线程执行
3. 每个线程记录当前处理图片的index
4. 如中断，下次执行时首先获得所有线程记录的index中最小值，并直接从该index开始续传
5. 每次处理完图片，线程会将该图片sha1值记录到本地文件，便于去重（主要是中断情况下的去重）。每次程序开始也会读取该文件获得已处理图片的sha1（如果有）
6. 入库的图片来源有两个：指定文件夹里的所有图片或指定url列表csv文件中的所有图片url，默认的图片文件夹是/workspace/source/，默认的url列表csv文件是/workspace/urlSource


## 2.4 字段录入

> 调用人脸入库接口时，传入id, uri, tag, desc四个字段

1. 图片来源为文件夹时：
    | 字段   | 说明                 |
    | :---- | :--------            |
    | id    | 图片在系统中的路径      |
    | uri   | 图片base64值          |
    | tag   | 图片文件名去除扩展名后的第一个‘_‘字符之前的值。如’abc_123.jpg‘,则为abc；如’abcd.jpg‘,则为abcd；如’abcd_123_56‘,则为abcd   |
    | desc  | 图片文件名去除扩展名后的第一个‘_‘字符之后的值。如’abc_123.jpg‘,则为123；如’abcd.jpg‘,则为空；如’abcd_123_56‘,则为123_56   |

2. 图片来源为csv文件时，默认为每行格式 url，tag，desc。例：http://somesite.com/test.jpg,faceTag,faceDescription
    | 字段   | 说明                     |
    | :---- | :--------                |
    | id    | csv文件中的第一列值，即url  |
    | uri   | url指向的图片的base64值    |
    | tag   | csv文件中的第二列值，如果有  |
    | desc  | csv文件中的第三列值，如果有  |

## 2.5 使用说明

> 交付docker文件

1.  docker镜像为dbstorage
2.  图片来源为文件夹时，运行
    docker run --name dbstorage_tool -v [宿主机上dbstorage_tool.conf文件路径]:/workspace/dbstorage_tool.conf -v [宿主机上图片源文件夹路径]:/workspace/source/ -v [宿主机上运行结果目录]:/workspace/log/ dbstorage_tool
3.  图片来源为url列表文件时，运行
    docker run --name dbstorage_tool -v [宿主机上dbstorage_tool.conf文件路径]:/workspace/dbstorage_tool.conf -v [宿主机上url列表文件路径]:/workspace/urlSource -v [宿主机上运行结果目录]:/workspace/log/ dbstorage_tool
4.  程序支持中断续传，若终端，可执行如下指令续传
    docker start -i dbstorage_tool
5.  出错的文件都会输出到文件： [宿主机上运行结果目录]/error
6.  当前已处理完的图片数目会实时更新到文件：[宿主机上运行结果目录]/count
7.  dbstorage_tool.conf配置文件字段说明：

    | 字段                           | 类型      | 说明                                                                              |
    | :---------------------------- | :-------- | :--------------------------------------------------------------------------------|
    | load_image_from_folder        | boolean   | 导入图片来源有两个：true为从文件夹里导入，false为从url列表文件导入，需配合docker run指令使用 |
    | max_try_service_time          | int       | 对于每个图片，调用入库接口时的尝试数 |
    | max_try_download_time         | int       | 对于每个图片，从图片源url下载图片是时的总尝试数 |
    | task_config.thread_num        | int       | 该任务起用的线程数，可选，默认为20 |
    | feature_group_service         | object    | 入库服务 |  
    | feature_group_service.host    | string    | 入库服务地址 |  
    | feature_group_service.timeout | int       | 入库服务超时时间，单位为秒，0为不超时 |  
    | face_service                  | object    | 人脸服务 |  
    | face_service.host             | string    | 人脸服务地址 |  
    | face_service.timeout          | int       | 人脸服务超时时间，单位为秒，0为不超时 |  
    | group_name                    | string    | 入库的人脸库名称 |
    | task_config                   | object    | 图片入库的配置 |
    | task_config.pos_pitch         | int       | 允许的人脸最大俯仰角，超过则该人脸不会入库，可选，取值为[1, 90], 默认为空（即不限制） |
    | task_config.pos_yaw           | int       | 允许的人脸最大偏航角，超过该值则该人脸不会入库，可选，取值为[1, 90], 默认为空（即不限制） |
    | task_config.pos_rol           | int       | 允许的人脸最大翻滚角，超过该值则该人脸不会入库，可选，取值为[1, 90], 默认为空（即不限制） |
    | task_config.width             | int       | 允许的人脸最小宽度，低于该值则该人脸不会入库，可选，取值为大于0的整数，默认为空（即不限制） |
    | task_config.height            | int       | 允许的人脸最小高度，低于该值则该人脸不会入库，可选，取值为大于0的整数，默认为空（即不限制） |
    | task_config.mode        | int       | 入库时遇见错误图片时的处理方式，0为跳过继续执行，1为将任务停止，可选，默认为0 |

