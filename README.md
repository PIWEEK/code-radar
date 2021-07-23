<h1 align="center">
  <br>
  <img width="30"
       src="https://raw.githubusercontent.com/PIWEEK/code-radar/main/front/public/favicon.svg" alt="LOGO">
  Code Radar
</h1>

**CodeRadar** is a GIT repository analyzer that retrieves information about file activity and file "ownership". 

![Code radar screenshot](https://github.com/PIWEEK/code-radar/blob/main/resources/screenshots/penpot.png)

## Building

```bash
./build.sh
```

After the script is done a self-contained binary will be generated under `target/code-radar`

## Running

If your current path is a GIT repository just running `code-radar` will analyze de current's path repository.

```
code-radar
```

Alternatively, you can run the application with a remote repository. The application will clone this into memory (so no artifacts will be created in your system).

```
code-radar <GIT repository URL>
```

Another option is running the application with a local repository but in a different path location:

```
code-radar --local <Local repository path>
```

By default `code-radar` will choose a free open port. You can force the port with the flag `--port`

```
code-radar --port 8000
```

After the repository is analyzed you can navigate to the URL shown in the terminal to check the results.

```bash
2021/07/23 11:12:42 Init local repository .
2021/07/23 11:12:42 [START] Processing repository
2021/07/23 11:15:09 [END] Processing repository
2021/07/23 11:15:09 Listening on: http://localhost:44103
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
