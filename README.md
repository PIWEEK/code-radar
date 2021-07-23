<h1 align="center">
  <br>
  <img width="30"
       src="https://raw.githubusercontent.com/PIWEEK/code-radar/readme/front/public/favicon.svg" alt="LOGO">
  Code Radar
</h1>

**CodeRadar** is a GIT repository analyzer that retrieves information about file activity and file "ownership". 

![Code radar screenshot](https://github.com/PIWEEK/code-radar/blob/readme/resources/screenshots/penpot.png)

## Build and run

```bash
./build.sh
```

After the script is done a self-contained binary will be generated under `target/code-radar`

You can run the binary with:

```
code-radar <GIT repository URL>
```

Alternative you can run the application with a local repository with:

```
code-radar --local <Local repository path>
```

By default `code-radar` will choose a port open. You can force the port with the flag `--port`

```
code-radar --port 8000
```


## Development

For development you should run both client and server in different command lines.

### Client

```
cd front
npm install
npm run dev
```

### Server

```
cd back
go run cmd/code-radar/main.go
```

You alternative can use [Air](https://github.com/cosmtrek/air) to start a live-reload server.

```
cd back
air
```

## References

We took inspiration in the book [Your code as a crime scene](https://pragprog.com/titles/atcrime/your-code-as-a-crime-scene/) by [Adam Tornhill](https://www.adamtornhill.com/). If you're interested in the subject please throw your money at him.

## License

MIT License

Copyright (c) 2021 ΠWEEK

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
