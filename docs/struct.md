Structure
=========

### Сценарий основных действий в платформе

Действия с проектом (имя проекта должно быть уникально в рамках пользователя)
- Создать проект
- Получить список проектов
- Получить проект по id/name
- Обновить информацию по проекту по id/name
- Удалить проект по id/name


Действия с сервисом (имя сервиса должно быть уникально в рамках проекта).
- Создать сервис
- Получить список сервисов по project_id/name
- Получить сервис по id/name
- Обновить информацию по сервису по id/name
- Удалить сервис по id/name


### Дерево хранения данных в etc

```yaml
/lastbackend

/lastbackend/vendors/<vendor>: <vendor (github|bitbucket|gitlab) object>

# project information data layer
/lastbackend/projects/<project id>/meta: <project info object>
/lastbackend/projects/<project id>/services/<service id>/meta: <service info object>
/lastbackend/projects/<project id>/services/<service id>/config: <service config object>
/lastbackend/projects/<project id>/services/<service id>/source: <service sources object>
/lastbackend/projects/<project id>/services/<service id>/domains: <service domains object>
/lastbackend/projects/<project id>/services/<service id>/image: <image name>
/lastbackend/projects/<project id>/services/<service id>/containers/<container id>: <service container object>
/lastbackend/projects/<project id>/services/<service id>/builds/<build number>: <service build object>

# image information data layer
/lastbackend/images/<image id>/meta: <image info object>
/lastbackend/images/<image id>/source: <image source object>
/lastbackend/images/<image id>/builds/<build number>: <build object>

# helpers information data layer
/lastbackend/helper/projects/<name>: <project id>
/lastbackend/helper/projects/<project id>/services/<name>: <service id>
/lastbackend/helper/pods/<pod id>: <project id>:<service id>
```

### Структуры данных

Image info object
```json
{
  "name": "hub.lastbackend.com/lastbackend/hello-world",
  "description": "hello-world description",
  "created": "Wed Mar 01 2017 17:13:08 GMT+03:00",
  "updated": "Wed Mar 01 2017 17:13:08 GMT+03:00"
}
```

Image source object
```json
{
  "hub":"github.com",
  "owner":"lastbackend",
  "repo":"proxy",
  "branch":"master"
}
```

Build object
```json
{
  "commit": "a454517a3c5c657cc71548b874d023f2e2d8915b",
  "commitMessage": "Merge pull request #218",
  "committer": "unloop",
  "status": "failed",
  "message": "clone repo failed",
  "created": "Wed Mar 01 2017 17:13:08 GMT+03:00",
  "updated": "Wed Mar 01 2017 17:13:08 GMT+03:00"
}
```

Registry account object
```json
{
 "email": "",
 "password": "",
 "username": "",
 "host": "hub.lastbackend.com"
}
```

Vendor map
```json
{
  "github": {},
  "bitbucket": {},
  "gitlab": {},
  "slack": {}
}

```

Vendor github object
```json
{
  "vendor": "github",
  "host": "github.com",
  "username": "unloop",
  "email": "pastor.konstantin@gmail.com",
  "service_id":  "1877907",
  "token": "a4f29b4c1a7fa86e8b1f308ffec6feeea979f98e",
  "token_type": "bearer",
  "updated": "Tue Mar 14 2017 12:49:45 GMT+03:00",
  "created": "Tue Mar 14 2017 12:49:45 GMT+03:00"
}
```

Vendor bitbucket object
```json
{
  "vendor": "bitbucket",
  "host": "bitbucket.org",
  "username": "unloop",
  "email": "pastor.konstantin@gmail.com",
  "service_id":  "1877907",
  "token": "NqFQPmgsa4QQ9StW2R",
  "token_type": "bearer",
  "updated": "Tue Mar 14 2017 12:49:45 GMT+03:00",
  "created": "Tue Mar 14 2017 12:49:45 GMT+03:00"
}
```
Vendor gitlab object
```json
{
  "vendor": "githlab",
  "host": "gitlab.com",
  "username": "unloop",
  "email": "pastor.konstantin@gmail.com",
  "service_id":  "1877907",
  "token": "1f0af717251950dbd4d73154fdf0a474a5c5119adad999683f5b450c460726aa",
  "token_type": "bearer",
  "expiry": "Tue Mar 14 2017 12:49:45 GMT+03:00",
  "updated": "Tue Mar 14 2017 12:49:45 GMT+03:00",
  "created": "Tue Mar 14 2017 12:49:45 GMT+03:00"
}
```

Project info object
```json
{
  "name": "demo",
  "created": "Wed Mar 01 2017 17:13:08 GMT+03:00",
  "updated": "Wed Mar 01 2017 17:13:08 GMT+03:00"
}
```

Service info object
```json
{
  "name": "mysql",
  "created": "Wed Mar 01 2017 17:13:08 GMT+03:00",
  "updated": "Wed Mar 01 2017 17:13:08 GMT+03:00"
}
```

Service config object
```json
{
  "image": "library/lastbackend/proxy:latest",
  "name": "",
  "replicas": 2,
  "memory": 32,
  "ports": {},
  "env": {},
  "volumes": {}
}
```

Service domains object
```json
{
  "service.lbapp.in": true,
  "service.domain.com": false
}
```

Service container object
```json
{
  "id": "59e8bce5a3032034dd84339c64fec42a8084bba90cdac6115f9456e29f646015",
  "status": "running",
  "ports" : {
    "3306/TCP": 44536
  },
  "updated": "Wed Mar 01 2017 17:13:39 GMT+03:00",
  "created": "Wed Mar 01 2017 17:13:39 GMT+03:00"
}
```