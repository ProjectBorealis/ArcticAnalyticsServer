# ArcticAnalytics-Server

Basic server for saving performance analytics to a file.

### Dependencies
```
npm install express basic-auth csv-writer ajv dotenv sqlite3
```

### Run Server
```
node main.js
```

#### Notes
* Authorization timing is not working properly.

### Endpoints
* POST /v1/user/performance
* GET /v1/user/performance/csv
* GET /example (basic-auth example)