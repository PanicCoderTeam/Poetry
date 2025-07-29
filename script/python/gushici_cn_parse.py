import requests
import re
import json
from bs4 import BeautifulSoup
import mysql.connector

from openai import OpenAI

total = 0
def insert_mysql_poetry(json_data):
    # 2. 建立数据库连接
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='pi=3.1415',
        database='poetry'
    )
    cursor = conn.cursor()

    # # 4. 数据插入
    for item in json_data:
        insert_query = """
        INSERT INTO poetry (title, author, paragraphs, dynasty, paragraphs_tradition, title_tradition, author_tradition, notes, comment,translation, pinyin)
        VALUES (%s, %s, %s, %s, '','','','', '', '','')
        """
        # 转换嵌套数据为JSON字符串
        # paragraphs = json.dumps(item.get('paragraphs'), ensure_ascii=False)
        # notes = json.dumps(item.get('notes'),ensure_ascii=False)
        # if item.get('rhythmic') is None:
        #     item['rhythmic'] = ''
        cursor.execute(insert_query, (
            item.get('title'),
            item.get('author'),
            item.get('paragraphs'),
            item.get('dynasty'),
        ))
    conn.commit()
    conn.close()

def generateT(poetry):
    
    # 定义 API 的 URL
    url = 'https://open.hunyuan.tencent.com/openapi/v1/agent/chat/completions'

    # 定义请求头
    headers = {
        'X-Source': 'openapi',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer pFNzBgvfqcatZHSCnHmpH3rPmP7YPq1A'
    }
    if poetry['author'] is None:
        poetry['author'] = '佚名'
    if poetry['paragraphs'] is None:
        poetry['paragraphs'] = ''
    # 定义请求体
    data = {
        "assistant_id": "YeKXUwbBE0hg",
        "user_id": "username",
        "stream": False,
        "messages": [
            {
                "role": "user",
                "content": [
                    {
                        "type": "text",
                        "text": "作者:"+ poetry["author"]+ "，诗句："+ poetry["paragraphs"]+ ",诗名"+poetry["title"],
                    }
                ]
            }
        ]
    }

    # 将请求体转换为 JSON 格式的字符串
    json_data = json.dumps(data)

    # 发送 POST 请求
    response = requests.post(url, headers=headers, json=data)  # 使用 json 参数自动设置正确的 Content-Type
    # print(response.text)
    dataParse = json.loads(response.text)
    # print(dataParse)
    # print(dataParse['choices'])
    # print(dataParse['choices'][0]['message']['content'])
    returnData = json.loads(dataParse['choices'][0]['message']['content'])
    # print(returnData)
    poetry['notes']=returnData['note']
    poetry['translation'] = returnData['translation']
    poetry['comment'] = returnData['comment']
    poetry['pinyin'] = returnData['pinyin']
    return poetry

def extract_categories(url):
    """
    提取网页中的分类信息及其URL
    :param url: 目标网页URL
    :return: 分类字典 {分类名称: 分类URL}
    """
    try:
        # 发送HTTP请求
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(url, headers=headers, timeout=10)
        response.encoding = 'utf-8'
        
        # 解析HTML内容
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # 查找分类区域 - 根据网站结构调整选择器
        rigntDiv = soup.select('.right')
        category_section = rigntDiv[1].select('.sons')  # 常见class名
        category_dict_list = []
        # 提取分类链接
        for category in category_section:
            category_dict = {}
            titleDiv = category.find('div', class_="title")
            titleText = titleDiv.get_text(strip=True)
            category_dict['main_category'] = titleText
            
            category_item_list = category.find_all('a')
            a_list = []
            for ci in category_item_list:
                poetry_info = {"name":ci.get_text(),"url":'https://www.gushiwen.cn'+ci['href']}
                a_list.append(poetry_info)
            
            category_dict['category_list'] =  a_list
            print(category_dict)
            category_dict_list.append(category_dict)
        
        return category_dict_list
    
    except Exception as e:
        print(f"解析过程中出错: {e}")
        return {}

def parse_shijing_title(full_title):
    """解析诗经标题和朝代"""
    # 分割主标题和朝代
    if '〔' in full_title and '〕' in full_title:
        title_part = full_title.split('〔')[0].strip()
        dynasty_part = full_title.split('〔')[1].replace('〕', '').strip()
    else:
        title_part = full_title
        dynasty_part = None
    return title_part, dynasty_part


