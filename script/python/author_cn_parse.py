import requests
import re
import json
from bs4 import BeautifulSoup
import mysql.connector

def extract_author(url, d):
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
        rigntDiv = soup.select('.left')
        category_section = rigntDiv[1].select('.sonspic')  # 常见class名
        category_dict_list = []
        
        amore_div = rigntDiv[1].select('.amore')
        hasNext = False
        if len(amore_div) > 0:
            hasNext = amore_div[0].has_attr('style') == False
            print(amore_div)
            print(hasNext)
        # 提取分类链接
        for category in category_section:
            # print(category)
            category_dict = {}
            titleDiv = category.find('div', class_="cont")
            # titleText = titleDiv.get_text(strip=True)
            # category_dict['main_category'] = titleText
            
            category_item_list = category.find_all('p')
            # print(category_item_list)
            category_dict_list.append({"name":category_item_list[0].get_text(),"desc":category_item_list[1].get_text(), "dynasty": d})
            # a_list = []
            # for ci in category_item_list:
            #     a_list.append({"name":ci.get_text(),"url":'https://www.gushiwen.cn'+ci['href']})
            # category_dict['category_list'] =  a_list
            # print(category_dict)
            # category_dict_list.append(category_dict)
        
        return category_dict_list, hasNext
    
    except Exception as e:
        print(f"解析过程中出错: {e}")
        return {}


def insert_mysql_author(json_data) :
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
        INSERT INTO author (name, `desc`, dynasty)
        VALUES (%s, %s, %s)
        """
        if item.get('desc') is None:
            item.put('desc','')
        if item['dynasty'] is None:
            item['dynasty']=''
        cursor.execute(insert_query, (
            item.get('name'),
            item.get('desc'),
            item.get('dynasty'),
        ))
    conn.commit()
    conn.close()

if __name__ == "__main__":
    # 替换为您要解析的实际URL
    target_url = "https://www.gushiwen.cn/authors/Default.aspx?"
    dynasty =  ["先秦","两汉","魏晋","南北朝","隋代","唐代","五代","宋代","金朝","元代","明代","清代"]
    dynasty_author_list = []
    for d in dynasty:
        page = 1
        author_list = []
        if d == None:
            continue
        while True:
            author_category, hasNext = extract_author(target_url+"p="+str(page)+"&c="+str(d), d)
            author_list.extend(author_category)
            page=page+1
            if hasNext==False:
                break
        print("authorList:" + str(page)+":"+d)
        dynasty_author_list.append({"dynasty":d,"author_list":author_list})
        
    json_str = json.dumps(dynasty_author_list, indent=4, ensure_ascii=False)  # 生成格式化JSON字符串
    with open('data.json', 'w', encoding='utf-8') as f:
        f.write(json_str)
    
    for item in dynasty_author_list:
        insert_mysql_author(item.get('author_list'))


        
        
