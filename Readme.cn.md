# 简述

这是一个自动分配子网络(Subnet)的工具，这个工具是用来配合Kubernetes来完成网络设置的。考虑到越来越多的团队和公司开始基于`直接路由`的方式来完成Kubernetes最底层网络的设置工作，所以我将此工具开放出来。

采用`直接路由`的网络模式，在网络间传输的开销是非常小的，但这也会带来一定的人工成本。人工成本的点在于需要为每一个Kubernetes Minion Node来可用的子网络IP段。这也是Flannel等工具提供的好处之一。在这个工具的初期版本中，我们会尝试假设每一个K8S Minion Node都可以最多支持执行110个POD的运行，这也就意味着我们会为每个K8S Minion Node分配一个"/24"的子网段。

# 使用方式

此工具启动成功后，会开放一个TCP的8080端口用于支持Web API的访问。使用者只需要将K8S Minion Node的IP，和当前这个K8S Minion Node所在的业务环境(生产、测试)编号发给Web API, 那么Web API就会返回一个分配到此Minion Node IP上的子网段。需要明确说明的是，如果将相同的Minion Node IP和业务环境编号多次发给Web API, 那么Web API会根据组合条件来匹配，并返回相同的结果(已分配的IP网段)。

## 启动方式

第一次启动之前，请使用本仓储中的`mysql_database.sql`文件来初始化MYSQL数据库。

如果相关数据库表已经初始化完成，并且是第一次启动的话，则需要使用以下指令来初始化数据库表中的待分配IP段数据:

./tinydhcp-dockerip -i 这里填写待分配的网段(比如"192.168.1.1/8") -e 这里填写初始化IP段资源时，这些IP段资源所归属的业务环境ID -mysql "MYSQL用户名:MYSQL密码@tcp(MYSQL服务器IP地址:3306)/tiny_dhcp" -n true

如果是后续启动的话，则只需要按照如下方式来启动即可:

./tinydhcp-dockerip -i 这里填写待分配的网段(比如"192.168.1.1/8") -e 这里填写初始化IP段资源时，这些IP段资源所归属的业务环境ID -mysql "MYSQL用户名:MYSQL密码@tcp(MYSQL服务器IP地址:3306)/tiny_dhcp"

这里需要明确指出的是，无论-i参数后面给出的待分配网段是多大，最终都会被转化为若干个"/24"的子网段。

## HTTP 请求

```http
GET /ip?node-ip={NODE-IP}&env-id={ENV-ID}&owner={OWNER}&desc={DESC}
```

## HTTP 应答

```json
{
"error-id": 0,
"docker-ip": "192.202.10.1/24",
"reason": ""
}
```

## 请求参数

|参数名称|参数用途|参数类型|是否必选|默认值|
|---|---|---|---|---|
|node-ip|Kubernetes Minion Node IP|string|Y|N/A|
|env-id|用于区分业务区域，比如生产环境或者测试环境等等。业务区域的ID需要事先在数据库中初始化好|int16|Y|N/A|
|owner|当前K8S Minion Node主机资源拥有者，只是一个名称标示|string|N|""|
|desc|一些描述信息|string|N|""|
