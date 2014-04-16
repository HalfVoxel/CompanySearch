import urllib2
import urllib
import re
from bs4 import BeautifulSoup
from bs4 import NavigableString
import sys
import codecs

letter = sys.argv[1]#raw_input("Letter: ")
start = int(sys.argv[2])
end = start + int(sys.argv[3])

print ("Spider on " + letter  +":" + str(start) +"..."+str(end))
i = start
index = 0
any = True
reg = re.compile('<div class="column1NoCheckbox">\s*<h5>\s*(.*?)\s*</h5>\s*<h5>\s*(.*?)\s*</h5>\s*</div>\s*<div class="column2">\s*<h5>\s*(.*?)\s*</h5>\s*<h5>\s*(.*?)\s*</h5>\s*</div>\s*<div class="columnLinks">\s*<h4>\s*<a href="([^"]*?)">Visa mer</a>\s*</h4>"' )
output = codecs.open("/Users/arong/web/python/scraping/output/"+letter+"_"+str(start)+".csv","w","utf-8")

while any and i < end:
	any = False
	resp = urllib2.urlopen("http://www.foretagsfakta.se/foretag-ao/"+urllib.quote_plus(letter)+"/"+str(i))
	s = resp.read()

	soup = BeautifulSoup(s)

	total = (end-start)#re.search("av\s*(\d+)",soup.find("div", class_="PageNavigator-Pages").string, flags=re.MULTILINE).group(1)
	#itot = int(total)
	print ("Searching... " + letter+":" + str(i) + " of " + str(total) + " ( " + str(float(i-start)/total) + " )")
	#print(soup)
	

	for li in soup.find_all('li'):
		div1 = li.find("div", class_="column1NoCheckbox")

		if div1 is None:
			continue


		any = True
		id = li['id']
		nameLink = li.find("div",class_="columnName").h2.a
		name = nameLink.string.strip()

		moreUrl = nameLink["href"]
		#print (id)
		#print (nameLink["href"])
		#print (name)

		col1 = li.find("div", class_="column1NoCheckbox")
		tel = col1.find_all("h5")[0].string.strip()
		city = col1.find_all("h5")[1].string.strip()

		col2 = li.find("div", class_="column2")

		url = ""
		conts = col2.find_all("h5")
		if not isinstance(conts[0],NavigableString) and conts[0].a is not None:
			url = conts[0].a["href"]

		email = ""
		if not isinstance(conts[1],NavigableString) and conts[1].a is not None:
			email = conts[1].a["href"]

		col3 = li.find("div", class_="columnLinks")
		conts = col3.find_all("h4")
		mapUrl = ""
		if not isinstance(conts[1],NavigableString) and conts[1].a is not None:
			mapUrl = conts[1].a["href"]


		output.write ( str(index) + "," + "\"" + name + "\", \"" + id + "\", \"" + tel + "\", \"" + city + "\", \"" + url + "\", \"" + moreUrl + "\", \"" + email + "\", \"" + mapUrl + "\"" + "\n")
		index += 1
		#print (city)
		#print (url)
		#print (email)
		#print (mapUrl)

		#a = [x.string.strip() for x in div1.find_all("h5")]
		#print (a)
		#print (x.contents for x in div1.find_all("h5"))
		#print (x.contents for x in li.find("div", class_="column2").find_all("h5"))
		#print (x.contents for x in li.find("div", class_="columnLinks").find_all("h5"))

	#for match in reg.finditer(s):
	#	print(match)
	i += 1

output.close()