# Restful-API-Project

[![Go.Dev reference](https://img.shields.io/badge/gorm-reference-blue?logo=go&logoColor=blue)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/echo-reference-blue?logo=go&logoColor=blue)](https://github.com/labstack/echo)


# Table of Content

- [Description](#description)
- [How to use](#how-to-use)
- [Endpoints](#endpoints)
- [Credits](#credits)

# Description
project-base task alterra academy

# How to use
- Install Go and MySQL
- Clone this repository in your $PATH:
```
$ git clone https://github.com/darienkentanu/Restful-API-Project.git
```

To run this project first you must insert the following query to mysql

* CREATE DATABASE IF NOT EXISTS `restfulapiproject`;
* USE `restfulapiproject`;

* INSERT INTO `users` (`fullname`,`username`,`email`,`password`,`phone_number`,
`gender`, `address`, `role`, `created_at`)
VALUES ('admin', 'admin','admin@gmail.com', '$2a$14$jqZPvRBaylmWVCK4Rnh6tOpKTn3B/6PlNlT0HOSLYgtOSGg9z/BGG',
'081234567890','male','jakarta', 'admin',curdate());

* admin password = "password"

Run `main.go`
```
$ go run main.go
```


# Endpoints

| Method | Endpoint | Description| Authentication | Authorization
|:-----|:--------|:----------| :----------:| :----------:|
| POST  | /register | Register a new user | No | No
| POST | /login | Login existing user| No | No
|---|---|---|---|---|
| GET    | /users | Get list of all user | Yes | Yes
| PUT | /users | Update user profile | Yes | Yes
|---|---|---|---|---|
| GET   | /categories | Get category list | Yes | Yes
| POST   | /categories | Add category | Yes | Yes
| DELETE   | /categories/:id | Delete category by id | Yes | Yes
|---|---|---|---|---|
| GET | /products | Get list of all products | No | No
| GET | /products/:id | Get product by product id | No | No
| POST | /products | Add products by admin | Yes | Yes
| PUT | /products/:id | Update products by admin | Yes | Yes
| DELETE | /products/:id | Delete products by admin | Yes | Yes
|---|---|---|---|---|
| POST | /carts | Add products to cart | Yes | Yes
| GET | /carts | Get list of all cart item | Yes | Yes
| PUT | /cartitems/:id | Update cart item by id | Yes | Yes
| DELETE | /cartitems/:id | Delete cart item by id | Yes | Yes
|---|---|---|---|---|
| POST | /checkout | List of products checkout | Yes | Yes
|---|---|---|---|---|
| GET | /payments/:id | Get transaction status | Yes | Yes
|---|---|---|---|---|
| GET | /transactions | Get list of all transaction | Yes | Yes
| GET | /transactionreport?range={range} | Get transactions with range date | Yes | Yes
|---|---|---|---|---|

<br>

## Credits

- [Darien Kentanu](https://github.com/darienkentanu) (Author and maintainer)
- [Rizka Khairani](https://github.com/rizkakhairani) (Author and maintainer)
- [Adi Cipta Pratama](https://github.com/adicipta) (Author and maintainer)
