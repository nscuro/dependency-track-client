<html>

<head>
    <title>Project Report</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css"
        integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous">
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js"
        integrity="sha384-OgVRvuATP1z7JjHLkuOU7Xw704+h835Lr+6QL9UvYjZE3Ipu6Tp75j7Bh/kR0JKI"
        crossorigin="anonymous"></script>
</head>

<body>
    <div class="container-fluid">
        <h1>Executive Summary</h1>
        <div>
            The project dependends on a total of {{.Components | len}} components. {{.Findings | len}} vulnerabilities
            have been identified.
        </div>
        <h1>Components</h1>
        <table class="table">
            <thead>
                <tr>
                    <th scope="col">Name</th>
                    <th scope="col">Group</th>
                    <th scope="col">Version</th>
                </tr>
            </thead>
            <tbody>
                {{range .Components}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.Group}}</td>
                    <td>{{.Version}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>

</html>