export namespace entities {
	
	export class Client {
	    id: number;
	    name: string;
	    email: string;
	    phone: string;
	    address: string;
	    contact_person: string;
	    notes: string;
	    status: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Client(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.contact_person = source["contact_person"];
	        this.notes = source["notes"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Sort {
	    field: string;
	    order: string;
	
	    static createFrom(source: any = {}) {
	        return new Sort(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.field = source["field"];
	        this.order = source["order"];
	    }
	}
	export class Pagination {
	    page: number;
	    page_size: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new Pagination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	        this.total = source["total"];
	    }
	}
	export class ClientQueryParams {
	    id_in: number[];
	    name: string;
	    name_like: string;
	    email: string;
	    email_like: string;
	    phone: string;
	    phone_like: string;
	    address_like: string;
	    contact_person_like: string;
	    notes_like: string;
	    status: number;
	    status_in: number[];
	    // Go type: time
	    created_at_gte?: any;
	    // Go type: time
	    created_at_lte?: any;
	    // Go type: time
	    updated_at_gte?: any;
	    // Go type: time
	    updated_at_lte?: any;
	    // Go type: Pagination
	    pagination?: any;
	    sorts?: Sort[];
	
	    static createFrom(source: any = {}) {
	        return new ClientQueryParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id_in = source["id_in"];
	        this.name = source["name"];
	        this.name_like = source["name_like"];
	        this.email = source["email"];
	        this.email_like = source["email_like"];
	        this.phone = source["phone"];
	        this.phone_like = source["phone_like"];
	        this.address_like = source["address_like"];
	        this.contact_person_like = source["contact_person_like"];
	        this.notes_like = source["notes_like"];
	        this.status = source["status"];
	        this.status_in = source["status_in"];
	        this.created_at_gte = this.convertValues(source["created_at_gte"], null);
	        this.created_at_lte = this.convertValues(source["created_at_lte"], null);
	        this.updated_at_gte = this.convertValues(source["updated_at_gte"], null);
	        this.updated_at_lte = this.convertValues(source["updated_at_lte"], null);
	        this.pagination = this.convertValues(source["pagination"], null);
	        this.sorts = this.convertValues(source["sorts"], Sort);
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

