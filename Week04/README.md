作业

`按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。`



```
──
├── cmd
│   └── week04
│       └── main.go
├── deployments         // 部署脚本
├── go.mod
├── internal
│   ├── dao             // 数据访问层
│   ├── server          // 对外暴露的接口
│   │   ├── http        // http接口
│   │   └── rpc         // rpc接口
│   └── service         // 业务逻辑层
├── pkg                 // 可被引用的包
│   └── api             // 存放 proto 文件
├── README.md
└── test                // 存在测试文件或数据，比如 ".http" 文件
```

