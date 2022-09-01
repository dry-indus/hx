[![pipeline status](https://gitlab.com/lzw5399/hx/badges/master/pipeline.svg)](https://gitlab.com/lzw5399/hx/-/commits/master)

# 摸索一下基于gin开发的姿势

## 目标

是在基于符合go语言基本规范的基础上：
- 参考前辈们的项目结构&组件选择
- 在个人开发习惯上作倾斜


## 目前参考的项目

- https://github.com/LyricTian/gin-admin
- https://github.com/flipped-aurora/gin-vue-admin

## 项目结构

```
├─hx 
   ├─controller     （api方法层）
   ├─config         （配置文件&配置结构体存放的地方）
   ├─docs  	    （swagger文档目录）
   ├─global         （全局对象，里面的对象会在initialize里面初始化）
   ├─initialize     （初始化）
   ├─middleware     （中间件）
   ├─model          （模型）
   ├─router         （路由）
   ├─service         (业务逻辑层)
   ├─util	    （公共功能）

...可以添加的
   ├─db             （数据库脚本）
   ├─resource       （静态资源）
   ├─core           （如果要自定义http的话可以抽取到这部分）
   ├─view           （如果需要在项目中包含前端页面，可以添加这个）
```

## 目前使用到的组件

- Web框架
   > github.com/gin-gonic/gin
- ORM
   > github.com/jinzhu/gorm
- 配置加载
   > github.com/jinzhu/configor
- API文档
   > github.com/swaggo/gin-swagger
- CORS(gin中间件)
   > github.com/gin-contrib/cors
- 认证方式Jwt
   > github.com/dgrijalva/jwt-go
- 授权
   > github.com/casbin/casbin
- Log
   > github.com/op/go-logging
- Mapper
   > github.com/jinzhu/copier

## 其他没用到但是可选的组件

- 配置加载
   > github.com/spf13/viper
- Log
   > github.com/sirupsen/logrus