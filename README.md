一、完成的功能：基础功能: 

1、文件上传（支持多文件上传）和下载
a
2、文件分享

进阶功能:

1、登录注册
2、文件权限管理a
3、加密分享链接
4、一次性快传
                     

接口文档：

1、GET     注册接口：                     /register/:userID/:password

先select数据库user表查询是否有userID，没有则注册成功



2、GET      登录接口：                    /login/:userID/:password

登录成功后，会创建一个键为ID，值为userID的session，用于权限管理的一部分

同时也会创建一个以userID命名的文件夹作为用户专属的网盘存储空间

也创建了一个可视化的网盘目录 /userID



3、GET       加密分享接口：           /share/:file/:date   

参数 file     分享的文件名

参数 date  设置过期天数

会返回一个分享文件链接，和一个由四位随机的大小写字母和数字构成的提取码（参照百度网盘）



4、GET        重定向接口---满足条件则跳转到下载的接口：/sharing/:afterParams/:code   

 由（3、得到的链接后面跟上  /提取码） 就可以跳转到下载了

这里需要满足 提取码正确和在链接设置的有效时间内才可以下载



5、POST     上传一个文件的接口：/upLoad

参数 file  ，上传之前用户登录的时候创建的文件夹（以用户的userID命名的文件夹）



6、POST     上传多个文件的接口：/upLoads

参数 file[]



7、GET        下载文件的接口:         /downLoad/:fileName?code=  &ownerID= 

参数 fileName  需要下载的文件名

参数 code          提取码

参数 ownerID    该文件归属的用户的userID，下载别人网盘的文件才需要加上

这个接口 fileName必须，通过code和ownerID参数有无的判断来区分这是一个

加密的分享链接(含提取码)下载、从自己的网盘上下载和从别人的网盘上下载（这里有别人文件

上的权限控制） 

分享链接下载（这里一般是4、重定向接口跳转过来的），提取码错误不能下载，



8、POST      一次性快传接口：     /oneupload

参数 file   

上传后获得一个一次性下载链接的地址



9、GET         快传后下载的接口： /oneDown/:ownerID/:fileName

参数不用管，直接对应 上面得到的一次性下载链接即可跳转到下载或失效

该链接只能下载一次，之后失效



10、GET       用户更改自己云盘文件权限的接口 ：  /cFilePermit/:fileName/:public

参数 fileName 是上传了的文件名

public  0 代表公开   1 代表 私有，其他人不可以下载



​      







                          

