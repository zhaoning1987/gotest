
# 基本参数

## 输入图片格式
支持JPG、PNG、BMP、GIF

## 资源表示方式（URI）

通过统一方式定位、获取资源（图片、二进制数据等）

* HTTP，网络资源，形如：http://host/path、https://host/path
* FILE，本地文件，形如：file://path
* Data，Data URI Scheme形态的二进制文件，形如：data:application/octet-stream;base64,xxx。ps: 当前只支持前缀为data:application/octet-stream;base64,的数据

## 错误返回

| 错误码 | 描述 |
| :--- | :--- |
| 4000100 | 请求参数错误 |
| 4000201 | 图片地址格式不支持 |
| 4000202 | 图片地址有误 |
| 4240203 | 获取图片失败 |
| 4240204 | 获取图片超时 |
| 4150301 | 图片格式不支持 |
| 4000302 | 图片过大，长宽大于4999像素、或图片大小超过10M |

# API列表

## /v1/pulp (1.0.0)

> 用第三方服务和AtLab的剑皇服务做融合剑皇<br>
> 每次输入一张图片，返回其内容是否含色情信息<br>

*Request*

```
POST /v1/pulp  Http/1.1
Content-Type:application/json

{
	"data": 
		{
			"uri": "http://oayjpradp.bkt.clouddn.com/Audrey_Hepburn.jpg"
		}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"code": ,
	"message": "",
	"result": {
		"label":1,
		"score":0.987,
		"review":false
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |
| params.detail | bool | 是否显示详细信息；可选参数 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| code | int | 0:表示正确 |
| message | string | 结果描述信息 |
| result.label | int | 标签{0:色情，1:性感，2:正常} |
| result.score | float | 色情识别准确度 |
| result.review | bool | 是否需要人工review |
| result.confidences | list | 图片打标信息列表 |
| result.confidences.index | int | 类别编号, 即 0:pulp,1:sexy,2:norm |
| result.confidences.class | string | 图片内容鉴别结果，分为色情、性感或正常3类 |
| result.confidences.score | float32 | 将图片判别为某一类的准确度，取值范围0~1，1为准确度最高 |

## /v1/censor/pulp (1.0.0)

> 用第三方服务和AtLab的剑皇服务做融合剑皇判断是否违规<br>
> 每次输入一张图片，返回其内容是否违规<br>

*Request*

```
POST /v1/pulp/censor  Http/1.1
Content-Type:application/json

