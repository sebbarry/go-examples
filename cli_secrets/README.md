#Introduction to Crypto Packages in Golang.
> This application stores secret keys in an encrypted 'vault'
## Further understanding of CLI apps.
## Further understanding of Interface Chaining.



## Details: 
```
secret set <keyname> <keystring>
secret get <keyname> -key=<password>
# or
SECRET_KEY=<key> secret get <keyname>
```

```go
flag.String(encryptionKey)
var c secret.Client{
    EncryptionKey: "some-key"
}
c.Get("key") #returns "some-value"
```


```
./app -encryption_key="key"
```
