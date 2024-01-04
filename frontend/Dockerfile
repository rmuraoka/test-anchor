# ベースイメージ
FROM node:latest

# 作業ディレクトリの設定
WORKDIR /front

# 依存関係のインストール
COPY package*.json ./
RUN npm install

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN npm run build

# アプリケーションの起動
CMD ["npm", "start"]
