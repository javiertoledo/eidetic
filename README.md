eidetic (Work in Progress)
==========================

Yet another HTTP reverse proxy response recorder written in Go

The intent of this code is to speed up any third party API dependant testing processes acting as a reverse proxy and saving any response in a Redis database.

The main feature that make this different from other tools is that it is language/platform agnostic: Real HTTP calls are made without stubbing HTTP libraries. This makes it ideal for testing complex or multiplatform tools.

This is a Work in progress, but for me is already speeding up some tests. You can install it and give it a try with this commands:

    $ go get github.com/javiertoledo/eidetic
    $ eidetic http://whatever.your.api.dev.host.is:8080/with_some_base_route/

The `eidetic` command will hang your console and start listening on `http://localhost:8080` so it will intercept any incoming request and redirect it to the provided base URL.

To allow the cache work you should have a redis instance running and listening on `127.0.0.1:6379`. This can't be configured yet (Remember this is a WIP).

## Contributing

Your contributions or suggestions are welcome, just open an issue or send a pull request!

## License

The MIT License (MIT)

Copyright (c) 2014 Javier Toledo <javier@theagilemonkeys.com>

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

