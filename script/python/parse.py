import json
import mysql.connector
import os
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
        INSERT INTO poetry (title, author, paragraphs, rhythmic, notes, poetry_type, dynasty, paragraphs_tradition)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """
        # 转换嵌套数据为JSON字符串
        paragraphs = json.dumps(item.get('paragraphs'), ensure_ascii=False)
        notes = json.dumps(item.get('notes'),ensure_ascii=False)
        cursor.execute(insert_query, (
            item.get('title'),
            item.get('author'),
            paragraphs,
            item.get('rhythmic'),
            notes,
            item.get('poetry_type'),
            item.get('dynasty'),
            paragraphs,
        ))
    conn.commit()
    conn.close()
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
        cursor.execute(insert_query, (
            item.get('name'),
            item.get('desc'),
            item.get('dynasty'),
        ))
    conn.commit()
    conn.close()
def read_wudai_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data
def read_author_json(path, json_data):
    with open(path, 'r', encoding='utf-8') as f:
        data = json.load(f)
        json_data.append(data)
    for item in json_data:
        for data_item in item:
            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
            print(formatted_json)
def read_wudai_author_and_insert():
    json_data = []
    read_author_json("/root/data/chinese-poetry/五代诗词/nantang/authors.json", json_data)
    insert_mysql_author(json_data[0])

def read_wudai_sc_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_wudai_json("/root/data/chinese-poetry/五代诗词/huajianji", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        for data_item in data:
            data_item['poetry_type'] = '五代诗词'
            data_item['dynasty'] = '五代十国'
            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
            print(formatted_json)
        insert_mysql(data)

def read_yuanqu_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_wudai_json("/root/data/chinese-poetry/元曲", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        for data_item in data:
            data_item['poetry_type'] = '元曲'
            data_item['notes'] = ''
            data_item['rhythmic'] = ''
            data_item['dynasty'] = data_item['dynasty']
            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
            print(formatted_json)
        insert_mysql_poetry(data)

def read_song_author_and_insert():
    json_data = []
    read_author_json("/root/data/chinese-poetry/全唐诗/authors.song.json", json_data)
    insert_mysql_author(json_data[0])

def read_tang_author_and_insert():
    json_data = []
    read_author_json("/root/data/chinese-poetry/全唐诗/authors.tang.json", json_data)
    insert_mysql_author(json_data[0])

def read_tang_song_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if "author" in file :
                print(file+":author skip")
                continue
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['poetry_type'] = '全唐诗'
                            if 'song' in file:
                                data_item['dynasty'] = '宋'
                            else:
                                data_item['dynasty'] = '唐'
                            data_item['notes'] = ''
                            data_item['rhythmic'] = ''
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data

def read_tangshi_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_tang_song_json("/root/data/chinese-poetry/全唐诗", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)

def read_song_ci_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['dynasty'] = '宋'
                            data_item['poetry_type'] = '宋词'
                            if 'title' not in data_item:
                                data_item['title'] = ''
                            if 'notes' not in data_item:
                                data_item['notes'] = ''
                            if 'author' not in data_item:
                                data_item['author'] = ''
                            if 'rhythmic' not in data_item:
                                data_item['rhythmic'] = ''
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data

def read_song_ci_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_song_ci_json("/root/data/chinese-poetry/宋词", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)

def read_yuding_tang_song_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['poetry_type'] = '御定全唐诗'
                            data_item['dynasty'] = '唐' 
                            if 'title' not in data_item:
                                data_item['title'] = ''
                            if 'notes' not in data_item:
                                data_item['notes'] = ''
                            if 'author' not in data_item:
                                data_item['author'] = ''
                            if 'rhythmic' not in data_item:
                                data_item['rhythmic'] = ''
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data
def read_yuding_tangshi_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_yuding_tang_song_json("/root/data/chinese-poetry/御定全唐詩/json", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)


def read_cao_cao_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['poetry_type'] = '曹操诗集'
                            data_item['dynasty'] = '三国' 
                            if 'title' not in data_item:
                                data_item['title'] = ''
                            if 'notes' not in data_item:
                                data_item['notes'] = ''
                            if 'author' not in data_item:
                                data_item['author'] = '曹操'
                            if 'rhythmic' not in data_item:
                                data_item['rhythmic'] = ''
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data

def read_caocao_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_cao_cao_json("/root/data/chinese-poetry/曹操诗集", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)


def read_chu_ci_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['poetry_type'] = '楚辞'
                            data_item['dynasty'] = '三国' 
                            if 'title' not in data_item:
                                data_item['title'] = ''
                            if 'notes' not in data_item:
                                data_item['notes'] = ''
                            if 'author' not in data_item:
                                data_item['author'] = '曹操'
                            if 'rhythmic' not in data_item:
                                data_item['rhythmic'] = ''
                            if 'section' in data_item and 'rhythmic' not in data_item:
                                data_item['rhythmic'] = data_item['section']
                            if 'paragraphs' not in data_item and 'content' in data_item:
                                data_item['paragraphs'] = data_item['content']
                            if 'paragraphs' not in data_item:
                                data_item['paragraphs'] = '[]'
                            if 'paragraphs_tradition' not in data_item:
                                data_item['paragraphs_tradition'] = data_item['paragraphs']
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data
    
def read_chuci_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_chu_ci_json("/root/data/chinese-poetry/楚辞", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)


def read_shuimo_tangshi_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['poetry_type'] = '水墨唐诗'
                            data_item['dynasty'] = '唐' 

                            if 'title' not in data_item:
                                data_item['title'] = ''
                            if 'notes' not in data_item:
                                data_item['notes'] = ''
                            if 'author' not in data_item:
                                data_item['author'] = '曹操'
                            if 'rhythmic' not in data_item:
                                data_item['rhythmic'] = ''
                            if 'prologue' in data_item and 'notes' not in data_item:
                                data_item['notes'] = data_item['prologue']
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data
    
def read_shuimo_tangshi_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_chu_ci_json("/root/data/chinese-poetry/水墨唐诗", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)


def read_nalanxingde_json(directory, json_data):
    for root, dirs, files in os.walk(directory):  # 递归遍历目录[1,2,5](@ref)
        for file in files:
            if file.endswith(".json"):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        for data_item in data:
                            data_item['poetry_type'] = '纳兰性德'
                            data_item['dynasty'] = '清' 
                            data_item['paragraphs'] = data_item['para']
                            if 'title' not in data_item:
                                data_item['title'] = ''
                            if 'notes' not in data_item:
                                data_item['notes'] = ''
                            if 'author' not in data_item:
                                data_item['author'] = '纳兰性德'
                            if 'rhythmic' not in data_item:
                                data_item['rhythmic'] = ''
                            formatted_json = json.dumps(data_item, indent=4, ensure_ascii=False)
                            print(formatted_json)
                        json_data.append(data)
                except json.JSONDecodeError:
                    print(f"⚠️ 文件 {file_path} 格式无效")
    return json_data
    
def read_nalanxingde_and_insert():
    json_data = []
    ## 解析五代的诗词
    read_nalanxingde_json("/root/data/chinese-poetry/纳兰性德", json_data)
    # read_wudai_json("/root/data/chinese-poetry/五代诗词/nantang", json_data)
    for data in json_data:
        insert_mysql_poetry(data)

read_caocao_and_insert()