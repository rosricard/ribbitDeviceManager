note that you'll need to set the following env variables in your launch.json and/or trigger in your bash or docker environment depending on what your using to be able to run the application to run locally:
-dbconn. ie `export DSN_ENV="root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"`
-projectID
-apikey