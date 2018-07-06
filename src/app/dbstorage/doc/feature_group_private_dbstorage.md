<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [入库工具说明](#%E5%85%A5%E5%BA%93%E5%B7%A5%E5%85%B7%E8%AF%B4%E6%98%8E)
  - [依赖](#%E4%BE%9D%E8%B5%96)
  - [功能说明](#%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E)
  - [实现细节](#%E5%AE%9E%E7%8E%B0%E7%BB%86%E8%8A%82)
  - [使用说明](#%E4%BD%BF%E7%94%A8%E8%AF%B4%E6%98%8E)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 入库工具说明

## 依赖

> 该工具依赖feature_group_private的api接口，具体接口为

1. 创建group接口: /v1/face/groups/[group_name]
2. 人脸入库接口: /v1/face/groups/[group_name]/add

## 功能说明

> 读取指定文件夹下所有图片（包括子文件夹下图片）或指定文件中的url列表，调用接口入库，同时实现：

1. 去除重复文件，即重复文件不上传
2. 支持中断续传
3. 列出错误文件及错误原因（ex：非图片文件，图片损坏，图片不包含人脸，etc）

## 实现细节

> 多线程入库，且处理多线程下的中断续传（场景为海量文件，非大容量文件，重点为多线程模式下记录上次中断点）

1. 按filepath.walk()默认的方式获取文件夹内所有图片，按文件名字典排序。或逐行读取csv文件内的url信息
2. 将每个图片包装成job，扔给任务池，多线程执行
3. 每个线程记录当前处理图片的index
4. 如中断，下次执行时首先获得所有线程记录的index中最小值，并直接从该index开始续传
5. 每次处理完图片，线程会将该图片sha1值记录到本地文件，便于去重（主要是中断情况下的去重）。每次程序开始也会读取该文件获得已处理图片的sha1（如果有）
6. 入库的图片来源有两个：指定文件夹里的所有图片或指定url列表csv文件中的所有图片url，默认的图片文件夹是/workspace/source/，默认的url列表csv文件是/workspace/urlSource

## 字段录入

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

## 使用说明

> 交付docker文件

1.  docker镜像为dbstorage
2.  图片来源为文件夹时，运行
    docker run --name dbstorage -v [宿主机上dbstorage.conf文件路径]:/workspace/dbstorage.conf -v [宿主机上图片源文件夹路径]:/workspace/source/ -v [宿主机上运行结果目录]:/workspace/log/ dbstorage
3.  图片来源为url列表文件时，运行
    docker run --name dbstorage -v [宿主机上dbstorage.conf文件路径]:/workspace/dbstorage.conf -v [宿主机上url列表文件路径]:/workspace/urlSource -v [宿主机上运行结果目录]:/workspace/log/ dbstorage
4.  程序支持中断续传，若终端，可执行如下指令续传
    docker start -i dbstorage
5.  出错的文件都会输出到文件： [宿主机上运行结果目录]/errorlist
6.  当前已处理完的图片数目会实时更新到文件：[宿主机上运行结果目录]/process
7.  dbstorage.conf配置文件字段说明：

    | 字段                           | 类型      | 说明                                                                                |
    | :---------------------------- | :-------- | :--------------------------------------------------------------------------------- |
    | load_image_from_folder        | boolean   | 导入图片来源有两个：true为从文件夹里导入，false为从url列表文件导入，需配合docker run指令使用   |
    | service_host_url              | string    | 七牛服务host地址                                                                     |
    | http_timeout_in_millisecond   | int       | 访问七牛服务的客户端等待时间，单位毫秒，0为客户端无超时                                     |
    | max_try_service_time          | int       | 对于每个图片，调用接口入库图片时的总尝试数                                                |
    | max_try_download_time         | int       | 对于每个图片，从图片源url下载图片是时的总尝试数                                            |
    | thread_number                 | int       | 线程池中线程数目                                                                      |
    | job_pool_size                 | int       | 工作池大小                                                                           |
    | group_name                    | string    | 入库接口的group_name参数                                                              |

