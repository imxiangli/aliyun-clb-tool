### 用法：
#### 1、下载对应目标机器的二进制程序文件，或使用源码自行编译目标机器的
#### 2、修改 config/app.yaml 其中的秘钥配置
#### 3、示例：复制转发规则
```shell
# 示例：复制转发规则
./copyRules -sourceRegionId cn-shenzhen -sourceLoadBalancerId lb-wz9***lsn -sourceListener https:443 -targetRegionId cn-shenzhen -targetLoadBalancerId lb-wz9***lsn -targetListener http:8080
```