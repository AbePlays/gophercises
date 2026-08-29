package constants

var AdventureTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8" >
        <title>Choose Your Own Adventure</title>
        <style>
            *,*::before,*::after{box-sizing:border-box}html{font-size:16px;scroll-behavior:smooth}body{margin:0;padding:1.5rem;max-width:65ch;margin-inline:auto;font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;line-height:1.6;color:#1a1a1a;background:#fafafa}h1,h2,h3,h4,h5,h6{line-height:1.25;margin:1.5em 0 .5em;color:#111}h1{font-size:2.2rem}h2{font-size:1.7rem}h3{font-size:1.35rem}h4{font-size:1.15rem}p{margin:0 0 1.2em}a{color:#06c;text-decoration:underline;text-underline-offset:2px}a:hover{color:#049}ul,ol{margin:0 0 1.2em;padding-left:1.5em}li{margin-bottom:.4em}img{max-width:100%;height:auto;display:block}code{font-family:ui-monospace,"Cascadia Code","Source Code Pro",Menlo,Consolas,monospac
        </style>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}

        <ul>
            {{range .Options}}
                <li>
                    <a href="/{{.Chapter}}">{{.Text}}</a>
                </li>
            {{end}}
        </ul>
    </body>
</html>
`
