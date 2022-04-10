<body> 
<ol>
  {{ range . }}
  <li>
    <p><a href="/{{.Id}}" >电影名: {{ .Title }}</a></p>
  </li>
  {{ end }}
</ol>
</body>