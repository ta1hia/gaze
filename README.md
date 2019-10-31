# gaze

Chat server and client implemented using websockets. 

##### Start a server:
```
$ gaze serve --bind localhost:8844
```
##### Create a room:
```
$ curl -XPOST http://localhost:8844/catdog
```
TODO: move this into a cobra cmd

##### Join a room as a client:
```
$ gaze connect --room catdog --nick alphonso
```
