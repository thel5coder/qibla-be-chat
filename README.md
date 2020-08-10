# Qibla Chat Backend 1.0.0
### Version 1.0.0
### By Ivan Satyaputra

## Building

### Building requires a
[working Go 1.12+ installation](http://golang.org/doc/install).

### Main package
```
$ git clone https://repo.nusatek.id/qibla/backend/api/qibla-backend-chat.git
$ go mod download
$ go mod vendor
$ cd server
$ go run main.go
```

### Repository structure
```
files = Save files that related to this project
helper = Helper function that usually called in usecase
key = Credential file e.g. private key, etc
log = Log file
model = PgSQL query function
mongomodel = MongoDB query function
pkg = 3rd party & global function
server = Main package
├── bootstrap = Init middleware and routes
├── handler = Handler function to validate parameter inputed and handle response body
├── middleware = Route middleware
├── request = Request body struct
usecase = API logic flow
├── viewmodel = Struct of usecase response body
```

## Dependencies

NOTE: You must create .env file based on .env.example file that provided

APP_DEBUG=false : true/false, flag to debug app if panic happen  
APP_HOST=0.0.0.0:3000 : default port that app will running  
APP_LOCALE=en : default validator v9 default language  
APP_PRIVATE_KEY_LOCATION=../key/id_rsa : private key path for encrypting jwt payload  
APP_PRIVATE_KEY_PASSPHRASE= : private key passphrase  
APP_CORS_DOMAIN=http://127.0.0.1 : whitelist cors domain, can be setted by multiple domain by using coma separator. e.g. https://qibla.com,https://asia.qibla.com,etc  
APP_DEFAULT_LOCATION=Asia/Jakarta : app default location  

TOKEN_SECRET=jwtsecret : jwt string secret  
TOKEN_REFRESH_SECRET=jwtsecretrefresh : jwt string refresh secret  
TOKEN_EXP_SECRET=72 : jwt secret lifetime in hours  
TOKEN_EXP_REFRESH_SECRET=720 : jwt refresh secret lifetime in hours  

REDIS_HOST=127.0.0.1:6379 : redis connection  
REDIS_PASSWORD= : redis password  

DATABASE_HOST=127.0.0.1 : postgres ip host  
DATABASE_DB=qibla-backend-chat : postgres db name  
DATABASE_USER=postgres : postgres username  
DATABASE_PASSWORD= : postgres password  
DATABASE_PORT=5432 : postgres port  
DATABASE_SSL_MODE=disable : ssl mode, disable means no private key required for conenction  

LOG_DEFAULT=system : file/system, need to fill file path if log default is file  
LOG_FILE_PATH=../log/system.log : log file path  

FILE_MAX_UPLOAD_SIZE=10000000 : max upload size in bytes  
FILE_STATIC_FILE=../static : local public directory  
FILE_PATH=/qibla-backend-chat-bucket : subpath to access public directory  

AES_KEY=goinitsecret32bitsupersecret : secret key to encrypt sensitive data  

ODOO_ADMIN= : odoo username  
ODOO_PASSWORD= : odoo password  
ODOO_DATABASE= : odoo database name  
ODOO_URL= : odoo url  

PUSHER_APP_ID= : pusher app id  
PUSHER_KEY= : pusher access key  
PUSHER_SECRET= : pusher secret key  
PUSHER_CLUSTER= : pusher cluster  

MONGO_URL=mongodb://127.0.0.1:27017/qibla : mongoDB connection string  
MONGO_DB=qibla : mongoDB database name  

S3_URL= : s3 base url  
S3_ACCESS_KEY= : s3 access key  
S3_SECRET_KEY= : s3 secret key  
S3_BUCKET= : s3 bucket name  
S3_REGION= : s3 region  

## Deploy using Docker

NOTE: You must install docker and docker-compose first in your instance

### Docker
- How to install Docker in CentOS 8

```bash
sudo su
dnf config-manager --add-repo=https://download.docker.com/linux/centos/docker-ce.repo
dnf repolist -v
dnf install docker-ce --nobest
systemctl start docker
systemctl enable docker
```

### Docker compose
- How to install `docker-compose` in CentOS 8

```bash
dnf install curl
curl -L https://github.com/docker/compose/releases/download/1.26.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod 0744 /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```

### Run command below and your app will run on port 3000
```bash
bash docker_build.sh
```
