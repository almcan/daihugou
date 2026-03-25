import os
import requests
import time

# 保存先のフォルダ名
save_folder = "cards"

# フォルダが存在しない場合は自動で作成
os.makedirs(save_folder, exist_ok=True)

# 1から53まで順番にダウンロード
for i in range(1, 54):
    url = f"https://chicodeza.com/wordpress/wp-content/uploads/torannpu-illust{i}.png"
    
    # 番号からマークと数字を判定するロジック
    if 1 <= i <= 13:
        mark = "spade"
        number = i
    elif 14 <= i <= 26:
        mark = "clover"
        number = i - 13
    elif 27 <= i <= 39:
        mark = "diamond"
        number = i - 26
    elif 40 <= i <= 52:
        mark = "heart"
        number = i - 39
    elif i == 53:
        mark = "joker"
        number = "" # ジョーカーは数字なし
    
    # 新しいファイル名を作成 (例: spade_1.png, joker.png)
    if mark == "joker":
        filename = f"{mark}.png"
    else:
        filename = f"{mark}_{number}.png"
        
    filepath = os.path.join(save_folder, filename)
    
    try:
        print(f"元画像 illust{i}.png を 『{filename}』 として保存中...")
        response = requests.get(url, stream=True)
        response.raise_for_status() 
        
        # 画像をフォルダに保存
        with open(filepath, 'wb') as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
                
        # ※重要※ サーバー負荷軽減のため1秒待機
        time.sleep(1)
        
    except requests.exceptions.RequestException as e:
        print(f"エラーが発生しました（スキップします）: {url} - {e}")

print("すべてのダウンロードと名前の変更が完了しました！")