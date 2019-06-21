Here is my CHUCKNORRIS application! It utilizes two different web APIs to create a web service that fetches a random name (http://uinames.com/api/) and replaces Chuck Norris'
 name in a joke (http://api.icndb.com/jokes/random?firstName=John&lastName=Doe&limitTo=[nerdy]) with the random names.

 All that you have to do to use the service is run it (`go run main.go` from the root directory), open a new console and enter `curl http://localhost:5000`. You can also open an internet browser and navigate to `http://localhost:5000`.
 This will automatically return the corrected joke from the service into the console that you entered the `curl` command, or onto the webpage of the internet browser. Enter this command as many times as you'd like to fetch different jokes.

  In order to stop the application, press Command+C (Cancel) in the terminal that you ran the `go run main.go` command.

 This application was intended to be configuration-based, so each service is included in the .configs folder (config-e0.yaml) with configurable options. At this time, the only options are the URL and request timeout duration. Viper was used to read in the configurations by the application.

 Due to time constraints, there is currently no test coverage, but unit tests will be added shortly.

 ***If a request returns any status code other than 200, the application will wait 3 seconds and then retry. The request will be retried two times before finally returning the error. If it is a response code error, this will not kill the server or application, but application level errors will stop it.***

 I have also included a graceful shutdown for the application, along with a graceful shutdown of the server. The configurations are read and set in `main.go`, while all of the http handling is done in `jokeHandler()`. This prevents a new processor being created with each request, increasing performance issues.

 The nature of a gorilla mux router is to handle each request in parallel, so concurrency should not be an issue.

If you have any questions or comments, please feel free to contact me:
---------------------------------------------------------------------
 Developer: Stevan Cunningham
 Cell:      813-751-4088
 Email:     stevanc08@gmail.com
 --------------------------------------------------------------------