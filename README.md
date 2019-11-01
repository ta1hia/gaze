# gaze

Chat server and client implemented using websockets. 

##### Start a server:
```
$ gaze serve --bind localhost:8844
```

##### Join a room as a client:
```
$ gaze connect --room catdog --nick alphonso
```
This creates the `catdog` room if it doesn't already exist, and connects you to the room using `nick` as your username.
