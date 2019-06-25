#! /usr/bin/python3.7
# rootVIII
# Download/clone all of a user's public repositories
# Pass a the Github user's username with the -u option
# Usage: python git_clones.py -u <github username>
# Example: python git_clones.py -u rootVIII
#
from argparse import ArgumentParser
from sys import exit
from re import findall
from urllib.request import urlopen
from subprocess import call


class GitClones:
    def __init__(self, user):
        self.url = "https://github.com/%s?tab=repositories" % user
        self.git_clone = "git clone https://github.com/%s/" % user
        self.git_clone += "%s.git"
        self.user = user
        self.repos = []
        self.page = ''

    def get_repo_data(self):
        try:
            r = urlopen(self.url)
        except Exception:
            print("Unable to make request to %s's Github page" % self.user)
            exit(1)
        else:
            self.page = r.read().decode('utf-8')
            pattern = r"repository_nwo:%s/(.*)," % self.user
            for line in findall(pattern, self.page):
                yield line.split(',')[0]

    def get_repositories(self):
        self.repos = [repo for repo in self.get_repo_data()]
        return set(self.repos)

    def download(self, git_repos):
        for git in git_repos:
            cmd = self.git_clone % git
            try:
                call(cmd.split())
            except Exception as e:
                print(e)
                print('unable to download:%s\n ' % git)


if __name__ == "__main__":
    message = 'Usage: python git_clones.py -u <github username>'
    h = 'Github Username'
    parser = ArgumentParser(description=message)
    parser.add_argument('-u', '--user', required=True, help=h)
    d = parser.parse_args()
    clones = GitClones(d.user)
    repositories = clones.get_repositories()
    clones.download(repositories)
