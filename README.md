# zlsgo-app

[使用文档 https://docs.73zls.com/zls-go/#/](https://docs.73zls.com/zlsgo/#/89a1532c-15bd-495e-a427-30b3cff6e061)


## 开发

```bash
# 先编译再执行，首次执行会自动生成配置文件
# 配置文件和执行文件处于同一个目录
go build -o tmpApp && ./tmpApp
```


```bash
docker run --rm --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=666666 -e POSTGRES_USER=root -e POSTGRES_DB=zls -v $PWD/../tmp/postgres/data:/var/lib/postgresql/data -d postgres:15


docker run --rm --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=666666 -e MYSQL_DATABASE=zls  -v $PWD/../tmp/mariadb:/var/lib/mysql -d mariadb:10.5.5
```
