# 新航前端 HTML 版本（自动导出）

## 说明
这是从 Vue 前端自动导出的静态 HTML 版本，用于设计修改。

## 文件结构
- index.html — 首页
- about.html — 关于新航
- academics.html — 学术课程
- campus.html — 校园
- student-life.html — 学生生活
- apply.html — 招生报名
- styles.css — 全部样式
- script.js — 语言切换 + 导航 + 动画
- assets/ — 图片资源

## 使用方式
1. 直接用浏览器打开任意 .html 文件即可预览
2. 修改 HTML 内容后保存，刷新浏览器查看效果
3. 点击右上角"中文"按钮可切换中英文

## 双语文本格式
所有双语文本都使用这种结构:
```html
<span data-en="English text" data-zh="中文文本">English text</span>
```
- 默认显示 data-en 的内容
- 切换语言后显示 data-zh 的内容
- 修改时请同时更新 data-en 和 data-zh 属性

## 注意事项
- 这只是展示层，不包含后端逻辑（登录、报名提交等）
- 报名表仅保留 HTML 结构，表单提交功能需要后端支持
- 修改完成后请将改动后的 HTML 文件发回，我们会同步到 Vue 版本

## 导出时间
2026/6/26 01:48:12
