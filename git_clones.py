from argparse import ArgumentParser
from re import findall
from subprocess import call
from sys import exit, version_info


try:
    from urllib.request import urlopen
except ImportError:
    from urllib2 import Request, urlopen


# rootVIII


class GitClones:

    """
    Download/clone all of a user's PUBLIC source repositories.
    Pass the Github user's username with the -u option.

    Args:
        user (str): Github username

    Example:
        python git_clones.py -u <github username>
        python3 git_clones.py -u rootVIII

    Note:
        Compatible with Python2 & Python3.
        Pycodestyle validated.
    """

    def __init__(self, user):
        self.url = 'https://github.com/%s' % user
        self.url += '?&tab=repositories&q=&type=source'
        self.git_clone = 'git clone https://github.com/%s/%%s.git' % user
        self.user = user

    def http_get(self):
        if version_info[0] != 2:
            req = urlopen(self.url)
            return req.read().decode('utf-8')
        req = Request(self.url)
        request = urlopen(req)
        return request.read()

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
    parser.add_argument('-u', '--user',
                        required=True, help='Github Username')

    clones = GitClones(parser.parse_args().user)
    try:
        clones.download()
    except Exception as err:
        print('%s: %s' % (type(err).__name__, str(err)))
        exit(1)
