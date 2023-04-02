# renpy_translate
自动把renpy里的对话机器翻译成中文

# 使用
1. 下载unrpa，用来把game\下的.rpa文件解压
```
pip install unrpa
```
2. 下载[unrpyc](https://github.com/CensoredUsername/unrpyc)，根据release里的使用说明操作。
需要把un.rpy拷贝到game\下，运行游戏。最后得到.rpy文件。
3. 找到game\script目录，应该有很多.rpy文件。执行本工具。
```
./renpy_translate -input script的目录
```
工具会找到所有的对话，然后翻译成中文替换。翻译的结果会保存运行目录下的translate.json里。
4. 安装字体，把SourceHanSansLite.ttf放到game\fonts下。并创建文件game\tl\style.rpy，内容为：
```python
translate chinese python:
    gui.system_font = gui.main_font = gui.text_font = gui.name_text_font = gui.interface_text_font = gui.button_text_font = gui.choice_button_text_font = "fonts/SourceHanSansLite.ttf"
```