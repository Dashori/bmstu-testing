# Vault!

```docker-compose up -d```

```export VAULT_ADDR='http://127.0.0.1:8200'```

\*Тут нужно установить волт\*

```vault status``` -- должно быть "sealed": false

```vault operator init``` -- выплюнет 5 ключей и рут токен

для 3х ключей нужно сделать:
``` vault operator unseal <пароль>```

после этого ```vault status``` -- должно быть "sealed": true

```vault login <рут токен>``` -- атуентификация

```vault secrets enable -path=secret/ kv``` -- делаем папочку secret возможной для kv

```vault kv put secret/hello pasword=qwerty``` -- кладем pasword=qwerty по пути secret/hello

```vault kv get secret/hello``` -- получаем секрет по пути secret/hello

И в CI:
```
- export VAULT_ADDR='http://127.0.0.1:8200'
- export PASSWORD="$(vault kv get -field=password secret/emailFrom)"
```


### PROFIT!



Дальше команды, которые нужны если волт не локальный и нужно использовать jwt:

```vault policy write email  - <<EOF
# Policy name: email
#
# Read-only permission on 'secret/*'
path "secret/*" {
    capabilities = [ "read" ]
}
EOF
```


```vault write auth/jwt/role/email - <<EOF
{
  "role_type": "jwt",
  "policies": ["email"],
  "token_explicit_max_ttl": 120,
  "user_claim": "user_email",
  "bound_claims": {
    "project_id": "project_id",
    "ref": "*",
    "ref_type": "branch"
  }
}
EOF
```


```vault write auth/jwt/config \
    oidc_discovery_url="https://git.iu7.bmstu.ru" \
    bound_issuer="https://git.iu7.bmstu.ru"
```

Ну и с CI:
```- export VAULT_TOKEN="$(vault write -field=token auth/jwt/login role=email jwt=$CI_JOB_JWT)"```

полезные ссылки:
https://blog.ruanbekker.com/blog/2019/05/06/setup-hashicorp-vault-server-on-docker-and-cli-guide/

https://www.dmosk.ru/instruktions.php?object=vault-hashicorp

https://habr.com/ru/companies/nixys/articles/512754/
