# 依据来源

写条款与深挖机制时的上游材料。条款 `sources` 从本清单里挑与条款直接对应的深链；数据格式与写作口径见 `gdm/CLAUSES.md`。

## 官方与权威（英文）

- [Effective Go](https://go.dev/doc/effective_go)
- [Package names（官方博客）](https://go.dev/blog/package-names)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Organizing a Go module](https://go.dev/doc/modules/layout)
- [go.mod file reference](https://go.dev/doc/modules/gomod-ref)
- [Managing dependencies](https://go.dev/doc/modules/managing-dependencies)
- [Go 1.4 Release Notes：Internal packages](https://go.dev/doc/go1.4#internalpackages)
- [Go Test Comments](https://go.dev/wiki/TestComments)
- [Google Go Style 总览](https://google.github.io/styleguide/go/)
- [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
- [Google Go Style Best Practices](https://google.github.io/styleguide/go/best-practices)
- [Structured Logging with slog（官方博客）](https://go.dev/blog/slog)
- [log/slog](https://pkg.go.dev/log/slog)
- [Using go fix to modernize Go code（官方博客）](https://go.dev/blog/gofix)
- [Go 1.26 Release Notes](https://go.dev/doc/go1.26)
- [cmd/fix](https://pkg.go.dev/cmd/fix)
- [Google API 设计指南：错误](https://cloud.google.com/apis/design/errors)

## 安全编码（官方）

- [Avoiding SQL injection risk](https://go.dev/doc/database/sql-injection)
- [Using prepared statements](https://go.dev/doc/database/prepared-statements)
- [Directory-limited filesystem access（Go 1.24 发布说明）](https://go.dev/doc/go1.24#directory-limited-filesystem-access)
- [os（os.Root）](https://pkg.go.dev/os#Root)
- [os/exec](https://pkg.go.dev/os/exec)
- [crypto/rand](https://pkg.go.dev/crypto/rand)
- [math/rand](https://pkg.go.dev/math/rand)
- [crypto/tls（Config）](https://pkg.go.dev/crypto/tls#Config)
- [crypto/md5](https://pkg.go.dev/crypto/md5)
- [crypto/sha1](https://pkg.go.dev/crypto/sha1)
- [crypto/des](https://pkg.go.dev/crypto/des)
- [crypto/rc4](https://pkg.go.dev/crypto/rc4)
- [html/template](https://pkg.go.dev/html/template)
- [net/http（Server.ReadHeaderTimeout）](https://pkg.go.dev/net/http#Server.ReadHeaderTimeout)
- [crypto/subtle](https://pkg.go.dev/crypto/subtle)
- [crypto/hmac（Equal）](https://pkg.go.dev/crypto/hmac#Equal)

## 工具与检测

- [gosec 规则清单（RULES.md）](https://github.com/securego/gosec/blob/master/RULES.md)
- [revive 规则清单（RULES_DESCRIPTIONS.md）](https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md)
- [go-critic 检查项总览](https://go-critic.com/overview.html)

## 工具链

- [golangci-lint 总览](https://golangci-lint.run/)
- [golangci-lint：False Positives（nolint 指令）](https://golangci-lint.run/docs/linters/false-positives/)
- [golangci-lint：Linters Settings](https://golangci-lint.run/docs/linters/configuration/)
- [go 命令：Generate Go files by processing source](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source)
- [go 命令：Build constraints（文件名约定）](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [protoc 生成的 Go 代码（输出文件命名）](https://protobuf.dev/reference/go/go-generated/)
- [cmd/fix](https://pkg.go.dev/cmd/fix)
- [testify](https://github.com/stretchr/testify)
- [testifylint](https://github.com/Antonboom/testifylint)
- [cobra](https://github.com/spf13/cobra)
- [wire](https://github.com/google/wire)

## 中文镜像与翻译

- [高效 Go 编程（golang.ac.cn 镜像）](https://golang.ac.cn/doc/effective_go)
- [Go 代码审查评论（golang.ac.cn 镜像）](https://golang.ac.cn/wiki/CodeReviewComments)

## 社区标准

- [Uber Go Style Guide](https://github.com/uber-go/guide)
- [Uber Go 语言编码规范（中文翻译）](https://github.com/xxjwxc/uber_go_guide_cn)
- [Go 编码建议（中文电子书）](https://dablelv.github.io/go-coding-advice/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Proverbs](https://go-proverbs.github.io/)
- [腾讯云：Go 编码规范建议](https://developer.cloud.tencent.com/article/1911268?policyId=1003)
- [腾讯云：Go 语言开发规范指南](https://cloud.tencent.cn/developer/article/2512042)
- [字节 Go 开发规范（第三方转载笔记）](https://www.haveyb.com/article/3102)
- [Practical Go：写可维护 Go 程序的建议（Dave Cheney）](https://dave.cheney.net/practical-go/presentations/qcon-china.html)

## 企业级与项目级分级规范

- [ChainMaker 贡献代码管理规范及流程](https://docs.chainmaker.org.cn/v2.2.1/html/contribution/%E8%B4%A1%E7%8C%AE%E4%BB%A3%E7%A0%81%E7%AE%A1%E7%90%86%E8%A7%84%E8%8C%83%E5%8F%8A%E6%B5%81%E7%A8%8B.html)
- [Lerian Studio Go standards（quality）](https://raw.githubusercontent.com/LerianStudio/ring/main/dev-team/docs/standards/golang/quality.md)
- [qarium/goga-lang-conventions](https://raw.githubusercontent.com/qarium/goga-lang-conventions/refs/heads/0.0.x/golang/project.md)

## 阿里巴巴 Java 开发手册（形态参照）

- [阿里巴巴 Java 开发手册（嵩山版专题）](https://developer.aliyun.com/topic/java20)
- [Alibaba Java Coding Guidelines（英文 GitBook）](https://alibaba.github.io/Alibaba-Java-Coding-Guidelines/)
- [alibaba/p3c](https://github.com/alibaba/p3c)

## Agent 检索机制参照

- [JetBrains go-modern-guidelines](https://github.com/JetBrains/go-modern-guidelines)
- [samber/cc-skills-golang](https://github.com/samber/cc-skills-golang)
- [awesome-copilot Go instructions](https://awesome-copilot.github.com/instruction/go/)
