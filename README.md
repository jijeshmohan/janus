##Janus 

Janus is a fake rest api server which can be used for various purpose including frontend application development , testing etc.

### Features

* Easy to configure and use.
* Supports all http methods and REST resources.
* Supports CORS.
* Available in all platforms.


### Install

##### Get Binary Distribution

You can download the binary distribution from [here](https://github.com/jijeshmohan/janus/releases/tag/1.0.0)

##### Build from source

Make sure that the golang is installed in your system

```go get github.com/jijeshmohan/janus```

### How to run

To run janus, go to any directory and create a json file called ```config.json``` . This defines the REST endpoints which server need to expose. ( config.json file described in the next section.). Run the application in the same directory by typing `janus` in the terminal
 
### Basic Structure of a config.json file

The basic structure of config file shown below. A config file has follwing attributes
* port
* resources 
* urls

e.g config.json

```json
{
    "port": 8080, 
	"resources":[
		{
        	"name": "user",
       	 "headers": {
				"key": "value",
			}
        }

    ],
    "urls":[
    	{
        	"url": "direct/url/for/something",
       	    "method": "GET",
       	    "content_type": "application/json",
       	    "status": 200,
       	    "file": "./files/some.json",
       	    "headers": {
            	"key": "value",
                "key1": "value1"
            }
        }
    ]

}
```

**port** in configuration file defines in which port the server needs to run.This is an optional field and if not provided the default port is 8000.

You can define two types of REST endpoints in the configuration file.  

* Resources
* Urls

##### Resources 

This represent basic REST resource which will exposes all [standard methods](http://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods). Janus will look for a folder with the name of the resource in the same directory as routes.json for sending the data correspoding to the methods.

e.g:

```js
{
	"name": "user"
	"headers": {  // headers field is optional
		"key": "value"
	}
}
```

e.g for user resource , it wil look follwoing files to send the data

| Verb | Url | FILE | 
|--------|--------|---|
| GET       | /user       |  ./user/index.json  | if the file not available app will send a 404|
| POST       | /user       |  ./user/post.json  | if file is present it will send 201 with the content of the file otherwise 404|
|GET| /user/:item |./user/[any file name which is matching :item].json | if the file present, it will send the file with 200, otherwise 404. you can add item1.json, item2.json etc if you want to fake different get request
| PUT | /user/:item |  will use the same file specified above | if file is present it will send 200 with the content of the file otherwise 404. You can specify any number of files to match the get requst like  user1.json, admin.json etc.
| PATCH | /user/:item | will use the same file specified above | if file is present it will send 200 with the content of the file otherwise 404.
| DELETE | /user/:item | will use the same file specified above | if file is present it will send 200 with the content of the file otherwise 405. If you have specified item1.json, then it will send 200 for that particular item's request.

You can specify any header informations which need to send along with the methods. This is optional field

##### Urls
Urls section is for specifying individual urls which can't qualify for a standard REST resource methods.

This section gives more freedom in terms for defining HTTP Methods and content type and files. 

A single url representasion showed below

```js
 {
   "url": "/admin/user/enable", //mandatory 
   "method": "GET", // mandatory 
   "content_type": "application/json", // optional. default to application/json; charset=utf-8
   "status": 200, // optional , default to 200
   "file": "./files/some.json" // optional, if not specified , the response will be empty string. if specified it should be a valid file.
   "headers": {  //Optional
    "key": "value",
    "key1": "value1"
   }
}
```

Note: Url also support dynamic url like ```/admin/{user}/enable``` which can match to any  user name like ```/admin/user1/enable``` or ```/admin/anotheruser/enable``` .

### TODO

* Different authentication supports.
* Admin UI for configuration and adding url at runtime.
* Logging

### License

The MIT License.
