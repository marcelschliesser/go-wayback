import urllib.request
import urllib.parse
import logging
import base64
import json
import time
import os

logging.basicConfig(format='%(asctime)s - %(message)s', level=logging.INFO)

def return_archive_urls_for_download(host: str):
    """
    Return urls from xxl-angeln.de on archive.org from 2012 - 2013.
    This years are the last valid entries for xxl-angeln.de.

        Parameters:
            host (str): Host for archive urls.

        Returns:
            List[str]: All related urls for the host.
    """

    url = "https://web.archive.org/cdx/search/cdx?"
    params = {
        "url": host,
        "matchType": "domain",
        "filters": "statuscode:200",
        "output": "json",
        "from": "2012",
        "to": "2013"
    }
    url = url + urllib.parse.urlencode(params)
    logging.info(f'Request-URL: {url}')
    with urllib.request.urlopen(url=url, timeout=900) as req:
        urls = req.read().decode("utf-8")
    urls = json.loads(urls)
    logging.info(f'URLs to download: {len(urls)}')
    return urls


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
    logging.info('Mimetype-Folder-Structure created.')


def partition(lst):
    unique_urls = list(set([i[2] for i in lst]))
    for url in unique_urls:
        yield filter(lambda i: i[2] == url, lst)


def main():
    urls = return_archive_urls_for_download("xxl-angeln.de")
    mimetypes(urls)


if __name__ == "__main__":
    main()