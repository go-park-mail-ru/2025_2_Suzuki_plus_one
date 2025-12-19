# AWS S3 with Cloud.ru

Классное [s3](https://cloud.ru/products/evolution-object-storage) от сбербанка

Пароли лежат в [bitwarden](https://send.bitwarden.com/#FN9Mch_MCE-FprOtACe0fA/aYijwQR5GhPemfOuAWF2QQ).

Пароль от bitwarden в личке тг.

## ЧаВО

0. [Войти](#log-in-cloudru) в UI
1. [Установить](#aws-cli) aws cli
2. [Ввести](#aws-creds) токены, которые я сгенерил ([или через файлик](https://cloud.ru/docs/s3e/ug/topics/tools__aws-cli?source-platform=Evolution#id1))
3. [Загрузить файл](#common-cmds) в бакет

## Log in Cloud.ru

## s3 UI url

[Console](https://console.cloud.ru/spa/svp/s3-storage?tabName=buckets&customerId=8fbf5ed3-c5bf-42c6-9630-6b3978d7f98e&projectId=39a1fd2e-4e1f-49a7-ab69-5ffd1d27b243)

## Creds

[Bitwarden](https://send.bitwarden.com/#FN9Mch_MCE-FprOtACe0fA/aYijwQR5GhPemfOuAWF2QQ)

## Aws cli

[Installation](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html#getting-started-install-instructions)

[Sber tutorial](https://cloud.ru/docs/s3e/ug/topics/tools__aws-cli)

### Linux

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

### Cmds

#### Aws creds

[Bitwarden](https://send.bitwarden.com/#FN9Mch_MCE-FprOtACe0fA/aYijwQR5GhPemfOuAWF2QQ)

#### Common cmds

- upload file

  ```bash
  KEY=dima.png
  BUCKET=bucket-test
  aws s3api put-object --endpoint-url https://s3.cloud.ru --bucket $BUCKET --key $KEY --body $KEY
  ```

- list current buckets

  ```bash
  KEY=dima.png
  BUCKET=bucket-test
  aws s3api list-buckets --endpoint-url https://s3.cloud.ru --query "Buckets[].Name"
  ```

- list files in bucket

  ```bash
  BUCKET=bucket-test
  aws s3api list-objects --endpoint-url https://s3.cloud.ru --bucket $BUCKET --query "Contents[].Key"
  ```

- download file

  ```bash
  KEY=dima.png
  BUCKET=bucket-test
  aws s3api get-object --endpoint-url https://s3.cloud.ru --bucket $BUCKET --key $KEY $
  ```

#### Additional cmds

- create new bucket

  ```bash
  KEY=dima.png
  BUCKET=bucket-test
  aws s3api create-bucket --endpoint-url https://s3.cloud.ru --bucket $BUCKET
  ```

- [WARN] remove all data from the bucket

  ```bash
  KEY=dima.png
  BUCKET=bucket-test
  aws s3 rm s3://$BUCKET --recursive --endpoint-url https://s3.cloud.ru
  ```
