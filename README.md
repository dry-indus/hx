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
- API文档
   > github.com/swaggo/gin-swagger
- CORS(gin中间件)
   > github.com/gin-contrib/cors
- Log
   > github.com/op/go-logging

## 构建文档
Since swag 1.7.9 we are allowing registration of multiple endpoints into the same server.

Generate documentation for merchant/v1 endpoints
```shell
sudo swag i -g api/merchant/v1/api.go --exclude ./controller/userctr  --instanceName mv1
```

Generate documentation for user/v1 endpoints
```shell
sudo swag i -g api/user/v1/api.go --exclude ./controller/merchantctr  --instanceName uv1
```

Run example
```shell
    go run main.go
```

merchant v1 swagger here [http://localhost:7777/swagger/mv1/index.html](http://localhost:7777/swagger/mv1/index.html)
user v1 swagger here [http://localhost:7777/swagger/uv1/index.html](http://localhost:7777/swagger/uv1/index.html)

## 部署Drone
```shell
docker run --volume=/data/drone/data:/data --env=DRONE_GITHUB_CLIENT_ID=56f626d61deb34fdb3ed --env=DRONE_GITHUB_CLIENT_SECRET=81ab9e02986d58cf6657f713236c8286196ba852 --env=DRONE_RPC_SECRET=b62a5214790682873063d6176c1e2004 --env=DRONE_SERVER_HOST=drone.mik888.com --env=DRONE_SERVER_PROTO=http --publish=3080:80 --env=DRONE_USER_CREATE=username:root,admin:true --restart=always --detach=true --name=drone drone/drone:2
docker run -d -v /var/run/docker.sock:/var/run/docker.sock -e DRONE_RPC_PROTO=http -e DRONE_RPC_HOST=drone.mik888.com:3080 -e DRONE_RPC_SECRET=b62a5214790682873063d6176c1e2004 -e DRONE_RUNNER_CAPACITY=2 -e DRONE_RUNNER_NAME=runner-docker -e TZ="Asia/Shanghai" --publish=3000:3000  --restart always --name drone-runner drone/drone-runner-docker:1
```