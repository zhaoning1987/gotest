入库工具说明

功能说明:
1: 读取指定文件夹下所有图片（包括子文件夹下图片），调用 /v1/face/groups/<*>/add 接口获得特征值并入库
2: 去除重复文件，即重复文件不上传
3: 支持中断续传
4: 列出错误文件及错误原因（ex：非图片文件，图片损坏，图片不包含人脸，etc）

实现细节:
1: 按filepath.walk()默认的方式获取所有图片，按文件名字典排序
2: 将每个图片包装成job，扔给任务池，多线程执行
3: 每个线程记录当前处理图片的index
4: 如中断，下次执行时首先获得所有线程记录的index中最小值，并直接从该index开始续传
5: 每次处理完图片，线程会将该图片md5值记录到本地文件，便于去重（主要是中断情况下的去重）。每次程序开始也会读取该文件获得已处理图片的md5（如果有）

使用说明：
1: 程序位于 ～/ava/platform/src/qiniu.com/ava/argus/feature_group_private/dbstorage 下
2: 程序唯一入参：配置文件地址。例：执行 ./dbstorage -f dbstorage.conf
3: 配置文件说明
    {
        "log_path": "./log/",                                       //程序执行过程中log存储地址，包含中断信息，错误图片信息，已处理图片的md5信息等
        "image_dir_path": "/Users/zhaoning/Documents/testImage/",   //图片源文件夹地址，内包含待处理的图片文件（图片源方式之一）
        "image_list_file": "/Users/zhaoning/Desktop/imageList",     //图片源文件地址，包含待处理的所有文件url（图片源方式之二）
        "use_image_dir_path": true,                                 //图片来源：true为使用参数image_dir_path，false为使用参数image_list_file
        "qiniu_host_url": "http://100.100.58.85:6125",              //七牛服务host地址
        "http_timeout_in_millisecond": 0,                           //访问七牛服务的客户端等待时间，单位毫秒，0为客户端无超时
        "max_try_service_time": 1,                                  //对于每个图片，调用接口入库图片时的总尝试数
        "max_try_download_time": 1,                                 //对于每个图片，从图片源url下载图片是时的总尝试数
        "thread_number": 20,                                        //线程池中线程数目
        "job_pool_size": 100,                                       //工作池大小
        "group_name": "group_name"                                  //入库接口的group_name参数
    }