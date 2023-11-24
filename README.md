# AWS Lambda response streaming

> At time of im writing this aws gateway does not support response streaming

I was trying to implement golang response streaming on golang but i hit a multiple walls face first. On the official documentation it shows how to build it from the scratch using the runtime API. On golang i wasnt able to find a way to write an http 1.1 request from client to to server option using transfer encoding chunk. But after reading through the all of the documentation i was not able to implement it. If you know how to implement it please let me know i would love to talk about this topic and networking in general.

<br>
 So i will talk about how to implement it in 2 ways but runtime option is which i failed.

## AWS lambda runtime API - failed option its just a little bit of my journey how i started first

So how aws lambda works in general. You have your application and lambda starts up a isolated container for your lambda function if your using the aws provided container your code needs to have bootstrap executable which is the starting trigger for your function if its a container that you made yourself its free of choice. Btw `{AWS_LAMBDA_RUNTIME_API}` is in the enviromnet variable. documentation is in [here.](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html) <br>

- `/2018-06-01/runtime/invocation/next` get the request context header contains the request id body includes the params
- Do some processing the request
- `2018-06-01/runtime/invocation/{requestID}/response` to send the response to back to the runtime api

If you want more detail its best read the documenation.<br>

Now the streaming part here everything goes to chaos there is few things to do use the runtime response.
When sending the response you need to set these headers `Lambda-Runtime-Function-Response-Mode` to `streaming`, `Transfer-Encoding` to `chunked` and make it http 1.1 compliant but i could not find a way to send it in one request as streaming no configuration on the console to set any variable to streaming only on the function URL. if you find a way to make it work please let me know my email is : puntsag0609@gmail.com

## Using aws lambda web adabter

the original github link is [here.](https://github.com/awslabs/aws-lambda-web-adapter) They have some awesome documentation here i just played around the fields little bit. <br>

So how does this work? What it does is its kind works as a proxy between the aws lambda runtime and your web server. also this makes more testable than the node js that aws built in method because now you have an option to run in your local machine and use it similar to normal way. Benefits are many but felt little bit slower than using the runtime option but thats a trade im willing to make in this case. <br>

**The things you need make it work**

You will need to include the docker image in your docker file.

```
    // basically this line here
    COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.7.1 /lambda-adapter /opt/extensions/lambda-adapter
```

also on top of this set on your `template.yaml` under the resource properties

```
    Environment:
            Variables:
            AWS_LWA_INVOKE_MODE: RESPONSE_STREAM
```

these are the things you need to set the basic needs for it to run but i highly encourage you to play around more with the fields to get more about it also take a look at the original github repo they have some really good examples.

```
 // to deploy
    sam build
    sam deploy -g

```
