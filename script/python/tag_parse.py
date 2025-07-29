import json
import mysql.connector
import os

import requests
import re
from bs4 import BeautifulSoup
import mysql.connector

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


def insert_mysql_tag(json_data):
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
        INSERT INTO tag (name, parent_tag, level)
        VALUES (%s, %s, %s)
        """
        # 转换嵌套数据为JSON字符串
        # paragraphs = json.dumps(item.get('paragraphs'), ensure_ascii=False)
        # notes = json.dumps(item.get('notes'),ensure_ascii=False)
        # if item.get('rhythmic') is None:
        #     item['rhythmic'] = ''
        cursor.execute(insert_query, (
            item.get('name'),
            item.get('parent_tag'),
            item.get('level'),
        ))
    conn.commit()
    conn.close()



def insert_mysql_poetry_tag(json_data):
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
        INSERT INTO poetry_tag (poetry_id, tag, min_tag_id)
        VALUES (%s, %s, %s)
        """
        # 转换嵌套数据为JSON字符串
        # paragraphs = json.dumps(item.get('paragraphs'), ensure_ascii=False)
        # notes = json.dumps(item.get('notes'),ensure_ascii=False)
        # if item.get('rhythmic') is None:
        #     item['rhythmic'] = ''
        cursor.execute(insert_query, (
            item['poetry_id'],
            item['tag'],
            item['min_tag_id'],
        ))
    conn.commit()
    conn.close()
def readContent():
    content = ''
    with open('poetry.json', 'r', encoding='utf-8') as file:
        content = file.read()
    data = json.loads(content)
    category_list = data[0]['category_list']
    tag_list = []
   
    tag_info_list = select_mysql_tag()
    tag_map ={}
    for tag_info in tag_info_list:
        if len(tag_info['parent_tag']) > 0:
            tag_map[tag_info['parent_tag']+"-"+tag_info['name']]=tag_info['id']
    # print(tag_map)
    poetry_tag_list = select_mysql_poetry_tag()
    poetry_tag_map = {}
    # print(poetry_tag_list)
    for poetry_tag in poetry_tag_list:
        poetry_tag_map[str(poetry_tag['poetry_id'])+"-"+poetry_tag['tag']] = poetry_tag['id']
    # print(poetry_tag_map)
    for category in category_list:
        for sub_category in category['sub_categories']:
            # print(category['name'])
            # print(sub_category['sub_category_type'])
            poetry_tag = []
            for poetry in sub_category['poetry_list']:
                poetry_list = select_mysql_poetry( poetry['author'],poetry['title'])
                # print(len(poetry_list))
                if len(poetry_list) > 1:
                    print("程序暂停，等待用户输入")
                    input("按回车键继续...")
                    print("程序继续执行")
                if len(poetry_list) == 1:
                    tagStr=category['name']+"-"
                    if sub_category['sub_category_type'] and sub_category['sub_category_type'] != 'none':
                        tagStr+=sub_category['sub_category_type']
                    else:
                        tagStr+=category['name']
                    
                    p_t_key = str(poetry_list[0]['id'])+"-"+tagStr
                    if p_t_key not in poetry_tag_map:
                        print("p_t_key:"+p_t_key)
                        poetry_tag.append({'poetry_id':poetry_list[0]['id'],'tag':tagStr,'min_tag_id':tag_map[tagStr]})
                        poetry_tag_map[p_t_key]=1
                        # print('poetry_id:'+str(poetry_list[0]['id'])+" ; tag:"+ tagStr)
                # print(poetry_tag)
            if len(poetry_tag) > 0:
                print(len(poetry_tag))
                insert_mysql_poetry_tag(poetry_tag)
    # unique_data = [eval(i) for i in {str(d) for d in tag_list}]
    # print(unique_data)
    # insert_mysql_tag(unique_data)
    
def select_mysql_poetry(author, title):
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='pi=3.1415',
        database='poetry'
    )
    query_result = []
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM poetry where title = %s and author = %s", (title,author))  # 执行查询
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
  
def select_mysql_tag():
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='pi=3.1415',
        database='poetry'
    )
    query_result = []
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM tag")  # 执行查询
    results = cursor.fetchall()  # 获取全部结果
    for row in results:
        query_result.append({
            'id':row[0],
            'name':row[1],
            'parent_tag':row[2],
            'level':row[3],
        })
        # print(row)
    cursor.close()
    conn.close()
    return query_result


def select_mysql_poetry_tag():
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='pi=3.1415',
        database='poetry'
    )
    query_result = []
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM poetry_tag")  # 执行查询
    results = cursor.fetchall()  # 获取全部结果
    for row in results:
        query_result.append({
            'id':row[0],
            'poetry_id':row[1],
            'tag':row[2],
            'min_tag_id':row[3],
        })
        # print(row)
    cursor.close()
    conn.close()
    return query_result

# 使用示例
if __name__ == "__main__":
    readContent()
   