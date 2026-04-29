-- Integration test seed data for mysql-align
-- Target: MySQL 5.7

-- fixture_clean: positive control — must NOT appear in any check result
CREATE SCHEMA IF NOT EXISTS fixture_clean DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_clean.users (
    id   INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(64),
    PRIMARY KEY (id)
) ENGINE=InnoDB ROW_FORMAT=DYNAMIC;

-- fixture_nopk: triggers CheckNoPK (table with no primary key)
CREATE SCHEMA IF NOT EXISTS fixture_nopk DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_nopk.events (
    id      INT,
    payload TEXT
) ENGINE=InnoDB;

-- fixture_utf8: triggers CheckCharSet with Severity=ERROR (charset=utf8)
CREATE SCHEMA IF NOT EXISTS fixture_utf8 DEFAULT CHARACTER SET utf8;
CREATE TABLE fixture_utf8.t1 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;

-- fixture_latin1: triggers CheckCharSet with Severity=WARNING (charset=latin1)
CREATE SCHEMA IF NOT EXISTS fixture_latin1 DEFAULT CHARACTER SET latin1;
CREATE TABLE fixture_latin1.t1 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;

-- fixture_engine: triggers CheckEngine (MyISAM and MEMORY tables)
CREATE SCHEMA IF NOT EXISTS fixture_engine DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_engine.t_myisam (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=MyISAM;
CREATE TABLE fixture_engine.t_mem    (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=MEMORY;

-- fixture_rowfmt: triggers CheckRowFormat (COMPACT and REDUNDANT row formats)
CREATE SCHEMA IF NOT EXISTS fixture_rowfmt DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_rowfmt.t_compact   (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB ROW_FORMAT=COMPACT;
CREATE TABLE fixture_rowfmt.t_redundant (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB ROW_FORMAT=REDUNDANT;

-- fixture_fkdup: used by TestCheckFKDuplication.
-- MySQL 5.7 InnoDB enforces unique FK names per schema at the engine level,
-- so duplicate FK constraint names cannot be created — the check always returns
-- empty on 5.7 (which is the correct behavior: no duplicates exist to find).
-- Distinct names are used here to prevent a seed error that would abort the file.
CREATE SCHEMA IF NOT EXISTS fixture_fkdup DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_fkdup.parent_a (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;
CREATE TABLE fixture_fkdup.parent_b (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;
CREATE TABLE fixture_fkdup.child_a (
    id   INT NOT NULL,
    a_id INT,
    PRIMARY KEY (id),
    CONSTRAINT fk_to_parent_a FOREIGN KEY (a_id) REFERENCES fixture_fkdup.parent_a(id)
) ENGINE=InnoDB;
CREATE TABLE fixture_fkdup.child_b (
    id   INT NOT NULL,
    b_id INT,
    PRIMARY KEY (id),
    CONSTRAINT fk_to_parent_b FOREIGN KEY (b_id) REFERENCES fixture_fkdup.parent_b(id)
) ENGINE=InnoDB;

-- fixture_views: triggers CheckViewDeprecated (view with GROUP BY ... ASC)
CREATE SCHEMA IF NOT EXISTS fixture_views DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_views.t_base (a INT, b INT) ENGINE=InnoDB;
CREATE VIEW fixture_views.v_bad AS
    SELECT a, COUNT(*) AS c FROM fixture_views.t_base GROUP BY a DESC;

-- fixture_routines: triggers CheckRoutineSyntaxDeprecated and CheckRoutineFunctionDeprecated
-- Single-statement bodies avoid needing DELIMITER.
CREATE SCHEMA IF NOT EXISTS fixture_routines DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_routines.t_base (x INT, y INT) ENGINE=InnoDB;

CREATE PROCEDURE fixture_routines.p_bad_groupby()
    SELECT x, COUNT(*) FROM fixture_routines.t_base GROUP BY x DESC;

CREATE FUNCTION fixture_routines.f_decode(ciphertext BLOB, pass_str VARCHAR(64))
    RETURNS BLOB DETERMINISTIC RETURN DECODE(ciphertext, pass_str);

CREATE PROCEDURE fixture_routines.p_compress()
    SELECT COMPRESS('hello world');

-- fixture_rows: exact row counts used by TestReconcileRow
CREATE SCHEMA IF NOT EXISTS fixture_rows DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE fixture_rows.t_a (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;
CREATE TABLE fixture_rows.t_b (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB;
INSERT INTO fixture_rows.t_a (id) VALUES (1),(2),(3);
INSERT INTO fixture_rows.t_b (id) VALUES (1),(2),(3),(4),(5);
