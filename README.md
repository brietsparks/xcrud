This is an example CRUD data layer with supporting concerns
such as testing, migrations, errors/logging, and a CLI.  

## Usage
The data layer can be accessed via a standalone CLI or via Go code.

### CLI
1. clone the repo: ```git clone https://github.com/brietsparks/xcrud```

2. in the repo, run: ```go install``` (creates a binary in GOPATH)
    
3. in the repo, create a .env file with the variables for the db connection:

    ```
    DB_HOST=my-db-host.com
    DB_USER=postgres
    DB_PORT=5432
    DB_NAME=database_name
    DB_PASSWORD=password1234
    ```

4. Now you can run commands:

    **Apply schema to the database:**

    ```
    xcrud migrate up
    ```
    
    **Create a user:**
    
    ```
    xcrud resources user:create --FirstName Bo --LastName Peep
    ```
    
    Output: `{"id":1,"firstName":"Bo","lastName":"Peep"}`
    
    **Get a user:**
    
    ```
    xcrud resources user:get 1
    ```

    Output: `{"id":1,"firstName":"Bo","lastName":"Peep"}`

    **Update a user:**
    
    ```
    xcrud resources user:update 1 --LastName Jackson
    ```
    
    **Create a group:**
    
    ```
    xcrud resources group:create --Name groupA
    ```
    
    Output: `{"id":1,"name":"groupA"}`
    
    **Add a user to a group:**
    
    ```
    xcrud resources group:add-user --GroupId 1 --UserId 1
    ```
    
    **Get users by group ID:**
    
    ```
    xcrud resources users:get --GroupId 1
    ```
    
    Output: `[{"id":1,"firstName":"Bo","lastName":"Peep"}]`

    **Get groups by user ID:**
       
    ```
    xcrud resources groups:get --UserId 1
    ```
    
    Output: `[{"id":1,"name":"groupA"}]`
    
    **Remove a user from a group:**
    
    ```
    xcrud resources group:remove-user --GroupId 1 --UserId 1
    ```
    
### Go

1. install: ```go get -u github.com/brietsparks/xcrud```

2. see [example usage code](https://github.com/brietsparks/xcrud/blob/master/example/example.go)
