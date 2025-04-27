import pandas as pd

# 读取Excel文件
file_path = 'your_file.xlsx'  # 替换为你的Excel文件路径
data = pd.read_excel(file_path)

# 显示数据
print(data.head())  # 查看前几行数据

# 访问特定列
print(data['列名'])  # 替换'列名'为你的实际列名

# 转换为列表或字典
data_list = data.values.tolist()
data_dict = data.to_dict('records')