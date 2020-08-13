import imaplib
import email
import re

email ="sfdsfsdfd@gmail.com"
password="sdsdfafsdf"
folderName="data"


#connection to the email server
mail = imaplib.IMAP4_SSL('imap.gmail.com')
mail.login(email,password)
mail.list()
# Out: list of "folders" aka labels in gmail.
mail.select(folderName, readonly=True) # connect to inbox.

result, data = mail.uid('search', None, "ALL") # search and return uids instead

ids = data[0] # data is a list.
id_list = ids.split() # ids is a space separated string
latest_email_uid = data[0].split()[-1]

result, data = mail.uid('fetch', latest_email_uid, '(RFC822)') # fetch the email headers and body (RFC822) for the given ID

raw_email = data[0][1] # here's the body, which is raw headers and html and body of the whole email
link_pattern = re.compile("https:\/\/kajabi-storefronts-production.s3.amazonaws.com(.*)(b>)")
search = link_pattern.search(raw_email.decode('utf-8'))
print (raw_email)
print ("***********************")