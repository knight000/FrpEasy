export namespace frpc {
	
	export class TemplateDisplayInfo {
	    NamePattern: string;
	    Protocol: string;
	    LocalPorts: string;
	    RemotePorts: string;
	    RawContent: string;
	
	    static createFrom(source: any = {}) {
	        return new TemplateDisplayInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.NamePattern = source["NamePattern"];
	        this.Protocol = source["Protocol"];
	        this.LocalPorts = source["LocalPorts"];
	        this.RemotePorts = source["RemotePorts"];
	        this.RawContent = source["RawContent"];
	    }
	}

}

export namespace models {
	
	export class Service {
	    id: string;
	    name: string;
	    protocol: string;
	    local_ip: string;
	    local_port: number;
	    remote_port: number;
	    use_encryption: boolean;
	    use_compression: boolean;
	    advanced_config: string;
	    is_advanced: boolean;
	    display_ports?: string;
	
	    static createFrom(source: any = {}) {
	        return new Service(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.protocol = source["protocol"];
	        this.local_ip = source["local_ip"];
	        this.local_port = source["local_port"];
	        this.remote_port = source["remote_port"];
	        this.use_encryption = source["use_encryption"];
	        this.use_compression = source["use_compression"];
	        this.advanced_config = source["advanced_config"];
	        this.is_advanced = source["is_advanced"];
	        this.display_ports = source["display_ports"];
	    }
	}
	export class LogEntry {
	    id: string;
	    timestamp: number;
	    message: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = source["timestamp"];
	        this.message = source["message"];
	        this.type = source["type"];
	    }
	}
	export class Server {
	    id: string;
	    name: string;
	    address: string;
	    port: number;
	    token: string;
	    status: string;
	    enabled: boolean;
	    logs: LogEntry[];
	    uptime: number;
	
	    static createFrom(source: any = {}) {
	        return new Server(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.token = source["token"];
	        this.status = source["status"];
	        this.enabled = source["enabled"];
	        this.logs = this.convertValues(source["logs"], LogEntry);
	        this.uptime = source["uptime"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Preset {
	    id: string;
	    name: string;
	    servers: Server[];
	    services: Service[];
	
	    static createFrom(source: any = {}) {
	        return new Preset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.servers = this.convertValues(source["servers"], Server);
	        this.services = this.convertValues(source["services"], Service);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportResult {
	    preset?: Preset;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preset = this.convertValues(source["preset"], Preset);
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class ServerRuntime {
	    preset_id: string;
	    server_id: string;
	    process_pid: number;
	    config_path: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerRuntime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preset_id = source["preset_id"];
	        this.server_id = source["server_id"];
	        this.process_pid = source["process_pid"];
	        this.config_path = source["config_path"];
	    }
	}

}

