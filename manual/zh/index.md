# Golang 开发手册

个人维护的 Golang 开发手册。条款分 MUST / SHOULD / MAY 三级，目录与条款同源于 `gdm/data/`。

- 基准版本：Go 1.27.1
- 条款总数：113

## 目录

1. [一、编程规约](01-programming-conventions.md)
    - 命名规约
    - 常量与枚举规约
    - 格式化与代码风格规约
    - 函数与方法规约
    - 数据结构规约
    - 并发规约
    - 控制语句规约
    - 注释规约
2. [二、错误与日志](02-errors-and-logging.md)
    - 错误处理规约
    - 错误码规约
    - 日志规约
3. [三、单元测试规约](03-unit-testing.md)
    - 测试命名与结构
    - 表驱动与用例设计
    - 断言与测试数据
    - Mock 与测试替身
    - 覆盖率的使用
4. [四、依赖与工程结构](04-dependencies-and-layout.md)
    - 分层规约
    - 包组织规约
    - 依赖管理规约
    - 项目布局规约
5. [五、接口与设计规约](05-interfaces-and-design.md)
    - 接口定义与落位
    - 组合与复用
    - 依赖注入与装配
    - context 传递
    - 结构体安全形状
6. [六、性能规约](06-performance.md)
7. [七、安全规约](07-security.md)
8. [八、静态检查规约](08-static-analysis.md)
    - 配置与运行
    - 必开规则
    - 推荐规则
9. [九、工具规约](09-tooling.md)
    - 命令行
    - 测试工具
    - 代码生成
10. [十、附录](10-appendix.md)
