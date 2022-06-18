#include <algorithm>
#include <iostream>
#include <stdexcept>
#include <utility>
#include <curl/curl.h>
#include <stdlib.h>


using namespace std;

/*
rootVIII 2019-2021

Build commands if not using CMakeLists.txt:

MacOS:
c++ git_clones.cpp -o git_clones -lcurl -std=c++11

Linux:
g++ git_clones.cpp -o git_clones -lcurl -std=c++11
*/


class GitClones {

private:
    string base_url, clone_url, username;

    static size_t write_func(void *ptr, size_t size, size_t nmemb, string* data) {
        data->append((char*) ptr, size * nmemb);
        return size * nmemb;
    }

    string http_get() {
        CURL *curl = curl_easy_init();
        if (!(curl)) {
            throw runtime_error("Failed: CURL init NULL");
        }
        string response;
        long http_status_code;
        curl_easy_setopt(curl, CURLOPT_URL, clone_url.c_str());
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_func);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);
        curl_easy_perform(curl);
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_status_code);
        curl_easy_cleanup(curl);
        if (http_status_code != 200)
            throw runtime_error("Failed to retrieve repositories, ensure username is valid");
        return response;
    };

public:
    explicit GitClones(const string& user) {
        base_url = "https://github.com/" + user;
        clone_url = base_url + "?&tab=repositories&q=&type=source";
        username = user;
    }

    void clone_repo(const string &repository) {
        string cmd = "git clone " + base_url + "/" + repository;
        system(cmd.c_str());
    }

    void clone_repos() {
        string html = http_get();
        string space = " ";
        string href = "href=\"/" + username + "/";
        auto start = 0u;
        auto end = html.find(space);
        while (end != string::npos) {
            string line = html.substr(start, end - start);
            if (line.find(href) != string::npos && count(line.begin(), line.end(), '/') == 2)
                clone_repo(line.substr(href.length(), line.length() - href.length() - 1));

            start = end + space.length();
            end = html.find(space, start);
        }
    }
};

int main(int argc, char *argv[]) {
    if (argc != 2) {
        cout << "Provide a single Github username as a command-line argument\n";
        exit(1);
    }

    GitClones gc(*++argv);

    try {
        gc.clone_repos();
    } catch (const exception &err) {
        cerr << "Error occurred: " << err.what() << endl;
        exit(1);
    }

    return 0;
}
