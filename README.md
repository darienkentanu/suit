# SUIT - (Sampah jadi dUIT)
Turn your trash into rewards

[![Go.Dev reference](https://img.shields.io/badge/gorm-reference-blue?logo=go&logoColor=blue)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/echo-reference-blue?logo=go&logoColor=blue)](https://github.com/labstack/echo)
![ERD](./erd.png)

# Table of Content
- [Description](#description)
- [How to use](#how-to-use)
- [Endpoints](#endpoints)
- [Credits](#credits)

# Description
This is final project of alterra academy

# How to use
- Install Go and MySQL or (install docker and docker-compose)
- Clone this repository in your $PATH:
```
$ git clone https://github.com/darienkentanu/suit
```
$ go run main.go or $ docker-compose up --build -d

```
after that import the file-insert.sql to your database

```


# Endpoints

| Method | Endpoint | Description| Authentication | Authorization
|:-----|:--------|:----------| :----------:| :----------:|
| POST  | /register | Register a new user | No | No
| POST | /login | Login existing user | No | No
|---|---|---|---|---|
| GET | /users | Get list of all user | Yes | Yes
| PUT | /users | Update user profile | Yes | Yes
|---|---|---|---|---|
| POST | /registeradmin | Register a new admin | No | No
| PUT | /admin | Update admin profile | Yes | Yes
|---|---|---|---|---|
| GET | /category | Get list of all category | No | No
| POST | /category | Add category by admin | Yes | Yes
| PUT | /category/:id | Update category by admin | Yes | Yes
| DELETE | /category/:id | Delete category by admin | Yes | Yes
|---|---|---|---|---|
| GET | /vouchers | Get list of all voucher | No | No
| POST | /voucher | Add voucher by admin | Yes | Yes
| PUT | /voucher/:id | Update voucher by admin | Yes | Yes
| DELETE | /voucher/:id | Delete voucher by admin | Yes | Yes
|---|---|---|---|---|
| GET | /uservouchers | Get list of all user voucher | Yes | Yes
| POST | /claimvoucher/:id | Claim voucher by voucher id | Yes | Yes
| POST | /redeemvoucher/:id | Redeem voucher by voucher id | Yes | Yes
|---|---|---|---|---|
| POST | /carts | Add category list to cart | Yes | Yes
| GET | /carts | Get list of all cart item | Yes | Yes
| PUT | /cartitems/:id | Update cart item by id | Yes | Yes
| DELETE | /cartitems/:id | Delete cart item by id | Yes | Yes
|---|---|---|---|---|
| POST | /checkout/:param | List of product on request pickup | Yes | Yes
|---|---|---|---|---|
| GET | /droppoint | Get list of all drop point | Yes | Yes
| POST | /droppoint | Add list of drop point | Yes | Yes
| PUT | /droppoint/:id | Update drop point | Yes | Yes
| DELETE | /droppoint/:id | Delete drop point | Yes | Yes
|---|---|---|---|---|
| GET | /transactions | Get list of all transaction | Yes | Yes
| GET | /transactionreport?range={range} | Get transactions with range date | Yes | Yes
|---|---|---|---|---|
| GET | /transactionsqty | Get total history of all transaction quantity | Yes | Yes
|---|---|---|---|---|


<br>

## Credits

- [Darien Kentanu](https://github.com/darienkentanu) (Author and maintainer)
- [Rizka Khairani](https://github.com/rizkakhairani) (Author and maintainer)
- [Adi Cipta Pratama](https://github.com/adicipta) (Author and maintainer)
