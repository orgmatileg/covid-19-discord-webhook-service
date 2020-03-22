# COVID-19 Discord Webhook Service Indonesia
Setiap orang saat ini ingin mengetahui perkembangan tentang COVID-19 yang terjadi di Indonesia, maka dari itu saya membuat Service ini untuk mengetahui perkembangan COVID-19 Melalui Discord Channel

## Project Goals
- Mengetahui angka terkonfirmasi
- Mengetahui angka dan persentasai telah sembuh
- Mengetahui angka dan persentasi kematian

## Requirement
- Docker

## How To
Aplikasi ini bisa berjalan di pc dengan cara build sendiri menggunakan go atau menggunakan Docker Image yang sudah saya sediakan.

Untuk menjalankan aplikasi ini menggunakan docker ada dibawah ini:

Saya asumsikan kalian bekerja di direktory ~/covid

- Buat file dengan nama config.yaml and copy paste kode dibawah ini:
```
# 1s = 1 Second
# 1m = 1 Minute
# 1h = 1 Hour
check_every: 1m
discord_webhook: "FILL YOUR DISCORD WEBHOOK URL"
```
Dont forget change discord_webhook with your webhook url

- Cara :
```
docker run -v ${pwd}/config.yaml:/app/config.yaml docker.pkg.github.com/orgmatileg/covid-19-discord-webhook-service-indonesia/covid-19-discord-webhook-service-indonesia:latest
```


# Author
- Luqmanul Hakim

## Referensi
- https://github.com/mathdroid/covid-19-api