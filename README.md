# Schedule_Command_module

* API Server

  * time template API
  * command API
    * header template API
  * header template API
  * header template API
  * schedule API
* Time Server

  1. 每秒監控一次schedule
  2. 時間符合則執行相對應的命令

## 啟動方式

#### 1. 使用docker-compose

沒有mysql container可以選擇此方式

1. 確定docker 有mysql image

   `docker pull mysql`
2. 修改docker-compose.yaml
   **Windows:** line:24 註解linux命令 取消line:25的註解
   **Linux:** 取消line:24註解 line:25註解windows命令
   line:24:

   ```
   dockerfile: deploy/api/linux/Dockerfile
   ```

   line25:

   ```
   dockerfile: deploy/api/windows/Dockerfile
   ```
3. docker-compose up

   `docker-compose up`

#### 2. 使用DockerFile啟動

有固定的mysql container可選擇此方式

1. ##### Windows系統

   1. 到mysql創建一個新的database "schedule"
      `create database schedule;`
   2. 創建需要的image
      api image:
      `docker build -t schedule:latest -f deploy/api/windows/Dockerfile .`
      migrate image:

      `docker build -t schedule-migrate:latest -f deploy/migrate/windows/Dockerfile .`
   3. 創建並啟動container

      1. run migrate container
         DB_HOST可指定特定的DB IP

         `docker run --name schedule-migrate --rm -e DB_HOST=host.docker.internal schedule-migrate:latest`
      2. run api container

         `docker run --name schedule -p 5487:5487 -e DB_HOST=host.docker.internal -v ${PWD}/deploy/api/dockerLog:/app/log schedule:latest`
2. ##### Linux系統

   1. 到mysql創建一個新的database "schedule"
      `create database schedule;`
   2. 創建需要的image
      api image:
      `docker build -t schedule:latest -f deploy/api/linux/Dockerfile .`
      migrate image:

      `docker build -t schedule-migrate:latest -f deploy/migrate/linux/Dockerfile .`
   3. 創建並啟動container

      1. run migrate container
         DB_HOST可指定特定的DB IP

         `docker run --name schedule-migrate --rm --network="host" schedule-migrate:latest`
      2. run api container

         `docker run --name schedule -p 5487:5487 --network="host" -v ${PWD}/deploy/api/dockerLog:/app/log schedule:latest`

# Log File

用docker啟動的程式log file在  deploy/api/dockerLog/

## Schedule DB schemas

![db_schedule.png](image/db_schedule.png?t=1660386742232)

## Swagger API

* 產生swagger document

  ```
  swag init --parseDependency --parseInternal --parseDepth 1 -d cmd/api
  ````
* swagger document url:

  http://{host}:{post}/swagger/index.html

![swagger.png](image/swagger.png)


asdfasdf
