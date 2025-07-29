# Database Migration Guide

## Penjelasan Error 1054: "Unknown column 'id_gedung' in 'where clause'"

Error ini terjadi karena ada **ketidaksesuaian antara kode aplikasi dan skema database MySQL**:

### Penyebab:
1. **Kode aplikasi** sudah direfactor untuk menggunakan nama baru:
   - `id_gedung` (bukan `id_setting`)
   - `gedung_id` (bukan `setting_id`)
   - table `gedung` (bukan `setting`)
   - table `hak_akses` (bukan `pengelola_gedung`)
   - table `torent` (bukan `data_torent`)

2. **Database schema** masih menggunakan nama lama:
   - Column `id_setting` 
   - Column `setting_id`
   - Table `setting`
   - Table `pengelola_gedung`
   - Table `data_torent`

### Solusi:

#### 1. Jalankan Migration SQL
Eksekusi file `migration.sql` di database MySQL Anda:

```bash
mysql -u [username] -p [database_name] < migration.sql
```

Atau copy-paste isi file `migration.sql` ke MySQL console.

#### 2. Verifikasi Migration
Setelah migration, verifikasi bahwa:
- Table `setting` → `gedung`
- Table `pengelola_gedung` → `hak_akses` 
- Table `data_torent` → `torent`
- Column `id_setting` → `id_gedung`
- Column `setting_id` → `gedung_id`

#### 3. Restart Aplikasi
Setelah database migration selesai, restart aplikasi Go:

```bash
go run main.go
```

### Perubahan yang Dilakukan dalam Kode:

#### Entity Changes:
- `entities/setting.go` → `entities/gedung.go`
- `entities/pengelola_gedung.go` → `entities/hak_akses.go`
- `entities/data_torent.go` → `entities/torent.go`

#### Field Name Changes:
- `IDSetting` → `IDGedung`
- `SettingID` → `GedungID`
- `PengelolaGedung` → `HakAkses`

#### API Endpoint Changes:
- `/setting/*` → `/gedung/*`
- `/pengelola-gedung/*` → `/hak-akses/*`
- `/data-torent/*` → `/torent/*`

#### SQL Query Updates:
- `WHERE id_setting = ?` → `WHERE id_gedung = ?`
- `JOIN setting s` → `JOIN gedung g`
- `ha.setting_id` → `ha.gedung_id`

### Testing
Setelah migration, test endpoint berikut:
1. `GET /gedung` - List semua gedung
2. `GET /hak-akses` - List semua hak akses
3. `GET /torent` - List semua torent
4. `GET /monitoring-data/air/{id}` - Get monitoring data
5. `POST /login` - Login user

### Backup Database
**PENTING**: Backup database sebelum menjalankan migration:

```bash
mysqldump -u [username] -p [database_name] > backup_before_migration.sql
```

Jika ada masalah, restore dengan:

```bash
mysql -u [username] -p [database_name] < backup_before_migration.sql
```
