export namespace main {
	
	export class dirNode {
	    files: string[];
	    subdirs: {[key: string]: dirNode};
	
	    static createFrom(source: any = {}) {
	        return new dirNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = source["files"];
	        this.subdirs = source["subdirs"];
	    }
	}

}

