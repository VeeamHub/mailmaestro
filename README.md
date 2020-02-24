# Mail Maestro

**Author:** Timothy Dewin

**Function:** Self Service Recovery Demo

Demo for self service recovery with VBO365 1.5 API Calls. Will do the authentication with an LDAP source (e.g. active directory), and then allow the user to restore it's own mails back to it's original location. 

## 📗 Documentation

Please download, modify maestroconf.json to reflect your environment. Change VBOMailBox to modify the selected user you want to a restore for.

Passing parameters will be done in the following order > Config > Overwrite with commandline parameters. Passwords can be passed via CLI or Config but will be interactively asked if not supplied. Passwords can not be empty

You need the GO compiler (referred to golang), set a GOPATH (workdirectory for GO) and Git. Might depend on the OS. Both Centos and Windows were tested and work.

For example Centos:

```bash
yum install git -y
yum install golang -y
export GOPATH=/usr/share/gopath
mkdir $GOPATH
```

Refer to : [https://golang.org/doc/install](https://golang.org/doc/install)

Download and compile (will pull in the dependencies)

```bash
go get -d github.com/veeamhub/mailmaestro
go install github.com/veeamhub/mailmaestro
```

Copy the binary and config and edit to reflect your env (for example add SSL keys)

```bash
export MAILMAESTROPATH=/usr/share/mailmaestro
mkdir $MAILMAESTROPATH
cp $GOPATH/src/github.com/veeamhub/mailmaestro/maestroconf.json  $MAILMAESTROPATH
cp $GOPATH/bin/mailmaestro $MAILMAESTROPATH
```

Running the code

```text
$MAILMAESTROPATH/mailmaestro -config $MAILMAESTROPATH/maestroconf.json
```

Hopefully you get similar results

Running from console:
![MailMaestro run from console](https://github.com/VeeamHub/mailmaestro/raw/master/githubmedia/run-screenshot.png)

Login screen (login with AD/LDAP credentials):
![MailMaestro login screen](https://github.com/VeeamHub/mailmaestro/raw/master/githubmedia/run-login.png)

Restoring an item:
![MailMaestro restore an item](https://github.com/VeeamHub/mailmaestro/raw/master/githubmedia/run-restore.png)


### Using SSL

If you want to go more advanced, you can use https by using openssl to generate keys

```text
openssl genrsa 2048 > $MAILMAESTROPATH/private.pem
openssl req -new -x509 -key $MAILMAESTROPATH/private.pem -out $MAILMAESTROPATH/req.pem -days 3650
```

(If you are running this on windows, you can use [http://gnuwin32.sourceforge.net/packages/openssl.htm](http://gnuwin32.sourceforge.net/packages/openssl.htm). Make sure `OPENSSL\_CONF` is properly set eg : `set OPENSSL\_CONF=C:\d\openssl\share\openssl.cnf`)

You can then edit maestroconf.json, or just supply via the cmdline e.g.

```text
$MAILMAESTROPATH/mailmaestro -config $MAILMAESTROPATH/maestroconf.json -localkey $MAILMAESTROPATH/private.pem -localcert $MAILMAESTROPATH/req.pem
```

You can now browse to `https://myip:4123`. Notice that the port remains the same, if both key and cert is supplied, the server starts in SSL mode automatically

### Final notes

Since the API are in beta, this might not work with final version. Also, this is a demo, please consider reviewing and updating the code if you want to use this in production.

## ✍ Contributions

We welcome contributions from the community! We encourage you to create [issues](https://github.com/VeeamHub/mailmaestro/issues/new/choose) for Bugs & Feature Requests and submit Pull Requests. For more detailed information, refer to our [Contributing Guide](CONTRIBUTING.md).

## 🤝🏾 License

* [MIT License](LICENSE)

## 🤔 Questions

If you have any questions or something is unclear, please don't hesitate to [create an issue](https://github.com/VeeamHub/mailmaestro/issues/new/choose) and let us know!
