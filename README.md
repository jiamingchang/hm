# 基于鸿蒙操作系统的配送药品智能机器人系统设计——后端

项目开源网址：https://gitee.com/openharmony-sig/online_event/tree/master/solution_student_challenge/%E5%9F%BA%E4%BA%8EOpenHarmony%E7%9A%84%E8%87%AA%E5%8A%A8%E9%85%8D%E8%8D%AF%E6%9C%BA%E5%99%A8%E4%BA%BA%E7%B3%BB%E7%BB%9F_%E9%99%88%E6%B1%89%E6%AD%A6/

## 1、开发指导

1）目录介绍

```go
├─api		  				//接口代码
│ 
├─middleware				//中间件
│ 
├─models				    //数据库交互代码
│ 
├─setting			    	//配置文件
│ 
└─main				    	//路由及主函数
```

2）后端开发

技术栈：gin、gorm

数据库：mysql

部署：docker、nginx

1、基础封装

```go
// GenerateToken 产生token的函数
func GenerateToken(userid, authority string)(string,error){
	...
}

// ParseToken 验证token的函数
func ParseToken(token string)(*Claims,error){
	...
}

// BindAndValid 表单数据验证函数
func BindAndValid(context *gin.Context, data interface{}) bool {
	...
}

// AuthRequired token验证中间件
func AuthRequired() gin.HandlerFunc {
	return func(context *gin.Context) {
		...
	}
}

// Authority 权限验证中间件
func Authority(auth ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		...	
    }
}

```

2、创建数据库对象

创建对象，让我们能通过gorm框架去操控数据库；使用标签，指定json格式以及表单验证

```go
// 药单
type Prescription struct {
	gorm.Model
	PrescriptionID string	`json:"prescription_id" validate:"required"`
	State string		    `json:"state" validate:"required"`

	BasicInfoId uint
	BasicInfo BasicInfo		`json:"basic_info" validate:"-"`
	Mifs []MedInfoList  	`gorm:"foreignKey:PrescriptionId"`
	UserId uint				`json:"user_id" validate:"required"`
	User User				`json:"user" validate:"-"`
}

// 药单基础信息 
type BasicInfo struct {
	gorm.Model
	RecordNumber string			`json:"record_number" validate:"required"`
	PrescriptionNumber string 	`json:"prescription_number" validate:"required"`
	OpenDate string 			`json:"open_date" validate:"required"`
	BedNumber int				`json:"bed_number" validate:"required"`
	Category string				`json:"category" validate:"required"`
	ClinicalDiagnosis string	`json:"clinical_diagnosis" validate:"required"`
	AuditDoctor string			`json:"audit_doctor" validate:"required"`
	DeploymentDoctor string		`json:"deployment_doctor" validate:"required"`
	CheckDoctor string			`json:"check_doctor" validate:"required"`
	Doctor string				`json:"doctor" validate:"required"`
}

// 药品信息
type MedInfoList struct {
	gorm.Model
	MedID string	`json:"med_id" validate:"required"`
	MedName string	`json:"med_name" validate:"required"`
	Dosage string	`json:"dosage" validate:"required"`
	UseType string	`json:"use_type" validate:"required"`
	TotalNum string	`json:"total_num" validate:"required"`
	PrescriptionId uint
}

// 用户表
type User struct {
	gorm.Model
	Name string		 `json:"name" form:"name" validate:"required"`
	Password string  `json:"password" form:"password" validate:"required,max=20,min=6"`
	TrueName string  `json:"true_name" form:"true_name" validate:"required"`
	Sex int          `json:"sex" form:"sex" validate:"required"`
	Age int          `json:"age" form:"age" validate:"required"`
	Phone string	 `json:"phone" form:"phone" validate:"required,len=11"`
	Authority string `gorm:"default:user"`
}

```

3、用户管理

```go
// Login 登录账号
func Login(name, password string) (user User, err error){
	err = db.Where("name = ? AND password = ?", name, password).First(&user).Error
	return
}
// AddUser 注册账号
func AddUser(user User) (err error) {
	err = db.First(&User{}, "name", user.Name).Error
	if err != nil{
		db.Create(&user)
	}
	return
}
// GetUserList 获取用户列表
func GetUserList(name string)(users []User){
	result := db.Model(&User{})
	if name != ""{
		result = result.Where("true_name", name)
	}
	result.Order("id desc").Where("authority != ?", "admin").Find(&users)
	return
}
// DeleteUser 注销账号
func (u *User)DeleteUser() error{
	return db.Delete(&User{}, u.ID).Error
}


```

4、处方管理

```go
// Addprescription 增加处方
func (p *Prescription)Addprescription(Ids []MedInfoList) (err error){
	user, _:= GetUser(p.UserId)
	p.User = user
	db.Create(&p)
	err = db.Model(&p).Association("Mifs").Append(&Ids)
	return
}
// Deleteprescription 删除处方
func Deleteprescription(id interface{}) error{
	var prescription Prescription
	result := db.First(&prescription, id)
	if result.RowsAffected == 0{
		return errors.New("处方不存在")
	}
	db.Select(clause.Associations).Delete(&BasicInfo{}, prescription.BasicInfoId)
	db.Select(clause.Associations).Delete(&MedInfoList{}, "prescription_id", prescription.ID)
	db.Select(clause.Associations).Delete(&prescription)
	return nil
}
// Editprescriptionv1 编辑处方（user）
func Editprescriptionv1(userid, id, state interface{})(err error){
	var p Prescription
	err = db.First(&p, id).Error
	if err!=nil {
		return errors.New("更改状态失败")
	}
	if p.UserId != userid {
		return errors.New("不是你的处方")
	}
	db.Model(&p).Update("state", state)
	return nil
}
// Editprescriptionv2 编辑处方（admin）
func Editprescriptionv2(id, state interface{})(err error){
	var p Prescription
	err = db.First(&p, id).Error
	if err!=nil {
		return errors.New("更改状态失败")
	}
	db.Model(&p).Update("state", state)
	return
}
// GetprescriptionListv1 获取处方（user）
func GetprescriptionListv1(userid, state interface{}) (prescriptionList []Prescription) {
	result := db.Order("id desc")
	if state !=""{
		result = result.Where("state", state)
	}
	result.Where("user_id", userid).Preload(clause.Associations).Find(&prescriptionList)
	return
}
// GetprescriptionListv2 获取处方（admin）
func GetprescriptionListv2(state string) (prescriptionList []Prescription) {
	result := db.Order("id desc")
	if state !=""{
		result = result.Where("state", state)
	}
	result.Preload(clause.Associations).Find(&prescriptionList)
	return
}
```

