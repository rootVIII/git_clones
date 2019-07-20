# rootVIII
# Download/clone all of a user's public source repositories
# Pass the Github user's username with the -u option
# Usage: python git_clones.py -u <github username>
# Example: python git_clones.py -u rootVIII
#
from argparse import ArgumentParser
from sys import exit, version_info
from re import findall
from subprocess import call
try:
    from urllib.request import urlopen
except ImportError:
    from urllib2 import Request, urlopen


class GitClones:
    def __init__(self, user):
        self.url = "https://github.com/%s" % user
        self.url += "?&tab=repositories&q=&type=source"
        self.git_clone = "git clone https://github.com/%s/%%s.git" % user
        self.user = user

    def http_get(self):
        if version_info[0] != 2:
            req = urlopen(self.url)
            return req.read().decode('utf-8')
        req = Request(self.url)
        request = urlopen(req)
        return request.read()

    def get_repo_data(self):
        try:
            response = self.http_get()
        except Exception:
            print("Unable to make request to %s's Github page" % self.user)
            exit(1)
        else:
            pattern = r"<a\s?href\W+%s/(.*)\"\s+" % self.user
            for line in findall(pattern, response):
                yield line.split('\"')[0]

    def get_repositories(self):
        return [repo for repo in self.get_repo_data()]

    def download(self, git_repos):
        for git in git_repos:
            cmd = self.git_clone % git
            try:
                call(cmd.split())
            except Exception as e:
                print(e)
                print("unable to download:%s\n" % git)


if __name__ == "__main__":
    message = 'Usage: python git_clones.py -u <github username>'
    h = 'Github Username'
    parser = ArgumentParser(description=message)
    parser.add_argument('-u', '--user', required=True, help=h)
    d = parser.parse_args()
    clones = GitClones(d.user)
    repositories = clones.get_repositories()
    clones.download(repositories)
