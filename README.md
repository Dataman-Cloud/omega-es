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
  vim omega-es.yaml.sample 
  ```
  ```
  docker run -d omage-es:latest
  ``` 
  or
  ```
  docker run -d omage-es:latest
  ```
  ```
  启动成功后 omega-es 默认监听5009 端口
  ```
  ```
  以上步骤不保证能跑起来
  ```
