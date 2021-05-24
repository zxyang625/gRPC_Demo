rm *.pem

# 1.生成CA的私钥和自签名证书
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "//C=CH/ST=SiChuan/L=ChengDu/O=Study/OU=Go/CN=Tony/emailAddress=Tony@email.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2.生成服务器的私钥和证书签名请求(CSR) #-nodes表示.pem不加密
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "//C=CH/ST=GuangDong/L=GuangZhou/O=Computer/OU=Go/CN=Harry/emailAddress=Harry@email.com"

# 3.使用CA的私钥签署服务器的CSR并生成签名证书
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's certificate signed by CA"
openssl x509 -in server-cert.pem -noout -text

echo "Verify ca-cert.pem server-cert.pem"
openssl verify -CAfile ca-cert.pem server-cert.pem

# 4.生成客户端的私钥和证书签名请求(CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "//C=CH/ST=FuJian/L=XiaMen/O=Client/OU=Gopher/CN=Alice/emailAddress=Alice@email.com"

# 5.用CA的私钥签署客户端的CSR并生成签名证书
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf

echo "client's certificate signed by CA"
openssl x509 -in client-cert.pem -noout -text

echo "Verify ca-cert.pem client-cert.pem"
openssl verify -CAfile ca-cert.pem client-cert.pem

