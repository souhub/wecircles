# README.md

## WECIRCLES

完成品はこちらです。→https://wecircles.net

スマホで使用されることを想定していますが、パソコンでもタブレットでも使用可能になってます。
※ロードバランサーの負荷分散の影響でプロフィール画像のアップロードのみ数回に一回失敗するエラーが発生します。
数回繰り返していただければできますので、その場合はそうしていただけると幸いです。

### 概要

コロナの影響で新歓ができない大学生向けのサークル運営者を中心とするサークル所属者と新入生を中心とするサークル探しをしている人たちを結びつける SNS です。
「サークル」というグループを作成したり、参加したりでき、その中でチャット機能等を用いて交流することができます。

### 制作の背景

私は大学でサークルの運営者として今年 1 年を過ごしてきました。
その経験からサークルの新歓に特化したプラットフォームが "With コロナ"時代のサークル運営には必要だと強く感じ、制作するに至りました。

今年はコロナ騒動で大学の新歓が禁止され大きな影響を受けました。私はサークル運営側としてサークルの HP を立ち上げたり、公式ラインを作ったり、YouTube チャンネルを作成したり、Zoom でオンライン新歓を行ったりして対応しました。
そして夏休みが明け 9 月。ついに新歓が解禁されました。新歓で 1 年生と実際にあって話してみると驚くべき発見がありました。自分ではオンラインで自分のサークルの魅力を普段の新歓ほどではないものの、十分伝えたと思っていたものの、HP や YouTube チャンネルや、Zoom のオンラインイベントなどのその存在そのものを知らなかった 1 年生が圧倒的多数でした。
すなわち、1 年生はサークルのことを知りたがっていて、サークルは自分たちのことを知ってほしいという両思いの関係にも関わらず、お互いの存在を認知できる術が現状存在しないことに気づきました。
このときサークルの新歓に特化したプラットフォームが "With コロナ"時代のサークル運営には必要だと強く感じ、制作するに至りました。

## 機能一覧

- サインアップ
- 退会
- ログイン・ログアウト
- 投稿の閲覧・投稿・編集・削除
- 投稿のサムネイル画像の設定・編集
- マイページ閲覧・編集
- プロフィール画像の設定・編集
- ユーザー基本情報の編集・更新
- サークルの作成・編集
- サークルへの入会・退会
- サークル加入者のみ投稿・閲覧可能なチャット機能
- サークル管理者のみ上のツイートを削除する権限を持ち、削除することが可能
- ツイッター連携機能
- ログ収集し、ターミナルと log.json に出力

## 使用技術

- Git:2.23.0
- Github
- Bootstrap:4.5.0
- Go:1.15.0
- MySQL:8.20.0
- Docker:19.03.13-beta2
- Docker Compose:1.27.0
- Visual Studio Code:1.49
- CircleCI (CI/CD)
- Terraform:0.13.3
- AWS
  - VPC
  - S3
  - RDS
  - ECR
  - ECS (Fargate)
  - ACM
  - SSM
  - Route53
  - Cloud Watch
- JavaScript
- CSS
- HTML5

## インフラ構成図

![インフラ構成図](/assets/infra_structures.png)
