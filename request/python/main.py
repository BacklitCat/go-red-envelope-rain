import json
import requests
import time
from multiprocessing import Process
from retrying import retry

host = "192.168.2.176"

loginUrl = "http://"+host+":7050/user/login"
rainUrl = "http://"+host+":7054/rain"
checkStatusUrl = rainUrl + "/checkStatus"
openEnvelopeUrl = rainUrl + "/openEnvelope"
headers = {
    "Content-Type": "application/json; charset=utf-8",
    "Host": host+":7050"
}

def retry_if_request_error(exception):
    # return isinstance(exception, (requests.HTTPError, requests.ConnectionError, requests.ConnectTimeout,
    #                               json.JSONDecodeError, requests.exceptions.ConnectionError,Exception))
    return True


class user(object):
    account = ""
    password = ""
    loginBody = ""
    code = 0
    status = True
    name = ""
    accessToken = ""
    accessExpire = 0
    refreshAfter = 0
    remaining = 0
    balance = 0

    def __init__(self, account, password):
        self.account = account
        self.password = password
        self.loginBody = json.dumps({"account": account, "password": password})

    def set_name(self, name):
        self.name = name

    def set_token_field(self, accessToken, accessExpire, refreshAfter):
        self.accessToken = accessToken
        self.accessExpire = accessExpire
        self.refreshAfter = refreshAfter

    def set_activity_field(self, remaining, balance):
        self.remaining = remaining
        self.balance = balance

    @retry(stop_max_attempt_number=10, stop_max_delay=20000, retry_on_exception=retry_if_request_error)
    def login(self):
        time.sleep(3.5)
        r = requests.post(loginUrl, self.loginBody, headers=headers).json()
        if r["code"] != 0:
            raise Exception("svc error:" + r["msg"])
        self.code = r["code"]
        self.set_name(r["data"]['name'])
        self.set_token_field(r["data"]['accessToken'],
                             r["data"]['accessExpire'],
                             r["data"]['refreshAfter'])

    @retry(stop_max_attempt_number=10, stop_max_delay=20000, retry_on_exception=retry_if_request_error)
    def check_status(self):
        time.sleep(3.5)
        h = headers
        h["Authorization"] = self.accessToken
        r = requests.post(checkStatusUrl, headers=h).json()
        self.code = r["code"]
        if r["code"] != 0:
            raise Exception("svc error:" + r["msg"])
        self.set_activity_field(r["data"]['remaining'],
                                r["data"]['balance'])

    @retry(stop_max_attempt_number=10, stop_max_delay=20000, retry_on_exception=retry_if_request_error)
    def open_envelope(self):
        time.sleep(3.5)
        h = headers
        h["Authorization"] = self.accessToken
        r = requests.post(openEnvelopeUrl, headers=h).json()
        self.code = r["code"]
        if r["code"] != 0:
            print(self.account + "执行open_envelope报错：", r)
            if r["msg"] == '用户访问接口过快':
                time.sleep(5)
            if r["msg"] == '用户无剩余活动次数':
                self.status = False
            return
        print(self.account + "参与活动，获得了", r["data"]['amount'], "分，余额", r["data"]['balance'],
              "分，还有", r["data"]['remaining'], "次机会")
        self.set_activity_field(r["data"]['remaining'],
                                r["data"]['balance'])

    def simulate(self):
        self.login()
        self.check_status()
        while self.remaining > 0 and self.status:
            if time.time() > self.refreshAfter:  # 检查凭证
                self.login()
                self.check_status()
            self.open_envelope()
        print(self.account + "模拟参与活动完成，共获得了", self.balance, "分")


if __name__ == "__main__":
    user_num = 1000  # 用户数量
    process_list = []
    for i in range(user_num):
        u = user(str(1006000 + i + 1), "123456")
        p = Process(target=u.simulate)
        p.start()
        process_list.append(p)

    for p in process_list:
        p.join()
