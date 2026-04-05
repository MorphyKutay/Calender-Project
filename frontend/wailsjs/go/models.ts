export namespace main {
	
	export class Event {
	    id: string;
	    date: string;
	    title: string;
	    color: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = source["date"];
	        this.title = source["title"];
	        this.color = source["color"];
	    }
	}

}

