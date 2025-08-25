# geeksaver
该工具仅供学习使用，可以将自己买的极客时间课程保存为本地的 markdown 文档，借助 AI 的总结能力，方便学习和复习。

### 支持的命令
```
login       极客时间登录态设置
config      配置信息
md          极客时间课程转 Markdown 并本地保存
```
### 使用方法
1. #### 先设置登陆态
在极客时间网页登陆完之后，打开f12，接着点击：应用程序 -> Cookie，就获取 GCID 和 GCESS 了
    ```
      geeksaver login --gcid 123456 --gcess abcdefg
    ```
2. #### 检查配置信息
    ```
      geeksaver config
    ```

3. #### 从购买的课程详情页，获取课程id
    例如：
    课程详情页：https://time.geekbang.org/column/intro/100029501
    课程ID：100029501

4. #### 开始保存！
    课程到 markdown 格式，默认保存在 `$HOME/geek-docs`
    ```
      geeksaver md --cid 100029501
    ```

5. #### 使用 markdown 工具打开课程文件夹开始阅读
    我用的 markdown 工具是 Obsidian。 