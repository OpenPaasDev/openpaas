CREATE TABLE IF NOT EXISTS datacenters(
		id TEXT PRIMARY KEY,
		Idegion TEXT
	);
CREATE TABLE IF NOT EXISTS server_groups(
		id TEXT PRIMARY KEY,
		dc_id TEXT,
		FOREIGN KEY(dc_id) REFERENCES datacenters(id)
	);

CREATE TABLE IF NOT EXISTS servers(
        id TEXT PRIMARY KEY,
        public_ip TEXT,
        private_ip TEXT,
        hostname TEXT,
        is_lb_target BOOLEAN,
        instance_type TEXT,
        server_group_id TEXT,
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

CREATE TABLE IF NOT EXISTS k3s_clusters(
	id TEXT PRIMARY KEY
	);

CREATE TABLE IF NOT EXISTS k3s_members(
		id TEXT PRIMARY KEY,
		server_id TEXT,
		state TEXT,
		role TEXT, -- server or agent
		cluster_id TEXT,
		FOREIGN KEY(id) REFERENCES servers(id),
		FOREIGN KEY(cluster_id) REFERENCES k3s_clusters(id)
	);
