# omega-es
omega elasticsearch service

### Building 


cd $GOPATH/src/github.com/Dataman-Cloud
git clone git@github.com:Dataman-Cloud/omega-es.git
cd omega-es
go build

### Run Standalone
  ```
  git clone git@github.com:Dataman-Cloud/omega-es.git
  ```
  ```
  cd omega-es/
  ```
  ```
  docker build -t omage-es .
  ```
  ```
  vim omega-es.yaml.sample  把该配置的都配置上 
  ```
  ```
  如果你上面配置中没有`127.0.0.1`的地址 那么用如下方式启动:

  docker run -d omage-es:latest
  ``` 
  ```
  如果你上面配置中有`127.0.0.1`的地址 那么需要以host模式启动: 

  docker run -d --net host omage-es:latest
  ```
  ```
  启动成功后 omega-es 默认监听5009 端口
  ```
  ```
  以上步骤不保证能跑起来 如果搞不定 请联系郭已钦
  ```