{
	"data": 
		{
			"uri": "http://oayjpradp.bkt.clouddn.com/Audrey_Hepburn.jpg"
		}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"suggestion": "pass",
	"result": {
		"label":"normal",
		"score":0.772529,
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| suggestion | string | pass/block/review |
| result.label | string | 标签{normal，pulp，sexy} |
| result.score | float | 色情识别准确度 |

## /v1/terror (1.0.0)

> 用检测暴恐识别和分类暴恐识别方法做融合暴恐识别<br>
> 每次输入一张图片，返回其内容是否含暴恐信息<br>

*Request*

```
POST /v1/terror  Http/1.1
Content-Type: application/json

{
	"data": {
		"uri": "http://oayjpradp.bkt.clouddn.com/Audrey_Hepburn.jpg"
	},
	"params": {
		"detail": true
	}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"code": 0,
	"message": "",
	"result": {
		"label":1,
		"score":0.987,
		"review":false
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |
| params.detail | bool | 是否显示详细信息；可选参数 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| code | int | 0:表示正确 |
| message | string | 结果描述信息 |
| result.label | int | 标签{0:正常，1:暴恐} |
| result.class | string | 标签（指定detail=true的情况下返回 |
| result.score | float | 暴恐识别准确度 |
| result.review | bool | 是否需要人工review |

## /v1/censor/terror (1.0.0)

> 用检测暴恐识别和分类暴恐识别方法做融合暴恐识别判断是否违规<br>
> 每次输入一张图片，返回其内容是否违规<br>

*Request*

```
POST /v1/terror/censor  Http/1.1
Content-Type: application/json

{
	"data": {
		"uri": "http://oayjpradp.bkt.clouddn.com/Audrey_Hepburn.jpg"
	}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"suggestion": "pass",
	"result": {
		"label":"normal",
		"score":0.9999995,
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| suggestion | string | pass/block/review |
| result.label | string | 标签 |
| result.score | float | 暴恐识别准确度 |

## /v1/face/search/politician (1.0.0)

> 政治人物搜索，对输入图片识别检索是否存在政治人物<br>

*Request*

```
POST /v1/face/search/politician Http/1.1
Content-Type: application/json

{
	"data": {
		"uri": "http://xx.com/xxx"
	}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"code": 0,
	"message": "",
	"result": {
		"review": True,
		"detections": [
			{
				"boundingBox":{
					"pts": [[1213,400],[205,400],[205,535],[1213,535]],
					"score":0.998
				},
				"value": {
					"name": "xx",
					"group": "Inferior Artist",
					"score":0.567,
					"review": True
				},
				"sample": {
					"url": "",
					"pts": [[1213,400],[205,400],[205,535],[1213,535]]
				}
			},
			{
				"boundingBox":{
					"pts": [[1109,500],[205,500],[205,535],[1109,535]],
					"score":0.98
				},
				"value": {
					"score":0.987,
					"review": False
				}
			}
		]
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| code | int | 0:表示正确 |
| message | string | 结果描述信息 |
| review | boolean | True或False,图片是否需要人工review, 只要有一个value.review为True,则此字段为True |
| boundingBox | map | 人脸边框信息 |
| boundingBox.pst | list[4] | 人脸边框在图片中的位置[左上，右上，右下，左下] |
| boundingBox.score | float | 人脸位置检测准确度 |
| value.name | string | 检索得到的政治人物姓名,value.score < 0.4时未找到相似人物,没有这个字段 |
| value.group | string | 人物分组信息，总共有7个组{'Domestic politician','Foreign politician','Sacked Officer(Government)','Sacked Officer (Enterprise)''Anti-China Molecule','Terrorist','Inferior Artist'} |
| value.review | boolean | True或False,当前人脸识别结果是否需要人工review |
| value.score | float | 0~1,检索结果的可信度, 0.35 <= value.score <=0.45 时 value.review 为True |
| sample | object | 该政治人物的示例图片信息，value.score < 0.4时未找到相似人物, 没有这个字段 |
| sample.url | string | 该政治人物的示例图片 |
| sample.pts | list[4] | 人脸在示例图片中的边框 |

## /v1/censor/politician (1.0.0)

> 政治人物搜索，对输入图片识别检索是否存在政治人物判断是否违规<br>

*Request*

```
POST /v1/face/search/politician/censor Http/1.1
Content-Type: application/json

{
	"data": {
		"uri": "http://xx.com/xxx"
	}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{	
	"suggestion": "pass",
	"result": {
		"label":"face",
		"faces": [
			{
				"boundingBox":{
					"pts": [[452,226],[1065,226],[1065,1159],[452,1159]],
					"score":0.9999962
				},
				"faces": [
					{
					"id": "",
					"name": "xx",
					"score":0.567,
					"group": "Inferior Artist"
				
				"sample": {
					"url": "",
					"pts": [[1213,400],[205,400],[205,535],[1213,535]]
						}
					},
				]
			},
		]
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| suggestion | string | pass/block/review |
| result.label | string | 标签 |
| result.faces | string | 人脸信息 |
| boundingBox | map | 人脸边框信息 |
| boundingBox.pst | list[4] | 人脸边框在图片中的位置[左上，右上，右下，左下] |
| boundingBox.score | float | 人脸位置检测准确度 |
| faces.faces | list | 检索得到的政治人物的信息 |
| faces.score | float | 0~1,检索结果的可信度, 0.35 <= value.score <=0.45 时 value.review 为True |
| faces.name | string | 检索得到的政治人物姓名,face.score < 0.4时未找到相似人物,没有这个字段 |
| faces.group | string | 人物分组信息，总共有7个组{'Domestic politician','Foreign politician','Sacked Officer(Government)','Sacked Officer (Enterprise)''Anti-China Molecule','Terrorist','Inferior Artist'} |
| faces.sample.url | string | 该政治人物的示例图片 |
| faces.sample.pts | list[4] | 人脸在示例图片中的边框 |

## /v1/face/sim (1.0.0)

> 人脸相似性检测<br>
> 若一张图片中有多个脸则选择最大的脸<br>

*Request*

```
POST /v1/face/sim  Http/1.1
Content-Type:application/json

{
	"data": [
		{
			"uri": "http://image2.jpeg"
		},
		{
			"uri": "http://image1.jpeg
		}
	]
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"code": 0,
	"message": "success",
	"result": {
		"faces":[{
				"score": 0.987,
				"pts": [[225,195], [351,195], [351,389], [225,389]]
			},
			{
				"score": 0.997,
				"pts": [[225,195], [351,195], [351,389], [225,389]]
			}], 
		"similarity": 0.87,
		"same": 0  
	}	
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| data | list | 两个图片资源地址 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| code | int | 0:表示处理成功；不为0:表示出错 |
| message | string | 描述结果或出错信息 |
| faces | list | 两张图片中选择出来进行对比的脸 |
| score | float | 人脸识别的准确度，取值范围0~1，1为准确度最高 |
| pts | list | 人脸在图片上的坐标 |
| similarity | float | 人脸相似度，取值范围0~1，1为准确度最高 |
| same | bool | 是否为同一个人 |

## /v1/face/detect (1.0.0)

> 检测人脸所在位置<br>
> 每次输入一张图片，返回所有检测到的脸的位置<br>

*Request*

```
POST /v1/face/detect  Http/1.1
Content-Type:application/json

{
	"data": {
			"uri": "http://image1.jpeg"
	}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"code": 0,
	"message": "",
	"result": {
        "detections": [
            {
                "bounding_box": {
                    "pts": [[268,212], [354,212], [354,320], [268,320]],
                    "score": 0.9998436
                }
            },
            {
                "bounding_box": {
                    "pts": [[159,309], [235,309], [235,408], [159,408]],
                    "score": 0.9997162
                }
            }
        ]
	}	
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| code | int | 0:表示处理成功；不为0:表示出错 |
| message | string | 描述结果或出错信息 |
| detections | list | 检测出的人脸列表 |
| bounding_box | map | 人脸坐标信息 |
| pts | list | 人脸在图片中的位置，四点坐标值 [左上，右上，右下，左下] 四点坐标框定的脸部 |
| score | float | 人脸的检测置信度，取值范围0~1，1为准确度最高 |

## /v1/image/censor (1.0.0)

> 图片审核<br>

*Request*

```
POST /v1/image/censor Http/1.1
Content-Type: application/json

{
	"data": {
		"uri": "http://oayjpradp.bkt.clouddn.com/Audrey_Hepburn.jpg"
	},
	"params": {
		"type": [
			"pulp",
			"terror",
			"politician",
			"terror-complex"
		],
		"detail": true
	}
}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"code": 0,
	"message": "",
	"result": {
		"label": 1,
		"score": 0.888,
		"details": [
			{
				"type": "pulp",
				"label": 1,
				"score": 0.8726697,
				"review": false
			},
			{
				"type": "terror",
				"label": 1,
				"class": <class>,
				"score": 0.6530496,
				"review": false
			},
			{
				"type": "politician",
				"label": 1,
				"score": 0.77954,
				"review": True,
				"more": [
					{
						"boundingBox":{
							"pts": [[1213,400],[205,400],[205,535],[1213,535]],
							"score":0.998
						},
						"value": {
							"name": "xx",
							"score":0.567,
							"review": True
						},
						"sample": {
							"url": "",
							"pts": [[1213,400],[205,400],[205,535],[1213,535]]
						}
					},
					{
						"boundingBox":{
							"pts": [[1109,500],[205,500],[205,535],[1109,535]],
							"score":0.98
						},
						"value": {
							"score":0.987,
							"review": False
						}
					}
				]
			},
			{
				"type": "terror-complex",
				"label": 1,
				"classes": <classes>,
				"score": 0.6530496,
				"review": false
			},
		]
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| uri | string | 图片资源地址 |
| params.type | string | 选择的审核类型，可选项：'pulp'/'terror'/'politician'/'terror-complex'；可选参数，不填表示全部执行 |
| params.detail | bool | 是否显示详细信息；可选参数 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| code | int | 0:表示正确 |
| message | string | 结果描述信息 |
| result.label | int | 是否违规，0：不违规；1：违规 |
| result.score | float | 是否违规置信度 |
| result.detail.type | string | 审核类型 |
| result.detail.label | int | 审核结果类别，具体看各类型 |
| result.detail.class | string | 详细类别，具体看各类型 |
| result.detail.classes | []string | 详细类别列表，具体看各类型 |
| result.detail.score | float | 审核结果置信度 |

## /v1/censor (1.0.0)

> 图片审核<br>

*Request*

```
POST /v1/censor/image/recognition Http/1.1
Content-Type: application/json

{
	    "datas":[{"uri":"http://odqrp5nr9.bkt.clouddn.com//cluster/1/1.jpg"}],
		"scenes": [
			"pulp",
			"terror",
			"politician",
			"terror-complex"
		]

}
```

*Response*

```
 200 ok
Content-Type:application/json

{
	"tasks":[{
	"code":200,
	"message":"OK",
	"suggestion":"pass",
	"scenes":{
		"politician":{
			"suggestion":"pass",
			"result":{
				"label":"face",
				"faces":[
					{
						"bounding_box":{
							"pts":[[452,226],[1065,226],[1065,1159],[452,1159]],
							"score":0.9999962
						}
					}
				]
			}
		},
		"pulp":{
					"suggestion":"pass",
					"result":{"label":"normal","score":0.999498}
				},
		"terror":{
						"suggestion":"pass",
						"result":{"label":"normal","score":0.9998622}
				}
			}
		}]
	}
}
```
*Request Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| datas | list | 多个图片资源地址和id |
| scenes  | list | 审核场景，可填入：'pulp'/'terror'/'politician'/'terror-complex'；可选参数，不填表示全部执行 |
| params.Scenes | map | 场景对应的配置信息,不填表示使用默认参数 |

*Response Params*

| 参数 | 类型 | 描述 |
| :--- | :---- | :--- |
| tasks | list | 每张图片的处理结果列表 |
| tasks.code | int | 200表示处理成功，否则返回错误码 |
| tasks.message | string | 结果描述信息 |
| tasks.suggestion | string | pass/block/review |
| tasks.scenes | map | 返回各个场景的审核结果 |

