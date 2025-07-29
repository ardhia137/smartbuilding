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

-- 8.1. Fix foreign key constraint issues
-- Temporarily disable foreign key checks
SET FOREIGN_KEY_CHECKS = 0;

-- Check and fix orphaned records in hak_akses table
-- Delete hak_akses records that reference non-existent users
DELETE FROM hak_akses WHERE user_id NOT IN (SELECT id FROM user);

-- Check and fix orphaned records that reference non-existent gedung
DELETE FROM hak_akses WHERE gedung_id NOT IN (SELECT id FROM gedung);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;

-- Verifikasi hasil migration
SHOW TABLES;
DESCRIBE gedung;
DESCRIBE torent;
DESCRIBE hak_akses;
DESCRIBE monitoring_data;
DESCRIBE monitoring_data_harian;

-- 9. Merge monitoring_data and monitoring_data_harian tables into monitoring_logs
-- First, create the new monitoring_logs table
CREATE TABLE monitoring_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    monitoring_name VARCHAR(255) NOT NULL,
    monitoring_value VARCHAR(50) NOT NULL,
    id_gedung INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_monitoring_logs_gedung (id_gedung),
    INDEX idx_monitoring_logs_name (monitoring_name),
    FOREIGN KEY (id_gedung) REFERENCES gedung(id)
);

-- Copy data from monitoring_data to monitoring_logs
INSERT INTO monitoring_logs (monitoring_name, monitoring_value, id_gedung, created_at, updated_at)
SELECT monitoring_name, monitoring_value, id_gedung, created_at, updated_at
FROM monitoring_data;

-- Copy data from monitoring_data_harian to monitoring_logs
INSERT INTO monitoring_logs (monitoring_name, monitoring_value, id_gedung, created_at, updated_at)
SELECT monitoring_name, monitoring_value, id_gedung, created_at, updated_at
FROM monitoring_data_harian;

-- Drop the old tables after data migration
-- DROP TABLE monitoring_data;
-- DROP TABLE monitoring_data_harian;

-- Verify the new table
DESCRIBE monitoring_logs;
SELECT COUNT(*) as total_records FROM monitoring_logs;
