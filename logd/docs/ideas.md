# Logd
  Logd provides a simple but extendable logging library using golangs logger as the fundation upon which it builds ontop.

## API
  The **Logd** api is rather simple and combines a few ideas to provide the flexibility it needs to match the use cases that it might be used for. For more expressive and useful purposes, most messages are stored as `LogdReport` to allow a richer content for logging formatters and backends

  - Formatters
    To provide high level of flexibility as to the output representation of Logd information, formatters provide a customer approach in transforming log information into the the appropriate format that makes sense. Formatters are more in the domain left to use developer to create or using the deault formatter that come within logd which relies basically on fmt for some go formatting goodness

    - Formatters API

      - `Format(*LogdReport,io.Writer)`
        providing a single function which takes the message and the necessary information required by the formatter,it transforms the input into an appropriate representation and feeds into the given writer


  - Backends
    To provide the highest level of flexibility, Logd contains the concept of backends which provide the output layers where Logd sends out it's reports to, multiple backends can be including to allow a sharding of log information to multiple endpoints (eg as a means of backup)

    - Backend API
      The backend api is rather simple and its just logs information out

      - `Write(*LogdReport)`


    ```go

      //creating a backend is simple,Logd provides a set of basic backends

      //file-based backend
      fileLog := logd.NewFileBackend("./logs/debug.logs",logd.DefautFormatter)

      //redisdb -based backend
      redisLog := log.NewRedisBackend(&RedisConfig{ Addr:"redis://:300",User: "wackom"},logd.DefaultFormatter)

      //standard output backend
      stdOut := logd.NewStreamBackend(os.Stdout,logd.DefaultFormatter)
      stderr := logd.NewStreamBackend(os.Stderr,logd.DefaultFormatter)

      //use the default backend
      logOut := logd.DefaultBackend() //defaults to a os.Stdout backend


      logOut.Write(&logd.LogdReport{
        Context: &Levels,
        File: ...,
        Line: ...,
        Message: ...,
        Args: []interface{},
      })

    ```


  - Levels/Context
    Levels or Context in Logd defines the difference set of accepted log information that would be coming out from an application logs, the levels actually have tied into them the backend mechanics which allows filtering out specific types of log information into separate buckets

    ```go

      debug := logd.NewContext({
        UID: 3,
        Name: "app.Debug",
        Backends: []Backends{stdOut,stderr},
      })

      //adds the backend instance into the list of current backends
      debug.AddBackend(fileLog)

      user := logd.NewContext({
        UID: 23,
        Name: "app.User",
        Backends: []Backends{stdOut,stderr},
      })

    ```


  -  Loggers
    Logd loggers combine all the previous features to provide a flexible and usable logging system versatile enough to handle the demand the loggers need to meet when in use with applications

    ```go

      log := logd.NewLogger()

      //add up the context we are using
      log.AddContext(debug)
      log.AddContext(user)

      log.Log("app.")
    ```
