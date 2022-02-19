# Simple webcrawler for a website. 

## Usage: 
```
go run src/main.go -w <url> -d <depth of links>
```


### Dumps out the file into XML:

```
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
  </url>
  <url>
    <loc>http://www.example.com/dogs</loc>
  </url>
</urlset>
```
