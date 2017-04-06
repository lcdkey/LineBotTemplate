// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

def gettotal():
	
	apikey={1HVzlSqwxT1f6vOult1W3qZuGhneOjZk5r83KrlvNIK8}
	getvalueurl='https://sheets.googleapis.com/v4/spreadsheets/{1}/values/sheet1!A:C?key=%s' % (apikey)
	res = requests.get(getvalueurl)
	data = res.content
	
	jsondata = json.loads(data)
	values = jsondata['values']
    #開關	
    count=0
    #以下3個變數作為加總的時候用的 	
    cash=0
	credit=0
	other=0
	for i in values:
		if count==0:
			count+=1
			continue
		else:
			#1代表cash
			if i[2]=='cash':
				cash+=int(i[1])
			#2代表credit
			elif i[2]=='credit':
				credit+=int(i[1])
			else:
				other+=int(i[1])
	outputvalue = 'Cash:{}, Credit:{}, \nOther:{}, Total:{}'.format(cash,credit,other,cash+credit+other)
	return outputvalue
"""
填寫google表單，進行記帳的動作
範例：1000,1,note 第1個參數為金額，第2個為種類：1表示cash、2表示credit、其他為other，第3個參數為備註。
"""
def editRecord(accountlist):
	actionurl = 'https://docs.google.com/forms/d/e/{sheet ID}/formResponse'
	formurl = 'https://docs.google.com/forms/d/e/{sheet ID}/viewform'
	r = requests.Session()
	res = r.get(formurl)
	soup = bs(res.content,'html.parser')
	inputlist = soup.find_all('input')
	textarealist = soup.find_all('textarea')
	test={}
	for i in inputlist:
		#print i['aria-label']+':'+i['value']
		try:
			if i['aria-label']=='金額'.decode("utf-8"):
				test.update({i['name']:accountlist[0]})
			elif i['aria-label']=='種類'.decode('utf-8'):
				if accountlist[1]=='1':
					inputkind = 'cash'
				elif accountlist[1]=='2':
					inputkind = 'credit'
				else:
					inputkind = 'other'
				test.update({i['name']:inputkind})
		except:
			test.update({i['name']:i['value']})
    #因為表單當中的備註欄位不是必填欄位，所以在記帳的時候，這個部分可有可無，所以需要判斷，有輸入的話才填入。		
    if len(accountlist)==3:	
		for j in textarealist:
			#print j['name']+':'+j['value']
			test.update({j['name']:accountlist[2]})
	res = r.post(actionurl,data=test)
	return res.status_code

app = Flask(__name__)

@app.route("/",methods=["GET"])
def index():
	return "hello world",200
	
@app.route("/callback",methods=["POST"])
def index2():
	temp = request.get_json()
	for i in temp['events']:
		ttt =  i['replyToken']
		print i['source']['userId']
		if i['message']['type']=='text':
			msg = i['message']['text']
		replyapi(ttt,msg)
	return "hello world",200
	
def replyapi(accesstoken,msg):
	channeltoken={Line channel Token}

	t = msg.encode('utf-8')
	inputlist = t.split(' ')
	count=0
	mat=[]
	if len(inputlist)>1:
		if inputlist[0]=='d' or inputlist[0]=='D':
			count,mat,outputstring = searchdictionary(count,inputlist)
	else:
		accountlist = inputlist[0].split(',')
		if inputlist[0]=='l' or inputlist[0]=='L':
			#進行記錄查詢的動作
			outputstring = gettotal()
			count=1
		elif len(accountlist)>1:
			#進行記帳的動作
			if editRecord(accountlist)==200:
				time.sleep(1)
				outputstring = '紀錄成功'
				count=1
			else:
				outputstring = '紀錄失敗'
				count=1
			#print 'hello'
		else:
			pat = re.compile(r".*(掰掰).*")
			print type(msg)
			mat = pat.findall(t)
	
	if len(mat)==0:
		if count==1:
			data={
			'replyToken':accesstoken,
			'messages':[{'type':'text','text':outputstring}]
			}
		else:
			data={
			'replyToken':accesstoken,
			'messages':[{'type':'text','text':'朕知道了'},{'type':'text','text':'可以退下了'}]
			}
	else:
		data={
		'replyToken':accesstoken,
		'messages':[{'type':'text','text':'慢走不送'}]
		}
	
	
	
	headers={
		'Content-Type':'application/json',
		'Authorization':'Bearer '+channeltoken
	}
	urladdress = 'https://api.line.me/v2/bot/message/reply'
	datajson = json.dumps(data)
	#print type(datajson)
	res = requests.post(urladdress,headers=headers,data=datajson)
	print res.status_code
	

if __name__=='__main__':
	app.run()