4、接口讲解

封装多个接口给应用端使用

![image-20220828160936682](C:\Users\86147\AppData\Roaming\Typora\typora-user-images\image-20220828160936682.png)

封装的接口，使用apipost进行测试

![image-20220828161112740](C:\Users\86147\AppData\Roaming\Typora\typora-user-images\image-20220828161112740.png)

5、docker部署

使用docker，方便快捷部署

docker-compose.yaml

```dockerfile
version: '3'
services:

  mysql:
    image: mysql:8
    container_name: mymysql
    restart: always
    privileged: true
    environment:
      TZ: Asia/Shanghai
      # 设置 root 用户密码
      MYSQL_ROOT_PASSWORD: mysql_hm1234
      # 新建数据库
      MYSQL_DATABASE: hmsql
      # 创建新的用户
      MYSQL_USER: hm
      MYSQL_PASSWORD: mysql_hm1234
    ports:
      - "3306:3306"
    volumes:
      - /home/mysql/lib:/var/lib/mysql
      - /home/mysql/log:/var/log/mysql
    command:
      --default-authentication-plugin=mysql_native_password
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_0900_ai_ci
    networks:
      - hm-server

  nginx:
    image: nginx
    container_name: mynginx
    restart: always
    privileged: true
    ports:
      - "80:80"
      - "443:443"
    volumes:
      # 挂载nginx目录
      - /home/nginx/share:/usr/share/nginx
      # 挂载nginx日志
      - /home/nginx/log:/var/log/nginx
      # 挂载nginx配置文件
      - /home/nginx/conf.d:/etc/nginx/conf.d
      - /home/nginx/ssl:/etc/nginx/ssl
    networks:
      - hm-server

  hm:
    container_name: myhm
    restart: always
    privileged: true
    build:
      context: ./
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    volumes:
      # 挂载代码，方便同步上传
      - /home/hm:/hm
    depends_on:
      - mysql
      - nginx
    # 等待mysql完全启动
    entrypoint: ["./wait-for-it.sh", "mymysql:3306", "--", "air"]

    networks:
      - hm-server

networks:
  hm-server:
```

dockerfile

```dockerfile
FROM golang:1.17

MAINTAINER fsr

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /hm

COPY . .

# docker build时
# 热更新
RUN go get -u github.com/cosmtrek/air \
    && go mod download
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8080

```

## 2、快速上手

### 方法一：本地运行

1、下载并安装goland

2、代码下载

```
git clone https://gitee.com/jiamingchang/hm.git
```

3、项目导入

```
打开goland，点击File->Open->代码路径
```

4、打开app.ini

```
更换你自己本地的配置
```

5、开启本地服务

```
点击run按钮
```



### 方法二：使用云服务器部署

1、安装docker、docker-compose、git

```go
apt-get update
apt install docker
apt install docker-compose
apt install git
```

2、项目导入

```go
cd home
git clone https://gitee.com/jiamingchang/hm.git
```

3、nginx配置

```go
// 1、将域名证书及私钥文件放入 /home/nginx/ssl
// 2、将配置文件命名为default.conf，放入 /home/nginx/conf.d
```

nginx配置文件示例

```go
server {
    listen       80;
    listen  [::]:80;
    
    # 请更换你的域名
    server_name  xxxx;
    
   	rewrite ^(.*)$  https://$host$1 permanent;
    
    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
    }
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }
}
server {
	listen	443 ssl http2;
	listen	[::]:443 ssl http2;
  
  	# 请更换你的域名
    server_name  xxxx;
	
    # 请更换你的域名证书和私钥文件路径
    ssl_certificate	/etc/nginx/ssl/xxxx;
  	ssl_certificate_key	/etc/nginx/ssl/xxxx;
    

	location / {
	    add_header 'Access-Control-Allow-Origin' '*';
    	add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
    	add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization';

    	if ($request_method = 'OPTIONS') {
        	return 204;
    	}

		client_max_body_size 1000m;

		proxy_set_header   Host                 $host;
		proxy_set_header   X-Real-IP            $remote_addr;
		proxy_set_header   X-Forwarded-For      $proxy_add_x_forwarded_for;
		proxy_set_header   X-Forwarded-Proto    $scheme;
       
        proxy_http_version 1.1;       
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection upgrade;
        
	    proxy_pass http://myhm:8080/;
    	proxy_buffer_size 32k;
    	proxy_buffers 4 64k;
    	proxy_busy_buffers_size 64k;
    	proxy_connect_timeout 300s;
    	proxy_send_timeout 300s;
    	proxy_read_timeout 300s;
	}
}
```

4、开启服务

```
cd /home/hm
docker-compose up
```

