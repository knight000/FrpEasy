export namespace main {
	
	export class ImportResult {
	    preset?: models.Preset;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preset = this.convertValues(source["preset"], models.Preset);
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

}

export namespace models {
	
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
	export class Service {
	    id: string;
	    name: string;
	    protocol: string;
	    localIp: string;
	    localPort: number;
	    remotePort: number;
	    useEncryption: boolean;
	    useCompression: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Service(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.protocol = source["protocol"];
	        this.localIp = source["localIp"];
	        this.localPort = source["localPort"];
	        this.remotePort = source["remotePort"];
	        this.useEncryption = source["useEncryption"];
	        this.useCompression = source["useCompression"];
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
	
	export class ServerRuntime {
	    presetId: string;
	    serverId: string;
	    processPid: number;
	    configPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerRuntime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.presetId = source["presetId"];
	        this.serverId = source["serverId"];
	        this.processPid = source["processPid"];
	        this.configPath = source["configPath"];
	    }
	}

}

