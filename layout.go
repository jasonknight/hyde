package main
func DefaultLayout() string {
    return `
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <title>{{$title}}</title>
            <meta name="description" content="{{$metadescription}}">
            <meta charset="utf-8">
            <meta name="author" content="{{$author}}">
            <META NAME="ROBOTS" CONTENT="INDEX, FOLLOW">
            <!--[if lt IE 9]>
                <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
            <![endif]-->
            ript src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.0/jquery.min.js"></script>
            <link rel="stylesheet" href="https://ajax.googleapis.com/ajax/libs/jquerymobile/1.4.5/jquery.mobile.min.css">
            <script src="https://ajax.googleapis.com/ajax/libs/jquerymobile/1.4.5/jquery.mobile.min.js"></script>
        </head>
        <body>
            <h1>{{$title}}</h1>
            {{.Content}}
        </body>
    </html>
    `
}