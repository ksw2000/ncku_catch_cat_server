# catch_cat_server

The server of https://github.com/ksw2000/catch_cat_flutter.

Powered by Golang

## Database

### cat_kind

+ `cat_kind_id` *int* **key** (auto-generated)
+ `name` *string* 
+ `description` *string*
+ `weight` *int*

### cat

+ `cat_id` *int* **key** (auto-generated)
+ `cat_kind_id` *int* **foreign key**
+ `lng` *float64* (經度)
+ `lat` *float64* (緯度)
+ `theme_id` *int* **foreign key**

### theme

+ `theme_id` *int* **key** (auto-generated)
+ `name` *string*
+ `thumbnail` *string* (thumbnail 使用內部連結)
+ `description` *string*

### user_cat

+ `user_id` *int* **foreign key**
+ `cat_id` *int* **foreign key**
+ `timing` *int*

### user

+ `user_id` *int64*
+ `salt` *string*
+ `password` *string* (加密)
+ `name` *string*
+ `profile` *string | null*
+ `email` *string*
+ `creating` *uint64*
+ `last_login` *uint64*
+ `last_lng` *float64* (使用者同意下才可存取)
+ `last_lat` *float64* (使用者同意下才可存取)
+ `verified` *boolean* (是否通過郵箱驗證)

### verify_email

+ `verify_id` **key** (auto-generated)
+ `user_id` *int64* **foreign key**
+ `email` *string*
+ `token` *string*
+ `expire` *int64* (token 過期時間)

### user_view (view table)

+ `user_id` *int64*
+ `cats` *uint* (捕獲貓的數量)
+ `score` *int64* (user_cat.user_id 計算權重，可進一步換算 level)

### friend

+ `friend_id` *int* **key** (auto-generated)
+ `user_id_src` int **foreign key**
+ `user_id_dest` int **foreign key**
+ `accepted` *boolean* (src 向 dest 邀請，dest 是否已接收)
+ `ban` *boolean* (src 向 dest 封鎖)

> 當 a 與 b 有好友關係時，雙向都要加入，刪除好友關係時雙向也都要刪除

## API

### user

```
/POST/register (向資料庫註冊新用戶)✅
    - password
    - confirm_password
    - email
    - name

檢查 name 字元數 < 10
檢查 password 是否等於 confirm_password
檢查 email 格式 ❌
(寄發 email 確認：太麻煩先跳過)  ❌
產生 uid (12位數數字), 產生時檢查是否可用
產生 salt (256位a-zA-Z0-0)
產生 hash 後的 password
寫入資料庫

HTTP 200 成功，但不符合規定
HTTP 201 成功，成功建立資源

retrun
	- error (string)
```

```
/POST/login (向資料庫查尋用戶)✅
	- passowrd
	- email
根據 email 查尋資料庫
比對密碼
若成功則寫入 session

HTTP 200 成功
HTTP 401 登入失敗

return
	- error (string)
	- session (session id)
	- name
	- uid
	- profile
	- email
	- verified
	- rank
	- cats (int)
```

```
/POST/logout (登出)✅
	- session

刪除 session
更新最後登入時間

HTTP 200 成功

return
	- error
```

```
/GET/me (查尋登入狀態)
	- session

查尋 session 如果已登入
查找資料庫

return
	- is_login
	- name
	- uid
	- profile
	- email
	- verified
	- rank
	- cats
```

```
/POST/modify_user_password (更新用戶密碼)
	- session
	- uid
    - original_password
    - password
    - confirm_password

檢查是否已登入
檢查 original_password 是否正確
檢查 password 是否等於 confirm_password
產生新的 salt (256位a-zA-Z0-0)
產生 hash 後的 password
寫入資料庫

return
	- ok (boolean)
	- error (string)
```

```
/POST/modify_user_name (更新用戶名)
	- session
	- uid
	- name

檢查是否已登入
寫入資料庫

return
	- ok (boolean)
	- error (string)
```

```
/POST/modify_user_email (更新 email)
	- uid
	- email

檢查是否已登入
(寄發 email 確認：太麻煩先跳過 → 直接寫進資料庫)
```

```
/GET/verify_email (確定更新 email)
	- token
	
檢查資料庫
修改資料庫
```

```
/POST/update_position
	- session
	- lat
	- lng
	
檢查是否登入

更新資料庫
```

### friend

```
/POST/friends_position
	- session

檢查是否登入
查尋朋友位置

return 
	- error
	- position_list
		- name (朋友名字)
		- uid (朋友 id)
		- lat
		- lng
```

```
/POST/friends_theme_rank (查尋某個主題中自己及朋友的分數)
	- session
	- theme_id

檢查是否登入
查尋朋友的貓貓清單

return
	- error
	- sorted_rank_list
		- name (朋友名字)
		- uid (朋友 id)
		- cats (朋友捕獲的貓咪數量)
```


```
/POST/friend/delete ✅
	- session
	- friend_uid

檢查是否登入
修改資料庫 src->dest 刪除
如果 dest->src 是封鎖狀態，不刪除

HTTP 401 (未登入)
HTTP 200 請求成功，修改沒成功
HTTP 201 成功修改

return
	- ok
```

```
/POST/friend/ban
	- session
	- ban_uid
	
檢查是否登入
修改資料庫
刪除雙方友誼關係
src->dest 關係改為封鎖狀態

HTTP 401 (未登入)
HTTP 200 請求成功，修改沒成功
HTTP 201 成功修改

return
	- ok
```

```
/POST/friend/agree ✅
	- session
	- friend_uid

檢查是否登入
修改資料庫 (允許該筆資料)，並且反過來增加一筆資料

HTTP 401 (未登入)
HTTP 200 請求成功，修改沒成功
HTTP 201 成功修改

return
	- ok
```

```
/POST/friend/decline ✅
	- session
	- friend_uid

檢查是否登入
修改資料庫 (刪除 friend -> me)

HTTP 401 (未登入)
HTTP 200 請求成功，修改沒成功
HTTP 201 成功修改

return
	- ok
```

```
/POST/friends/list ✅
	- session

檢查是否登入

HTTP 401 (未登入)
HTTP 200

return
	- error
	- list
		- name
		- uid
		- level
		- last_login
```

```
/POST/friends/inviting_me ✅
	- session

檢查是否登入

HTTP 401 (未登入)
HTTP 200

return
	- error
	- list
		- name
		- uid
		- level
		- last_login
```

```
/POST/friend/invite ✅
	- session
	- finding_uid

檢查是否登入
HTTP 401 如果沒有登入

如果該用戶被封鎖回傳找不到

HTTP 200 成功，但可能找不到
HTTP 201 成功，成功建立資源

return
	- error
```

### cat, theme

```
/GET/theme_list ✅

HTTP 200 成功

回傳主題
	- error
	- list
		- thumbnail
		- name
		- theme_id
```

```
/GET/theme
	- theme_id

檢查是否登入

回傳主題內容
	- error
	- cat_list
		- cat_id
		- cat_kind_id
		- weight (貓咪對應分數)
		- lng
		- lat
		- caught (是否已被使用者捕獲)
```

```
/POST/add_caught_cat
	- cat_id

檢查是否登入
修改資料庫(新增已抓到的貓)

return 
	- ok
	- error
```

```
/GET/my_caught_cat_kind (用來處理圖鑑)

檢查是否登入

return
	- cat_list
		- cat_kind_id
		- weight (貓咪對應分數)
		- name
		- description
```

