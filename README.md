## shulun
#### 最大公因子
#### miller-rabin素性测试
#### crt
一种同构映射
1. 数论四大定理:威尔逊定理、欧拉定理、孙子定理、费马小定理
2. 威尔逊定理：当且仅当p为素数时,(p-1)! ≡ -1(mod p)
#### rsa 问题
1. 大数分解实验
2. 教科书式的 rsa 实验
3. 随机生成两个大素数并计算它们的乘积
4. 欧拉函数的计算、公私钥的计算、模幂、模逆
5. 扩展到 DDH 问题、key-exchange-protocol
6. 扩展到 elgamal 加密
7. 扩展到 Paillier
    - Zn^2 与 Zn
8. 扩展到盲签名blind-signature
#### paillier 算法原理
#### 二次剩余假设
1. 当N=pq的分解形式未知的时候，目前还没有多项式时间的算法判定某个x是否是模N的二次剩余。
2. Goldwasser-Micali 是基于二次剩余假设的对单比特加密的公钥加密方案。
3. Rabin 方案，在不知道N的分解形式时，求合数模N平方根是困难的，已知N的分解形式时，求合数模N平方根是简单的。
#### 计算素数模的平方根
1. 素数模平方根
2. 合数模平方根，比如N=pq;在不知道分解形式的条件下计算合数模的平方根是困难。当知道分解形式的时候计算平方根是简单的
3. 导致一个单向函数，在未知分解时，计算平方值很容易，但是计算平方根就很困难。