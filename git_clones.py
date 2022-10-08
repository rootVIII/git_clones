from argparse import ArgumentParser
from re import findall
from subprocess import call
from sys import exit
from urllib.request import urlopen


"""
rootVIII
Download/clone all of a user's PUBLIC source repositories.
Pass the Github user's username with the -u option.
"""


class GitClones:
    def __init__(self, user):
        self.url = f'https://github.com/{user}'
        self.url += '?&tab=repositories&q=&type=source'
        self.git_clone = 'git clone https://github.com/%s/%%s.git' % user
        self.user = user

    def http_get(self):
        return urlopen(self.url).read().decode('utf-8')

    def get_repo_data(self):
        pattern = r"<a\s?href\W+%s/(.*)\"\s+" % self.user
        for line in findall(pattern, self.http_get()):
            yield line.split('\"')[0]

    def get_repos(self):
        return [repo for repo in self.get_repo_data()]

    def download(self):
        _ = [call((self.git_clone % git).split()) for git in self.get_repos()]


if __name__ == "__main__":
    message = 'Usage: python3 git_clones.py -u <github username>'
    parser = ArgumentParser(description=message)
    parser.add_argument('-u', '--user', required=True, help='Github Username')

    clones = GitClones(parser.parse_args().user)
    try:
        clones.download()
    except Exception as err:
        print(f'{type(err).__name__}: {str(err)}')
        exit(1)
