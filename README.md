# goice

元旦在家无聊写了个demo,一来是练练手,也抱着想搞一点好玩的东西。

目前Demo 都是一些简单的实现 可能还很不完善 包含:
- 简单的依赖注入,只要是注册的服务都可以自动注入
- Redis 的 Orm 现在只完成了 保存 和 获取同时由于Golang是静态语言,在反射的时候填充结构体时 必须要断言,但是 既然是框架 那不能写死是哪个结构体,现在的做法是 reflect_type
 事先把断言写好 然后返回,目前可以解决这个问题.但每次写这个文件比较麻烦而且机械,后期可以写个自动生成的工具,但可能由更好的办法 希望可以改进。
 
 
 ### 2021-01-03 18:27
 - 下面要把Redis的批量查找和 删除 修改写完
 - 后面 考虑 把 所有的查询都可以缓存起来,然后按照策略可以淘汰什么的。
 