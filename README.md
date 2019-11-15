# gaze

Chat server and client implemented using websockets. 

##### Start a server:
```
$ gaze serve --bind localhost:8844
```

##### Connect as a client:
```
$ gaze connect 
```
Once you're connected, you can join a room:
```
/join catdog 
```
This creates the `catdog` room if it doesn't already exist and connects you to the room. 
