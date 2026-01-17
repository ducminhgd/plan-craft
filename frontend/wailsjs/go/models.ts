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
	export class ClientListResponse {
	    data: Client[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ClientListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], Client);
	        this.total = source["total"];
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
	export class HumanResource {
	    id: number;
	    name: string;
	    title: string;
	    level: string;
	    status: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new HumanResource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.title = source["title"];
	        this.level = source["level"];
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
	export class HumanResourceListResponse {
	    data: HumanResource[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new HumanResourceListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], HumanResource);
	        this.total = source["total"];
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
	export class HumanResourceQueryParams {
	    id_in: number[];
	    name: string;
	    name_like: string;
	    title: string;
	    title_like: string;
	    level: string;
	    level_like: string;
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
	        return new HumanResourceQueryParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id_in = source["id_in"];
	        this.name = source["name"];
	        this.name_like = source["name_like"];
	        this.title = source["title"];
	        this.title_like = source["title_like"];
	        this.level = source["level"];
	        this.level_like = source["level_like"];
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
	export class ProjectResource {
	    id: number;
	    project_id: number;
	    human_resource_id: number;
	    role: string;
	    allocation: number;
	    // Go type: time
	    start_date?: any;
	    // Go type: time
	    end_date?: any;
	    notes: string;
	    status: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    project?: Project;
	    human_resource?: HumanResource;
	
	    static createFrom(source: any = {}) {
	        return new ProjectResource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.human_resource_id = source["human_resource_id"];
	        this.role = source["role"];
	        this.allocation = source["allocation"];
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.end_date = this.convertValues(source["end_date"], null);
	        this.notes = source["notes"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.project = this.convertValues(source["project"], Project);
	        this.human_resource = this.convertValues(source["human_resource"], HumanResource);
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
	export class Project {
	    id: number;
	    name: string;
	    description: string;
	    client_id: number;
	    // Go type: time
	    start_date?: any;
	    // Go type: time
	    end_date?: any;
	    status: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    client?: Client;
	    project_resources?: ProjectResource[];
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.client_id = source["client_id"];
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.end_date = this.convertValues(source["end_date"], null);
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.client = this.convertValues(source["client"], Client);
	        this.project_resources = this.convertValues(source["project_resources"], ProjectResource);
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
	export class ProjectListResponse {
	    data: Project[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], Project);
	        this.total = source["total"];
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
	export class ProjectQueryParams {
	    id_in: number[];
	    name: string;
	    name_like: string;
	    description_like: string;
	    client_id: number;
	    client_id_in: number[];
	    status: number;
	    status_in: number[];
	    // Go type: time
	    start_date_gte?: any;
	    // Go type: time
	    start_date_lte?: any;
	    // Go type: time
	    end_date_gte?: any;
	    // Go type: time
	    end_date_lte?: any;
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
	        return new ProjectQueryParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id_in = source["id_in"];
	        this.name = source["name"];
	        this.name_like = source["name_like"];
	        this.description_like = source["description_like"];
	        this.client_id = source["client_id"];
	        this.client_id_in = source["client_id_in"];
	        this.status = source["status"];
	        this.status_in = source["status_in"];
	        this.start_date_gte = this.convertValues(source["start_date_gte"], null);
	        this.start_date_lte = this.convertValues(source["start_date_lte"], null);
	        this.end_date_gte = this.convertValues(source["end_date_gte"], null);
	        this.end_date_lte = this.convertValues(source["end_date_lte"], null);
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
	
	export class ProjectResourceListResponse {
	    data: ProjectResource[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectResourceListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], ProjectResource);
	        this.total = source["total"];
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
	export class ProjectResourceQueryParams {
	    id_in: number[];
	    project_id: number;
	    project_id_in: number[];
	    human_resource_id: number;
	    human_resource_id_in: number[];
	    role: string;
	    role_like: string;
	    allocation_gte?: number;
	    allocation_lte?: number;
	    status: number;
	    status_in: number[];
	    // Go type: time
	    start_date_gte?: any;
	    // Go type: time
	    start_date_lte?: any;
	    // Go type: time
	    end_date_gte?: any;
	    // Go type: time
	    end_date_lte?: any;
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
	        return new ProjectResourceQueryParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id_in = source["id_in"];
	        this.project_id = source["project_id"];
	        this.project_id_in = source["project_id_in"];
	        this.human_resource_id = source["human_resource_id"];
	        this.human_resource_id_in = source["human_resource_id_in"];
	        this.role = source["role"];
	        this.role_like = source["role_like"];
	        this.allocation_gte = source["allocation_gte"];
	        this.allocation_lte = source["allocation_lte"];
	        this.status = source["status"];
	        this.status_in = source["status_in"];
	        this.start_date_gte = this.convertValues(source["start_date_gte"], null);
	        this.start_date_lte = this.convertValues(source["start_date_lte"], null);
	        this.end_date_gte = this.convertValues(source["end_date_gte"], null);
	        this.end_date_lte = this.convertValues(source["end_date_lte"], null);
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

