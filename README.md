# Bark

<img src="https://wx3.sinaimg.cn/mw690/0060lm7Tly1g0nfnjjxbbj30sg0sg757.jpg" width=200px height=200px />

[Bark](https://github.com/Finb/Bark) is an iOS App which allows you to push customed notifications to your iPhone.


## Table of Contents

   * [Bark](#bark)
      * [Installation](#installation)
        * [For Azure Function](#for-azure-function)
      * [Contributing to bark-server](#contributing-to-bark-server)
         * [Development environment](#development-environment)
      * [Update](#update)


## Installation

### For Azure Function

```sh
#1. build app
go build

#2. deploy to Azure Function.

#3. create Azure Cosmos DB.

#4. set Cosmos DB params in Function's Application Config
{
    "MONGODB_CONNECTION_STRING": "mongodb://xxxxxxx",
    "MONGODB_DATABASE": "base",
    "MONGODB_COLLECTION": "bark",
}
```

## Other Docs

### 中文:

- [https://day.app/2018/06/bark-server-document/](https://day.app/2018/06/bark-server-document/)
  

## Contributing to bark-server

### Development environment

This project requires at least the golang 1.12 version to compile and requires Go mod support.

- Golang 1.14
- GoLand 2019.3 or other Go IDE
- Docker(Optional)

## Update 

The push certificate embedded in the program expires on **`2020/01/30`**, please update the program after **`2019/12/01`**