flag = False
def extract_shiwen(url):
    # 发送HTTP请求
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
    }
    response = requests.get(url, headers=headers, timeout=10)
    response.encoding = 'utf-8'
    # 解析HTML内容
    soup = BeautifulSoup(response.text, 'html.parser')

    # 查找分类区域 - 根据网站结构调整选择器
    leftDiv = soup.select('.left')
    # print(leftDiv[1])
    category_section = leftDiv[1].select('.sons')  # 常见class名
    category_dict_list = []
    # print(category_section)
    secions = []
    if len(category_section) > 0:
        secions = category_section[0].select('.typecont')
    for category in secions:
        # print("00------0000")
        typeCont = category
        typeContMl = typeCont.find('div', class_='bookMl')
        spanList = typeCont.find_all('span')
        spanDictList = []
        sub_category_type = "none"
        if typeContMl:
            sub_category_type = typeContMl.find('strong').get_text()
        # global flag
        # if sub_category_type == '近代曲辞':
        #     flag = True
        # if flag == False:
        #     continue
        contType={'sub_category_type': sub_category_type}
        for span in spanList:
            
            # print("---------span--------")
            # print(span)
            author = span.get_text()
            pattern = r'\(.*\)'  
            result = re.findall(pattern, author)
            # print('author:'+str(result))
            if len(result) > 0:
                author = result[0]
            if span :
                a_elem = span.find('a')
                if a_elem and a_elem.has_attr('href'):
                    try:
                        poetry = extract_poetry('https://www.gushiwen.cn'+ a_elem['href'])
                        if poetry and poetry['title']:
                        # poetry={'author':'佚名','title':'击壤歌','content':'日出而作，日入而息。凿井而饮，耕田而食。帝力于我何有哉！'}
                            spanDictList.append(poetry)
                    except ZeroDivisionError as e:
                        poetry = extract_poetry('https://www.gushiwen.cn'+ a_elem['href'])
                        if poetry and poetry['title']:
                        # poetry={'author':'佚名','title':'击壤歌','content':'日出而作，日入而息。凿井而饮，耕田而食。帝力于我何有哉！'}
                            spanDictList.append(poetry) 
        # insert_mysql_poetry(spanDictList)
        global total
        total += len(spanDictList)
        print("Len:"+str(total))
        print("sub:"+sub_category_type)
        contType['poetry_list']=spanDictList
        # print(spanDictList)
        # print(contType)
        category_dict_list.append(contType)
    # print(category_dict_list)
    return category_dict_list

def extract_poetry(url) :
    # 发送HTTP请求
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
    }
    response = requests.get(url, headers=headers, timeout=300)
    response.encoding = 'utf-8'
    # 解析HTML内容
    soup = BeautifulSoup(response.text, 'html.parser')

    # 查找分类区域 - 根据网站结构调整选择器
    leftDiv = soup.select('.left') 
    # print(leftDiv)
    print(url)
    if len(leftDiv) <=1:
        print("error:"+ url)
    if len(leftDiv) > 1:
        category_section = leftDiv[1].select('.sons')  # 常见class名
        if len(category_section) > 0 :
            # print(category_section[0])
            # print("-----------------")
            # print(category_section[1])
            # print("----------------fff-------")
            conts = category_section[0].select('.cont')
            if len(conts) > 0:
                # print(conts)
                divs = conts[0].find_all('div')
                # print(divs[1])
                # print(divs[2])
                title = divs[1].find('h1')
                poetry = {}
                if title :
                    # print(title.get_text())
                    poetry['title'] = title.get_text()
                source = divs[1].find('p')
                if source:
                    # print(source.get_text())
                    author_info = source.get_text()
                    poetry['author'],poetry['dynasty']=parse_shijing_title(author_info)
                    
                conson = divs[1].find('div')
                if conson:
                    # print(conson.get_text())
                    poetry['paragraphs'] = conson.get_text().rstrip('\n').lstrip('\n')
        return poetry
                
            
# 使用示例
if __name__ == "__main__":
    # 替换为您要解析的实际URL
    target_url = "https://www.gushiwen.cn/shiwens/"
    
    shiwen_category = extract_categories(target_url)
    
    # print("提取到的分类信息：")
    for category in shiwen_category:
        print(category)
    
    ## 解析诗文
    for category in shiwen_category[0]['category_list']:
        # print('extractshiwen:' + category['url'])
        # if category['url']== "https://www.gushiwen.cn/gushi/yongwu.aspx" :
        #     flag = True
        # if flag == False or category['url']== "https://www.gushiwen.cn/gushi/yongwu.aspx":
        #     continue
        sub_categories = extract_shiwen(category['url'])
        category['sub_categories'] = sub_categories
    json_str = json.dumps(shiwen_category, indent=4, ensure_ascii=False)  # 生成格式化JSON字符串
    with open('poetry.json', 'w', encoding='utf-8') as f:
        f.write(json_str)
    

    # poetry = extract_poetry('https://www.gushiwen.cn/shiwenv_6ba04da3fdc1.aspx')
    # print(poetry)
    # # poetry={'author':'佚名','title':'击壤歌','content':'日出而作，日入而息。凿井而饮，耕田而食。帝力于我何有哉！'}
    # generateT(poetry)
    # print(poetry)