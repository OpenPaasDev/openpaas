CREATE TABLE IF NOT EXISTS datacenters(
		id TEXT PRIMARY KEY,
		name TEXT,
		region TEXT,
		availability_zone TEXT
	);
CREATE TABLE IF NOT EXISTS server_groups(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		dc_id TEXT,
		FOREIGN KEY(dc_id) REFERENCES datacenters(id)
	);

CREATE TABLE IF NOT EXISTS servers(
        id TEXT PRIMARY KEY,
        name TEXT,
        public_ip TEXT,
        private_ip TEXT,
        hostname TEXT,
        dc_name TEXT,
        is_lb_target BOOLEAN,
        instance_type TEXT,
        server_group_id INTEGER,
        FOREIGN KEY(server_group_id) REFERENCES server_groups(id)
    );

CREATE TABLE IF NOT EXISTS mounts(
		id TEXT PRIMARY KEY,
		server_id TEXT,
		local_path TEXT,
		mount_path TEXT,
		owner TEXT,
		FOREIGN KEY(server_id) REFERENCES servers(id)
	);

CREATE TABLE IF NOT EXISTS k3s_members(
		id TEXT PRIMARY KEY,
		server_id TEXT,
		state TEXT,
		role TEXT, -- server or agent
		FOREIGN KEY(id) REFERENCES servers(id)
	);
