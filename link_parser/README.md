# Link parser

This package needs some understanding of graph theory and trees.
Remember, html files are just trees... 

### Take: 

```
<a href="/dog">
  <span>Something in a span</span>
  Text not in a span
  <b>Bold text!</b>
</a>
```

### Output: 

```
Link{
  Href: "/dog",
  Text: "Something in a span Text not in a span Bold text!",
}
```
