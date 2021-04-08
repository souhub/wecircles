# README.md

## WECIRCLES

~~完成品はこちらです。→https://wecircles.net~~
<br/>公開停止しました。　※2021/4/8 追記<br/>

スマホで使用されることを想定していますが、パソコンでもタブレットでも使用可能になってます。

※ロードバランサーの負荷分散の影響でプロフィール画像のアップロードのみ数回に一回失敗するエラーが発生します。
数回繰り返していただければできますので、その場合はそうしていただけると幸いです。

### 概要

大学生向けの新歓のための SNS です。新入生とサークル運営者または新入生同士の出会いの場を提供し、サークルのイベントや練習へ参加するきっかけ作りをになることを目的としています。</br>
自分はサークルの副会長という立場だったことから新歓に精通しています。その経験から新歓を分解すると次の図のようにの 4 つのステップからなります。</br>

- 【発見】:サークルを見つける
- 【交流】:ブースで交流する
- 【参加】:イベントや練習に参加する
- 【入会】:正式にサークルに加入する</br>

WeCircles は上の 【発見】と【交流】 を行える場を提供し、新入生とサークルの出会いのきっかけ作りを行います。
![概要図](/assets/shinkan_flow.png)

### 制作の背景

私は大学でサークルの運営者として今年 1 年を過ごしてきました。
その経験からサークルの新歓に特化したプラットフォームが "With コロナ"時代のサークル運営には必要だと強く感じ、制作するに至りました。

今年はコロナ騒動で大学の新歓が禁止され大きな影響を受けました。私はサークル運営側としてサークルの HP を立ち上げたり、公式ラインを作ったり、YouTube チャンネルを作成したり、Zoom でオンライン新歓を行ったりして対応しました。
そして夏休みが明け 9 月。ついに新歓が解禁されました。新歓で 1 年生と実際にあって話してみると驚くべき発見がありました。自分ではオンラインで自分のサークルの魅力を普段の新歓ほどではないものの、十分伝わっていると思っていました。しかし、せっかく作った HP や YouTube チャンネルや、Zoom のオンラインイベントなどのその存在そのものを知らなかった 1 年生が圧倒的多数でした。
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
- ブースの作成・編集
- ブースへの参加・退出
- ブース参加者のみ投稿・閲覧可能なチャット機能
- ブース管理者のみ上のツイートを削除する権限を持ち、削除することが可能
- Twitter 連携機能
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
