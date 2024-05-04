ProjectSprint! Batch 2 Requirement Gathering Week 1 | CatsSocial

https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769

Meeting Notes
Perbedaan engineer yang bagus dan tidak adalah mereka yang aktif memperjelas requirements.

Use Golang, Postgres, run in 8080, no ORM/Query Generator RAW only.


QnA
1.	User sama Cat terpisah, satu user bisa punya banyak cat.
2.	Social login boleh?
MVP, ngeload-test, nanti Google marah.

Load test harus mensimulate banyak user.
Google ada limit based on IP.
Misal load test 10rb user, yang ada IP kita akan diblacklist.

Mungkin akan coba mock flow, phase berikutnya.
Integrasi ke mock API Google.A
3.	Ini harus punya kita catnya?
Owned punya kita
4.	Wujud infranya gimana?
SCP kita ke server backend kalian, server K6 load test hit ke BE.
 
Makannya ada banyak room showcase.
1 showcase 5-10 menit.

Awal-awal SCP dulu.
Ntar pakai Amazon ECS.
5.	ECS itu apa?
Amazon ECS sebuah cara biar kita bisa run aplikasi kita yang sudah dibungkus docker.
Ada beberapa tahap.

Ada w1, w2, w3, w4
W1: manual, upload scp
W2: manual, git clone, compile docker di server.
W3: otomatis CI/CD, push docker ke registry, Run di ECR (disetupin Mas Nanda)
W4: otomatis CI/CD, kalian setup semua sendiri di ECR.

W1, W2, W3, pelan-pelan dibiasakan interaksi ke deployment sampai W4.

Runner CI/CD kita yang menyiapkan.
 

Repo apa bebas, Gitlab boleh, Github boleh.

Unit test tidak perlu, belum perlu.

6.	Operasi delete cat? Kalau sudah match bagaimana?
Ekspektasinya yang match tetap muncul.
Match kedelete, tetap muncul hasil match-nya.

7.	Jumlah table users, cats, bebas?
Bebas.

8.	Perlu sofiate?
Table berapa, framework apa, bebas.
Yang penting Menuhin requirement.
9.	Boleh E6 batch sebelumnya?
Boleh
10.	Matchnya tetap ada, itu kalau diupdate?
Ketika sudah diupdate apakah bisa diedit? Tidak.
Match, kucing A, diupdate dari kucing B Namanya, di match juga update.
11.	Kalau sudah disapprove apakah bisa didelete?
Masih bisa.
Kucingnya bisa diedit didelete, sudah disapprove maupun belum.

12.	Kalau gender kucing diganti?
Sex is edited when cat is already matched.

13.	Kalau mau match?
Tidak boleh match kucing yang sudah match.
14.	Kalau udah match, misal ada beberapa kucing yang dimatch-kan, misal 3 orang.?
Once a match is approved, other match req that matches bot hthe issuer and the receiver cat’s, will get removed.
15.	Showcase ngapain?
Migrasi, test cases, load test.
Kemarin test casesnya ada

Git /nandanugg nandanugg (Nanda Nugraha) (github.com)

Di script akan tes beberapa cases.
Pertama akan ngambil Oauth, Invalid body, harus hati-hati invalid body, akan test semua. Akan dicoba semua kind invalid.
Semua API akan dibomb dengan invalid body.
Akan mencoba test rutenya dengan hal-hal yang tidak sesuai kontrak. Invalid payload.
Transaction, apakah historynya keluar.
Upload file.
Header kosong gimana.

Casesnya dibuatin.
16.	Kalau sudah match gaboleh ganti sex, kalau masih approval?
Sex is edited when cat is already requested to match.
Edit kelamin kucing bisa.
Datanya aja sih. Misal salah input data.

Kita asumsikan ada user error.
17.	Ras kucing selain di enum?
Tidak boleh.
18.	Jodohin kucing sendiri boleh tak?
Tak boleh sama.
19.	Image kucing upload perlu index url?
W2-W3-W4 gunakan s3 bucket.
20.	Image url ada allowed extension tak?
Yang penting url dulu. Kalau beneran upload harus upload dulu.
21.	Kucingnya bisa cerai tak?
Mau? Gausah.
Ada hidden requirement.
22.	Cases dishare kapan?
Kamis atau Jumat random.
Test cases biasanya keluar pas akhiran.
Sebenernya sudah bikin, tapi bisa jadi ada perubahan.
23.	Delete match hanya issuer?
Match can only be deleted by issuer.
24.	1 kucing bisa kawin berapa kali?
Iya
25.	Kalau poligami boleh?
Mau diubah? Ga usah.
Kalau count harus buat history, bisa diget ada filternya.
Emang erd bisa lihat di mana?
26.	Tidak jelas, memperjelas, sama nambah requirement itu beda.
Bener-bener memperjelas requirement, tapi tujuannya ga nambah kerjaan.
Senior mampu memperjelas requirement.
Tahu future risk, case tak tercover.
Poligami itu ada.

https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/PostgreSQL.Concepts.General.SSL.html
