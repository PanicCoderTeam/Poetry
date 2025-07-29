import json
import mysql.connector
import os

import requests
import re
from bs4 import BeautifulSoup
import mysql.connector
import threading
from openai import OpenAI

def update_mysql_poetry(poetry_info):
    # 2. 建立数据库连接
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='pi=3.1415',
        database='poetry'
    )
    cursor = conn.cursor()

    # # 4. 数据插入
    insert_query = """
    update poetry set notes = %s, comment = %s ,translation = %s, pinyin = %s where id = %s
    """
    # 转换嵌套数据为JSON字符串
    # paragraphs = json.dumps(item.get('paragraphs'), ensure_ascii=False)
    # notes = json.dumps(item.get('notes'),ensure_ascii=False)
    # if item.get('rhythmic') is None:
    #     item['rhythmic'] = ''
    cursor.execute(insert_query, (
        poetry_info['notes'],
        poetry_info['comment'],
        poetry_info['translation'],
        poetry_info['pinyin'],
        poetry_info['id']
    ))
    conn.commit()
    conn.close()


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

dataList = []

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
    print(poetry)
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
                        "text": "作者:"+ poetry["author"]+ "，诗句："+ poetry["paragraphs"]+ ",诗名"+poetry["title"]+". 整体生成的json必须完整! json 必须可解析，去掉多余换行符",
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
    dataParse = json.loads(response.text,strict=False)
    # print(dataParse)
    # print(dataParse['choices'])
    # print(dataParse['choices'][0]['message']['content'])
    try:
        returnData = json.loads(dataParse['choices'][0]['message']['content'],strict=False)
        # print(returnData)
        poetry['notes']=returnData['note']
        poetry['translation'] = returnData['translation']
        poetry['comment'] = returnData['comment']
        poetry['pinyin'] = json.dumps(returnData['pinyin'], indent=4, ensure_ascii=False)
    except Exception as e:
        print(dataParse)
        print(f"解析过程中出错: e{e}")
        global dataList
        poetry['content'] = dataParse['choices'][0]['message']['content']
        dataList.append(poetry)
    return poetry

def select_mysql_poetry(limit, offset):
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='pi=3.1415',
        database='poetry'
    )
    query_result = []
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM poetry where notes='' order by id limit %s offset %s", (limit,offset))  # 执行查询
    results = cursor.fetchall()  # 获取全部结果
    for row in results:
        query_result.append({
            'id':row[0],
            'title':row[1],
            'title_tradition':row[2],
            'paragraphs':row[3],
            'paragraphs_tradition':row[4],
            'author':row[5],
            'author_tradition':row[6],
            'dynasty':row[7],
            'notes':row[8],
            'comment':row[9],
            'translation':row[10],
            'pinyin':row[11],
        })
        # print(row)
    cursor.close()
    conn.close()
    return query_result

def task(poetry):
    poetry_info = generateT(poetry)
    print(poetry_info)
    update_mysql_poetry(poetry_info)

# 使用示例
if __name__ == "__main__":
    offset = 0
    while offset < 9280:
        poetry_list = select_mysql_poetry(9,offset)
        offset+=9
        threads = []
        print(offset)
        for poetry in poetry_list:
            # poetry_info = poetry
            print(str(poetry['id'])+" "+poetry['notes'])
            if poetry['notes'] == '':
                # poetry_info['notes']='三日：古代风俗，新媳妇婚后三日须下厨房做饭菜。入厨下：到厨房去。洗手：是表示恭敬和认真的意思。作羹汤：煮饭烧菜。未谙：不熟悉。姑：婆婆（丈夫的母亲）。食性：口味。遣：让。小姑：丈夫的妹妹。'
                # poetry_info['translation']='婚后第三天来到厨房，洗净双手开始熬汤。还不熟悉婆婆的口味，先叫小姑来尝一尝。'
                # poetry_info['pinyin']= json.dumps([{'paragraph': '三日入厨下', 'result': 'sān rì rù chú xià'}, {'paragraph': '洗手作羹汤', 'result': 'xǐ shǒu zuò gēng tāng'}, {'paragraph': '未谙姑食性', 'result': 'wèi ān gū shí xìng'}, {'paragraph': '先遣小姑尝', 'result': 'xiān qiǎn xiǎo gū cháng'}], indent=4, ensure_ascii=False)
                # poetry_info['comment']='这首诗通过描写新嫁娘初入夫家时的一个生活细节，生动地刻画了新嫁娘聪慧机敏的形象。她深知不熟悉婆婆口味的情况下贸然行事可能不佳，于是巧妙地先让小姑尝羹汤，以此来推测婆婆的喜好。短短二十个字，将新嫁娘小心谨慎又心思细腻的形象跃然纸上，富有生活情趣，同时也反映了当时的民俗风情。'
                threads.append(threading.Thread(target=task, args=(poetry,)))

        offset -= len(threads)
        for t in threads:
            t.start()
        for t in threads:
            t.join()
    with open('dataMap.json', 'w', encoding='utf-8') as f:
        f.write(json.dumps(dataList))
        
                

