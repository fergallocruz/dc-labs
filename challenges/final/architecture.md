Architecture Document
=====================
Contributors: María Fernanda Gallo Cruz & Sebastián Woolfolk

---------> API ENDPOINTS: <---------------

-   /login
    -
    curl --location --request GET 'http://localhost:8080/login?user=fer&password=gallo' --HEADER 'Authorization: application/x-www-form-urlencoded'
-   /workloads/test
    - 
    curl --location --request GET 'http://localhost:8080/workloads/test' --header 'Authorization: Bearer ZmVyOmdhbGxv'
-   /workloads/filter
    -
    
    curl --location --request POST 'http://localhost:8080/workloads/filter?workload-id=wii-filter&filter=bw' --header 'Authorization: Bearer ZmVyOmdhbGxv' --form 'data=@/Users/mariafernandagallocruz/Downloads/ice_cream.png'
-   /status
    -
    curl --location --request GET 'http://localhost:8080/status' --header 'Authorization: Bearer ZmVyOmdhbGxv'
    curl --location --request GET 'http://localhost:8080/status/WORKERNAME' --header 'Authorization: Bearer ZmVyOmdhbGxv'
-   /logout
    -
    curl --location --request GET 'http://localhost:8080/logout' --header 'content-type:multipart/form-data' --header 'Authorization: Bearer ZmVyOmdhbGxv'
- /results/:workloadsID
    -
    curl --location --request GET 'http://localhost:8080/results/my-filter'  --header 'Authorization: Bearer ZmVyOmdhbGxv'
