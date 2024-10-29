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

# 本地证书
RHILEX自签证书对，主要用来保护前后端交互需要加密的数据。

## 说明
1. **生成密钥对**
   - 输入: 无
   - 处理: 使用RSA算法生成一对密钥（公钥和私钥）
   - 输出: 公钥（public_key.pem）和私钥（private_key.pem）

2. **加密过程**
   - 输入: 原始数据（plaintext.txt）和公钥（public_key.pem）
   - 处理: 使用公钥对原始数据进行加密
   - 输出: 加密后的数据（ciphertext.bin）

3. **解密过程**
   - 输入: 加密后的数据（ciphertext.bin）和私钥（private_key.pem）
   - 处理: 使用私钥对加密后的数据进行解密
   - 输出: 解密后的数据（decrypted_text.txt）
