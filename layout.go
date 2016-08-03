package main
func DefaultLayout() string {
    return `
    <html>
        <head>
            <title>{{.Title}}</title>
        </head>
        <body>
            {{.Content}}
        </body>
    </html>
    `
}