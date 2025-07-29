-- Migration script untuk mengubah nama table dan column sesuai dengan refactoring kode
-- Jalankan script ini di MySQL database untuk menyelaraskan database dengan kode baru

-- 1. Rename table 'setting' menjadi 'gedung'
ALTER TABLE setting RENAME TO gedung;

-- 2. Rename table 'data_torent' menjadi 'torent'
ALTER TABLE data_torent RENAME TO torent;

-- 3. Rename table 'pengelola_gedung' menjadi 'hak_akses'
ALTER TABLE pengelola_gedung RENAME TO hak_akses;

-- 4. Update column names di table 'monitoring_data'
ALTER TABLE monitoring_data CHANGE id_setting id_gedung INT NOT NULL;

-- 5. Update column names di table 'monitoring_data_harian'
ALTER TABLE monitoring_data_harian CHANGE id_setting id_gedung INT NOT NULL;

-- 6. Update column names di table 'hak_akses' (formerly pengelola_gedung)
ALTER TABLE hak_akses CHANGE setting_id gedung_id INT NOT NULL;

-- 7. Update column names di table 'torent' (formerly data_torent)
ALTER TABLE torent CHANGE id_setting id_gedung INT NOT NULL;

-- 8. Update foreign key constraints jika ada (adjust sesuai dengan constraint yang ada)
-- ALTER TABLE monitoring_data DROP FOREIGN KEY IF EXISTS fk_monitoring_data_setting;
-- ALTER TABLE monitoring_data ADD CONSTRAINT fk_monitoring_data_gedung FOREIGN KEY (id_gedung) REFERENCES gedung(id);

-- ALTER TABLE monitoring_data_harian DROP FOREIGN KEY IF EXISTS fk_monitoring_data_harian_setting;
-- ALTER TABLE monitoring_data_harian ADD CONSTRAINT fk_monitoring_data_harian_gedung FOREIGN KEY (id_gedung) REFERENCES gedung(id);

-- ALTER TABLE hak_akses DROP FOREIGN KEY IF EXISTS fk_hak_akses_setting;
-- ALTER TABLE hak_akses ADD CONSTRAINT fk_hak_akses_gedung FOREIGN KEY (gedung_id) REFERENCES gedung(id);

-- ALTER TABLE torent DROP FOREIGN KEY IF EXISTS fk_torent_setting;
-- ALTER TABLE torent ADD CONSTRAINT fk_torent_gedung FOREIGN KEY (id_gedung) REFERENCES gedung(id);

-- Verifikasi hasil migration
SHOW TABLES;
DESCRIBE gedung;
DESCRIBE torent;
DESCRIBE hak_akses;
DESCRIBE monitoring_data;
DESCRIBE monitoring_data_harian;
