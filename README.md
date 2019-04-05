# Cloud build Badge

![build_badge](https://storage.googleapis.com/cloud-build-result/star-integrations-ci/github_star-integrations_cloud-build-badge-go/feature/build-badge/badge.svg)

## 使い方

1. CloudBuildサービスアカウントにプロジェクト編集者の権限付与
1. CloudFunctionsAPIを有効化
1. CloudStorageの{BucketName}に `success.svg` と `failure.svg` を配置
1. {BucketName}をbuildenv.yamlの `BUILD_RESULT_BADGE_BUCKET` に記載
1. `gcloud builds submit . --project=${ProjectID}` を実行
1. CloudStorageの `{BucketName}/{projectID}/{RepositoryName}/{Branch}/badge.svg` にビルド結果を配置

### 画像キャッシュ

バッジ画像がキャッシュされビルド結果が即時反映されない時は `Cache-Control` を設定する。  
元のバッジ画像に設定すればそれがビルド結果ファイルに反映される。

* コンソール画面からメタデータを編集する。 [ヘルプ](https://cloud.google.com/storage/docs/viewing-editing-metadata?hl=ja#edit)

* gsutilコマンドで設定する

    ```bash
    ex)
    gsutil setmeta -h "Cache-Control:no-cache" gs://{BucketName}/success.svg
    gsutil setmeta -h "Cache-Control:no-cache" gs://{BucketName}/failure.svg
    ```