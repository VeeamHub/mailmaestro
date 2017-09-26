# Superedit
## VeeamHub
Veeamhub projects are community driven projects, and are not created by Veeam R&D nor validated by Veeam Q&A. They are maintained by community members which might be or not be Veeam employees. 

## Distributed under MIT license
Copyright (c) 2016 VeeamHub

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


## Project Notes
**Author:** Timothy Dewin

**Function:** Self Service Recovery Demo


Demo for self service recovery with VBO365 1.5 API Calls. Will do the authentication with an LDAP source (e.g. active directory), and then allow the user to restore it's own mails back to it's original location. Please download, modify maestroconf.json to reflect your environment. Change VBOMailBox to modify the selected user you want to a restore for.

Passing parameters will be done in the following order > Config > Overwrite with commandline parameters. Passwords can be passed via CLI or Config but will be interactively asked if not supplied. Passwords can not be empty


Since the API are in beta, this might not work with final version. Also, this is a demo, please consider reviewing and updating the code if you want to use this in production. 