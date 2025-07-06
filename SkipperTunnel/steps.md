1. Create flags and validate them

2. Create a ping or stbalish a connection with the localhost provided by the flag

    Make that via ping with timing, so it goes every time.Duration pinging the localhost

3. Stablish connection with skipper


4. Send the request_tunnel creation to the proxy


5. See if the response is acknowledge, if not show the error. we could type the erros for example.
    Type the response as: 
        Error connecting to skipperProxy, try in a bit
        Subdomain already in used
        Something else could be?


6. If it was acknowledge, show the user that he can now use his subdomain.skipper.lat

7. allow communication with the http client, etc. Reactor pattern taht redirects to goroutines

8. goruotines must receive the data, serilize the  request, make the request, recieve the reponse, serilaize he reponsse. Create a header with the response, send it to the connection.


9. Gracefull shutdown is a must

10. in the case of the gouroutines for now i think the only ones that need to be running are, the goroutine with the PING to the localhost. Cloud be another to Ping the connection to the Proxy. and the workerpool




# more:
    Show a type of dashboard, mount a server web or something showing the requests made
