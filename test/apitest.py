# integration test to ensure API behaves according to spec
# requires requests python package: run "pip install requests"

# by default, assumes server is running on localhost:8080
# specify different target url in command line args
# ex: python3 apitest.py http://localhost:5000

import requests
import sys
from datetime import datetime
import dateutil.parser

url = "http://localhost:8080"
username = "admin"
password = "back-challenge"

if (len(sys.argv) >= 2):
    url = sys.argv[1]

#tests return error if they fail. If they pass they return None

def createMessage(name, email, text):
    message = {'name' : name,
               'email' : email,
               'text' : text}
    r = requests.post(url=url+"/messages", json=message)
    return r

#id is a string uuid
def getMessage(id):
    r = requests.get(url=url+"/messages/"+id, auth=(username, password))
    return r

def getMessages():
    r = requests.get(url=url+"/messages", auth=(username, password))
    return r

def changeMessage(id, text):
    message = {'text' : text}
    r = requests.put(url=url+"/messages/"+id, json=message, auth=(username, password))
    return r

def testPost():
    message = {'name':'testuser1',
            'email':'testemail1',
            'text':'test1'}
    r = requests.post(url=url+"/messages", json=message)
    if (r.status_code != requests.codes.created):
        return "expected status created, recieved code "+str(r.status_code)

def checkCode(r, expectedCode):
    if(r.status_code != expectedCode):
        return "expected status code "+str(expectedCode)+", recieved status code "+str(r.status_code)

#test create new message
def testCreateAndGetMessage():
    r = createMessage('testname','testemail','test text')
    err = checkCode(r, requests.codes.created)
    if err != None:
        return err
    originalMessage = r.json()

    r = getMessage(originalMessage['id'])
    err = checkCode(r, requests.codes.ok)
    retrievedMessage = r.json()
    if err != None:
        return err
    fields = ['id', 'name', 'email', 'text', 'creation_time']
    for x in range(len(fields)):
        if(originalMessage[fields[x]] != retrievedMessage[fields[x]]):
            return "original message "+fields[x]+"="+originalMessage[fields[x]]+" but get request retrieved "+retrievedMessage[fields[x]]

#test list all messages anti-chronologically
def testGetMessages():
    for x in range(10):
        r = createMessage('testname'+str(x), 'testemail'+str(x), 'test text '+str(x))
        err = checkCode(r, requests.codes.created)
        if err != None:
            return err
    r = getMessages()
    messageList = r.json()
    for x in range(1, len(messageList)):
        t0 = dateutil.parser.parse(messageList[x-1]['creation_time'])
        t1 = dateutil.parser.parse(messageList[x]['creation_time'])
        if(t0 < t1):
            return "message :"+str(messageList[x-1])+" located before message: "+str(messageList[x])

#test update the text of message by ID
def testUpdateMessage():
    r = createMessage('testname','testemail','test text')
    err = checkCode(r, requests.codes.created)
    if err != None:
        return err
    originalMessage = r.json()

    r = changeMessage(originalMessage['id'], 'update test text')
    err = checkCode(r, requests.codes.ok)
    if err != None:
        return err
    updatedMessage = r.json()

    r = getMessage(originalMessage['id'])
    err = checkCode(r, requests.codes.ok)
    retrievedMessage = r.json()
    if err != None:
        return err
    if(originalMessage['id'] != retrievedMessage['id']):
        return "updated message id does not match original"
    if(originalMessage['name'] != retrievedMessage['name']):
        return "updated message name does not match original"
    if(originalMessage['email'] != retrievedMessage['email']):
        return "updated message email does not match original"
    if(retrievedMessage['text'] != 'update test text'):
        return "message text not updated properly: expected update test text, got "+retrievedMessage['text']
    if(originalMessage['creation_time'] != retrievedMessage['creation_time']):
        return "updated message creation time does not match original"

#test authentication
def testAuthentication():
    r = createMessage('testname','testemail','test text')
    err = checkCode(r, requests.codes.created)
    if err != None:
        return err
    originalMessage = r.json()

    r = requests.get(url=url+"/messages/"+originalMessage['id'])
    err = checkCode(r, requests.codes.unauthorized)
    if err != None:
        return err

    r = requests.get(url=url+"/messages")
    err = checkCode(r, requests.codes.unauthorized)
    if err != None:
        return err

    message = {'text' : 'hahaha im changing other users messages'}
    r = requests.put(url=url+"/messages/"+originalMessage['id'])
    err = checkCode(r, requests.codes.unauthorized)
    if err != None:
        return err

def main():
    tests = [
                ['testCreateMessage', testCreateAndGetMessage],
                ['testGetMessages', testGetMessages],
                ['testUpdateMessage', testUpdateMessage],
                ['testAuthentication', testAuthentication],
            ]
    for x in range(len(tests)):
        err = tests[x][1]()
        if (err != None):
            print(tests[x][0]+" failed")
            print(err)
            return
        else:
            print(tests[x][0]+" passed")

main()
