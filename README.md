# SUIT - (Sampah jadi dUIT)

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
$ git clone https://github.com/darienkentanu/suit
```

To run this project first you must insert the following query to mysql

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
| GET | /category | Get list of all category | No | No
| POST | /category | Add category by admin | Yes | Yes
| PUT | /category/:id | Update category by admin | Yes | Yes
| DELETE | /category/:id | Delete category by admin | Yes | Yes
|---|---|---|---|---|
| GET | /voucher | Get list of all voucher | No | No
| POST | /voucher | Add voucher by admin | Yes | Yes
| PUT | /voucher/:id | Update voucher by admin | Yes | Yes
| DELETE | /voucher/:id | Delete voucher by admin | Yes | Yes
|---|---|---|---|---|
| POST | /carts | Add category list to cart | Yes | Yes
| GET | /carts | Get list of all cart item | Yes | Yes
| PUT | /cartitems/:id | Update cart item by id | Yes | Yes
| DELETE | /cartitems/:id | Delete cart item by id | Yes | Yes
|---|---|---|---|---|
| POST | /checkout/:param | list of product on request pickup | Yes | Yes
|---|---|---|---|---|
| GET | /transactions | Get list of all transaction | Yes | Yes
| GET | /transactionreport?range={range} | Get transactions with range date | Yes | Yes
|---|---|---|---|---|
| POST | /redeemvoucher/:id | Redeem voucher by ID | Yes | Yes
<br>

## Credits

- [Darien Kentanu](https://github.com/darienkentanu) (Author and maintainer)
- [Rizka Khairani](https://github.com/rizkakhairani) (Author and maintainer)
- [Adi Cipta Pratama](https://github.com/adicipta) (Author and maintainer)
