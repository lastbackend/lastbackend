Database storage structure.

As main storage engine we decided to use etcd v3 database. Etcd key,value storage is created to store most critical data
in modern distributed apps.
 
Main storage structure

###### system storage

/system/templates
```yaml
    categories:
      uuid:
        name:
    templates:
       uuid:
        name:
```

/system/hooks/<token>
```yaml
    token:
    image:
```

/system/clusters
```yaml
    - name:
      region:
      uri:
      ssl:
        ca:
        key:
        crt:
```

###### user storage


/users/<username>/account/info
```yaml
    username:
    email:
    gravatar:
```

/users/<username>/account/balance
```yaml
    balance:
```

/users/<username>/account/billing/debit/<date>
```yaml
    
```

/users/<username>/account/billing/credit/<date>
```yaml
    
```

/users/<username>/account/security/

/users/<username>/account/security/access
```yaml
    password:
    salt:
    token:
```
/users/<username>/account/security/keys/<name>
```yaml
  key:
    key:
    name:
```

/users/<username>/account/integrations/github
```yaml
  github:
    token:
```

/users/<username>/account/integrations/bitbucket
```yaml
  bitbucket:
    token:
    token-refresh:
    token-lifetime:
```

/users/<username>/account/integrations/gitlab
```yaml
  gitlab:
    token:
    token-refresh:
    token-lifetime:
```

/users/<username>/profile
```yaml
  profile:
    firstname:
    lastname:
```

/users/<username>/organization/<username>
```yaml
  role:
```

/users/<username>/projects/<uuid>
```yaml
  project:
    name:
    desc:
    namespace:
```

/users/<username>/images/<image>/<tag>/info
```yaml
  image:
    name:
    desc:
```

/users/<username>/images/<image>/<tag>/builds/<uuid>
```yaml
  build:
    number:
    status:
      step:
      cancelled:
      message:
      error:
      updated:
    source:
      hub:
      repo:
      tag:
      commit:
        commit:
        committer:
        author:
        message:
      auth:
        token:
    image:
      registry:
        host:
        token:
      repo:
      tag:
    request:
      type:
      owner:
    created:
    updated:
```

/organizations/

/organizations/<username>/profile
```yaml
  profile:
    name:
    desc:
```

/organizations/<username>/members

/organizations/<username>/members/<username>
```yaml
  role:
```

/organizations/<username>/projects/<uuid>
```yaml
  project:
    name:
    desc:
    namespace:
```

/organizations/<username>/images/<image>/<tag>/info
```yaml
  project:
    name:
    desc:
```

/organizations/<username>/images/<image>/<tag>/builds/<uuid>
```yaml
  build:
    number:
    status:
      step:
      cancelled:
      message:
      error:
      updated:
    source:
      hub:
      repo:
      tag:
      commit:
        commit:
        committer:
        author:
        message:
      auth:
        token:
    image:
      registry:
        host:
        token:
      repo:
       tag:
    request:
      type:
      owner:
    created:
    updated:
```