curl https://localhost:8888/TicketType \
--cacert ./certificate/kaadda/ca.pem \
--cert ./certificate/kaadda/client.pem --cert-type PEM \
--key ./certificate/kaadda/client.key --key-type PEM

#--cacert ./certificate/kaadda/ca.pem



#openssl x509 -noout -modulus -in ./certificate/kaadda/client.pem | openssl md5
#openssl rsa -noout -modulus -in ./certificate/kaadda/client.key | openssl md5