# geeksaver

[![Release](https://github.com/bynow2code/geeksaver/actions/workflows/release.yml/badge.svg)](https://github.com/bynow2code/geeksaver/actions/workflows/release.yml)

该工具仅供学习使用，可将已购买的极客时间课程保存为本地 Markdown 文档，并借助 AI 总结能力梳理内容，方便后续学习与复习。

## 支持的命令

| 命令      | 功能描述                           |
|---------|--------------------------------|
| login   | 配置极客时间登录状态（设置 GCID/GCESS）      |
| config  | 查看当前工具的配置信息（含登录态、保存路径等）        |
| md      | 将指定极客时间课程转换为 Markdown 格式并保存到本地 |
| version | 输出版本信息                         |

## 安装方式

1. ### 使用 go install 安装（适合有 Go 环境的用户，无版本信息）

   默认安装路径为系统 `GOPATH/bin` 目录，若已将该路径添加到系统环境变量，可直接全局调用工具。

    ```
    # 查看 GOPATH 路径
    go env GOPATH
    
    # 安装 geeksaver
    go install github.com/bynow2code/geeksaver@latest
    
    # 查看当前配置
    geeksaver config
    ```

2. ### 从 Releases 页面下载二进制文件（推荐，包含版本信息）
    1. 访问项目 GitHub 仓库的「Releases」页面；
    2. 根据本地操作系统（Windows/macOS/Linux）和硬件架构（amd64/arm64），下载对应的二进制压缩包（如
       `geeksaver-darwin-amd64.tar.gz` 对应 macOS 64 位系统）；
    3. 解压压缩包，将得到的二进制文件（如 `geeksaver` 或 `geeksaver.exe`）添加到系统环境变量，或直接在解压目录通过终端调用。

## 完整使用步骤

1. #### 配置登录态（必做：获取极客时间登录凭证 GCID/GCESS）
   需先在极客时间网页端登录账号，再通过浏览器开发者工具获取登录凭证：
    1. 打开极客时间官网，使用账号密码/验证码完成登录；
    2. 按下 `F12` 打开浏览器开发者工具，切换到「应用程序」（Application）选项卡；
    3. 在左侧导航栏展开「Cookie」，点击极客时间域名（如 `time.geekbang.org`）；
    4. 在右侧 Cookie 列表中，找到 GCID 和 GCESS 两项，复制对应的「值」（Value）。
    ```
    # 执行登录态配置命令（替换为实际获取的 GCID 和 GCESS 值）
    geeksaver login --gcid "你的GCID值（如：abc123xyz）" --gcess "你的GCESS值（如：def456uvw）"
    ```

2. #### 检查配置信息（确认登录态信息是否正确）
   执行命令查看当前配置，确认登陆信息和保存路径等信息是否正确：
    ```
    geeksaver config
    ```

3. #### 获取课程 ID（从课程详情页 URL 中提取）
   已购买课程的 ID 可从课程详情页 URL 末尾获取：
    - 示例：课程详情页 URL 为 https://time.geekbang.org/column/intro/100029501
    - 提取的课程 ID：100029501（即 URL 中 intro/ 后的数字串）

4. #### 转换并保存课程为 Markdown
   工具会自动下载课程所有章节内容，转换为 Markdown 格式后保存到本地，默认保存路径为 `$HOME/geek-docs`（Windows 系统路径为
   `C:\Users\你的用户名\geek-docs`）。

5. #### 打开并阅读课程文档
   使用任意 Markdown 阅读器打开保存目录即可学习，推荐工具：
    - **Obsidian**：支持目录管理、双向链接，适合构建课程知识体系，方便复习；

## 注意事项

- 该工具仅用于个人学习，请遵守极客时间用户协议，不得用于商业用途或非法传播；
- 若保存过程中出现「登录失效或用户未购买课程」等提示，需重新执行 `login` 命令更新 `GCID` 和 `GCESS`；
