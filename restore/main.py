import urllib.request
import urllib.parse
import threading
import base64
import time
import json
import os


def archive_urls_to_json_file(host: str) -> None:
    """
    Save all known urls from a host on archive.org in a json file.
    Save local because of a slow Archive Server.
    2012 - 2013 are the last valid entries for these urls.

        Parameters:
            host (str): A host for getting the archive urls.

        Returns:
            None
    """

    base = "https://web.archive.org/cdx/search/cdx?"
    params = {
        "url": host,
        "matchType": "domain",
        "filters": "statuscode:200",
        "output": "json",
        "from": "2012",
        "to": "2013"
    }
    url = base + urllib.parse.urlencode(params)
    req = urllib.request.urlopen(url)
    res = req.read()
    j = json.loads(res.decode("utf-8"))
    with open('urls_to_restore.json', 'w', encoding='utf-8') as f:
        json.dump(j, f)
    print(f"Saved all URLs in File. URL Count: {len(j)}")


def download(timestamp, original, mimetype):
    base = f'http://web.archive.org/web/{timestamp}if_/{original}'
    encode_url = base64.b64encode(original.encode())
    remaining_download_tries = 20
    while remaining_download_tries > 0:
        try:
            urllib.request.urlretrieve(
                base, f'data/{mimetype}/{encode_url.decode("utf-8")}.{timestamp}.{mimetype.split("/")[-1]}')
            print(f"successfully downloaded: {original}")
            time.sleep(0.5)
        except:
            print(f"error downloading: {original}")
            time.sleep(1)
            remaining_download_tries -= 1
            continue


def mimetypes(lst):
    mimetypes = list(set([i[3] for i in lst]))
    for m in mimetypes:
        if not os.path.exists(f'data/{m}'):
            os.makedirs(f'data/{m}')


def partition(lst):
    unique_urls = list(set([i[2] for i in lst]))
    for url in unique_urls:
        yield filter(lambda i: i[2] == url, lst)


def main():
    archive_urls_to_json_file("xxl-angeln.de")
    # for i in lst[1:11]:
    #    print(i)
    #    threading.Thread(target=download, args=(i[1], i[2], i[3], )).start()


if __name__ == "__main__":
    # with open('archive_urls.json') as f:
    #    lst = json.load(f)
    # archive_urls()
    # mimetypes()
    # for i in partition(lst):
    #    print(list(i))
    main()


#
# def images():
#  req = requests.get('http://web.archive.org/cdx/search/cdx?url=xxl-angeln.de&matchType=host&from=2010&to=2013&filter=statuscode:200&filter=mimetype:image/.*&output=json&collapse=digest')
#  res = req.json()
#
#  #base = f'http://web.archive.org/web/{date}/{url}'
#
#  for i in res[1:]:
#    url = f'http://web.archive.org/web/{i[1]}if_/{i[2]}'
#    urllib.request.urlretrieve(url, f"images/{i[1]}-{url.split('/')[-1]}")
#
#
# def fangmeldung():
#  req = requests.get('http://web.archive.org/cdx/search/cdx?url=http://www.xxl-angeln.de/angel_praxis/fangmeldungen/*&filter=statuscode:200&filter=mimetype:text/html&output=json&fl=original,timestamp&collapse=urlkey')
#  res = req.json()
#  print(len(res))
#  for i in res[1:]:
#    if 'offset' not in i[0]:
#      if 'species' not in i[0]:
#        fname = i[0].replace('http://www.xxl-angeln.de','').replace('http://xxl-angeln.de','').replace(':80','').replace('/','-').replace('-angel_praxis-fangmeldungen-','')
#        url = f'http://web.archive.org/web/{i[1]}if_/{i[0]}'
#        try:
#          urllib.request.urlretrieve(url, f"html/{fname}-{i[1]}.html")
#        except:
#          print(url)
#    #break
#
# def parse():
#  all = []
#  for fname in os.listdir('fangmeldungen')[:250]:
#    data = {}
#    with open(f'fangmeldungen/{fname}', encoding='UTF-8') as raw:
#      soup = BeautifulSoup(raw, "html.parser")
#
#      if "Du hast keine Berechtigung, diesen Beitrag einzusehen." not in raw.read():
#
#        posted = soup.find("p", class_="posted")
#        if posted:
#          data['posted'] = posted.find("strong").get_text()
#          data['from'] = posted.find("a").get_text()
#          try:
#            data['id'] = posted.find("a")['href'].split('(id)/')[1]
#          except:
#            pass
#
#        rating = soup.find("li", class_="current-rating")
#        if rating:
#          data['rating'] = rating.get_text().split(' ')[1]
#
#        image = soup.find("a", class_="thickbox", href=True)
#        if image:
#          data['image'] = image['href']
#
#        views = soup.find("div", class_="lineLinks")
#        if views:
#          if views.find('b'):
#            data['views'] = views.find('b').get_text()
#
#        description = soup.find("div", class_="content_tabcontent")
#        if description:
#          if description.find_next_sibling('p'):
#            data['description'] = description.get_text()
#
#        title = soup.find("h1", class_="title")
#        if title:
#          data['title'] = title.get_text()
#
#
#        for part in soup.select('div[class*="linedata"]'):
#          i = re.sub(r"\s\s+", '', part.get_text()).replace('\n','').replace(' :',':').split(':',1)
#          data[i[0]] = i[1]
#        if "Köder" not in data:
#          data['Köder'] = "Nicht vorhanden"
#        all.append(data)
#
#  with open('test.csv', 'w', newline='', encoding="UTF-8") as f:
#    w = csv.DictWriter(f, all[0].keys())
#    w.writeheader()
#    for i in all:
#      w.writerow(i)
#
#
#
#
# def profile():
#  req = requests.get('http://web.archive.org/cdx/search/cdx?url=http://www.xxl-angeln.de/content/view/profile/*&filter=statuscode:200&filter=mimetype:text/html&output=json&fl=original,timestamp&collapse=urlkey')
#  res = req.json()
#  print(len(res))
#  for i in res[1:]:
#    fname = i[0].replace('http://www.xxl-angeln.de','').replace('http://xxl-angeln.de','').replace(':80','').replace('/','-').replace('-content-view-profile-','')
#    url = f'http://web.archive.org/web/{i[1]}if_/{i[0]}'
#    try:
#      urllib.request.urlretrieve(url, f"profile/{fname}-{i[1]}.html")
#    except:
#      print(url)
#
# fangmeldung()
# parse()
# profile()
