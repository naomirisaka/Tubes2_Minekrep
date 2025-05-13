# <p align = center>Tugas Besar 2 IF2211 Strategi Algoritma</p>
# <p align = center>Pemanfaatan Algoritma BFS dan DFS dalam Pencarian Recipe pada Permainan Little Alchemy 2</p>

![image](https://github.com/user-attachments/assets/397a1b88-0a0c-4415-9d62-bcfb6bebf7be)

### Kelompok 28: Minekrep
| Nama | NIM |
|------|-----|
| Indah Novita Tangdililing | 13523047 |
| Bevinda Vivian | 13523120 |
| Naomi Risaka Sitorus | 13523122 |

## Deskripsi
Web Little Alchemy 2 Recipe Finder mampu membantu pengguna mencari tahapan-tahapan atau resepyang perlu dilakukan untuk membentuk sebuah elemen target. Jumlah resep untuk membentuk elemen yang dicari bisa berjumlah satu atau lebih. Pencarian multiple recipes akan diproses secara multithreading demi efisiensi waktu. Pencarian dilkaukan menggunakan algoritma BFS atau BFS sesuai pilihan pengguna dengan prinsip setiap elemen harus dapat dibentuk dari keempat elemen dasar, yaitu fire, air, water, dan earth.

## Fitur
### 1. Scraping 
Mengekstrak resep dari https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2) saat website pertama kali website dijalankan dan hasilnya disimpan dalam suatu file JSON.
### 2. Pencarian Resep Secara BFS
Mencari rute untuk menghasilkan elemen target (simpul tujuan) dengan membangun graf secara iteratif atau per tingkat. 
### 3. Pencarian Resep Secara DFS
Mencari rute untuk menghasilkan elemen target (simpul tujuan) dengan membangun graf secara rekursif hingga mencapai target atau simpul daun yang tidak bisa diekspansi lagi.
### 4. Visualisasi Pohon dan Statistik Pencarian
Menampilkan pohon resep hasil pencarian elemen target serta statistik pencarian berupa waktu yang diperlukan dan jumlah simpul yang dikunjungi.

## Struktur
```bash
├── docs
│   └── Minekrep.pdf
├── src
│   ├── backend
│   │   ├── api                   # Mengirim informasi dari backend ke frontend
│   │   │    └── handlers.go
│   │   ├── scraper               # Mengekstrak data resep dari website ke JSON
│   │   │    └── scrap.go
│   │   ├── searchalgo            # Algoritma pencarian resep
│   │   │    ├── bfs.go
│   │   │    ├── bidirectional.go
│   │   │    └── dfs.go
│   │   ├── utilities             # Struktur data dan helper function  untuk pencarian
│   │   |    ├── models.go
│   │   |    └── utils.go
│   │   └── Dockerfile            # Menciptakan container backend untuk Docker
│   ├── frontend
│   │   ├── minekrep
│   │   │   ├── public
│   │   │   │   └── images
│   │   │   ├── src
│   │   │   │   ├── app
│   │   │   │   │   ├── favicon.ico
│   │   │   │   │   ├── globals.css
│   │   │   │   │   ├── layout.js
│   │   │   │   │   ├── page.js     # Landing page
│   │   │   │   │   ├── profile     # Halaman profil
│   │   │   │   │   │   └── page.js
│   │   │   │   │   └── search      # Halaman pencarian
│   │   │   │   │       └── page.js
│   │   │   │   ├── components
│   │   │   │   │   ├── AlgorithmSelector.jsx
│   │   │   │   │   ├── MinecraftButton.jsx
│   │   │   │   │   ├── NavBar.jsx
│   │   │   │   │   ├── RecipeVisualizer.jsx
│   │   │   │   │   ├── SearchForm.jsx
│   │   │   │   │   └── TeamProfile.jsx
│   │   │   │   ├── styles
│   │   │   │   │   └── globals.css
│   │   │   │   ├── utils
│   │   │   │   |   └── api.js      # Memproses informasi dari backend di frontend
│   │   │   │   ├── Dockerfile      # Menciptakan container backend untuk Docker
│   │   │   │   └── ...
│   │   │   └── package-lock.json
│   └── docker-compose.yml
├── .gitignore
└── README.md
```

## Requirements
- Go Language v1.23 atau lebih baru
- Node JS v20 atau lebih baru
- Docker Desktop (jika ingin menjalankan dengan Docker)
   
## Cara Menjalankan (Tanpa Docker)
1. Clone repository ini dengan menjalankan perintah di bawah ini pada terminal IDE yang mendukung Go:
   ```sh
   git clone https://github.com/naomirisaka/Tubes2_Minekrep.git
2. Buka folder hasil clone di IDE.
3. Pindah ke directory backend dengan:
   ```sh
   cd src/backend
4. Jalankan backend dengan:
    ```sh
    run go main.go
5. Buka terminal baru di IDE dan pindah ke directory backend dengan:
   ```sh
   cd src/frontend/minekrep
6. Inisialisasi frontend dengan:
   ```sh
   npm i
7. Jalankan frontend dengan:
   ```sh
   npm run dev
8. Akses website melalui link yang diberikan di terminal frontend, seperti `http://localhost:3000`

## Cara Menjalankan (Dengan Docker)
1. Clone repository ini dengan menjalankan perintah di bawah ini pada terminal IDE yang mendukung Go:
   ```sh
   git clone https://github.com/naomirisaka/Tubes2_Minekrep.git
2. Buka folder hasil clone di IDE.
3. Pindah ke directory src dengan:
   ```sh
   cd src
4. Jalankan Docker dengan:
    ```sh
    docker-compose up --build
5. Akses website melalui link yang diberikan oleh container frontend, seperti `http://localhost:3000`
