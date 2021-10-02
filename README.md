# Introduction

My first `go` program uses to help my children to learn orthograph.
It is very basic and ugly (well... it's my first go program...).
In addition, I hate `HTML` so the the `HTML` content is also ugly :p

# Way of working

It exposes a web server on port 8080 and ask you to upload a file with a list of word. Then it use the browser voice API to say the word, and you have to correctly write it. Right now, it only works with french word, but it can easily be changed by updating the `javascript` file located in `www/js/golearn.js`.
