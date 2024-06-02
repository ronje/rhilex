<!--
 Copyright (C) 2024 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

**提交指南**

### 1. 检查你的更改

在提交代码之前，确保你的更改是正确且完整的。运行测试用例以确保代码没有引入新的问题。

### 2. 分解更改

如果你的更改涉及多个逻辑功能，考虑将它们拆分成独立的提交。这样可以使代码审查更加简单，并且更容易理解每个提交的目的。

### 3. 编写清晰的提交消息

每个提交都应该有一个清晰、简洁的提交消息。提交消息应该包括以下三个部分：

- **类型（Type）**：描述你的更改类型。常见的类型包括：`feat`（新功能）、`fix`（修复bug）、`docs`（文档更新）、`style`（代码格式化等）、`refactor`（重构代码）、`test`（添加测试）、`chore`（构建过程或辅助工具的变动）等。
- **范围（Scope）**（可选）：描述你的更改影响的范围，比如模块、功能等。如果更改涉及多个范围，可以使用通配符 `*` 代替。
- **主题（Subject）**：简要描述你的更改。建议使用中文描述，不超过50个字符，并且不要以句号或其他标点符号结尾。

例如：

```
feat(API): 添加用户认证功能
fix(Database): 修复用户查询缺少username属性的bug
docs(README): 更新安装说明
```

### 4. 提交代码

当你准备好提交代码时，使用Git提交命令提交你的更改。确保你的提交消息符合上述的规范。

```
git add .
git commit -m "feat(API): 添加用户认证功能"
```

### 5. 推送更改

如果你的更改是在一个分支上，推送你的更改到远程仓库。

```
git push origin your-branch
```

### 6. 更新你的分支

如果你的分支是基于远程仓库的一个主要分支（如 `master`），确保在推送更改之前，将主要分支的最新更改合并到你的分支中。

```
git pull origin master
```

### 7. 进行代码审查

等待团队成员对你的提交进行审查。在审查过程中，可能会提出改进建议，根据反馈进行相应的修改。

### 8. 合并更改

一旦你的更改通过审查，将它们合并到主要分支中。

```
git checkout master
git merge your-branch
```

### 9. 完成

恭喜你，你的更改已经成功地提交到主要分支中！

