## Requirements
  - Golang 
  - Go Module
  - Postgresql Database
  - gorm

## Installing
  - Use Go Modules, please read https://blog.golang.org/using-go-modules
  - Install Dependecies
     ```console
     $ go mod tidy
     ```
  - Fill your local .env configuration
  - Run Project
     ```console
     $ go run .
     ```

## How to test
```
curl --location --request POST 'localhost:8080/upload' \
--header 'Cookie: gosessionid=ef529dfc-8104-4ad9-aa35-ec92f5a79930; gosessionid=ef529dfc-8104-4ad9-aa35-ec92f5a79930; gosessionid=ef529dfc-8104-4ad9-aa35-ec92f5a79930' \
--form 'upload-file=@"/path/to/file"'
```
```
curl --location --request POST 'localhost:8080/bulk-upload' \
--header 'Cookie: gosessionid=ef529dfc-8104-4ad9-aa35-ec92f5a79930; gosessionid=ef529dfc-8104-4ad9-aa35-ec92f5a79930; gosessionid=ef529dfc-8104-4ad9-aa35-ec92f5a79930' \
--form 'upload-file=@"/path/to/file"'
```

import Curl above to postman and use the .csv file in the folder csv-file
