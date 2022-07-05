# Project Layout

Designing the project's layout

1. create a models folder so that we can store all the models structs
2. create a handlers folder with the handler.go file containing all the handlers handle any incoming HTTP requests
  - This code creates a Config{} struct with the MongoDB collection and context instances encapsulated. 
  - We must define a NewApp() function so that we can create an instance from the Config{} struct.
  - All handlers are a Config{} method. The handlers thus have access to all the variables of the struct such as the database connection because it is a method of the Config{} type.
3. From the main.go file, we'll provide all the database credentials and connect to the MongoDB server.
4. Then, we must create a global variable `app` to access the endpoints handlers.
5. Finally, we use the app variable to access the handler for each HTTP endpoint.